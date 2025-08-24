package bot

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/genians/endpoint-lab-slack-bot/pkg/slack"
)

type EventHandler interface {
	HandleCommandEvent(payload *slack.SlashCommandEventPayload)
	HandleInteractiveEvent(payload *slack.InteractiveEventPayload)
}

type Bot struct {
	client      *slack.Client
	handler     EventHandler
	eventCh     chan slack.SlackEvent
	reconnectCh chan struct{}
}

func NewBot(appToken string, botToken string, handler EventHandler) *Bot {
	return &Bot{
		client:      &slack.Client{AppToken: appToken, BotToken: botToken},
		handler:     handler,
		eventCh:     make(chan slack.SlackEvent, 1),
		reconnectCh: make(chan struct{}, 1),
	}
}

func (b *Bot) Run(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return nil
		}

		resp, err := b.client.GetWebSocketURL(ctx)
		if err != nil {
			return err
		}

		conn, _, err := websocket.DefaultDialer.Dial(resp.URL, nil)
		if err != nil {
			return err
		}

		connCtx, connCancel := context.WithCancel(ctx)
		go func() {
			<-b.reconnectCh
			connCancel()
		}()
		go func() {
			<-connCtx.Done()
			_ = conn.Close()
		}()

		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			b.runReadLoop(connCtx, conn)
			connCancel()
			wg.Done()
		}()
		go func() {
			b.runEventLoop(connCtx, conn)
			connCancel()
			wg.Done()
		}()
		wg.Wait()
	}
}

func (b *Bot) runReadLoop(ctx context.Context, conn *websocket.Conn) {
	ctx, cancel := context.WithCancel(ctx)

	ping := make(chan struct{}, 1)
	defaultHandler := conn.PingHandler()
	conn.SetPingHandler(func(appData string) error {
		ping <- struct{}{}
		return defaultHandler(appData)
	})

	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		defer ticker.Stop()
		// 서버로부터 일정시간동안 핑을 받지 못하면 네트워크 연결이 끊어진것으로 판단한다.
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				slog.Warn("ping timeout")
				cancel()
				return
			case <-ping:
				ticker.Reset(3 * time.Minute)
			}
		}
	}()

	retry := 0
	for {
		if ctx.Err() != nil {
			slog.Info("stopping read loop")
			return
		}

		if retry >= 3 {
			cancel()
			slog.Error("failed to read after retrying 3 times, exiting")
			return
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			retry += 1
			slog.Error("failed to read websocket", slog.Any("error", err))
			time.Sleep(time.Second * 1)
			continue
		}
		retry = 0

		event, err := slack.UnmarshalSlackEvent(data)
		if err != nil {
			slog.Error("failed to unmarshal slack event", slog.Any("error", err))
			continue
		}

		b.eventCh <- event
	}
}

func (b *Bot) runEventLoop(ctx context.Context, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping event loop")
			return
		case event := <-b.eventCh:
			switch event.EventType() {
			case slack.EventTypeHello:
				e := event.(*slack.HelloEvent)
				slog.Info("receive hello event", slog.Int("count", e.ConnCount))
			case slack.EventTypeDisconnect:
				e := event.(*slack.DisconnectEvent)
				slog.Info("receive disconnect event", slog.String("reason", e.Reason))
				select {
				case b.reconnectCh <- struct{}{}:
				default:
					slog.Warn("reconnectCh is full, skipping reconnect signal")
				}
			case slack.EventTypeSlashCommand:
				e := event.(*slack.SlashCommandEvent)
				// event.AcceptsResponsePayload 값에 따라 socket으로 응답을 할 수도 있지만,
				// 각각 로직을 따로 구분하면 복잡성이 증가하므로 핸들러가 처리하도록 로직을 통일한다.
				go b.handler.HandleCommandEvent(&e.Payload)

				response := map[string]any{
					"envelope_id": e.EnvelopeID,
				}
				if err := conn.WriteJSON(response); err != nil {
					slog.Error("failed to command response", slog.Any("error", err))
				}
			case slack.EventTypeInteractive:
				e := event.(*slack.InteractiveEvent)
				// event.AcceptsResponsePayload 값에 따라 socket으로 응답을 할 수도 있지만,
				// 각각 로직을 따로 구분하면 복잡성이 증가하므로 핸들러가 처리하도록 로직을 통일한다.
				go b.handler.HandleInteractiveEvent(&e.Payload)

				response := map[string]any{
					"envelope_id": e.EnvelopeID,
				}
				if err := conn.WriteJSON(response); err != nil {
					slog.Error("failed to interactive response", slog.Any("error", err))
				}
			}
		}
	}
}

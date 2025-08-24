package slack

import (
	"encoding/json"
	"errors"
)

type EventType string

const (
	EventTypeHello        EventType = "hello"
	EventTypeDisconnect   EventType = "disconnect"
	EventTypeSlashCommand EventType = "slash_commands"
	EventTypeInteractive  EventType = "interactive"
)

type SlackEvent interface {
	EventType() EventType
}

type HelloEvent struct {
	Type      EventType      `json:"type"`
	ConnInfo  ConnectionInfo `json:"connection_info"`
	ConnCount int            `json:"num_connections"`
	DebugInfo DebugInfo      `json:"debug_info"`
}

func (e *HelloEvent) EventType() EventType {
	return e.Type
}

type ConnectionInfo struct {
	AppID string `json:"app_id"`
}

type DebugInfo struct {
	Host           string `json:"host"`
	Started        string `json:"started"`
	BuildNumber    int    `json:"build_number"`
	ConnectionTime int    `json:"approximate_connection_time"`
}

type DisconnectEvent struct {
	Type      EventType `json:"type"`
	Reason    string    `json:"reason"`
	DebugInfo DebugInfo `json:"debug_info"`
}

func (e *DisconnectEvent) EventType() EventType {
	return e.Type
}

type SlashCommandEvent struct {
	Type                   EventType                `json:"type"`
	EnvelopeID             string                   `json:"envelope_id"`
	AcceptsResponsePayload bool                     `json:"accepts_response_payload"`
	Payload                SlashCommandEventPayload `json:"payload"`
}

func (e *SlashCommandEvent) EventType() EventType {
	return e.Type
}

type InteractiveEvent struct {
	Type                   EventType               `json:"type"`
	EnvelopeID             string                  `json:"envelope_id"`
	AcceptsResponsePayload bool                    `json:"accepts_response_payload"`
	Payload                InteractiveEventPayload `json:"payload"`
}

func (i *InteractiveEvent) EventType() EventType {
	return i.Type
}

func UnmarshalSlackEvent(data []byte) (SlackEvent, error) {
	var raw = struct {
		Type EventType `json:"type"`
	}{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var unmarshal = func(data []byte, value SlackEvent) (SlackEvent, error) {
		if err := json.Unmarshal(data, value); err != nil {
			return nil, err
		} else {
			return value, err
		}
	}

	switch raw.Type {
	case EventTypeHello:
		e := HelloEvent{}
		return unmarshal(data, &e)
	case EventTypeDisconnect:
		e := DisconnectEvent{}
		return unmarshal(data, &e)
	case EventTypeSlashCommand:
		e := SlashCommandEvent{}
		return unmarshal(data, &e)
	case EventTypeInteractive:
		e := InteractiveEvent{}
		return unmarshal(data, &e)
	default:
		return nil, errors.New("undefined type: " + string(raw.Type))
	}
}

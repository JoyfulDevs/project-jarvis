package slack

import (
	"github.com/genians/endpoint-lab-slack-bot/pkg/slack/blockkit"
)

// reference
// https://api.slack.com/interactivity/slash-commands
// https://api.slack.com/reference/interaction-payloads/block-actions

// 슬래시 커맨드 이벤트 발생시 전달받는 데이터.
type SlashCommandEventPayload struct {
	// Your Slack app's unique identifier.
	// Use this in conjunction with request signing to verify context for inbound requests.
	AppID string `json:"api_app_id"`
	// The ID of the user who triggered the command.
	UserID string `json:"user_id"`
	// The command that was entered to trigger this request.
	// This value can be useful if you want to use a single Request URL
	// to service multiple slash commands, as it allows you to tell them apart.
	Command string `json:"command"`
	// This is the part of the slash command after the command itself,
	// and it can contain absolutely anything the user might decide to type.
	// It is common to use this text parameter to provide extra context for the command.
	Text string `json:"text"`
	// A temporary webhook URL that you can use to generate message responses.
	ResponseURL string `json:"response_url"`
	// A short-lived ID that will allow your app to open a modal.
	TriggerID string `json:"trigger_id"`

	TeamID       string `json:"team_id"`
	ChannelID    string `json:"channel_id"`
	EnterpriseID string `json:"enterprise_id"`
}

// 상호작용 이벤트를 발생시킨 사용자 정보.
type InteractiveUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	TeamID   string `json:"team_id"`
}

// 상호작용 이벤트 발생시 전달받는 데이터.
type InteractiveEventPayload struct {
	// Helps identify which type of interactive component sent the payload.
	// An interactive element in a block will have a type of `block_actions`,
	// whereas an interactive element in a message attachment will have a type of `interactive_message`.
	Type string `json:"type"`
	// A short-lived ID that can be used to open modals.
	TriggerID string `json:"trigger_id"`
	// 상호작용 이벤트를 발생시킨 사용자 정보.
	User      InteractiveUser      `json:"user"`
	Channel   InteractiveChannel   `json:"channel"`
	Container InteractiveContainer `json:"container"`
	// Contains data from the specific interactive component that was used.
	// App surfaces can contain blocks with multiple interactive components,
	// and each of those components can have multiple values selected by users.
	Actions []InteractiveAction `json:"actions"`
	// A short-lived webhook that can be used to send messages in response to interactions.
	ResponseURL string `json:"response_url"`
}

type InteractiveContainer struct {
	Type        string `json:"type"`
	MessageTs   string `json:"message_ts"`
	ChannelID   string `json:"channel_id"`
	IsEphemeral bool   `json:"is_ephemeral"`
}

// 상호작용 이벤트를 발생시킨 액션에 대한 데이터.
type InteractiveAction struct {
	Type           blockkit.ElementType      `json:"type"`
	Timestamp      float64                   `json:"action_ts,string"`
	ActionID       string                    `json:"action_id"`
	BlockID        string                    `json:"block_id"`
	Value          string                    `json:"value"`
	SelectedOption InteractiveSelectedOption `json:"selected_option"`
}

type InteractiveChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type InteractiveSelectedOption struct {
	Value string              `json:"value"`
	Text  blockkit.TextObject `json:"text"`
}

type InteractiveResponseType string

const (
	// 채널에 포함된 모두가 볼 수 있도록 메시지를 전달한다.
	InChannel InteractiveResponseType = "in_channel"
	// 상호작용한 사용자만 볼 수 있도록 메시지를 전달한다.
	Ephemeral InteractiveResponseType = "ephemeral"
)

// 이벤트에 대한 응답으로 전달하는 데이터.
type InteractiveResponsePayload struct {
	// 채널 모두가 볼 수 있도록 메시지를 전달할지, 상호작용한 사용자만 볼 수 있도록 메시지를 전달할지 선택한다.
	// 기본값은 상호작용한 사용자만 볼 수 있도록 설정된다.
	ResponseType InteractiveResponseType `json:"response_type,omitempty"`
	// blocks 필드를 사용하지 않을 경우에는 이 필드를 메시지의 본문으로 사용한다.
	// blocks 필드에 문제가 있는 경우 대신 표시되므로 사용하는 것이 권장된다.
	Text string `json:"text"`
	// text 필드를 사용하지 않을 경우에는 이 필드를 메시지의 본문으로 사용한다.
	Blocks []blockkit.SlackBlock `json:"blocks,omitempty"`
	// The ID of another un-threaded message to reply to.
	ThreadTimestamp string `json:"thread_ts,omitempty"`
	// 원본 메시지를 삭제할지 여부를 결정한다.
	DeleteOriginal bool `json:"delete_original"`
	// 원본 메시지를 대체할지 여부를 결정한다.
	ReplaceOriginal bool `json:"replace_original"`
}

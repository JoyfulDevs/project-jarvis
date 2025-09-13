package slack

import (
	"strconv"
	"strings"

	"github.com/joyfuldevs/project-jarvis/pkg/slack/blockkit"
)

type APIResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type GetWebSocketURLResponse struct {
	APIResponse

	URL string `json:"url"`
}

type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type ChannelType string

const (
	PublicChannel      ChannelType = "public_channel"
	PrivateChannel     ChannelType = "private_channel"
	DirectMessage      ChannelType = "im"
	GroupDirectMessage ChannelType = "mpim"
)

type ListChannelsRequest struct {
	// Paginate through collections of data by setting the cursor parameter to a
	// next_cursor attribute returned by a previous request's response_metadata.
	// Default value fetches the first "page" of the collection. See pagination for more detail.
	Cursor string `json:"cursor,omitempty"`
	// Set to true to exclude archived channels from the list.
	ExcludeArchived bool `json:"exclude_archived"`
	// The maximum number of items to return.
	// Fewer than the requested number of items may be returned,
	// even if the end of the list hasn't been reached. Must be an integer under 1000.
	Limit int `json:"limit,omitempty"`
	// encoded team id to list channels in, required if token belongs to org-wide app.
	TeamID string `json:"team_id,omitempty"`
	// Mix and match channel types by providing a comma-separated list of
	// any combination of `public_channel`, `private_channel`, `mpim`, `im`
	Types ChannelType `json:"types,omitempty"`
}

func (r *ListChannelsRequest) URLParams() string {
	var builder strings.Builder
	builder.WriteString("?")
	if len(r.Cursor) > 0 {
		builder.WriteString("cursor=")
		builder.WriteString(r.Cursor)
		builder.WriteString("&")
	}
	if r.ExcludeArchived {
		builder.WriteString("exclude_archived=true")
		builder.WriteString("&")
	}
	if r.Limit > 0 {
		builder.WriteString("limit=")
		builder.WriteString(strconv.Itoa(r.Limit))
		builder.WriteString("&")
	}
	if len(r.TeamID) > 0 {
		builder.WriteString("team_id=")
		builder.WriteString(r.TeamID)
		builder.WriteString("&")
	}
	if len(r.Types) > 0 {
		builder.WriteString("types=")
		builder.WriteString(string(r.Types))
	}
	return builder.String()
}

type ListChannelsResponse struct {
	APIResponse

	Channels []ConversationObject `json:"channels"`
	Metadata ResponseMetadata     `json:"response_metadata"`
}

type MessageParseType string

const (
	MessageParseTypeNone     MessageParseType = "none"
	MessageParseTypeMarkdown MessageParseType = "mrkdwn"
	MessageParseTypeFull     MessageParseType = "full"
)

type PostMessageRequest struct {
	// An encoded ID that represents a channel, private group, or IM channel to send the message to.
	Channel string `json:"channel"`
	// How this field works and whether it is required depends on other fields you use in your API call.
	Text string `json:"text,omitempty"`
	// A JSON-based array of structured blocks, presented as a URL-encoded string.
	Blocks []blockkit.SlackBlock `json:"blocks,omitempty"`
	// URL to an image to use as the icon for this message.
	IconURL string `json:"icon_url,omitempty"`
	// Emoji to use as the icon for this message. Overrides icon_url.
	IconEmoji string `json:"icon_emoji,omitempty"`
	// Find and link user groups. No longer supports linking individual users.
	// use syntax shown in Mentioning Users instead.
	LinkNames bool `json:"link_names,omitempty"`
	// Disable Slack markup parsing by setting to false. Enabled by default.
	Markdown bool `json:"mrkdwn,omitempty"`
	// Change how messages are treated.
	//
	// By default, URLs will be hyperlinked. Set parse to `none` to remove the hyperlinks.
	// The behavior of parse is different for text formatted with `mrkdwn`.
	//
	// By default, or when parse is set to `none`, `mrkdwn` formatting is implemented.
	// To ignore mrkdwn formatting, set parse to `full`.
	Parse MessageParseType `json:"parse,omitempty"`
	// Used in conjunction with thread_ts and indicates whether reply should be made visible to everyone
	// in the channel or conversation. Defaults to false.
	ReplyBroadcast bool `json:"reply_broadcast,omitempty"`
	// Provide another message's ts value to make this message a reply.
	// Avoid using a reply's ts value. use its parent instead.
	ThreadTimestamp float64 `json:"thread_ts,omitempty,string"`
	// Pass true to enable unfurling of primarily text-based content.
	UnfurlLinks bool `json:"unfurl_links,omitempty"`
	// Pass false to disable unfurling of media content.
	UnfurlMedia bool `json:"unfurl_media,,omitempty"`
	// Set your bot's user name.
	Username string `json:"username,omitempty"`
}

type PostMessageResponse struct {
	APIResponse

	Channel   string  `json:"channel"`
	Timestamp float64 `json:"ts,string"`
}

type ListMessagesRequest struct {
	// Conversation ID to fetch history for.
	Channel string `json:"channel"`
	// Paginate through collections of data by setting the `cursor` parameter to a
	// `next_cursor` attribute returned by a previous request's `response_metadata`.
	Cursor string `json:"cursor,omitempty"`
	// Return all metadata associated with this message.
	IncludeAllMetadata bool `json:"include_all_metadata"`
	// Include messages with latest or oldest timestamp in results.
	// Ignored unless either timestamp is specified.
	Inclusive bool `json:"inclusive"`
	// The maximum number of items to return.
	// Fewer than the requested number of items may be returned,
	// even if the end of conversation history hasn't been reached.
	// Maximum of 999.
	Limit int `json:"limit,omitempty"`
	// Only messages before this Unix timestamp will be included in result.
	Latest float64 `json:"latest,string,omitempty"`
	// Only messages after this Unix timestamp will be included in result.
	Oldest float64 `json:"oldest,string,omitempty"`
}

type ListMessagesResponse struct {
	APIResponse

	HasMore  bool `json:"has_more"`
	Messages []struct {
		AppID     string  `json:"app_id"`
		User      string  `json:"user"`
		Text      string  `json:"text"`
		Timestamp float64 `json:"ts,string"`
	} `json:"messages"`
	Metadata ResponseMetadata `json:"response_metadata"`
}

type ListRepliesRequest struct {
	Channel   string  `json:"channel"`
	Timestamp float64 `json:"ts,string"`
}

type ListRepliesResponse struct {
	APIResponse

	HasMore  bool `json:"has_more"`
	Messages []struct {
		AppID     string  `json:"app_id"`
		User      string  `json:"user"`
		Text      string  `json:"text"`
		Timestamp float64 `json:"ts,string"`
	} `json:"messages"`
	Metadata ResponseMetadata `json:"response_metadata"`
}

type UserProfile struct {
	// The user's title.
	Title string `json:"title,omitempty"`
	// The display name the user has chosen to identify themselves by in their workspace profile.
	// Do not use this field as a unique identifier for a user, as it may change at any time.
	// Instead, use id and team_id in concert.
	DisplayName string `json:"display_name,omitempty"`
	// The user's first and last name. Updating this field will update first_name and last_name.
	// If only one name is provided, the value of last_name will be cleared.
	RealName  string `json:"real_name,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

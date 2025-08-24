package slack

// A conversation object contains information about a channel-like thing in Slack.
// It might be a public channel, a private channel, a direct message,
// a multi-person direct message, or a huddle.
// You'll find all of these objects throughout the Conversations API.
type ConversationObject struct {
	// The conversation ID.
	ID string `json:"id"`
	// Indicates the name of the channel-like thing, without a leading hash sign.
	Name string `json:"name"`
	// A normalized string of the name.
	NormalizedName string `json:"name_normalized"`
	// (unix-timestamp seconds)
	// Timestamp of when the conversation was created.
	Created int `json:"created"`
	// The ID of the member that created this conversation.
	Creator string `json:"creator"`
	// The timestamp, in milliseconds, when the channel settings were updated.
	// for example, the "topic" or "description" of the channel changed.
	Updated int `json:"updated"`
	// Provides information about the channel topic.
	Topic struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"topic"`
	// Provides information about the channel purpose.
	Purpose struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"purpose"`
	// A list of prior names the channel has used.
	PreviousNames []string `json:"previous_names"`
	// Indicates a conversation is archived, frozen in time.
	IsArchived bool `json:"is_archived"`
	// Indicates whether a conversation is a channel.
	// Private channels created before March 2021 (with IDs that begin with G) will return false,
	// and is_group will be true instead. Use is_private to determine whether a channel is private or public.
	IsChannel bool `json:"is_channel"`
	// Means the channel is the workspace's "general" discussion channel (even if it may not be named #general).
	// That might be important to your app because almost every user is a member.
	IsGeneral bool `json:"is_general"`
	// Means the channel is a private channel created before March 2021. is_private will also be true.
	IsGroup bool `json:"is_group"`
	// Means the conversation is a direct message between two distinguished individuals or a user and a bot.
	// is_private will also be true
	IsDM bool `json:"is_im"`
	// Indicates whether the user, bot user or Slack app associated with the token making
	// the API call is itself a member of the conversation.
	IsMember bool `json:"is_member"`
	// Represents an unnamed private conversation between multiple users. is_private will also be true.
	IsGroupDM bool `json:"is_mpim"`
	// Means the conversation is privileged between two or more members.
	// Ensure that you meet their privacy expectations.
	IsPrivate bool `json:"is_private"`
	// Means the conversation can't be written to by the user performing the API call.
	IsReadOnly bool `json:"is_read_only"`
	// Means the conversation is in some way shared between multiple workspaces.
	// Look for is_ext_shared and is_org_shared to learn which kind it is, and if that matters, act accordingly.
	IsShared bool `json:"is_shared"`
	// Means the conversation can't be written to by the user performing the API call,
	// except to reply to messages in the channel.
	IsThreatOnly bool `json:"is_thread_only"`
}

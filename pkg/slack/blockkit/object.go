package blockkit

// Reference
// https://api.slack.com/reference/block-kit/composition-objects

type TextObjectType string

const (
	TextTypePlainText TextObjectType = "plain_text"
	TextTypeMarkdown  TextObjectType = "mrkdwn"
)

// Defines an object containing some text.
type TextObject struct {
	// The formatting to use for this text object. Can be one of plain_text or mrkdwn.
	Type TextObjectType `json:"type"`
	// The text for the block.
	// This field accepts any of the standard text formatting markup when type is mrkdwn.
	// The minimum length is 1 and maximum length is 3000 characters.
	Text string `json:"text"`
	// Indicates whether emojis in a text field should be escaped into the colon emoji format.
	// This field is only usable when type is plain_text.
	Emoji bool `json:"emoji,omitempty"`
	// When set to false,
	// URLs will be auto-converted into links, conversation names will be link-ified,
	// and certain mentions will be automatically parsed.
	//
	// When set to true,
	// Slack will continue to process all markdown formatting and manual parsing strings,
	// but it won’t modify any plain-text content.
	Verbatim bool `json:"verbatim,omitempty"`
}

// Defines a dialog that adds a confirmation step to interactive elements.
type ConfirmDialogObject struct {
	// A `plain_text` text object that defines the dialog's title.
	// Maximum length for this field is 100 characters.
	Title TextObject `json:"title"`
	// A `plain_text` text object that defines the explanatory text that appears in the confirm dialog.
	// Maximum length for the text in this field is 300 characters.
	Text TextObject `json:"text"`
	// A `plain_text` text object to define the text of the button that confirms the action.
	// Maximum length for the text in this field is 30 characters.
	Confirm TextObject `json:"confirm"`
	// A `plain_text` text object to define the text of the button that cancels the action.
	// Maximum length for the text in this field is 30 characters.
	Deny TextObject `json:"deny"`
	// Defines the color scheme applied to the confirm button.
	Style ButtonStyle `json:"style,omitempty"`
}

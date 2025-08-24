package blockkit

import (
	"encoding/json"
)

// Reference
// https://api.slack.com/reference/block-kit/blocks

type BlockType string

const (
	BlockTypeSection BlockType = "section"
	BlockTypeHeader  BlockType = "header"
	BlockTypeAction  BlockType = "actions"
	BlockTypeDivider BlockType = "divider"
)

type SlackBlock interface {
	BlockType() BlockType
}

// Displays a larger-sized text block.
type HeaderBlock struct {
	// A unique identifier for a block. If not specified, one will be generated.
	// Maximum length for this field is 255 characters.
	// block_id should be unique for each message and each iteration of a message.
	// If a message is updated, use a new block_id.
	BlockID string `json:"block_id,omitempty"`
	// The text for the block, in the form of a `plain_text` text object.
	// Maximum length for the text in this field is 150 characters.
	Text TextObject `json:"text"`
}

func NewHeaderBlock(text string) *HeaderBlock {
	return &HeaderBlock{
		Text: TextObject{
			Type:  TextTypePlainText,
			Text:  text,
			Emoji: true,
		},
	}
}

func (h *HeaderBlock) BlockType() BlockType {
	return BlockTypeHeader
}

func (h *HeaderBlock) MarshalJSON() ([]byte, error) {
	raw := struct {
		HeaderBlock
		Type BlockType `json:"type"`
	}{
		HeaderBlock: *h,
		Type:        h.BlockType(),
	}
	return json.Marshal(raw)
}

// Displays text, possibly alongside block elements.
type SectionBlock struct {
	// A unique identifier for a block. If not specified, one will be generated.
	// Maximum length for this field is 255 characters.
	// block_id should be unique for each message and each iteration of a message.
	// If a message is updated, use a new block_id.
	BlockID string `json:"block_id,omitempty"`
	// The text for the block, in the form of a text object.
	// Minimum length for the text in this field is 1 and maximum length is 3000 characters.
	// This field is not required if a valid array of fields objects is provided instead.
	Text *TextObject `json:"text,omitempty"`
	// Required if no text is provided.
	// An array of text objects. Any text objects included with fields will be rendered
	// in a compact format that allows for 2 columns of side-by-side text.
	// Maximum number of items is 10. Maximum length for the text in each item is 2000 characters.
	Fields []TextObject `json:"fields,omitempty"`
	// One of the compatible element objects noted above.
	// Be sure to confirm the desired element works with section.
	Accessory SlackBlockElement `json:"accessory,omitempty"`
	// Whether or not this section block's text should always expand when rendered.
	// If false or not provided, it may be rendered with a 'see more' option to expand and show the full text.
	Expand bool `json:"expand,omitempty"`
}

func (s *SectionBlock) BlockType() BlockType {
	return BlockTypeSection
}

func (s *SectionBlock) MarshalJSON() ([]byte, error) {
	raw := struct {
		SectionBlock
		Type BlockType `json:"type"`
	}{
		SectionBlock: *s,
		Type:         s.BlockType(),
	}
	return json.Marshal(raw)
}

// Holds multiple interactive elements.
type ActionBlock struct {
	// A unique identifier for a block. If not specified, one will be generated.
	// Maximum length for this field is 255 characters.
	// block_id should be unique for each message and each iteration of a message.
	// If a message is updated, use a new block_id.
	BlockID string `json:"block_id,omitempty"`
	// An array of interactive element objects - buttons, select menus, overflow menus, or date pickers.
	// There is a maximum of 25 elements in each action block.
	Elements []SlackBlockElement `json:"elements,omitempty"`
}

func (a *ActionBlock) BlockType() BlockType {
	return BlockTypeAction
}

func (a *ActionBlock) MarshalJSON() ([]byte, error) {
	raw := struct {
		ActionBlock
		Type BlockType `json:"type"`
	}{
		ActionBlock: *a,
		Type:        a.BlockType(),
	}
	return json.Marshal(raw)
}

type DividerBlock struct {
}

func (d *DividerBlock) BlockType() BlockType {
	return BlockTypeDivider
}

func (d *DividerBlock) MarshalJSON() ([]byte, error) {
	raw := struct {
		Type BlockType `json:"type"`
	}{
		Type: d.BlockType(),
	}
	return json.Marshal(raw)
}

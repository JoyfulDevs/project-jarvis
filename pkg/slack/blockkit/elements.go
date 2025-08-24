package blockkit

import (
	"encoding/json"
)

// Reference
// https://api.slack.com/reference/block-kit/block-elements

type ElementType string

const (
	ElementTypeImage  ElementType = "image"
	ElementTypeButton ElementType = "button"
	ElementTypeSelect ElementType = "static_select"
)

// Block elements can be added to certain app surfaces and used within certain block types.
type SlackBlockElement interface {
	ElementType() ElementType
}

// SelectElement 타입 추가
type SelectElement struct {
	ActionID    string              `json:"action_id"`
	Type        ElementType         `json:"type"`
	Placeholder TextObject          `json:"placeholder"`
	Options     []OptionBlockObject `json:"options"`
}

type OptionBlockObject struct {
	Text  TextObject `json:"text"`
	Value string     `json:"value"`
}

func (s *SelectElement) ElementType() ElementType {
	return ElementTypeSelect
}

func (s *SelectElement) MarshalJSON() ([]byte, error) {
	raw := struct {
		ActionID    string              `json:"action_id"`
		Type        ElementType         `json:"type"`
		Placeholder TextObject          `json:"placeholder"`
		Options     []OptionBlockObject `json:"options"`
	}{
		ActionID:    s.ActionID,
		Type:        s.Type,
		Placeholder: s.Placeholder,
		Options:     s.Options,
	}
	return json.Marshal(raw)
}

// Displays an image as part of a larger block of content.
type ImageElement struct {
	// A plain-text summary of the image. This should not contain any markup.
	AltText string `json:"alt_text"`
	// The URL for a publicly hosted image. You must provide either an image_url or slack_file.
	// Maximum length for this field is 3000 characters.
	ImageURL string `json:"image_url"`
}

func (i *ImageElement) ElementType() ElementType {
	return ElementTypeImage
}

func (i *ImageElement) MarshalJSON() ([]byte, error) {
	raw := struct {
		ImageElement
		Type ElementType `json:"type"`
	}{
		ImageElement: *i,
		Type:         i.ElementType(),
	}
	return json.Marshal(raw)
}

// Allows users a direct path to performing basic actions.
type ButtonElement struct {
	// An identifier for this action.
	// You can use this when you receive an interaction payload to identify the source of the action.
	// Should be unique among all other action_ids in the containing block.
	// Maximum length is 255 characters.
	ActionID string `json:"action_id,omitempty"`
	// A text object that defines the button's text.
	// Can only be of type: `plain_text`. text may truncate with ~30 characters.
	// Maximum length for the text in this field is 75 characters.
	Text TextObject `json:"text"`
	// Decorates buttons with alternative visual color schemes.
	// Use this option with restraint.
	Style ButtonStyle `json:"style,omitempty"`
	// The value to send along with the interaction payload. Maximum length is 2000 characters.
	Value string `json:"value,omitempty"`
	// A URL to load in the user's browser when the button is clicked.
	// Maximum length is 3000 characters. If you're using url,
	// you'll still receive an interaction payload and will need to send an acknowledgement response.
	URL string `json:"url,omitempty"`
	// A label for longer descriptive text about a button element.
	// This label will be read out by screen readers instead of the button text object.
	// Maximum length is 75 characters.
	Label string `json:"accessibility_label,omitempty"`
}

func (b *ButtonElement) ElementType() ElementType {
	return ElementTypeButton
}

func (b *ButtonElement) MarshalJSON() ([]byte, error) {
	raw := struct {
		ButtonElement
		Type ElementType `json:"type"`
	}{
		ButtonElement: *b,
		Type:          b.ElementType(),
	}
	return json.Marshal(raw)
}

package block

import (
	"encoding/json"
	"fmt"
)

// A list block (unordered)
type ListItem struct {
	BaseBlock
}

func (l *ListItem) Decode(enc interface{}) error {
	return l.decodeBlock(enc)
}

func (l *ListItem) AsHtml() string {
	return fmt.Sprintf("<li>%s</li>", l.FormatHtmlText())
}

func (l *ListItem) ParentHtmlTag() string {
	return "ul"
}

// A list block (ordered)
type OrderedListItem struct {
	BaseBlock
}

func (l *OrderedListItem) Decode(enc interface{}) error {
	return l.decodeBlock(enc)
}

func (l *OrderedListItem) AsHtml() string {
	return fmt.Sprintf("<li>%s</li>", l.FormatHtmlText())
}

func (l *OrderedListItem) ParentHtmlTag() string {
	return "ol"
}

func (l *OrderedListItem) MarshalJSON() ([]byte, error) {
	type Alias OrderedListItem
	return json.Marshal(&struct {
		HTML string `json:"html"`
		*Alias
	}{
		HTML:  l.AsHtml(),
		Alias: (*Alias)(l),
	})
}

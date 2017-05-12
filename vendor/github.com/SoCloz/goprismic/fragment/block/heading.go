package block

import (
	"encoding/json"
	"fmt"
)

// A heading (h1, h2, h3, h4, h5, h6) block
type Heading struct {
	BaseBlock
}

func (h *Heading) Decode(enc interface{}) error {
	return h.decodeBlock(enc)
}

func (h *Heading) AsHtml() string {
	switch h.Type {
	case "heading1":
		return fmt.Sprintf("<h1>%s</h1>", h.FormatHtmlText())
	case "heading2":
		return fmt.Sprintf("<h2>%s</h2>", h.FormatHtmlText())
	case "heading3":
		return fmt.Sprintf("<h3>%s</h3>", h.FormatHtmlText())
	case "heading4":
		return fmt.Sprintf("<h4>%s</h4>", h.FormatHtmlText())
	case "heading5":
		return fmt.Sprintf("<h5>%s</h5>", h.FormatHtmlText())
	case "heading6":
		return fmt.Sprintf("<h6>%s</h6>", h.FormatHtmlText())
	}
	return ""
}

func (h *Heading) ParentHtmlTag() string {
	return ""
}

func (h *Heading) MarshalJSON() ([]byte, error) {
	type Alias Heading
	return json.Marshal(&struct {
		HTML string `json:"html"`
		*Alias
	}{
		HTML:  h.AsHtml(),
		Alias: (*Alias)(h),
	})
}

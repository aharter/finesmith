package fragment

import (
	"encoding/json"
	"fmt"

	"github.com/SoCloz/goprismic/fragment/link"
)

// A link fragment
type Link struct {
	Link link.Link
}

func (l *Link) Decode(t string, enc interface{}) error {
	var err error
	l.Link, err = link.Decode(t, enc)
	return err
}

func (l *Link) AsText() string {
	return l.Link.GetUrl()
}

func (l *Link) AsHtml() string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", l.Link.GetUrl(), l.Link.GetText())
}

func (l *Link) ResolveLinks(r link.Resolver) {
	l.Link.Resolve(r)
}

func (l *Link) MarshalJSON() ([]byte, error) {
	switch t := l.Link.(type) {
	default:
		return json.Marshal(&struct {
			Children map[string]Interface `json:"children"`
			HTML     string               `json:"html"`
			Text     string               `json:"text"`
			URL      string               `json:"url"`
		}{
			Children: l.Children,
			HTML:     l.AsHtml(),
			Text:     l.AsHtml(),
			URL:      l.Link.GetUrl(),
		})
	case *link.DocumentLink:
		return json.Marshal(&struct {
			Children map[string]Interface `json:"children"`
			HTML     string               `json:"html"`
			Text     string               `json:"text"`
			URL      string               `json:"url"`
			UID      string               `json:"uid"`
		}{
			Children: l.Children,
			HTML:     l.AsHtml(),
			Text:     l.AsHtml(),
			URL:      l.Link.GetUrl(),
			UID:      t.Document.UID,
		})
	}
}

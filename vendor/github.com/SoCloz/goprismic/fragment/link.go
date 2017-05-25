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
	marshalLink := func(l *Link) ([]byte, error) {
		return json.Marshal(&struct {
			HTML string `json:"html"`
			Text string `json:"text"`
			URL  string `json:"url"`
		}{
			HTML: l.AsHtml(),
			Text: l.AsText(),
			URL:  l.Link.GetUrl(),
		})
	}

	switch t := l.Link.(type) {
	default:
		return json.Marshal(t)
	case *link.WebLink:
		return marshalLink(l)
	case *link.MediaLink:
		return marshalLink(l)
	case *link.DocumentLink:
		return json.Marshal(&struct {
			HTML string `json:"html"`
			Text string `json:"text"`
			URL  string `json:"url"`
			UID  string `json:"uid"`
		}{
			HTML: l.AsHtml(),
			Text: l.AsText(),
			URL:  l.Link.GetUrl(),
			UID:  t.Document.UID,
		})
	}

}

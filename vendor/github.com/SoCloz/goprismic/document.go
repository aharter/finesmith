package goprismic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SoCloz/goprismic/fragment"
	"github.com/SoCloz/goprismic/fragment/block"
	"github.com/SoCloz/goprismic/fragment/link"
)

type DocumentLink struct {
	*link.DocumentLink
	Children fragment.Fragments
}

func (l *DocumentLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		HTML     string             `json:"html"`
		Text     string             `json:"text"`
		URL      string             `json:"url"`
		UID      string             `json:"uid"`
		Children fragment.Fragments `json:"children"`
	}{
		HTML:     fmt.Sprintf("<a href=\"%s\">%s</a>", l.GetUrl(), l.GetText()),
		Text:     l.GetText(),
		URL:      l.GetUrl(),
		UID:      l.Document.UID,
		Children: l.Children,
	})
}

// A document is made of fragments of various types
type Document struct {
	Id        string        `json:"id"`
	Type      string        `json:"type"`
	Href      string        `json:"href"`
	UID       string        `json:"uid"`
	Tags      []string      `json:"tags"`
	Slugs     []string      `json:"slugs"`
	Fragments fragment.Tree `json:"data"`
	URL       string        `json:"url,omitempty"`
}

// Returns the document slug
func (d *Document) GetSlug() string {
	return d.Slugs[0]
}

// Tests if the document has a slug
func (d *Document) HasSlug(slug string) bool {
	for _, v := range d.Slugs {
		if v == slug {
			return true
		}
	}
	return false
}

// Resolves links
func (d *Document) ResolveLinks(r link.Resolver) {
	frags, found := d.Fragments[d.Type]
	if !found {
		return
	}
	for _, list := range frags {
		for k := range list {
			list[k].ResolveLinks(r)
		}
	}
	d.URL = r(d.AsDocumentLink())
}

// Returns the list of fragments of a certain name
func (d *Document) GetFragments(field string) (fragment.List, bool) {
	frags, found := d.Fragments[d.Type]
	if !found {
		return nil, false
	}
	f, found := frags[field]
	return f, found
}

// Returns the nth fragment of a certain name
func (d *Document) GetFragmentAt(field string, index int) (fragment.Interface, bool) {
	frags, found := d.GetFragments(field)
	if !found || len(frags) <= index {
		return nil, false
	}
	return frags[index], true
}

// Returns the first fragment of a certain name
func (d *Document) GetFragment(field string) (fragment.Interface, bool) {
	return d.GetFragmentAt(field, 0)
}

// Returns an image fragment (the first found)
func (d *Document) GetImageFragment(field string) (*fragment.Image, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	i, ok := f.(*fragment.Image)
	if !ok {
		return nil, false
	}
	return i, true
}

// Returns a structured text fragment (returns the first found)
func (d *Document) GetStructuredTextFragment(field string) (*fragment.StructuredText, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.StructuredText)
	if !ok {
		return nil, false
	}
	return st, true
}

// Returns the list of blocks of a structured text fragment
func (d *Document) GetStructuredTextBlocks(field string) ([]block.Block, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.StructuredText)
	if !ok {
		return nil, false
	}
	return []block.Block(*st), true
}

// Returns a color fragment (the first found)
func (d *Document) GetColorFragment(field string) (*fragment.Color, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	c, ok := f.(*fragment.Color)
	if !ok {
		return nil, false
	}
	return c, true
}

// Returns a color value (the first found)
func (d *Document) GetColor(field string) (string, bool) {
	c, found := d.GetColorFragment(field)
	if !found {
		return "", false
	}
	return string(*c), true
}

// Returns a number fragment (the first found)
func (d *Document) GetNumberFragment(field string) (*fragment.Number, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	n, ok := f.(*fragment.Number)
	if !ok {
		return nil, false
	}
	return n, true
}

// Returns a number value (the first found)
func (d *Document) GetNumber(field string) (float64, bool) {
	n, found := d.GetNumberFragment(field)
	if !found {
		return float64(0), false
	}
	return float64(*n), true
}

// Returns a text fragment (the first found)
func (d *Document) GetTextFragment(field string) (*fragment.Text, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	t, ok := f.(*fragment.Text)
	if !ok {
		return nil, false
	}
	return t, true
}

// Returns a text value (the first found)
func (d *Document) GetText(field string) (string, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return "", false
	}
	return f.AsText(), true
}

// Returns the boolean representation of a fragment (the first found)
func (d *Document) GetBool(field string) (bool, bool) {
	t, found := d.GetText(field)
	return (t == "yes" || t == "true"), found
}

// Returns a date fragment (the first found)
func (d *Document) GetDateFragment(field string) (*fragment.Date, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	t, ok := f.(*fragment.Date)
	if !ok {
		return nil, false
	}
	return t, true
}

// Returns a date value (the first found)
func (d *Document) GetDate(field string) (time.Time, bool) {
	t, found := d.GetDateFragment(field)
	if !found {
		return time.Time{}, false
	}
	return time.Time(*t), true
}

// Returns a link fragment (the first found)
func (d *Document) GetLinkFragment(field string) (*fragment.Link, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	l, ok := f.(*fragment.Link)
	if !ok {
		return nil, false
	}
	return l, true
}

// Returns a geopoint fragment (returns the first found)
func (d *Document) GetGeoPointFragment(field string) (*fragment.GeoPoint, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	gp, ok := f.(*fragment.GeoPoint)
	if !ok {
		return nil, false
	}
	return gp, true
}

func (d *Document) AsDocumentLink() *DocumentLink {
	docLink := &link.DocumentLink{
		Document: struct {
			Id   string
			Type string
			Slug string
			UID  string
		}{
			Id:   d.Id,
			Type: d.Type,
			Slug: d.GetSlug(),
			UID:  d.UID,
		},
		Url:      d.URL,
		IsBroken: false,
	}

	l := &DocumentLink{DocumentLink: docLink, Children: make(fragment.Fragments)}
	for _, parentPageResult := range d.Fragments {
		for name, fragmentList := range parentPageResult {
			l.Children[name] = fragmentList
		}
	}

	return l
}

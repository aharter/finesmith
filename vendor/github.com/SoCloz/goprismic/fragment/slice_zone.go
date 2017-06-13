package fragment

import (
	"fmt"

	"github.com/SoCloz/goprismic/fragment/block"
	"github.com/SoCloz/goprismic/fragment/link"
)

type Slice interface {
	Interface
	GetFirstParagraph() (*block.Paragraph, bool)
	GetFirstImage() (*block.Image, bool)
}

// SliceZone fragment is a list of Slice fragments
type SliceZone []Slice

func NewSliceZone(enc interface{}) (*SliceZone, error) {
	dec, ok := enc.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%#v is not a slice", enc)
	}
	sz := make(SliceZone, 0, len(dec))
	return &sz, nil
}

func (sz *SliceZone) Decode(_ string, enc interface{}) error {
	dec := enc.([]interface{})
	for _, v := range dec {
		dec, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%+v is not a map", v)
		}
		var slice Slice
		val, ok := dec["repeat"].(bool)
		if ok && val {
			panic("Composite slice not implemented.")
		} else {
			slice = new(SimpleSlice)
		}
		slice.Decode("", v)
		*sz = append(*sz, slice)
	}
	return nil
}

// AsHtml returns the html formatted fragment content
func (sz SliceZone) AsHtml() string {
	html := ""
	for _, v := range sz {
		html += v.AsHtml()
	}
	return html
}

// AsText returns fragment content as plain text.
func (sz SliceZone) AsText() string {
	text := ""
	for _, v := range sz {
		if text != "" {
			text += "\n"
		}
		text += v.AsText()
	}
	return text
}

// Returns the first paragraph fragment
func (sz SliceZone) GetFirstParagraph() (*block.Paragraph, bool) {
	for _, v := range sz {
		p, ok := v.GetFirstParagraph()
		if ok {
			return p, true
		}
	}
	return nil, false
}

// Returns the first image fragment
func (sz SliceZone) GetFirstImage() (*block.Image, bool) {
	for _, v := range sz {
		i, ok := v.GetFirstImage()
		if ok {
			return i, true
		}
	}
	return nil, false
}

func (sz SliceZone) ResolveLinks(r link.Resolver) {
	for _, v := range sz {
		v.ResolveLinks(r)
	}
}

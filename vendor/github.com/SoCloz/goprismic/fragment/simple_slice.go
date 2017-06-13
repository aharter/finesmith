package fragment

import (
	"fmt"

	"github.com/SoCloz/goprismic/fragment/block"
	"github.com/SoCloz/goprismic/fragment/link"
)

// SimpleSlice wraps a child (fragment or group)
type SimpleSlice struct {
	SliceType string      `json:"sliceType"`
	Label     string      `json:"label"`
	Child     interface{} `json:"children"`
}

func (ss *SimpleSlice) Decode(t string, enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		fmt.Printf("%+v is not a map", enc)
		return fmt.Errorf("%+v is not a map", enc)
	}

	ss.SliceType = dec["slice_type"].(string)
	ss.Label = ""
	if v, ok := dec["slice_label"]; ok && v != nil {
		ss.Label = v.(string)
	}

	var err error
	child := dec["value"].(map[string]interface{})
	if ss.Child, err = ExtDecode(child["type"].(string), child["value"]); err != nil {
		return err
	}
	return nil
}

// AsHtml returns the html formatted fragment content
func (ss *SimpleSlice) AsHtml() string {
	classes := "slice"
	if ss.Label != "" {
		classes = fmt.Sprintf("%s %s", classes, ss.Label)
	}
	childHTML := ""
	if html, ok := ss.Child.(block.Block); ok {
		childHTML = html.AsHtml()
	}
	if html, ok := ss.Child.(Interface); ok {
		childHTML = html.AsHtml()
	}

	return fmt.Sprintf("<div data-slicetype=\"%s\" class=\"%s\">%s</div>", ss.SliceType, classes, childHTML)
}

// Formats the fragment content as text
func (ss *SimpleSlice) AsText() string {
	if block, ok := ss.Child.(block.Block); ok {
		return block.AsText()
	}
	if i, ok := ss.Child.(Interface); ok {
		return i.AsText()
	}
	return ""
}

// Returns the first paragraph fragment
func (ss *SimpleSlice) GetFirstParagraph() (*block.Paragraph, bool) {
	if p, ok := ss.Child.(*block.Paragraph); ok {
		return p, true
	}
	if st, ok := ss.Child.(*StructuredText); ok {
		return st.GetFirstParagraph()
	}
	if s, ok := ss.Child.(Slice); ok {
		return s.GetFirstParagraph()
	}

	return nil, false
}

// Returns the first image fragment
func (ss *SimpleSlice) GetFirstImage() (*block.Image, bool) {
	if b, ok := ss.Child.(*block.Image); ok {
		return b, true
	}
	if st, ok := ss.Child.(*StructuredText); ok {
		return st.GetFirstImage()
	}
	if s, ok := ss.Child.(Slice); ok {
		return s.GetFirstImage()
	}

	return nil, false
}

func (ss *SimpleSlice) ResolveLinks(r link.Resolver) {
	if b, ok := ss.Child.(block.Block); ok {
		b.ResolveLinks(r)
	}
	if i, ok := ss.Child.(Interface); ok {
		i.ResolveLinks(r)
	}
}

package block

import (
	"bytes"
	"fmt"
	"html"

	"github.com/SoCloz/goprismic/fragment/link"
	"github.com/SoCloz/goprismic/fragment/span"
)

// Common block properties
type BaseBlock struct {
	Type  string               `json:"type"`
	Text  string               `json:"text"`
	Spans []span.SpanInterface `json:"spans"`
}

func (b *BaseBlock) AsText() string {
	return b.Text
}

func contains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// FormatHtmlText adds HTML tags to the text and returns it.
func (b *BaseBlock) FormatHtmlText() string {
	t := html.EscapeString(b.Text)

	// Create mapping of start/end index and spans for easier lookup
	endTags := make(map[int][]span.SpanInterface)
	beginTags := make(map[int][]span.SpanInterface)
	for _, s := range b.Spans {
		if _, ok := endTags[s.GetEnd()]; !ok {
			endTags[s.GetEnd()] = make([]span.SpanInterface, 0)
		}
		if _, ok := beginTags[s.GetStart()]; !ok {
			beginTags[s.GetStart()] = make([]span.SpanInterface, 0)
		}

		endTags[s.GetEnd()] = append(endTags[s.GetEnd()], s)
		beginTags[s.GetStart()] = append(beginTags[s.GetStart()], s)
	}

	// Create mapping between Span and UTF-8 offsets
	offsets := make([]int, len(t))
	index := 0
	for k := range t {
		offsets[k] = index
		index++
	}

	var buffer bytes.Buffer

	for i, r := range t {
		if v, contains := endTags[offsets[i]]; contains {
			// Reverse iteration to close tags in correct order
			for k := len(v) - 1; k >= 0; k-- {
				buffer.WriteString(v[k].HtmlEndTag())
			}
		}

		if v, contains := beginTags[offsets[i]]; contains {
			for _, currentSpan := range v {
				buffer.WriteString(currentSpan.HtmlBeginTag())
			}
		}

		buffer.WriteRune(r)
	}

	// Close end tags
	if v, contains := endTags[len(t)]; contains {
		for k := len(v) - 1; k >= 0; k-- {
			buffer.WriteString(v[k].HtmlEndTag())
		}
	}
	return buffer.String()
}

func (b *BaseBlock) decodeBlock(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", enc)
	}
	if v, found := dec["type"]; found {
		b.Type = v.(string)
	}
	if v, found := dec["text"]; found {
		b.Text = v.(string)
	}
	if v, found := dec["spans"]; found {
		dec2, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("%+v is not a slice", dec2)
		}
		b.Spans = make([]span.SpanInterface, 0, len(dec2))
		for _, v := range dec2 {
			dec3, ok := v.(map[string]interface{})
			if ok {
				var s span.SpanInterface
				switch dec3["type"] {
				case "strong":
					s = new(span.Strong)
				case "em":
					s = new(span.Em)
				case "hyperlink":
					s = new(span.Hyperlink)
				default:
					panic(fmt.Sprintf("Unknown span type %s", dec3["type"]))
				}
				err := s.Decode(v)
				if err == nil {
					b.Spans = append(b.Spans, s)
				}
			}
		}
	}
	return nil
}

// Resolves links
func (b *BaseBlock) ResolveLinks(r link.Resolver) {
	for _, v := range b.Spans {
		v.ResolveLinks(r)
	}
}

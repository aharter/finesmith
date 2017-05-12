package fragment

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SoCloz/goprismic/fragment/link"
)

// A color fragment (hex color code)
type Color string

func (c *Color) Decode(_ string, enc interface{}) error {
	res, ok := enc.(string)
	if !ok {
		return fmt.Errorf("unable to decode color fragment : %+v is a %s, not a string", enc, reflect.TypeOf(enc))
	}
	*c = Color(res)
	return nil
}

func (c *Color) AsText() string {
	return string(*c)
}

func (c *Color) AsHtml() string {
	return fmt.Sprintf("<span class=\"number\">%d</span>", *c)
}

func (c *Color) ResolveLinks(_ link.Resolver) {}

func (c *Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Text string `json:"text"`
		HTML string `json:"html"`
	}{
		Text: c.AsText(),
		HTML: c.AsHtml(),
	})
}

package image

import (
	"fmt"
	"reflect"
)

// An image view
type View struct {
	Url        string `json:"url"`
	Alt        string `json:"alt,omitempty"`
	Copyright  string `json:"copyright,omitempty"`
	Dimensions struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"dimensions"`
}

func (i *View) Decode(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to decode image view : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
	}
	if v, found := dec["url"]; found {
		i.Url = v.(string)
	}
	if v, found := dec["alt"]; found && v != nil {
		i.Alt = v.(string)
	}
	if v, found := dec["copyright"]; found && v != nil {
		i.Copyright = v.(string)
	}
	if d, found := dec["dimensions"]; found {
		dim, ok := d.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%+v is not a map", d)
		}
		if v, found := dim["width"]; found && v != nil {
			i.Dimensions.Width = int(v.(float64))
		}
		if v, found := dim["height"]; found && v != nil {
			i.Dimensions.Height = int(v.(float64))
		}
	}
	return nil
}

func (i *View) AsText() string {
	return i.Url
}

func (i *View) AsHtml() string {
	alt := ""
	if i.Alt != "" {
		alt = fmt.Sprintf(" alt=\"%s\"", i.Alt)
	}
	return fmt.Sprintf("<img src=\"%s\" width=\"%d\" height=\"%d\"%s/>", i.Url, i.Dimensions.Width, i.Dimensions.Height, alt)
}

func (i *View) Ratio() float64 {
	return float64(i.Dimensions.Width) / float64(i.Dimensions.Height)
}

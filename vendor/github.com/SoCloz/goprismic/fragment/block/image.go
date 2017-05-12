package block

import (
	"encoding/json"

	"github.com/SoCloz/goprismic/fragment/image"
)

// An image block
type Image struct {
	BaseBlock
	View *image.View
}

func (i *Image) Decode(enc interface{}) error {
	i.View = new(image.View)
	err := i.View.Decode(enc)
	if err != nil {
		return err
	}
	return i.decodeBlock(enc)
}

func (i *Image) AsHtml() string {
	return i.View.AsHtml()
}

func (i *Image) ParentHtmlTag() string {
	return ""
}

func (i *Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	return json.Marshal(&struct {
		HTML string `json:"html"`
		*Alias
	}{
		HTML:  i.AsHtml(),
		Alias: (*Alias)(i),
	})
}

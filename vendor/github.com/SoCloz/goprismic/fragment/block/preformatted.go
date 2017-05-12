package block

import (
	"encoding/json"
	"fmt"
)

// A preformatted block
type Preformatted struct {
	BaseBlock
}

func (p *Preformatted) Decode(enc interface{}) error {
	return p.decodeBlock(enc)
}

func (p *Preformatted) AsHtml() string {
	return fmt.Sprintf("<pre>%s</pre>", p.FormatHtmlText())
}

func (p *Preformatted) ParentHtmlTag() string {
	return ""
}

func (p *Preformatted) MarshalJSON() ([]byte, error) {
	type Alias Preformatted
	return json.Marshal(&struct {
		HTML string `json:"html"`
		*Alias
	}{
		HTML:  p.AsHtml(),
		Alias: (*Alias)(p),
	})
}

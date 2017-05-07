package controller

import "encoding/json"

var (
	TextDescription text
	RefDescription  ref
	ImgDescription  img
)

type IDescriptionType interface {
	Id() uint8
	String() string
	GoString() string
}

type text struct{ IDescriptionType }

func (t text) Id() uint8        { return 0x01 }
func (t text) String() string   { return "text" }
func (t text) GoString() string { return "description.text" }

type ref struct{ IDescriptionType }

func (r ref) Id() uint8        { return 0x02 }
func (r ref) String() string   { return "ref" }
func (r ref) GoString() string { return "description.ref" }

type img struct{ IDescriptionType }

func (i img) Id() uint8        { return 0x03 }
func (i img) String() string   { return "img" }
func (i img) GoString() string { return "description.img" }

type Description struct {
	Type IDescriptionType `json:"type"`
	Text string           `json:"text"`
	Url  string           `json:"url"`
	Data []byte           `json:"-"`
}

type representation struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Url  string `json:"url"`
}

func (d *Description) MarshalJSON() ([]byte, error) {
	r := representation{
		Type: d.Type.String(),
		Text: d.Text,
		Url:  d.Url,
	}

	return json.Marshal(&r)
}

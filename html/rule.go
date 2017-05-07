package html

import (
	htmlnet "golang.org/x/net/html"

	"fmt"
	"strings"
)

type Rule struct {
	Tag   string
	Attrs map[string]string
}

func (r Rule) Equal(element interface{}) bool {

	var (
		elementData        *string
		elementAttributes  []htmlnet.Attribute
		elementTypeIsValid bool
	)

	switch element.(type) {
	case *htmlnet.Node:
		elementData = &element.(*htmlnet.Node).Data
		elementAttributes = element.(*htmlnet.Node).Attr
		elementTypeIsValid = element.(*htmlnet.Node).Type == htmlnet.ElementNode
	case *htmlnet.Token:
		elementData = &element.(*htmlnet.Token).Data
		elementAttributes = element.(*htmlnet.Token).Attr
		elementTypeIsValid = true
	default:
		panic(fmt.Sprintf("Unknown element type:%#v", element))
	}

	if elementTypeIsValid &&
		strings.Compare(*elementData, r.Tag) == 0 {

		if len(r.Attrs) == 0 {
			return true
		}

		var countChecks int

		for i, _ := range elementAttributes {
			attr := &elementAttributes[i]

			if etalonValue, ok := r.Attrs[attr.Key]; !ok {
				continue
			} else if strings.Compare(attr.Val, etalonValue) == 0 {
				countChecks += 1
			} else {
				break
			}
		}

		return countChecks == len(r.Attrs)
	}

	return false
}

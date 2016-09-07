package parser

import (
	"golang.org/x/net/html"
	"strings"
)

type Rule struct {
	Tag   string
	Attrs map[string]string
}

func (r Rule) EqualToken(token *html.Token) bool {

	if strings.Compare(token.Data, r.Tag) == 0 {
		if len(r.Attrs) == 0 {
			return true
		}

		var countChecks int

		for i, _ := range token.Attr {
			attr := &token.Attr[i]
			etalonValue, ok := r.Attrs[attr.Key]
			if !ok {
				continue
			}

			if strings.Compare(attr.Val, etalonValue) == 0 {
				countChecks += 1
			} else {
				break
			}
		}

		return countChecks == len(r.Attrs)
	}

	return false
}

func (r Rule) EqualNode(node *html.Node) bool {

	if node.Type == html.ElementNode &&
		strings.Compare(node.Data, r.Tag) == 0 {

		if len(r.Attrs) == 0 {
			return true
		}

		var countChecks int

		for i, _ := range node.Attr {
			attr := &node.Attr[i]
			etalonValue, ok := r.Attrs[attr.Key]
			if !ok {
				continue
			}

			if strings.Compare(attr.Val, etalonValue) == 0 {
				countChecks += 1
			} else {
				break
			}
		}

		return countChecks == len(r.Attrs)
	}

	return false
}

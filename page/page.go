package page

import (
	"bytes"
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

var httpClient = new(http.Client)

type Page struct {
	Req    *http.Request
	Client *http.Client
	Data   []byte
	Root   *Node
}

func New(req *http.Request) *Page {
	return &Page{
		Req: req,
	}
}

func (p *Page) Download() error {
	p.Data = nil

	client := p.Client
	if client == nil {
		client = httpClient
	}

	resp, err := client.Do(p.Req)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed get url '%s': %s", p.Req.URL.String(), err.Error()))
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if pos := strings.Index(ct, ";"); pos != -1 {
		ct = string([]rune(ct)[:pos])
	}

	utf8, err := charset.NewReader(resp.Body, ct)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed decode result of the page '%s': %s", p.Req.URL.String(), err.Error()))
	}

	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed read body of the page '%s': %s", p.Req.URL.String(), err.Error()))
	}

	p.Data = body

	return nil
}

func (p Page) Tokenizer() *html.Tokenizer {
	reader := bytes.NewReader(p.Data)
	return html.NewTokenizer(reader)
}

func TagAttr(attrs []html.Attribute, attrName string) (val string) {

	for i, _ := range attrs {
		a := &attrs[i]
		if strings.Compare(a.Key, attrName) == 0 {
			val = a.Val
			break
		}
	}

	return
}

func (p Page) Nodes(rule ...*Rule) []*Node {

	if p.Root == nil {
		p.Parse()
	}

	return Nodes(p.Root, rule...)
}

func Nodes(root *Node, rule ...*Rule) []*Node {

	results := make([][]*Node, len(rule), len(rule))

	for ruleIndex, ruleValue := range rule {
		ruleResults := make([]*Node, 0, 5)

		if ruleIndex == 0 {
			ruleResults = find(root, ruleValue)
		} else {
			parentsList := results[ruleIndex-1]
			if len(parentsList) > 0 {

				for _, node := range parentsList {
					if list := find(node, ruleValue); len(list) > 0 {
						ruleResults = append(ruleResults, list...)
					}
				}
			}
		}

		if len(ruleResults) == 0 {
			break
		} else {
			results[ruleIndex] = ruleResults
		}
	}

	var retval []*Node

	end := results[len(rule)-1]
	if end == nil {
		retval = make([]*Node, 0, 0)
	} else {
		retval = end
	}

	return retval
}

func find(parent *Node, r *Rule) (retval []*Node) {

	for _, child := range parent.Children {

		if r.Equal(child) {
			retval = append(retval, child)
		} else if len(retval) == 0 {
			if list := find(child, r); len(list) > 0 {
				retval = append(retval, list...)
			}
		}
	}

	return
}

func (p *Page) Parse() {

	tokenizer := p.Tokenizer()
	lastParent := NewNode("html", nil)

	var lastTextNode *Node

	for {
		t := tokenizer.Next()
		if t == html.ErrorToken {
			break // End of the document
		}

		switch {
		case t == html.EndTagToken:

			t := tokenizer.Token()

			if !strings.EqualFold(lastParent.Tag, t.Data) {
				// if tag not closed, we find real parent node
				for lastParent.Parent != nil {
					lastParent = lastParent.Parent
					if strings.EqualFold(lastParent.Tag, t.Data) {
						break
					}
				}
				if lastTextNode != nil {
					lastTextNode.Parent = lastParent
				}
			}

			if lastTextNode != nil {
				lastTextNode.Parent.Children = append(lastTextNode.Parent.Children, lastTextNode)
				lastTextNode = nil
			}

			if lastParent.Parent != nil {
				lastParent = lastParent.Parent
			}

		case t == html.SelfClosingTagToken:

			t := tokenizer.Token()

			node := NewNode(t.Data, t.Attr)
			node.Parent = lastParent
			lastParent.Children = append(lastParent.Children, node)

		case t == html.StartTagToken:

			t := tokenizer.Token()
			if lastTextNode != nil {
				lastTextNode.Parent.Children = append(lastTextNode.Parent.Children, lastTextNode)
				lastTextNode = nil
			}

			node := NewNode(t.Data, t.Attr)
			node.Parent = lastParent
			lastParent.Children = append(lastParent.Children, node)

			if node.Tag != "br" && node.Tag != "hr" && node.Tag != "img" {
				lastParent = node
			}

		case t == html.TextToken:
			t := tokenizer.Token()

			if len(strings.TrimSpace(t.Data)) > 0 {
				lastTextNode = NewNode(TEXT_TAG, make([]html.Attribute, 0, 0))
				lastTextNode.Text = t.Data
				lastTextNode.Parent = lastParent
			}
		}
	}

	p.Root = lastParent
}

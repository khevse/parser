package parser

import (
	"bytes"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
)

func DownloadHtml(host string, url string) []byte {

	resp, err := http.Get(strings.Join([]string{host, url}, ""))
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()

	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		panic("Encoding error:" + err.Error())
	}

	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		panic("IO error:" + err.Error())
	}

	return body
}

func BytesToTokenizer(data *[]byte) *html.Tokenizer {

	r := bytes.NewReader(*data)
	return html.NewTokenizer(r)
}

func BytesToHtmlDom(data *[]byte) (*html.Node, error) {

	r := bytes.NewReader(*data)
	return html.Parse(r)
}

func TagAttr(attrs *[]html.Attribute, attrName string) string {

	var val string

	for _, a := range *attrs {
		if strings.Compare(a.Key, attrName) == 0 {
			val = a.Val
			break
		}
	}

	return val
}

func FindNodes(parent *html.Node, r *Rule) []*html.Node {

	var nodes_list []*html.Node

	for node := parent.FirstChild; node != nil; node = node.NextSibling {

		if node.Type != html.ElementNode {
			continue
		}

		if r.EqualNode(node) {
			nodes_list = append(nodes_list, node)
		} else {
			nodes_list = FindNodes(node, r)
			if len(nodes_list) > 0 {
				break
			}
		}
	}

	return nodes_list
}

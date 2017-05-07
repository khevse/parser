package html

import (
	htmlnet "golang.org/x/net/html"
	"golang.org/x/net/html/charset"

	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

func DownloadHtml(host string, url string) []byte {

	resp, err := http.Get(host + url)
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

func BytesToTokenizer(data []byte) *htmlnet.Tokenizer {

	r := bytes.NewReader(data)
	return htmlnet.NewTokenizer(r)
}

func BytesToHtmlDom(data []byte) (*htmlnet.Node, error) {

	r := bytes.NewReader(data)
	return htmlnet.Parse(r)
}

func TagAttr(attrs []htmlnet.Attribute, attrName string) string {

	for i, _ := range attrs {
		a := &attrs[i]
		if strings.Compare(a.Key, attrName) == 0 {
			return a.Val
		}
	}

	return ""
}

func FindNodes(parent *htmlnet.Node, r *Rule) []*htmlnet.Node {

	nodesList := make([]*htmlnet.Node, 0, 0)

	for node := parent.FirstChild; node != nil; node = node.NextSibling {

		if node.Type != htmlnet.ElementNode {
			continue
		}

		if r.Equal(node) {
			nodesList = append(nodesList, node)
		} else {
			nodesList = FindNodes(node, r)
			if len(nodesList) > 0 {
				break
			}
		}
	}

	return nodesList
}

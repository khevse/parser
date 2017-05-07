package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/khevse/parser/page"
)

func TestParseTOC(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	req, err := http.NewRequest(http.MethodGet, "https://www.invitro.ru/analizes/for-doctors/", nil)
	if err != nil {
		t.Error(err)
	}

	sourceData, _ := ioutil.ReadFile("../test_data/toc.html")
	pageData := page.Page{
		Req:  req,
		Data: sourceData,
	}

	toc, err := parseTOC(&pageData)
	if err != nil {
		t.Error(err)
		return
	}

	srcList := []struct {
		Href   string
		Name   string
		Result bool
	}{
		{
			"/analizes/for-doctors/137/",
			"Гематологические исследования",
			true,
		},
		{
			"/analizes/for-doctors/140/",
			"Биохимические исследования",
			true,
		},
		{
			"/analizes/for-doctors/152/",
			"Гормональные исследования",
			true,
		},
	}

	if len(toc) != len(srcList) {
		t.Error("Fail:", len(toc), len(srcList))
		return
	}

	for i, _ := range srcList {
		item := toc[i].(*TOCItem)
		src := srcList[i]

		if src.Href != item.Href() {
			t.Error("Wrong href:", src.Href, item.Href())
		} else if src.Name != item.Name {
			t.Error("Wrong name:", src.Name, item.Name)
		}
	}
}

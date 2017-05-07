package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/khevse/parser/engine"
	"github.com/khevse/parser/page"
)

func TOCHandler(url string) ([]engine.ITOCItem, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	src := page.New(req)
	if err := src.Download(); err != nil {
		return nil, err
	} else if res, err := parseTOC(src); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func parseTOC(src *page.Page) ([]engine.ITOCItem, error) {
	// Example of the source data:
	// <div class="group-name-brd" style="border-left:4px solid #f3209d;">
	//     <a href="/analizes/for-doctors/140/" title="Биохимические исследования">Биохимические исследования</a>
	// </div>

	nodesList := src.Nodes(
		page.NewRule("ul", "class", "marked_list", "class", "analyses_list"),
		page.NewRule("li", "class", `an_(\d+)`),
		page.NewRule("a"),
	)

	itemsList := make([]engine.ITOCItem, 0, len(nodesList))

	for _, node := range nodesList {

		href := page.TagAttr(node.Attr, "href")
		var number int

		urlFormat := strings.Replace(src.Req.URL.RequestURI()+"/%d/", "//", "/", -1)
		if n, err := fmt.Sscanf(href, urlFormat, &number); err != nil || n == 0 {
			return nil, errors.New(fmt.Sprintln("Failed group number:", href, urlFormat))
		}

		item := TOCItem{
			Url:    href,
			Name:   node.GetText(),
			Number: fmt.Sprintf("%d", number),
		}

		if len(item.Url) == 0 || len(item.Name) == 0 || len(item.Number) == 0 {
			return nil, errors.New(fmt.Sprintln("Failed group data:", item))
		}

		itemsList = append(itemsList, &item)
	}

	if len(itemsList) == 0 {
		return nil, errors.New("TOC is empty.")
	}

	return itemsList, nil
}

package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/khevse/parser/engine"
	"github.com/khevse/parser/page"
)

type TOCItem struct {
	engine.ITOCItem

	Url    string
	Name   string
	Number string
	Data   *Target
}

func (t TOCItem) Href() string {
	return t.Url
}

func (t *TOCItem) GetChildren() []engine.ITarget {
	return dataToList(t.Data)
}

func (t *TOCItem) SetChildren(url string) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	} else {
		src := page.New(req)

		if err := src.Download(); err == nil {
			if res, err := parseTOCItem(t.Name, src); err == nil {
				if res == nil {
					log.Fatal(t.Url)
				}
				t.Data = res
			}
		}
	}
}

func dataToList(parent *Target) []engine.ITarget {

	retval := make([]engine.ITarget, 0, 1)

	if parent.Parent != nil {
		retval = append(retval, parent)
	}

	for _, item := range parent.Children {
		list := dataToList(item)

		if len(list) > 0 {
			if (cap(retval) - len(retval)) < len(list) {
				newLen := len(retval) + len(list)

				tmp := make([]engine.ITarget, len(retval), newLen*2)
				copy(tmp, retval)
				retval = tmp
			}

			retval = append(retval, list...)
		}
	}

	return retval
}

var (
	groupCheckerRule   = page.NewRule("th", "class", `price_content__head`)
	elementCheckerRule = page.NewRule("td", "class", `(^|\s)elem_list_(\d+)($|\s)`)
)

func parseTOCItem(groupName string, src *page.Page) (*Target, error) {
	// Example of the source data:
	// <body>
	// 	...
	//		<div id="catalog-section-analiz">
	//			<table class="data-table">
	//			  <tbody>
	//				<tr>...</tr>
	//				<tr>...</tr>
	//			  </tbody>
	//			</table>
	//		</div>
	// 	...
	// </body>

	nodesList := src.Nodes(
		page.NewRule("table", "class", "price_content"),
		page.NewRule("tbody"),
		page.NewRule("tr"),
	)

	tdRule := page.NewRule("td")
	thRule := page.NewRule("th", "class", "price_content__head")

	root := &Target{Name: groupName}
	var currentParent **Target = &root

	for _, rowNode := range nodesList {

		if columns := rowNode.Find(thRule); len(columns) > 0 {
			if group, err := readGroupRow(columns); err != nil {
				return nil, err
			} else {
				for {
					if (*currentParent).Level < group.Level {
						break
					}

					newParent := (*currentParent).Parent
					if newParent == nil {
						return nil, errors.New("Wrong parent")
					} else {
						currentParent = &newParent
					}
				}

				group.Parent = *currentParent
				group.Parent.Children = append(group.Parent.Children, group)
				currentParent = &(group.Parent.Children[len(group.Parent.Children)-1])
			}

		} else if columns := rowNode.Find(tdRule); len(columns) > 0 {
			elem, err := readElementRow(columns)
			if err != nil {
				return nil, err
			}

			elem.Parent = *currentParent
			(*currentParent).Children = append((*currentParent).Children, elem)
		}
	}

	return root, nil
}

func readElementRow(columns []*page.Node) (*Target, error) {
	// Example of the source data:
	// <tr>
	//     <td>1500</td>
	//     <td><a href="/analizes/for-doctors/140/29559/">Антиоксидантный статус Total antioxidant status, TAS</a></td>
	//     <td>
	//         4680
	//     </td>
	//     <td>
	//         <a href="#" item_id="29559" class="addToCart" title="Добавить в корзину"></a>
	//     </td>
	// </tr>

	var (
		nodeNum    *page.Node
		nodeName   *page.Node
		nodeAmount *page.Node
	)

	if len(columns) != 4 {
		msg := fmt.Sprintf("Wrong item data: %#v", columns)
		log.Printf(msg)
		return nil, errors.New(msg)
	} else if nodeName = columns[1]; len(nodeName.Children) != 1 {
		msg := fmt.Sprintf("Wrong item data: %#v", nodeName)
		log.Printf(msg)
		return nil, errors.New(msg)
	} else if nodeAmount = columns[2]; len(nodeAmount.Children) != 1 {
		msg := fmt.Sprintf("Wrong item data: %#v", nodeAmount)
		log.Printf(msg)
		return nil, errors.New(msg)
	} else if nodeName = nodeName.Children[0]; len(nodeName.Children) != 1 {
		msg := fmt.Sprintf("Wrong item data: %#v", nodeName)
		log.Printf(msg)
		return nil, errors.New(msg)
	} else {
		nodeNum = columns[0]
	}

	var (
		itemUrl    *url.URL
		itemAmount float64
	)

	if val, err := url.Parse(page.TagAttr(nodeName.Attr, "href")); err != nil {
		log.Printf("Wrong item url:", err)
		return nil, err
	} else {
		itemUrl = val
	}

	if val, err := strconv.ParseFloat(strings.TrimSpace(nodeAmount.GetText()), 10); err != nil {
		log.Printf("Wrong item amount:", err)
		return nil, err
	} else {
		itemAmount = val
	}

	item := Target{
		Name:      nodeName.GetText(),
		PriceCode: nodeNum.GetText(),
		Amount:    itemAmount,
		URL:       itemUrl,
	}

	if item.Amount == 0 {
		return nil, errors.New(fmt.Sprintf("Failed element amount: %+v", item))
	}

	return &item, nil
}

var groupNumRepalacer = strings.NewReplacer("h", "", "H", "")

func readGroupRow(columns []*page.Node) (*Target, error) {
	// Example of the source data:
	// <tr>
	//     <th colspan="4" class="price_content__head">
	//         <h2>Биохимические исследования</h2>
	//     </th>
	// </tr>

	var groupNode *page.Node

	if len(columns) != 1 {
		msg := fmt.Sprintf("Wrong group data: %#v", columns)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	children := columns[0].Children

	var child *page.Node
	for _, c := range children {
		if strings.HasPrefix(c.Tag, "h") || strings.HasPrefix(c.Tag, "H") {
			if child != nil {
				msg := "Wrong group data:"
				for _, c := range children {
					msg += fmt.Sprintf("\n - %#v", c)
				}
				log.Println(msg)
				return nil, errors.New(msg)
			}

			child = c
		}
	}

	if child == nil {
		msg := "Wrong group data:"
		for _, c := range columns[0].Children {
			msg += fmt.Sprintf("\n - %#v", c)
		}
		log.Println(msg)
		return nil, errors.New(msg)
	}

	groupNode = child

	level := groupNode.Tag
	level = groupNumRepalacer.Replace(level)

	var groupLevel int

	if val, err := strconv.ParseInt(level, 10, 64); err != nil {
		msg := fmt.Sprintf("Wrong group level '%v': %s ", level, err.Error())
		log.Println(msg)
		return nil, errors.New(msg)
	} else {
		groupLevel = int(val)
	}

	item := Target{
		Name:  strings.TrimSpace(groupNode.GetText()),
		Level: groupLevel,
	}

	if len(item.Name) == 0 {
		msg := fmt.Sprintf("Failed group data: %v", item)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	return &item, nil
}

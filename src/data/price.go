package data

import (
	"../parser"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

type Group struct {
	Href   string
	Name   string
	Number string
}

type Price struct {
	Parent   *Price
	Level    uint8
	Name     string
	Number   string
	Amount   float32
	Children []*Price
}

func ReadPrice(groupNumber string, groupName string, data []byte) (price []*Price, err error) {
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

	var rows []*html.Node
	rows, err = findRows(&data)
	if err != nil {
		return
	}

	groupChecker := func(class string) bool {
		return strings.Index(class, "name") != -1 && strings.Index(class, "sec_") != -1
	}

	elementChecker := func(class string) bool {
		return strings.Index(class, "elem_list") != -1
	}

	// Read price

	price = []*Price{&Price{Name: groupName, Number: groupNumber}}
	var currentParent **Price = &price[0]

	for i, _ := range rows {

		isGroup := checkRowType(rows[i], groupChecker)
		isElement := !isGroup && checkRowType(rows[i], elementChecker)

		if isGroup {
			group := readGroup(rows[i])

			for {
				if currentParent == nil || (*currentParent).Level < group.Level {
					break
				}

				newParent := (*currentParent).Parent
				if newParent == nil {
					currentParent = nil
				} else {
					currentParent = &(*currentParent).Parent
				}
			}

			if currentParent == nil {
				price = append(price, group)
				currentParent = &price[len(price)-1]
			} else {
				group.Parent = *currentParent
				ChildrenList := &(*currentParent).Children
				*ChildrenList = append(*ChildrenList, group)
				currentParent = &(*ChildrenList)[len(*ChildrenList)-1]
			}

		} else if isElement {
			elem := readElement(rows[i])

			if currentParent == nil {
				price = append(price, elem)
			} else {
				elem.Parent = (*currentParent)
				(*currentParent).Children = append((*currentParent).Children, elem)
			}

		}
	}

	return
}

func findRows(data *[]byte) (rows []*html.Node, err error) {
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

	doc, errDoc := parser.BytesToHtmlDom(data)
	if errDoc != nil {
		err = errors.New("Failed price html: " + errDoc.Error())
		return
	}

	CatalogSelection := parser.FindNodes(doc, &parser.Rule{
		Tag:   "div",
		Attrs: map[string]string{"id": "catalog-section-analiz"},
	})
	if len(CatalogSelection) != 1 {
		err = errors.New("Not found catalog selection")
		return
	}

	CatalogTable := parser.FindNodes(CatalogSelection[0], &parser.Rule{
		Tag:   "table",
		Attrs: map[string]string{"class": "data-table"},
	})
	if len(CatalogTable) != 1 {
		err = errors.New("Not found catalog table:\n" + string(*data))
		return
	}

	CatalogTableData := parser.FindNodes(CatalogTable[0], &parser.Rule{
		Tag:   "tbody",
		Attrs: map[string]string{},
	})
	if len(CatalogTableData) != 1 {
		err = errors.New("Not found catalog table data:\n" + string(*data))
		return
	}

	rows = parser.FindNodes(CatalogTableData[0], &parser.Rule{
		Tag:   "tr",
		Attrs: map[string]string{},
	})

	return
}

func ReadPriceGroups(data *[]byte) []Group {
	// Example of the source data:
	// <div class="group-name-brd" style="border-left:4px solid #f3209d;">
	//     <a href="/analizes/for-doctors/140/" title="Биохимические исследования">Биохимические исследования</a>
	// </div>

	const (
		ATTR_TITLE string = "title"
		ATTR_HREF  string = "href"
	)

	rule_group_data_wrapper := parser.Rule{
		Tag:   "div",
		Attrs: map[string]string{"class": "group-name-brd"},
	}

	rule_group_data := parser.Rule{
		Tag:   "a",
		Attrs: map[string]string{},
	}

	var (
		groups_list      []Group
		openTagGroupName bool
	)

	tokenizer := parser.BytesToTokenizer(data)

	for {
		t := tokenizer.Next()
		if t == html.ErrorToken {
			break // End of the document
		}

		switch {
		case t == html.EndTagToken:
			openTagGroupName = false
		case t == html.StartTagToken:
			t := tokenizer.Token()

			if !openTagGroupName && rule_group_data_wrapper.EqualToken(&t) {
				openTagGroupName = true

			} else if openTagGroupName && rule_group_data.EqualToken(&t) {

				href := parser.TagAttr(&t.Attr, ATTR_HREF)
				var number int
				fmt.Sscanf(href, "/analizes/for-doctors/%d/", &number)

				elem := Group{
					Href:   href,
					Name:   parser.TagAttr(&t.Attr, ATTR_TITLE),
					Number: fmt.Sprintf("%d", number),
				}

				if len(elem.Href) == 0 || len(elem.Name) == 0 || len(elem.Number) == 0 {
					panic("Failed group data.")
				}

				groups_list = append(groups_list, elem)
			}
		}
	}

	if len(groups_list) == 0 {
		panic("Groups list is empty.")
	}

	return groups_list
}

func PrintPrice(price []*Price) {

	for i, _ := range price {
		fmt.Println(price[i])
		PrintPrice(price[i].Children)
	}
}

func readElement(row *html.Node) *Price {
	// Example of the source data:
	// <tr>
	// 	<td class="number elem_list_0">93</td>
	// 	<td class="name elem_list_0">
	// 		<a href="/analizes/for-doctors/479/2860/">Группа крови</a>
	// 	</td>
	// 	<td class="price elem_list_0">
	//		360
	// 	</td>
	// <tr>

	var elem *Price = &Price{}

	for td := row.FirstChild; td != nil; td = td.NextSibling {
		if td.Type != html.ElementNode || td.Data != "td" {
			continue
		}

		class := parser.TagAttr(&td.Attr, "class")

		if strings.Index(class, "elem_list_") != -1 {

			if strings.Index(class, "number") != -1 {
				elem.Number = strings.Trim(td.FirstChild.Data, " ")
			} else if strings.Index(class, "name") != -1 {

				for node := td.FirstChild; node != nil; node = node.NextSibling {
					if node.Type == html.ElementNode && node.Data == "a" {
						elem.Name = strings.Trim(node.FirstChild.Data, " ")
						break
					}
				}

			} else if strings.Index(class, "price") != -1 {

				var amount float32
				for node := td.FirstChild; node != nil; node = node.NextSibling {
					if node.Type == html.TextNode {
						val := node.Data
						val = strings.Replace(val, "\n", "", -1)
						val = strings.Replace(val, "\r", "", -1)
						val = strings.Replace(val, " ", "", -1)

						fmt.Sscanf(val, "%f", &amount)
						if amount > 0 {
							break
						}
					}
				}

				elem.Amount = amount
			}
		}

	}

	if elem.Amount == 0 {
		panic(fmt.Sprintf("Failed price: %v", elem))
	}

	return elem
}

func readGroup(row *html.Node) *Price {
	// Example of the source data:
	// <tr>
	// 	<td class="name sec_3" colspan="3">
	// 		<a name="479"></a>
	// 		<span class="name_gr_3">Групповая принадлежность крови</span>
	// 	</td>
	// </tr>

	var elem *Price = nil

	for td := row.FirstChild; td != nil; td = td.NextSibling {
		if td.Type != html.ElementNode || td.Data != "td" {
			continue
		}

		class := parser.TagAttr(&td.Attr, "class")
		class = strings.Trim(strings.Replace(class, "name", "", -1), " ")

		var groupLevel int
		fmt.Sscanf(class, "sec_%d", &groupLevel)

		elem = &Price{
			Level:    uint8(groupLevel),
			Children: []*Price{},
		}

		for node := td.FirstChild; node != nil; node = node.NextSibling {
			if node.Type != html.ElementNode {
				continue
			}

			if node.Data == "a" {
				if val := parser.TagAttr(&node.Attr, "name"); len(val) > 0 {
					elem.Number = strings.Trim(val, " ")
				}
			} else if node.Data == "span" {
				elem.Name = strings.Trim(node.FirstChild.Data, " ")
			}
		}

		if len(elem.Number) == 0 || len(elem.Name) == 0 {
			panic(fmt.Sprintf("Failed group data: %v", elem))
		}
	}

	if elem == nil {
		panic("Bad group data.")
	}

	return elem
}

func checkRowType(row *html.Node, fn func(string) bool) bool {

	var retval bool

	for td := row.FirstChild; td != nil; td = td.NextSibling {
		if td.Type != html.ElementNode || td.Data != "td" {
			continue
		}

		class := parser.TagAttr(&td.Attr, "class")
		if retval = fn(class); retval {
			break
		}
	}

	return retval
}

package page

import (
	"testing"
)

func TestRuleEqual(t *testing.T) {

	testsList := []struct {
		Rule   *Rule
		Result bool
	}{
		{
			Rule:   NewRule("span", "id", "elem_28932"),
			Result: true,
		},
		{
			Rule:   NewRule("span", "class", "icon_span"),
			Result: true,
		},
		{
			Rule:   NewRule("span", "id", "elem_28932", "class", "icon_span"),
			Result: true,
		},
		{
			Rule:   NewRule("span", "id", "elem_28932", "class", `^[a-z]{4}_span`),
			Result: true,
		},
		{
			Rule:   NewRule("span", "id", `name`, "id", `(^|\s)elem_(\d+)($|\s)`, "class", `^[a-z]{4}_span$`, "class", `^icon_[a-z]{4}$`),
			Result: true,
		},
		{
			Rule:   NewRule("a", "href", "javascript:...."),
			Result: true,
		},
		{
			Rule:   NewRule("a", "href", "javascript"),
			Result: false,
		},
	}

	sourse_data := `
          <td class="price elem_list_0">
              <span id="name elem_28932" class="icon_span">
                  <a href="javascript:...."/>
              </span>
              PRICE-> 1020
          </td>`

	page := Page{
		Data: []byte(sourse_data),
	}
	page.Parse()

	for testIndex, test := range testsList {
		var node *Node
		if test.Rule.Tag == "span" {
			node = page.Root.Children[0].Children[0]
		} else if test.Rule.Tag == "a" {
			node = page.Root.Children[0].Children[0].Children[0]
		}

		result := test.Rule.Equal(node)

		if result != test.Result {
			t.Errorf("Wrong result for rule №%d: %v != %v", testIndex, test.Result, result)
		}
	}
}

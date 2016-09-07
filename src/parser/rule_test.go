package parser

import (
	"bytes"
	"golang.org/x/net/html"
	"testing"
)

func TestRuleEqual(t *testing.T) {

	sourse_data := `
          <td class="price elem_list_0">
              <span id="elem_28932" class="icon_span">
                  <a href="javascript:...."></a>
              </span>
              PRICE-> 1020
          </td>`

	r := bytes.NewReader([]byte(sourse_data))
	tokenizer := html.NewTokenizer(r)

	rule_1 := Rule{
		Tag:   "span",
		Attrs: map[string]string{"id": "elem_28932"},
	}

	rule_2 := Rule{
		Tag:   "span",
		Attrs: map[string]string{"class": "icon_span"},
	}

	rule_3 := Rule{
		Tag:   "span",
		Attrs: map[string]string{"id": "elem_28932", "class": "icon_span"},
	}

	rule_4 := Rule{
		Tag:   "a",
		Attrs: map[string]string{"href": "javascript:...."},
	}

	rule_5 := Rule{
		Tag:   "a",
		Attrs: map[string]string{"href": ""},
	}

	var (
		result_1 bool
		result_2 bool
		result_3 bool
		result_4 bool
		result_5 bool
	)

	for {
		t := tokenizer.Next()
		if t == html.ErrorToken {
			break // End of the document
		}

		switch {
		case t == html.StartTagToken:
			t := tokenizer.Token()

			result_1 = result_1 || rule_1.EqualToken(&t)
			result_2 = result_2 || rule_2.EqualToken(&t)
			result_3 = result_3 || rule_3.EqualToken(&t)
			result_4 = result_4 || rule_4.EqualToken(&t)
			result_5 = result_5 || rule_5.EqualToken(&t)
		}
	}

	if !result_1 {
		t.Error("Not found rule №1")
	}
	if !result_2 {
		t.Error("Not found rule №2")
	}
	if !result_3 {
		t.Error("Not found rule №3")
	}
	if !result_4 {
		t.Error("Not found rule №4")
	}
	if result_5 {
		t.Error("Failed rule rule №5")
	}
}

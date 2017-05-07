package html

import (
	htmlnet "golang.org/x/net/html"

	"bytes"
	"testing"
)

func TestRuleEqual(t *testing.T) {

	sourseData := `
          <td class="price elem_list_0">
              <span id="elem_28932" class="icon_span">
                  <a href="javascript:...."></a>
              </span>
              PRICE-> 1020
          </td>`

	r := bytes.NewReader([]byte(sourseData))
	tokenizer := htmlnet.NewTokenizer(r)

	rule1 := Rule{
		Tag:   "span",
		Attrs: map[string]string{"id": "elem_28932"},
	}

	rule2 := Rule{
		Tag:   "span",
		Attrs: map[string]string{"class": "icon_span"},
	}

	rule3 := Rule{
		Tag:   "span",
		Attrs: map[string]string{"id": "elem_28932", "class": "icon_span"},
	}

	rule4 := Rule{
		Tag:   "a",
		Attrs: map[string]string{"href": "javascript:...."},
	}

	rule5 := Rule{
		Tag:   "a",
		Attrs: map[string]string{"href": ""},
	}

	var (
		result1 bool
		result2 bool
		result3 bool
		result4 bool
		result5 bool
	)

	for {
		t := tokenizer.Next()
		if t == htmlnet.ErrorToken {
			break // End of the document
		}

		switch {
		case t == htmlnet.StartTagToken:
			t := tokenizer.Token()

			result1 = result1 || rule1.Equal(&t)
			result2 = result2 || rule2.Equal(&t)
			result3 = result3 || rule3.Equal(&t)
			result4 = result4 || rule4.Equal(&t)
			result5 = result5 || rule5.Equal(&t)
		}
	}

	if !result1 {
		t.Error("Not found rule №1")
	}
	if !result2 {
		t.Error("Not found rule №2")
	}
	if !result3 {
		t.Error("Not found rule №3")
	}
	if !result4 {
		t.Error("Not found rule №4")
	}
	if result5 {
		t.Error("Failed rule rule №5")
	}
}

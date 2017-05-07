package page

import (
	"net/http"
	"strings"
	"testing"
)

func TestPage(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "https://google.com:443/maps", nil)
	if err != nil {
		t.Error(err)
	}

	p := New(req)
	if err := p.Download(); err != nil {
		t.Error(err)
	} else if p.Data == nil || len(p.Data) == 0 {
		t.Error("Fail")
	} else {
		p.Parse()
	}
}

func TestFindNodes(t *testing.T) {

	sourceData := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "https://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="https://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml">

<head>
</head>

<body>
  <!-- comment -->
  <div class="view_brd">
    <div id="view">
    	<div id="item_1" />
    	<div id="item_2" />
    	<div id="item_3" />
    </div>
    <div id="view">
    	<div id="item_4">
    		<div id="item_4_1"/>
    		divText1
    	</div>
    	<div id="item_5"> text my</div>
    	<div id="item_6"></div>
    </div>
    <div id="view_test">
    	<span id="item_7">
    	divText2
    </div>
  </div>
</body>
    `
	p := new(Page)
	p.Data = []byte(sourceData)

	if list := p.Nodes(NewRule("div", "id", "view")); len(list) != 2 {
		t.Error("Fail", list)
	} else {
		for _, node := range list {
			if node.Tag != "div" || node.Attr[0].Val != "view" {
				t.Error("Fail", *node)
			}
		}
	}

	if list := p.Nodes(NewRule("div", "id", "view"), NewRule("div", "id", `(^|\s)item_(\d+)($|\s)`)); len(list) != 6 {
		for _, node := range list {
			t.Error("Fail", node.Tag, node.Attr)
		}
	}

	if list := p.Nodes(NewRule("div", "id", "item_4")); len(list) != 1 {
		for _, node := range list {
			t.Error("Fail", node.Tag, node.Attr)
		}
	} else {
		if strings.Index(list[0].GetText(), "divText1") == -1 {
			t.Error("Wrong text", list[0].Children)
		}
	}

	if list := p.Nodes(NewRule("div", "id", "view_test")); len(list) != 1 {
		for _, node := range list {
			t.Error("Fail", node.Tag, node.Attr)
		}
	} else {
		if children := list[0].Children; len(children) != 2 {
			t.Error("Fail", children)
		} else if strings.Index(list[0].GetText(), "divText2") == -1 || len(children[0].Text) > 0 || len(children[1].Text) == 0 {
			t.Error("Wrong text", list[0].Children)
		}
	}

	if list := p.Nodes(NewRule("div", "id", "view"), NewRule("div", "id", `(^|\s)item_(\d+)_(\d+)($|\s)`)); len(list) != 1 {
		for _, node := range list {
			t.Error("Fail", node.Tag, node.Attr)
		}
	}
}

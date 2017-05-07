package html

import (
	"testing"
)

func TestFindNodes(t *testing.T) {

	sourceData := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "https://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="https://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml">

<head>
</head>

<body>
  <div class="view_brd">
    <div id="view">
    </div>
  </div>
</body>
    `
	data := []byte(sourceData)
	doc, err := BytesToHtmlDom(data)
	if err != nil {
		t.Error("Failed price html:", err.Error())
	}

	result := FindNodes(doc, &Rule{
		Tag:   "div",
		Attrs: map[string]string{"id": "view"},
	})

	if len(result) != 1 {
		t.Error("Fail", result)
	}
}

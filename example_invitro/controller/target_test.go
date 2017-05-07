package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/khevse/parser/page"
)

func TestTarget(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	req, err := http.NewRequest(http.MethodGet, "https://www.invitro.ru/analizes/for-doctors/156/28932/", nil)
	if err != nil {
		t.Error(err)
		return
	}

	sourceData, _ := ioutil.ReadFile("../test_data/target.html")
	pageData := page.Page{
		Data: sourceData,
		Req:  req,
	}

	parseTarget(&pageData)
}

func TestParseText(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	sourceData, _ := ioutil.ReadFile("../test_data/target.html")
	pageData := page.Page{
		Data: sourceData,
	}

	tabsDesList := pageData.Nodes(
		page.NewRule("div", "class", `tabs_block`),
		page.NewRule("div", "class", `tab_cont`),
	)

	rootNode := tabsDesList[0]

	testData := []struct {
		Node *page.Node
		Desc []Description
	}{
		{
			Node: rootNode.Children[0],
			Desc: []Description{
				Description{Type: TextDescription, Url: "", Text: "Кровь - это жидкая ткань, выполняющая различные функции, в том числе, транспорта кислорода и питательных веществ к органам и тканям и выведения из них шлаковых продуктов. Состоит из плазмы и форменных элементов: эритроцитов, лейкоцитов и тромбоцитов."},
			},
		},
		{
			Node: rootNode.Children[2],
			Desc: []Description{
				Description{Type: ImgDescription, Url: "/images/gvozd1.jpg", Text: ""},
			},
		},
		{
			Node: rootNode.Children[4],
			Desc: []Description{
				Description{Type: TextDescription, Url: "", Text: "Общий анализ крови в лаборатории ИНВИТРО включает в себя определение концентрации гемоглобина, количества эритроцитов, лейкоцитов и тромбоцитов, величины гематокрита и эритроцитарных индексов (MCV, RDW, MCH, MCHC). Общий анализ -"},
				Description{Type: RefDescription, Url: "http://www.invitro.ru/analizes/for-doctors/137/2852/", Text: "см. тест № 5"},
				Description{Type: TextDescription, Url: "", Text: ", Лейкоцитарная формула -"},
				Description{Type: RefDescription, Url: "https://www.invitro.ru/analizes/for-doctors/156/28933/", Text: "см. тест № 911"},
				Description{Type: TextDescription, Url: "", Text: ", СОЭ -"},
				Description{Type: RefDescription, Url: "http://www.invitro.ru/analizes/for-doctors/137/2854/", Text: "см. тест № 139"},
				Description{Type: TextDescription, Url: "", Text: "."},
			},
		},
	}

	for _, src := range testData {
		res := parseText(src.Node)

		if len(src.Desc) != len(res) {
			t.Errorf("'%q'", res)
		} else {

			for i, _ := range src.Desc {
				if src.Desc[i].Type != res[i].Type {
					t.Errorf("'%s' != '%s'", src.Desc[i], res[i])
				} else if src.Desc[i].Url != res[i].Url {
					t.Errorf("'%s' != '%s'", src.Desc[i], res[i])
				} else if src.Desc[i].Text != res[i].Text {
					t.Errorf("'%s' != '%s'", src.Desc[i], res[i])
				}
			}
		}

	}
}

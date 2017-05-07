package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/khevse/parser/page"
)

func TestReadPrice(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	sourceData, _ := ioutil.ReadFile("../test_data/price.html")
	pageData := page.Page{
		Data: sourceData,
	}

	price, err := parseTOCItem("-", &pageData)
	if err != nil {
		t.Error(err)
		return
	}

	tocItem := TOCItem{
		Data: price,
	}
	list := tocItem.GetChildren()
	for _, i := range list {
		if i == nil {
			t.Error(i)
		}
	}

	if len(price.Children) != 1 {
		t.Error("Failed price", price.Children)
		return
	}

	mainGroup := price.Children[0]
	if err = chechGroup(mainGroup, "Биохимические исследования", 2, 3); err != nil {
		t.Error(err.Error())
		return
	}

	{
		if err = chechChild(mainGroup.Children[0], "1500", "Антиоксидантный статус Total antioxidant status, TAS", 4680, mainGroup); err != nil {
			t.Error(err.Error())
		}
	}
	{
		group := mainGroup.Children[1]
		if err = chechGroup(group, "Глюкоза и метаболиты углеводного обмена", 3, 4); err != nil {
			t.Error(err.Error())
		} else if err = chechChild(group.Children[0], "16", "Глюкоза (в крови)", 260, group); err != nil {
			t.Error(err.Error())
		} else if err = chechChild(group.Children[1], "", "Прием врача-гастроэнтеролога, к.м.н. повторный", 1800, group); err != nil {
			t.Error(err.Error())
		} else if err = chechChild(group.Children[2], "109", "Глюкоза (в моче)", 265, group); err != nil {
			t.Error(err.Error())
		} else if err = chechChild(group.Children[3], "ГТТ", "Глюкозотолерантный тест с определением глюкозы в венозной крови натощак и после нагрузки через 2 часа", 765, group); err != nil {
			t.Error(err.Error())
		}
	}

	{
		group := mainGroup.Children[2]
		if err = chechGroup(group, "Белки и аминокислоты", 3, 1); err != nil {
			t.Error(err.Error())
		} else if err = chechChild(group.Children[0], "95110", "Альбумин/креатинин-соотношение в разовой порции мочи  (Отношение альбумина к креатинину в разовой порции мочи)  (Albumin-to-creatinine ratio, ACR, random urine)", 750, group); err != nil {
			t.Error(err.Error())
		}
	}
}

// TODO remove
func printPrice(price *Target, level int) {

	fmt.Println(strings.Repeat(" ", level), fmt.Sprintf("%+v", price))

	for i, _ := range price.Children {
		printPrice(price.Children[i], level+1)
	}
}

func chechGroup(val *Target, name string, level int, countChildren int) error {

	if val.Name != name {
		return errors.New(fmt.Sprintf("Failed group name: '%v' != '%v'", val.Name, name))
	} else if val.Level != level {
		return errors.New(fmt.Sprintf("Failed group level: '%v' != '%v'", val.Level, level))
	} else if val.Amount != 0 {
		return errors.New("Failed group amount")
	} else if len(val.Children) != countChildren {
		return errors.New(fmt.Sprintf("Failed count children: '%v' != '%v' (%q)", len(val.Children), countChildren, val.Children))
	}

	return nil
}

func chechChild(val *Target, priceCode string, name string, amount float64, parent *Target) error {

	if val.PriceCode != priceCode {
		return errors.New(fmt.Sprintf("Failed child price code: '%v' != '%v'", val.PriceCode, priceCode))
	} else if val.Name != name {
		return errors.New(fmt.Sprintf("Failed child name: '%v' != '%v'", val.Name, name))
	} else if val.Amount != amount {
		return errors.New(fmt.Sprintf("Failed child amount: '%v' != '%v'", val.Amount, amount))
	} else if val.Parent != parent {
		return errors.New(fmt.Sprintf("Failed child parent: '%v' != '%v'", val.Parent, parent))
	} else if val.Level != 0 {
		return errors.New("Failed child level")
	} else if len(val.Children) != 0 {
		return errors.New("Failed count children")
	}

	return nil
}

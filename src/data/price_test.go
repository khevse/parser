package data

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestReadPrice(t *testing.T) {

	sourceData := readFile("price.html")
	price, _ := ReadPrice("", "", sourceData)

	root := price[0].Children
	if len(root) != 1 {
		t.Error("Failed price")
	}

	var err error

	main_group := root[0]
	if err = chechGroup(main_group, "500", "Анализ крови", 2, 2); err != nil {
		t.Error(err.Error())
	}

	group_1 := main_group.Children[0]
	if err = chechGroup(group_1, "479", "Групповая принадлежность крови", 3, 2); err != nil {
		t.Error(err.Error())
	}

	if err = chechChild(group_1.Children[0], "93", "Группа крови", 360, group_1); err != nil {
		t.Error(err.Error())
	}
	if err = chechChild(group_1.Children[1], "94", "Резус-принадлежность", 360, group_1); err != nil {
		t.Error(err.Error())
	}

	group_2 := main_group.Children[1]
	if err = chechGroup(group_2, "480", "Аллоиммунные антитела", 3, 1); err != nil {
		t.Error(err.Error())
	}

	if err = chechChild(group_2.Children[0], "140", "Аллоиммунные антитела", 675, group_2); err != nil {
		t.Error(err.Error())
	}
}

func chechGroup(val *Price, number string, name string, level uint8, countChildren int) error {

	if val.Number != number {
		return errors.New(fmt.Sprintf("Failed group number: '%v' != '%v'", val.Number, number))
	}

	if val.Name != name {
		return errors.New(fmt.Sprintf("Failed group name: '%v' != '%v'", val.Name, name))
	}

	if val.Level != level {
		return errors.New(fmt.Sprintf("Failed group level: '%v' != '%v'", val.Level, level))
	}

	if val.Amount != 0 {
		return errors.New("Failed group amount")
	}

	if len(val.Children) != countChildren {
		return errors.New(fmt.Sprintf("Failed count children: '%v' != '%v'", len(val.Children), countChildren))
	}

	return nil
}

func chechChild(val *Price, number string, name string, amount float32, parent *Price) error {

	if val.Number != number {
		return errors.New(fmt.Sprintf("Failed child number: '%v' != '%v'", val.Number, number))
	}

	if val.Name != name {
		return errors.New(fmt.Sprintf("Failed child name: '%v' != '%v'", val.Name, name))
	}

	if val.Amount != amount {
		return errors.New(fmt.Sprintf("Failed child amount: '%v' != '%v'", val.Amount, amount))
	}

	if val.Parent != parent {
		return errors.New(fmt.Sprintf("Failed child parent: '%v' != '%v'", val.Parent, parent))
	}

	if val.Level != 0 {
		return errors.New("Failed child level")
	}

	if len(val.Children) != 0 {
		return errors.New("Failed count children")
	}

	return nil
}

func TestGroups(t *testing.T) {

	sourceData := readFile("price_groups.html")
	groups_list := ReadPriceGroups(&sourceData)

	countGroups := 0
	foundGroup1 := false
	foundGroup2 := false

	for i, _ := range groups_list {
		group := &groups_list[i]

		foundGroup1 = foundGroup1 ||
			(strings.Compare(group.Href, "/analizes/for-doctors/137/") == 0 &&
				strings.Compare(group.Name, "Гематологические исследования") == 0)

		foundGroup2 = foundGroup2 ||
			(strings.Compare(group.Href, "/analizes/for-doctors/140/") == 0 &&
				strings.Compare(group.Name, "Биохимические исследования") == 0)

		countGroups += 1
	}

	if !foundGroup1 {
		t.Error("Not found group № 1")
	}
	if !foundGroup2 {
		t.Error("Not found group № 2")
	}
	if countGroups != 2 {
		t.Error("Wrong number")
	}
}

func readFile(name string) []byte {

	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err.Error())
	}

	return data
}

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NikSmith/cache"
	"github.com/khevse/parser/engine"
	"github.com/khevse/parser/example_invitro/model"
	"github.com/khevse/parser/page"
)

var filesCache = cache.New(time.Hour, time.Minute)

type Target struct {
	engine.ITarget

	DBID      interface{}
	Parent    *Target
	Level     int
	Name      string
	PriceCode string
	Amount    float64
	URL       *url.URL
	Children  []*Target
	Sections  map[string][]*Description
}

func (t *Target) Save(price *model.ModelPrice, desc *model.ModelDesc, files *model.ModelFiles) error {

	if t.Level > 0 {
		var parentId interface{}
		if t.Parent != nil {
			parentId = t.Parent.DBID
		}

		id, err := price.AddGroup(parentId, t.Level, t.Name)
		if err != nil {
			log.Println(t.Href(), err)
			return err
		}

		t.DBID = id

	} else {

		id, err := price.AddItem(t.Parent.DBID, t.Name, t.PriceCode, t.Amount)
		if err != nil {
			log.Println(t.Href(), err)
			return err
		}

		t.DBID = id
		for name, data := range t.Sections {

			for _, part := range data {
				if part.Type != ImgDescription {
					continue
				}

				if fileName, err := files.Add(part.Data); err != nil {
					log.Println(t.Href(), part.Url, err)
					part.Text = ""
				} else {
					part.Text = fileName
				}
			}

			if jd, err := json.Marshal(data); err != nil {
				log.Println(t.Href(), err)
				return err
			} else {
				_, err := desc.AddItem(t.DBID, name, string(jd))
				if err != nil {
					log.Println(t.Href(), err)
					return err
				}
			}
		}
	}

	return nil
}

func (t Target) Href() string {
	if t.URL == nil {
		return ""
	} else {
		return t.URL.String()
	}
}

func (t *Target) SetDescription(url string) {

	isGroup := len(url) == 0

	if isGroup {
		return // Nothing do
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	} else {
		src := page.New(req)
		if err := src.Download(); err != nil {
			log.Println(err)
		} else if res, err := parseTarget(src); err != nil {
			log.Println(err)
		} else {
			t.Sections = res
		}
	}
}

func parseTarget(src *page.Page) (map[string][]*Description, error) {

	rootRule := page.NewRule("div", "class", `col_mid`)

	tabsNamesList := src.Nodes(
		rootRule,
		page.NewRule("div", "class", `tabs`, "class", `small_tabs`),
		page.NewRule("li"),
	)

	tabsDesList := src.Nodes(
		rootRule,
		page.NewRule("div", "class", `tabs_block`),
		page.NewRule("div", "class", `tab_cont`),
	)

	if len(tabsNamesList) != len(tabsDesList) {
		msg := fmt.Sprintf("Tabs names not equal tabs description '%s': %d != %d", src.Req.URL.String(), len(tabsNamesList), len(tabsDesList))
		log.Println(msg)
		return nil, errors.New(msg)
	}

	retval := make(map[string][]*Description)

	for i, _ := range tabsDesList {
		header := tabsNamesList[i].GetText()
		section := parseText(tabsDesList[i])

		for _, item := range section {
			if item.Type == ImgDescription {
				item.Url = engine.ChildURL(engine.RootURL(src.Req.URL), item.Url)

				if dt, err := getFileData(item.Url); err != nil {
					return nil, err
				} else {
					item.Data = dt
				}
			}
		}

		retval[header] = section
	}

	return retval, nil
}

func getFileData(u string) ([]byte, error) {

	if dt := filesCache.Get(u); dt != nil {
		return dt.([]byte), nil
	}

	res, err := http.Get(u)
	if err != nil {
		msg := fmt.Sprintf("Failed get file data '%s': %s", u, err.Error())
		log.Println(msg)
		return nil, errors.New(msg)
	}

	defer res.Body.Close()

	if dt, err := ioutil.ReadAll(res.Body); err != nil {
		msg := fmt.Sprintf("Failed get file data '%s': %s", u, err.Error())
		log.Println(msg)
		return nil, errors.New(msg)
	} else {
		filesCache.Set(u, dt)
		return dt, nil
	}
}

func parseText(node *page.Node) []*Description {

	var partsList []*Description

	appendToList := func(item *Description) {
		if len(partsList) > 0 && item.Type == TextDescription && partsList[len(partsList)-1].Type == item.Type {
			partsList[len(partsList)-1].Text += " " + item.Text
		} else {
			partsList = append(partsList, item)
		}
	}

	for _, child := range node.Children {
		for _, item := range parseText(child) {
			appendToList(item)
		}
	}

	if txt := strings.TrimSpace(node.Text); len(txt) > 0 {
		if node.Parent.Tag == "a" {
			item := &Description{
				Text: txt,
				Type: RefDescription,
				Url:  page.TagAttr(node.Parent.Attr, "href"),
			}

			partsList = append(partsList, item)
		} else {
			item := &Description{
				Text: txt,
				Type: TextDescription,
			}

			appendToList(item)
		}

	} else if node.Tag == "img" {

		item := &Description{
			Type: ImgDescription,
			Url:  page.TagAttr(node.Attr, "src"),
		}

		partsList = append(partsList, item)
	}

	return partsList
}

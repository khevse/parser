package engine

import (
	"strconv"
	"testing"

	"github.com/khevse/parser/workers"
)

const (
	COUNT_TARGETS   = 1000
	COUNT_TOC_ITEMS = 1000
)

// TARGET

type Target struct {
	ITarget

	ParantId int
	Id       int
	Url      string
	DescUrl  string
}

func (t Target) Href() string {
	return t.Url
}

func (t *Target) SetDescription(url string) {
	t.DescUrl = url
}

// TOC ITEM

type TOCItem struct {
	ITOCItem

	Id       int
	Url      string
	Children []ITarget
}

func (t TOCItem) Href() string {
	return t.Url
}

func (t *TOCItem) SetChildren(url string) {
	t.Children = make([]ITarget, 0)

	for i := 0; i < COUNT_TARGETS; i++ {
		t.Children = append(t.Children, &Target{
			ParantId: t.Id,
			Id:       i,
			Url:      t.Href() + "/t" + strconv.Itoa(i),
		})
	}
}

func (t *TOCItem) GetChildren() []ITarget {
	return t.Children
}

// TEST

func TestEngine(t *testing.T) {

	toc := make([]ITOCItem, 0)
	for i := 0; i < COUNT_TOC_ITEMS; i++ {
		toc = append(toc, &TOCItem{
			Id:  i,
			Url: "/i" + strconv.Itoa(i),
		})
	}

	tocHandler := func(url string) ([]ITOCItem, error) {
		return toc, nil
	}

	d := new(workers.Dispatcher)
	d.Run()
	defer d.Close()

	e, err := NewEngine("http://example.com", d)
	if err != nil {
		t.Error(err)
	} else if res, err := e.Parse(tocHandler); err != nil {
		t.Error(err)
	} else {
		var count int
		for item := range res {
			target := item.(*Target)

			tartgetUrl := "http://example.com/i" + strconv.Itoa(target.ParantId) + "/t" + strconv.Itoa(target.Id)

			if target.DescUrl != tartgetUrl {
				t.Errorf("'%s' != '%s'", tartgetUrl, target.DescUrl)
			}

			count++
		}

		if count != COUNT_TOC_ITEMS*COUNT_TARGETS {
			t.Error("Fail:", count)
		}
	}

}

func BenchmarkEngine(b *testing.B) {

	toc := make([]ITOCItem, 0)
	for i := 0; i < COUNT_TOC_ITEMS; i++ {
		toc = append(toc, &TOCItem{
			Id:  i,
			Url: "/i",
		})
	}

	tocHandler := func(url string) ([]ITOCItem, error) {
		return toc, nil
	}

	d := new(workers.Dispatcher)
	d.Run()
	defer d.Close()

	e, err := NewEngine("http://example.com", d)
	if err != nil {
		b.Error(err)
	} else if res, err := e.Parse(tocHandler); err != nil {
		b.Error(err)
	} else {
		var count int
		for range res {
			count++
		}

		if count != COUNT_TOC_ITEMS*COUNT_TARGETS {
			b.Error("Fail:", count)
		}
	}
}

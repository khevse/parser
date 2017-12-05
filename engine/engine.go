package engine

import (
	"bytes"
	"errors"
	"net/url"
	"strings"
	"sync"

	"github.com/khevse/parser/workers"
)

type TOCHandler func(url string) ([]ITOCItem, error)

type Engine struct {
	TocURL *url.URL // TOC - table of content

	rootURL    string
	dispatcher *workers.Dispatcher
	once       sync.Once
}

func NewEngine(tocURL string, dispatcher *workers.Dispatcher) (*Engine, error) {
	u, err := url.Parse(tocURL)
	if err != nil {
		return nil, err
	} else if dispatcher == nil {
		return nil, errors.New("Dispatcher is empty.")
	}

	e := Engine{
		dispatcher: dispatcher,
		TocURL:     u,
	}

	return &e, nil
}

func (e *Engine) Parse(tocHandler TOCHandler) (retvalChan <-chan ITarget, retvalErr error) {

	items, err := tocHandler(e.TocURL.String())
	if err != nil {
		retvalErr = err
		return
	}

	channel := make(chan ITarget)
	retvalChan = channel

	go func() {
		wg := new(sync.WaitGroup)

		for item := range e.parseTOCItems(wg, items).GetData() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for target := range e.parseTargets(wg, item.(ITOCItem)).GetData() {
					channel <- target.(ITarget)
				}
			}()
		}

		wg.Wait()
		close(channel)
	}()

	return
}

func (e *Engine) parseTargets(wg *sync.WaitGroup, tocItem ITOCItem) *buffer {

	var list []ITarget = tocItem.GetChildren()

	buf := newBuffer(len(list))

	if len(list) == 0 {
		return buf
	}

	wg.Add(len(list))

	go func() {
		d := e.dispatcher

		for _, item := range list {
			href := item.Href()
			val := taskTarget{
				URL:    e.childURL(href),
				Target: item,
				Result: buf.Queue,
				wg:     wg,
			}
			d.JobQueue <- &val
		}
	}()

	return buf
}

func (e *Engine) parseTOCItems(wg *sync.WaitGroup, list []ITOCItem) *buffer {

	buf := newBuffer(len(list))
	if len(list) == 0 {
		return buf
	}

	wg.Add(len(list))

	go func() {
		d := e.dispatcher
		for _, item := range list {
			href := item.Href()
			val := taskTOCItem{
				URL:     e.childURL(href),
				TOCItem: item,
				Result:  buf.Queue,
				wg:      wg,
			}
			d.JobQueue <- &val
		}
	}()

	return buf
}

func (e *Engine) childURL(href string) string {
	e.internalInit()

	return ChildURL(e.rootURL, href)
}

func (e *Engine) internalInit() {
	e.once.Do(func() {
		e.rootURL = RootURL(e.TocURL)
	})
}

func RootURL(src *url.URL) string {
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	buf.WriteString(src.Scheme)
	buf.WriteString("://")
	buf.WriteString(src.Hostname())

	if len(src.Port()) > 0 {
		buf.WriteString(":")
		buf.WriteString(src.Port())
	}

	return buf.String()
}

func ChildURL(rootURL, href string) string {

	bufLen := len([]byte(rootURL)) + len([]byte(href)) + 1 // +1 - symbol'/'
	bufData := make([]byte, 0, bufLen)

	buf := bytes.NewBuffer(bufData)
	buf.WriteString(rootURL)

	if !strings.HasPrefix(href, "/") {
		buf.WriteString("/")
	}

	buf.WriteString(href)

	return buf.String()
}

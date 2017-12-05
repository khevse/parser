package engine

import (
	"sync"

	"github.com/khevse/parser/workers"
)

type taskTOCItem struct {
	workers.Doer

	URL     string
	TOCItem ITOCItem
	Result  chan<- interface{}
	wg      *sync.WaitGroup
}

func (t *taskTOCItem) Do() {
	t.TOCItem.SetChildren(t.URL)
	t.Result <- t.TOCItem
	t.wg.Done()
}

type taskTarget struct {
	workers.Doer

	URL    string
	Target ITarget
	Result chan<- interface{}
	wg     *sync.WaitGroup
}

func (t *taskTarget) Do() {
	t.Target.SetDescription(t.URL)
	t.Result <- t.Target
	t.wg.Done()
}

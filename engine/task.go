package engine

import "github.com/khevse/parser/workers"

type taskTOCItem struct {
	workers.Doer

	URL     string
	TOCItem ITOCItem
	Result  chan<- interface{}
}

func (t *taskTOCItem) Do() {
	t.TOCItem.SetChildren(t.URL)
	t.Result <- t.TOCItem
}

type taskTarget struct {
	workers.Doer

	URL    string
	Target ITarget
	Result chan<- interface{}
}

func (t *taskTarget) Do() {
	t.Target.SetDescription(t.URL)
	t.Result <- t.Target
}

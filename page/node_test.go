package page

import (
	"fmt"
	"testing"
)

func TestFindNode(t *testing.T) {

	srcList := []struct {
		SrcNode  *Node
		ResNodes []*Node
		Rule     *Rule
	}{
		{
			SrcNode: &Node{
				Tag:      "div",
				Children: []*Node{&Node{Tag: "a"}},
			},
			ResNodes: []*Node{&Node{Tag: "a"}},
			Rule:     NewRule("a"),
		},
	}

	for i, src := range srcList {

		res := src.SrcNode.Find(src.Rule)
		if len(res) != len(src.ResNodes) {
			t.Error(fmt.Sprintf("%d. Wrong result len: %d != %d", i, len(res), len(src.ResNodes)))
			continue
		}

		for itemIndex, _ := range src.ResNodes {
			resNode := res[itemIndex]
			srcNode := src.ResNodes[itemIndex]

			if resNode.Tag != srcNode.Tag || resNode.Text != srcNode.Text || len(resNode.Attr) != len(srcNode.Attr) {
				t.Error(i, ". Wrong node:", *resNode)
			}
		}
	}
}

package page

import "golang.org/x/net/html"

const TEXT_TAG = "text"

type Node struct {
	Parent   *Node
	Tag      string
	Attr     []html.Attribute
	Text     string
	Children []*Node
}

func NewNode(tag string, attr []html.Attribute) *Node {
	return &Node{
		Tag:      tag,
		Attr:     attr,
		Children: make([]*Node, 0, 5),
	}
}

func (n Node) Find(rule *Rule) []*Node {

	retval := make([]*Node, 0, 10)

	for _, child := range n.Children {
		if rule.Equal(child) {
			retval = append(retval, child)
		}
	}

	return retval
}

func (n *Node) NodesWithText() []*Node {

	retval := make([]*Node, 0, 1)
	if len(n.Text) > 0 {
		retval = append(retval, n)
	}

	for _, child := range n.Children {
		retval = append(retval, child.NodesWithText()...)
	}

	return retval
}

func (n *Node) GetText() string {

	textNodes := make([]*Node, 0, len(n.Children))
	for _, node := range n.Children {
		if node.Tag == TEXT_TAG {
			textNodes = append(textNodes, node)
		}
	}

	if len(textNodes) != 1 {
		return ""
	} else {
		return textNodes[0].Text
	}
}

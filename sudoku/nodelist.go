package sudoku

import (
	"container/list"
	// "fmt"
)

type NodeList struct {
	list.List
}

func (nl *NodeList) Append(node *Node) {
	if node.e != nil {
		// panic(fmt.Errorf("node is not allowed to append"))
		return
	}

	e := nl.PushBack(node)
	node.e = e
}

func (nl *NodeList) Remove(node *Node) {
	if node.e == nil {
		return
	}

	nl.List.Remove(node.e)
	node.e = nil
}

func (nl *NodeList) Insert(node *Node) {
	if node.e != nil {
		// panic(fmt.Errorf("node is not allowed to append"))
		return
	}

	e := nl.PushFront(node)
	node.e = e
}

package sudoku

import (
	"container/list"
	"fmt"
)

const (
	DefaultBits int = 0x1FF
)

const (
	X = iota
	Y = iota
	Z = iota
	N = iota
)

func ntb(n int) int {
	return 1 << (n - 1)
}

type Node struct {
	x    int
	y    int
	bits [N + 1]int
	v    int

	n int
	e *list.Element
}

func (node *Node) Set(value int) {
	bit := ntb(value)
	node.bits[X] = bit
	node.bits[Y] = bit
	node.bits[Z] = bit
	node.bits[N] = bit
	node.v = value
	node.n = 1

}

func (node *Node) Init(x, y int) {
	node.x = x
	node.y = y
	node.bits[X] = DefaultBits
	node.bits[Y] = DefaultBits
	node.bits[Z] = DefaultBits
	node.bits[N] = DefaultBits
	node.n = 9
	node.e = nil
	node.v = 0
}

func (node *Node) AppendBit(what int, bit int) {
	/*
		fmt.Printf("[AppendBit][%d,%d]{%d}<%9b, %9b, %9b, %9b>:%9b\n",
			node.x, node.y, what, node.bits[X], node.bits[Y],
			node.bits[Z], node.bits[N], bit)
	*/
	if node.v != 0 {
		return
	}

	node.bits[what] |= bit
	bits := node.bits[N]
	node.bits[N] = node.bits[X] & node.bits[Y] & node.bits[Z]
	bits ^= node.bits[N]
	if bits == 0 {
		return
	}

	if bits == bit {
		node.n += 1
		return
	}

	// panic(fmt.Errorf("AppendBit 不唯一: <%9b, %9b>", bits, bit))
}

func (node *Node) ClearBit(what int, bit int) bool {
	/*
		fmt.Printf("[ClearBit][%d,%d]{%d}<%9b, %9b, %9b, %9b>:%9b\n",
			node.x, node.y, what, node.bits[X], node.bits[Y],
			node.bits[Z], node.bits[N], bit)
	*/
	if node.v != 0 {
		return true
	}

	if node.bits[what]&bit == 0 {
		// fmt.Println("[ClearBit]Failed")
		return false
	}
	if node.bits[N] == bit {
		// fmt.Println("[ClearBit]Failed")
		return false
	}
	node.bits[what] ^= bit

	if node.bits[N]&bit != 0 {
		if node.n == 1 {
			// panic(node)
			return false
		}
		node.bits[N] ^= bit
		node.n -= 1
	}

	return true
}

func (node *Node) Show() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Printf("e: %p\n", node.e)
	fmt.Printf("row: %d\n", node.x)
	fmt.Printf("col: %d\n", node.y)
	fmt.Printf("n: %d\n", node.n)
	fmt.Printf("v: %d\n", node.v)
	fmt.Printf("X: %9b\n", node.bits[X])
	fmt.Printf("Y: %9b\n", node.bits[Y])
	fmt.Printf("Z: %9b\n", node.bits[Z])
	fmt.Printf("N: %9b\n", node.bits[N])
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}

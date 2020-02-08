package sudoku

import (
	"container/list"
	"fmt"
)

type Result [9][9]int

type Sudoku struct {
	maps [9][9]Node
	left [9]NodeList
}

func (sdk *Sudoku) AppendNode(what, i, j, bit int) {
	node := &sdk.maps[i][j]

	if node.v != 0 {
		return
	}

	n := node.n
	node.AppendBit(what, bit)

	if n == node.n {
		return
	}

	ol := &sdk.left[n-1]
	ol.Remove(node)

	nl := &sdk.left[node.n-1]
	nl.Append(node)
}

func (sdk *Sudoku) AppendX(i, bit int) {
	for j := 0; j < 9; j++ {
		sdk.AppendNode(X, i, j, bit)
	}
}

func (sdk *Sudoku) AppendY(j, bit int) {
	for i := 0; i < 9; i++ {
		sdk.AppendNode(Y, i, j, bit)
	}
}

func (sdk *Sudoku) AppendZ(i, j, bit int) {
	x := i / 3 * 3
	y := j / 3 * 3

	for nx := 0; nx < 3; nx++ {
		for ny := 0; ny < 3; ny++ {
			sdk.AppendNode(Z, x+nx, y+ny, bit)
		}
	}
}

func (sdk *Sudoku) Append(i, j, bit int) {
	sdk.AppendX(i, bit)
	sdk.AppendY(j, bit)
	sdk.AppendZ(i, j, bit)
}

func (sdk *Sudoku) UnSet(node *Node, value int) {
	if value < 1 || value > 9 {
		// panic(node)
		return
	}

	// fmt.Printf("[UnSet][%d, %d]:%d\n", node.x, node.y, value)
	sdk.Append(node.x, node.y, ntb(value))
}

func (sdk *Sudoku) Set(node *Node, value int) bool {
	if value < 1 || value > 9 {
		// panic(node)
		return false
	}

	// fmt.Printf("[Set][%d, %d]:%d\n", node.x, node.y, value)
	if node.e != nil {
		sdk.left[node.n-1].Remove(node)
	}

	node.Set(value)
	//sdk.left[0].Append(node)

	bit := ntb(value)
	return sdk.Clear(node.x, node.y, bit)
}

func (sdk *Sudoku) Clear(i, j, bit int) bool {
	if !sdk.ClearX(i, bit) {
		return false
	}
	if !sdk.ClearY(j, bit) {
		return false
	}
	if !sdk.ClearZ(i, j, bit) {
		return false
	}
	return true
}

func (sdk *Sudoku) ClearNode(what, i, j, bit int) bool {
	node := &sdk.maps[i][j]

	if node.v != 0 {
		return true
	}

	n := node.n
	if !node.ClearBit(what, bit) {
		return false
	}

	if n == node.n {
		return true
	}

	ol := &sdk.left[n-1]
	ol.Remove(node)
	//fmt.Printf("[ClearNode][Remove<%d, %d>][List%d:%d]\n", node.x, node.y, n-1, ol.Len())

	nl := &sdk.left[node.n-1]
	nl.Append(node)
	//fmt.Printf("[ClearNode][Append<%d, %d>][List%d:%d]\n", node.x, node.y, node.n-1, nl.Len())
	return true
}

func (sdk *Sudoku) ClearX(i, bit int) bool {
	for j := 0; j < 9; j++ {
		if !sdk.ClearNode(X, i, j, bit) {
			return false
		}
	}
	return true
}

func (sdk *Sudoku) ClearY(j, bit int) bool {
	for i := 0; i < 9; i++ {
		if !sdk.ClearNode(Y, i, j, bit) {
			return false
		}
	}
	return true
}

func (sdk *Sudoku) ClearZ(i, j, bit int) bool {
	x := i / 3 * 3
	y := j / 3 * 3

	for nx := 0; nx < 3; nx++ {
		for ny := 0; ny < 3; ny++ {
			if !sdk.ClearNode(Z, x+nx, y+ny, bit) {
				return false
			}
		}
	}
	return true
}

func New(inputs *Result) *Sudoku {
	var sdk Sudoku
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			node := &sdk.maps[i][j]
			node.Init(i, j)
			sdk.left[8].Append(node)
		}
	}

	for i, rows := range inputs {
		for j, value := range rows {
			if value == 0 {
				continue
			}
			sdk.Set(&sdk.maps[i][j], value)
		}
	}
	return &sdk
}

func (sdk *Sudoku) Result() *Result {
	r := &Result{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			n := &sdk.maps[i][j]
			r[i][j] = n.v
		}
	}

	return r
}

func (sdk *Sudoku) Execute() bool {
	var next *list.Element
	for {
		length := 0
		for n := 0; n < 9; n++ {
			link := &sdk.left[n]
			length += link.Len()
			// fmt.Printf("[Execute][Link %d][Size %d]\n", n, link.Len())
			for e := link.Front(); e != nil; e = next {
				next = e.Next()
				node, err := e.Value.(*Node)
				if !err {
					//panic(node)
					continue
				}
				if node.v != 0 {
					continue
				}
				// 出队，并缓存
				link.Remove(node)
				back := *node

				bits := node.bits[N]
				for i := 1; i <= 9; i++ {
					if ntb(i)&bits == 0 {
						continue
					}
					if !sdk.Set(node, i) {
						sdk.UnSet(node, i)
						continue
					}
					//sdk.Show()
					if !sdk.Execute() {
						sdk.UnSet(node, i)
						continue
					}
					return true
				}
				// 全部不对，还原，入队
				*node = back
				link.Insert(node)
				// fmt.Printf("[Execute][Retry][Link%d][Size %d]\n", n, link.Len())
				// node.Show()
				return false
			}
		}
		if length == 0 {
			break
		}
	}
	return true
}

func (sdk *Sudoku) Show() {
	var result [27][27]int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			node := &sdk.maps[i][j]
			for k := 0; k < 9; k++ {
				if node.bits[N]&(1<<k) == 0 {
					continue
				}
				x := k / 3
				y := k % 3
				result[i*3+x][j*3+y] = k + 1
			}
		}
	}
	fmtstr := "\n-------------------------------------------------------------------------------"
	for i := 0; i < 27; i++ {
		if i%9 == 0 {
			fmt.Printf(fmtstr)
		} else if i%3 == 0 && i > 0 {
			fmt.Println("")
		}
		fmt.Println("")
		for j := 0; j < 27; j++ {
			if j%3 == 0 {
				fmt.Printf("| ")
			}
			if j%9 == 0 && j > 0 {
				fmt.Printf(" | ")
			}
			if result[i][j] == 0 {
				fmt.Printf("  ")
			} else {
				fmt.Printf("%d ", result[i][j])
			}
			if j == 26 {
				fmt.Printf("|")
			}
		}
		if i == 26 {
			fmt.Println(fmtstr)
		}
	}
	// sdk.Debug()
}

func (sdk *Sudoku) Debug() {
	for {
		fmt.Printf("Debug:> ")
		var x, y int
		if _, err := fmt.Scanln(&x, &y); err != nil {
			break
		}
		if x < 1 || x > 9 || y < 1 || y > 9 {
			continue
		}
		node := &sdk.maps[x-1][y-1]
		node.Show()
	}
}

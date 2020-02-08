package main

import (
	"fmt"
	"github.com/YangHaoGG/sudoku/sudoku"
	"github.com/gin-gonic/gin"
	"os"
)

func load(path string) *sudoku.Result {
	fd, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("Open File Failed"))
	}
	defer fd.Close()

	var r sudoku.Result
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Fscan(fd, &r[i][j])
		}
	}

	return &r
}

func posting(c *gin.Context) {
	inputs := &sudoku.Result{}
	if err := c.ShouldBindJSON(inputs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	sdk := sudoku.New(inputs)
	ok := sdk.Execute()
	if !ok {
		c.JSON(400, gin.H{"error": "execute failed"})
		return
	}
	c.JSON(200, gin.H{"data": sdk.Result()})
}

func main() {
	route := gin.Default()
	route.POST("/", posting)
	route.Run("127.0.0.1:8080")
}

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type node struct {
	puzzle [][]int
	children []*node
}

//every operation is maded with a square matrix so i omitted this information

func createMatrix(size int) [][]int {
	m := make([][]int, size)

	for i :=0; i < size; i++ {
		m[i] = make([]int, size)
	}

	return m
}

func printMatrix(m [][]int) {

	size := len(m)

	for i := 0; i < size; i++ {
		fmt.Println(m[i], " ")
	}

	fmt.Println("")

}

func populateMatrix(m [][]int) {

	size := len(m)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			m[i][j] = i * size + j + 1
		}
	}

	//in puzzle15 the last number have value 0
	m[size-1][size-1] = 0

}

func generateRandomMove() int {
	//generate a random move beetween left, right, up, down
	//for simplify the code left => 0, right => 1, up => 2, down => 3
	return rand.Intn(4)
}

func findZero(m [][]int) (int, int) {
	//zero in puzzle15 is the empty space so is crucial to know where it is

	size := len(m)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if m[i][j] == 0 {
				return i,j
			}
		}
	}

	return 0,0
}

func findNumber(m [][]int, n int) (int, int) {
	size := len(m)

	if n == size*size { return findZero(m) }

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if m[i][j] == n {
				return i,j
			}
		}
	}

	return 0,0
}

func positionGoal(n int, size int) (int, int) {

	if n == 0 { return size-1,size-1 }

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if n == i*size+j+1 {
				return i,j
			}
		}
	}

	return 0,0;
}

func swap(a *int, b *int) {
	
	c := *b
	*b = *a
	*a = c

}

func isLegalMove(i int, j int, move int, size int) bool {

	//check if the move is legal

	switch move {
	case 0:
		if j > 0 {return true}
	case 1:
		if j < size - 1 {return true}
	case 2:
		if i > 0 {return true}
	case 3:
		if i < size - 1 {return true}
	}

	return false
}

func makeMove(move int, m [][]int) {

	//left move mean you move the left cell by looking the 0 on the position of the 0 cell, and the same for the others move

	i, j := findZero(m)

	if isLegalMove(i, j, move, len(m)) {

		switch move {

		case 0:
			swap(&m[i][j-1], &m[i][j])
		case 1:
			swap(&m[i][j+1], &m[i][j])
		case 2:
			swap(&m[i-1][j], &m[i][j])
		case 3:
			swap(&m[i+1][j], &m[i][j])

		}
	}
}

func randomizePuzzle(m [][]int) {

	sizeShuffle := 1000

	for i := 0; i < sizeShuffle; i++ {
		move := generateRandomMove()
		makeMove(move, m)
	}

}

func createPuzzle(size int) [][]int {

	matrix := createMatrix(size)
	populateMatrix(matrix)
	randomizePuzzle(matrix)
	return matrix
}

func findLegalMove(puzzle [][]int) (int, int, int, int) {
	
	l, r, u, d := 0, 0, 0, 0	//left, right, up, down

	size := len(puzzle)
	i, j := findZero(puzzle)

	if j > 0 {l = 1}
	if j < size-1 {r = 1}
	if i > 0 {u = 1}
	if i < size-1 {d = 1}

	return l, r, u, d


}

func createNode(puzzle [][]int) *node {
	n := new(node)

	//soft copy for matrix
	size := len(puzzle)
	n.puzzle = createMatrix(size)
	for i :=0; i < size; i++ {
		for j:=0; j < size; j++ {
			n.puzzle[i][j] = puzzle[i][j]
		}
	}

	return n
}

func createChildren(n *node) {
	l, r, u, d := findLegalMove(n.puzzle)
	n.children = make([]*node, 0)

	if l == 1 {
		child := createNode(n.puzzle)
		makeMove(0, child.puzzle)
		n.children = append(n.children, child)
	}

	if r == 1 {
		child := createNode(n.puzzle)
		makeMove(1, child.puzzle)
		n.children = append(n.children, child)
	}

	if u == 1 {
		child := createNode(n.puzzle)
		makeMove(2, child.puzzle)
		n.children = append(n.children, child)
	}

	if d == 1 {
		child := createNode(n.puzzle)
		makeMove(3, child.puzzle)
		n.children = append(n.children, child)
	}

}

func createNodeTree(n *node, levels int) {

	if levels != 0 {

		for i := 0; i < len(n.children); i++ {
			createChildren(n.children[i])
		}

		levels--

		for i := 0; i < len(n.children); i++ {
			createNodeTree(n.children[i], levels)
		}
	}
}

func printTree(n *node, levels int) {

	if n.children != nil {

		fmt.Printf("---------- \n Level %d \n---------- \n", levels)

		for i:=0; i < len(n.children); i++ {
			fmt.Println(manhattanDistance(n.children[i].puzzle))
			printMatrix(n.children[i].puzzle)
		}

		levels--

		for i:=0; i < len(n.children); i++ {
			printTree(n.children[i], levels)
		}
	}
	
}

func absDiffInt(x, y int) int {
	if x < y {
	   return y - x
	}
	return x - y
 }


func manhattanDistance(puzzle [][]int) int {

	distance := 0
	size := len(puzzle)

	for i := 1; i <= size*size; i++ {

		p, q := findNumber(puzzle, i)
		m, n := positionGoal(i, size)

		distance += absDiffInt(p, m) + absDiffInt(q, n)
	}

	return distance
}

func findCandidate(n *node, candidate *node) {

	if n != nil {

		for i:=0; i < len(n.children); i++ {
			
			childrenDistance := manhattanDistance(n.children[i].puzzle)
			candidateDistance := manhattanDistance(candidate.puzzle)
			
			if candidateDistance >= childrenDistance {
				*candidate = *n.children[i]
			}

		}

		if n.children != nil {
			
			for i:=0; i < len(n.children); i++ {
				findCandidate(n.children[i], candidate)
			}
		} 
		
	}


}

func isCorrect(puzzle [][]int) bool {
	size := len(puzzle)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			
			if puzzle[i][j] == 0 && i == size-1 && j == size-1 {
				return true
			}
			
			if puzzle[i][j] != i*size+j+1 {
				return false
			}
		}
	}

	return true
}

func resolvePuzzle(n *node) {
	levels := 8
	for !isCorrect(n.puzzle) {
		createChildren(n)
		createNodeTree(n, levels)
		
		//printTree(n, levels)

		c := createNode(n.puzzle)
		findCandidate(n, c)
		c.children = nil
		n = c

		fmt.Println("candidato scelto:")
		printMatrix(n.puzzle)

	}

	printMatrix(n.puzzle)
	
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func main() {
	
	size := 3
	puzzle := createPuzzle(size)

	root := createNode(puzzle)

	resolvePuzzle(root)


}
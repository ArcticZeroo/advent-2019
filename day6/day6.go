package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"util/datafile"
)

type node struct {
	id       string
	children []*node
	parent   *node
}

func (n node) Depth() int {
	depth := 0
	cur := n.parent
	for cur != nil {
		depth++
		cur = cur.parent
	}
	return depth
}

type treeNodeMap map[string]*node

func (tree treeNodeMap) GetNode(id string) *node {
	if tree[id] == nil {
		tree[id] = &node{id, []*node{}, nil}
	}
	return tree[id]
}

func (tree treeNodeMap) GetRoot() *node {
	for _, node := range tree {
		if node.parent == nil {
			return node
		}
	}
	return nil
}

func sumDepth(n *node, depth int) int {
	total := depth
	for _, child := range n.children {
		total += sumDepth(child, depth+1)
	}
	return total
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func createTree() treeNodeMap {
	tree := treeNodeMap{}

	file := datafile.Open("advent-2019/day6.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		pieces := strings.Split(text, ")")
		parentID := pieces[0]
		childID := pieces[1]
		parent := tree.GetNode(parentID)
		child := tree.GetNode(childID)
		child.parent = parent
		parent.children = append(parent.children, child)
	}

	return tree
}

func part1() {
	tree := createTree()

	root := tree.GetRoot()
	if root == nil {
		log.Fatal("Root not found")
	}

	fmt.Println(sumDepth(root, 0))
}

func part2() {
	tree := createTree()

	source := tree.GetNode("YOU")
	dest := tree.GetNode("SAN")

	sourceDepth := source.Depth()
	destDepth := dest.Depth()

	delta := absInt(sourceDepth - destDepth)
	taller, shorter := source, dest
	if destDepth > sourceDepth {
		taller, shorter = dest, source
	}

	if delta != 0 {
		for i := 0; i < delta; i++ {
			taller = taller.parent
		}
	}

	fmt.Println("Taller depth:", taller.Depth())
	fmt.Println("Shorter depth:", shorter.Depth())

	for taller != shorter {
		taller = taller.parent
		shorter = shorter.parent
		delta += 2
	}

	fmt.Println(delta - 2)
}

func main() {
	//part1()
	part2()
}

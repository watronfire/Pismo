package main

import (
	"bufio"
	"errors"
	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/tree"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

func SliceToMap(slice []string) map[string]struct{} {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	return set
}

func ContainsMap(s map[string]struct{}, q string) bool {
	_, ok := s[q]
	return ok
}

func Contains(s []string, q string) bool {
	for _, i := range s {
		if i == q {
			return true
		}
	}
	return false
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func LoadTree(path string) *tree.Tree {
	var (
		t   *tree.Tree
		err error
		f   *os.File
	)

	if f, err = os.Open(path); err != nil {
		panic(err)
	}
	t, err = newick.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}
	return t
}

func Phylosor(t *tree.Tree, commA, commB map[string]struct{}) float64 {
	blA := 0.0
	blB := 0.0
	blBoth := 0.0

	for _, tip := range t.Tips() {
		if ContainsMap(commA, tip.Name()) {
			blA += GetEdgeLength(tip)
			for ancestor := range AncestorIter(tip) {
				if Contains(ancestor.Comments(), "A") {
					break
				}
				ancestor.AddComment("A")
				ancestorEdge := GetEdgeLength(ancestor)
				blA += ancestorEdge
				if Contains(ancestor.Comments(), "B") {
					blBoth += ancestorEdge
				}
			}
		} else if ContainsMap(commB, tip.Name()) {
			blB += GetEdgeLength(tip)
			for ancestor := range AncestorIter(tip) {
				if Contains(ancestor.Comments(), "B") {
					break
				}
				ancestor.AddComment("B")
				ancestorEdge := GetEdgeLength(ancestor)
				blB += ancestorEdge
				if Contains(ancestor.Comments(), "A") {
					blBoth += ancestorEdge
				}
			}
		}
	}
	return blBoth / (0.5 * (blA + blB))
}

func AncestorIter(node *tree.Node) <-chan *tree.Node {
	currentNode := node
	ch := make(chan *tree.Node)
	go func() {
		for true {
			parentNode, err := currentNode.Parent()
			if err != nil {
				close(ch)
				break
			}
			currentNode = parentNode
			ch <- parentNode
		}
	}()
	return ch
}

func GetEdgeLength(node *tree.Node) float64 {
	edge, err := node.ParentEdge()
	if err != nil {
		return 0
	}
	return edge.Length()
}

func Init() (treeLoc, commA, commB string) {
	// Assign flags
	flag.StringVar(&treeLoc, "tree", "None", "locations of phylogeny")
	flag.StringVar(&commA, "commA", "None", "file containing taxa of community A")
	flag.StringVar(&commB, "commB", "None", "file containing taxa of community B")

	// Parse flags
	flag.Parse()

	// Assert flags are required
	if treeLoc == "None" {
		log.Fatal(errors.New("A tree is required"))
	}
	if commA == "None" {
		log.Fatal(errors.New("taxa in community A must be specified"))
	}
	if commB == "None" {
		log.Fatal(errors.New("taxa in community B must be specified"))
	}

	return treeLoc, commA, commB
}

func main() {
	treeLoc, commA, commB := Init()

	var t = LoadTree(treeLoc)

	commASlice, err := readLines(commA)
	if err != nil {
		panic(err)
	}
	commBSlice, err := readLines(commB)
	if err != nil {
		panic(err)
	}

	commATaxa := SliceToMap(commASlice)
	commBTaxa := SliceToMap(commBSlice)

	got := Phylosor(t, commATaxa, commBTaxa)
	println(got)
}

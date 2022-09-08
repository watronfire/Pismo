package main

import (
	"math"
	"testing"
)

var THRESHOLD = 0.01

func TestPhylosor(t *testing.T) {
	var tree = LoadTree("res/test.tree")
	commASlice, err := readLines("res/commA.txt")
	if err != nil {
		panic(err)
	}
	commBSlice, err := readLines("res/commB.txt")
	if err != nil {
		panic(err)
	}

	commATaxa := SliceToMap(commASlice)
	commBTaxa := SliceToMap(commBSlice)

	got := Phylosor(tree, commATaxa, commBTaxa).GetPhylosor()
	want := 0.2659961328235229
	if math.Abs(got-want) > THRESHOLD {
		t.Errorf("got %f want %f", got, want)
	}
}

func TestLargePhylosor(t *testing.T) {
	var tree = LoadTree("res/real/SanDiego_CA-BritishColumbia_CAN.tree")
	commASlice, err := readLines("res/real/comm1.txt")
	if err != nil {
		panic(err)
	}
	commBSlice, err := readLines("res/real/comm2.txt")
	if err != nil {
		panic(err)
	}
	commATaxa := SliceToMap(commASlice)
	commBTaxa := SliceToMap(commBSlice)

	got := Phylosor(tree, commATaxa, commBTaxa).GetPhylosor()
	want := 0.083859
	if math.Abs(got-want) > THRESHOLD {
		t.Errorf("got %f want %f", got, want)
	}
}

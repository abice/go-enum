package cso

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	src := `1,1.2,abc,"abc",{},{1,true},[1],[1,,],{,,}`
	node, err := Read(strings.NewReader(src))
	if err != nil {
		t.Errorf("Read error: %v", err)
	} else {
		fmt.Fprintf(os.Stdout, "numChild: %d\n", node.NumChild())
		node.output(os.Stdout)
		fmt.Fprintf(os.Stdout, "\n")
	}
}

func TestReadAll(t *testing.T) {
	for _, src := range []string{
		`1,1.2,abc,"abc",{},{1,true},[1],[1,,],{,,}`,
		`1,1.2,abc,"abc",{},{1,true},[1],[1,,],{,,},[]
1,1.2,abc,"abc",`,
	} {
		nodes, err := ReadAll(strings.NewReader(src))
		if err != nil {
			t.Errorf("Read error: %v", err)
		} else {
			for i, node := range nodes {
				fmt.Fprintf(os.Stdout, "%d node numChild: %d\n", i, node.NumChild())
				node.output(os.Stdout)
				fmt.Fprintf(os.Stdout, "\n")
			}
		}
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	ortfodb "github.com/ortfo/db"
	"gopkg.in/yaml.v3"
)

func main() {
	contents, err := os.ReadFile("tags.yaml")
	if err != nil {
		fmt.Printf("while reading file %s: %s", "tags.yaml", err)
	}

	var nodes []yaml.Node
	err = yaml.Unmarshal(contents, &nodes)
	if err != nil {
		fmt.Printf("while loading %s as YAML: %s", "tags.yaml", err)
	}

	var tag ortfodb.Tag
	nodes[0].Decode(&tag)
	spew.Dump(tag)
}

package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	// ortfodb "github.com/ortfo/db"
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

	for _, n := range nodes {
		m := make(map[string]string)
		var k string
		for i, v := range n.Content {
			if i%2 == 0 {
				k = v.Value
			} else {
				m[k] = fmt.Sprintf("%v", v.Kind)
			}
		}

		spew.Dump(m)

		// var tag ortfodb.Tag
		// spew.Dump(n)
	}
}

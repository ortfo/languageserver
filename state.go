package languageserver

import (
	"fmt"
	"os"
	"strings"

	"go.lsp.dev/protocol"
	"gopkg.in/yaml.v3"
)

type DescriptionFile struct {
	contents    string
	lines       []string
	cursor      protocol.Position
	frontmatter *yaml.Node
}

func (d DescriptionFile) CurrentLine() string {
	return d.lines[d.cursor.Line]
}

func (d DescriptionFile) InFrontmatter() (closestKey string, closestNode *yaml.Node, found bool) {
	if d.frontmatter.Line > int(d.cursor.Line) {
		return "", nil, false
	}

	if len(d.frontmatter.Content) == 0 {
		return "", nil, false
	}

	key := ""
	for i, node := range d.frontmatter.Content {
		if i%2 == 0 {
			key = node.Value
			continue
		}

		nextNode := d.frontmatter.Content[i+1]
		if nextNode == nil {
			return "", nil, false
		}
		if isAfter(d.cursor, positionOf(nextNode)) {
			switch node.Kind {
			case yaml.SequenceNode:
				return key, closestNodeInArray(node, d.cursor), true
			case yaml.ScalarNode:
				return key, node, true
			default:
				return "", nil, false
			}
		}
	}
	return "", nil, false
}

func positionOf(node *yaml.Node) protocol.Position {
	return protocol.Position{
		Line:      uint32(node.Line),
		Character: uint32(node.Column),
	}
}

func isAfter(a protocol.Position, b protocol.Position) bool {
	return comparePositions(a, b) > 0
}

func comparePositions(a protocol.Position, b protocol.Position) int {
	if a.Line == b.Line {
		return int(a.Character) - int(b.Character)
	}
	return int(a.Line) - int(b.Line)
}

func closestNodeInArray(node *yaml.Node, pos protocol.Position) *yaml.Node {
	for i, child := range node.Content {
		next := node.Content[i+1]
		if next == nil {
			return child
		}
		if isAfter(pos, positionOf(next)) {
			return child
		}
	}
	return nil
}

func CurrentFile(params protocol.TextDocumentPositionParams) (DescriptionFile, error) {
	contentsRaw, err := os.ReadFile(params.TextDocument.URI.Filename())
	if err != nil {
		return DescriptionFile{}, fmt.Errorf("while reading file at %s: %w", params.TextDocument.URI.Filename(), err)
	}
	contents := string(contentsRaw)
	return DescriptionFile{
		contents: contents,
		lines:    strings.Split(contents, "\n"),
		cursor:   params.Position,
	}, nil
}

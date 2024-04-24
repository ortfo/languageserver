package languageserver

import (
	"fmt"
	"os"
	"strings"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type DescriptionFile struct {
	contents          string
	lines             []string
	cursor            protocol.Position
	frontmatter       *yaml.Node
	frontmatterEndsAt protocol.Position
}

func (d DescriptionFile) CurrentLine() string {
	return d.lines[d.cursor.Line]
}

func (d DescriptionFile) InFrontmatter() (closestKey string, closestNode *yaml.Node, found bool) {
	logger.Debug("InFrontmatter with", zap.Any("d", d))
	if isAfter(d.cursor, d.frontmatterEndsAt) {
		return "", nil, false
	}

	if len(d.frontmatter.Content) == 0 {
		return "", nil, false
	}

	key := ""
	for i, node := range d.frontmatter.Content {
		if i%2 == 0 {
			logger.Debug("InFrontmatter:key", zap.Any("node", node), zap.Int("i", i))
			key = node.Value
			continue
		}

		var nextNode *yaml.Node
		if len(d.frontmatter.Content) <= i+1 {
			nextNode = nil
		} else {
			nextNode = d.frontmatter.Content[i+1]
		}

		logger.Debug("InFrontmatter:value", zap.Any("node", node), zap.Int("i", i), zap.Any("nextNode", nextNode))
		if nextNode == nil || isAfter(positionOf(nextNode), d.cursor) {
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

func CurrentFile(params protocol.TextDocumentPositionParams) (DescriptionFile, error) {
	contentsRaw, err := os.ReadFile(params.TextDocument.URI.Filename())
	if err != nil {
		return DescriptionFile{}, fmt.Errorf("while reading file at %s: %w", params.TextDocument.URI.Filename(), err)
	}
	contents := string(contentsRaw)
	var frontmatter yaml.Node
	frontmatterRaw, frontmatterBoundaryLineNumber := extractFrontmatter(contents)
	if err := yaml.Unmarshal([]byte(frontmatterRaw), &frontmatter); err != nil {
		return DescriptionFile{}, fmt.Errorf("while parsing frontmatter: %w", err)
	}

	frontmatter = *frontmatter.Content[0]

	return DescriptionFile{
		contents:    contents,
		lines:       strings.Split(contents, "\n"),
		cursor:      params.Position,
		frontmatter: &frontmatter,
		frontmatterEndsAt: protocol.Position{
			Line:      uint32(frontmatterBoundaryLineNumber),
			Character: 0,
		},
	}, nil
}

func extractFrontmatter(contents string) (string, int) {
	lines := strings.Split(contents, "\n")
	if len(lines) == 0 {
		return "", 0
	}
	for i, line := range lines {
		if YAMLSeparator.MatchString(line) && i > 0 {
			return strings.Join(lines[:i+1], "\n"), i
		}
	}
	return "", 0
}

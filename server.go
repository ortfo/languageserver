package languageserver

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	ortfodb "github.com/ortfo/db"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
	"gopkg.in/yaml.v3"
)

var YAMLSeparator = regexp.MustCompile(ortfodb.PatternYAMLSeparator)

type DescriptionFile struct {
	contents string
	lines    []string
	cursor   protocol.Position
}

func (d DescriptionFile) YAMLArrayElement() any {
	_, value := d.YAMLMapping()
	switch value := value.(type) {
	case []any:
		commasCountBeforeCursor := strings.Count(d.CurrentLine()[:d.cursor.Character], ",")
		if len(value) > commasCountBeforeCursor {
			return value[commasCountBeforeCursor]
		}
		return nil
	default:
		return nil
	}
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

func (d DescriptionFile) CurrentLine() string {
	return d.lines[d.cursor.Line]
}

func (d DescriptionFile) YAMLMapping() (key string, value any) {
	if !d.InFrontmatter() {
		return "", nil
	}

	var mapping map[string]interface{}
	yaml.Unmarshal([]byte(d.CurrentLine()), &mapping)
	for k, v := range mapping {
		return k, v
	}

	return "", nil
}

func (d DescriptionFile) YAMLKey() string {
	k, _ := d.YAMLMapping()
	return k
}

func (d DescriptionFile) InFrontmatter() bool {
	for _, line := range d.lines[d.cursor.Line:] {
		if YAMLSeparator.MatchString(line) {
			return true
		}
	}
	return false
}

type Handler struct {
	protocol.Server
	config ortfodb.Configuration
	tags   []yaml.Node
}

func (h *Handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	config, err := ortfodb.NewConfiguration("ortfodb.yaml")
	if err != nil {
		return nil, fmt.Errorf("while loading ortfodb configuration from ./ortfodb.yaml: %w", err)
	}

	err = h.LoadTags()
	if err != nil {
		return nil, fmt.Errorf("while loading tags from repository: %w", err)
	}

	h.config = config
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			DefinitionProvider: true,
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "ortfols",
			Version: ortfodb.Version,
		},
	}, nil
}

func (h *Handler) Definition(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	file, err := CurrentFile(params.TextDocumentPositionParams)
	if err != nil {
		return []protocol.Location{}, fmt.Errorf("while getting current file: %w", err)
	}

	if file.InFrontmatter() {
		key, value := file.YAMLMapping()
		switch key {
		case "tags":
			loc, err := h.DefinitionLocationOfTag(value.(string))
			if err != nil {
				return []protocol.Location{}, fmt.Errorf("while getting location of tag %s: %w", value, err)
			}

			return []protocol.Location{loc}, nil
		}
	}
	return []protocol.Location{}, nil
}

func (h *Handler) DefinitionLocationOfTag(name string) (protocol.Location, error) {
	for _, tagNode := range h.tags {
		var tag ortfodb.Tag
		err := tagNode.Decode(&tag)
		if err != nil {
			return protocol.Location{}, fmt.Errorf("while decoding tag node %s : %w", tagNode, err)
		}

		if tag.ReferredToBy(name) {
			return protocol.Location{
				URI: uri.File(h.config.Tags.Repository),
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(tagNode.Line),
						Character: uint32(tagNode.Column),
					},
					End: protocol.Position{
						Line:      uint32(tagNode.Line),
						Character: uint32(tagNode.Column),
					},
				},
			}, nil
		}
	}

	return protocol.Location{}, fmt.Errorf("tag %q not found in repository", name)
}

func (h *Handler) LoadTags() error {
	contents, err := os.ReadFile(h.config.Tags.Repository)
	if err != nil {
		return fmt.Errorf("while reading file %s: %w", h.config.Tags.Repository, err)
	}

	err = yaml.Unmarshal(contents, &h.tags)
	if err != nil {
		return fmt.Errorf("while parsing %s as YAML: %w", h.config.Tags.Repository, err)
	}

	return nil
}

package languageserver

import (
	"context"
	"fmt"
	"os"
	"regexp"

	ortfodb "github.com/ortfo/db"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
	"gopkg.in/yaml.v3"
)

var YAMLSeparator = regexp.MustCompile(ortfodb.PatternYAMLSeparator)

type Handler struct {
	protocol.Server
	config       ortfodb.Configuration
	tags         []yaml.Node
	technologies []yaml.Node
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

	if key, node, inside := file.InFrontmatter(); inside {
		switch key {
		case "tags":
			pos, err := DefinitionLocationOf[ortfodb.Tag](node.Value, h.tags)
			return []protocol.Location{
				{
					URI: uri.File(h.config.Tags.Repository),
					Range: protocol.Range{
						Start: pos,
						End:   pos,
					},
				},
			}, err
		case "made with":
			pos, err := DefinitionLocationOf[ortfodb.Technology](node.Value, h.technologies)
			return []protocol.Location{
				{
					URI: uri.File(h.config.Technologies.Repository),
					Range: protocol.Range{
						Start: pos,
						End:   pos,
					},
				},
			}, err
		}
	}
	return []protocol.Location{}, nil
}

type referrable interface {
	ReferredToBy(string) bool
}

func DefinitionLocationOf[T referrable](name string, repo []yaml.Node) (protocol.Position, error) {
	for _, tagNode := range repo {
		var tag T
		err := tagNode.Decode(&tag)
		if err != nil {
			return protocol.Position{}, fmt.Errorf("while decoding tag node %s : %w", tagNode, err)
		}

		if tag.ReferredToBy(name) {
			return protocol.Position{
				Line:      uint32(tagNode.Line),
				Character: uint32(tagNode.Column),
			}, nil
		}
	}

	return protocol.Position{}, fmt.Errorf("tag %q not found in repository", name)
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

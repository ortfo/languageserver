package languageserver

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type referrable interface {
	ReferredToBy(string) bool
	DisplayName() string
	URLFriendlyName() string
}

func ReferrableDescription(item referrable, description string) protocol.MarkupContent {
	return protocol.MarkupContent{
		Kind: protocol.Markdown,
		Value: heredoc.Docf(`# %s (%s)
			%s
		`, item.DisplayName(), fmt.Sprintf("`%s`", item.URLFriendlyName()), description),
	}
}

func FindInRepository[T referrable](name string, kind string, repo []yaml.Node) (*yaml.Node, *T, error) {
	logger.Debug("InDefintionLocationOf", zap.String("name", name), zap.Any("repo", repo))
	for _, tagNode := range repo {
		var tag T
		err := tagNode.Decode(&tag)
		if err != nil {
			return nil, nil, fmt.Errorf("while decoding node %v : %w", tagNode, err)
		}

		logger.Debug("InDefinitionLocationOf:ReferredToBy?", zap.String("name", name), zap.Any("tag", tag))
		if tag.ReferredToBy(name) {
			return &tagNode, &tag, nil
		}
	}

	return nil, nil, fmt.Errorf("%s %q not found in repository", kind, name)
}

func LoadRepository(at string) ([]yaml.Node, error) {
	contents, err := os.ReadFile(at)
	if err != nil {
		return []yaml.Node{}, fmt.Errorf("while reading file %s: %w", at, err)
	}

	var tags []yaml.Node
	err = yaml.Unmarshal(contents, &tags)
	if err != nil {
		return []yaml.Node{}, fmt.Errorf("while parsing %s as YAML: %w", at, err)
	}

	return tags, nil
}

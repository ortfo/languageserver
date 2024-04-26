package languageserver

import (
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func positionOf(node *yaml.Node) protocol.Position {
	return protocol.Position{
		Line:      uint32(node.Line) - 1,
		Character: uint32(node.Column) - 1,
	}
}

func endPositionOf(node *yaml.Node) protocol.Position{
	return protocol.Position{
		Line:      uint32(node.Line) - 1,
		Character: uint32(node.Column) - 1 + uint32(len(node.Value)),
	}
}

func isAfterCursor(sequenceStyle yaml.Style, node *yaml.Node, cursor protocol.Position) bool {
	logger.Debug("isAfterCursor", zap.Any("node", node), zap.Any("cursor", cursor))
	if sequenceStyle == yaml.FlowStyle {
		logger.Debug("isAfterCursor:flow style", zap.Any("cursor", cursor), zap.Any("node pos", positionOf(node)))
		// look character by character if the sequence is written in flow style
		return isAfter(positionOf(node), cursor)
	} else {
		logger.Debug("isAfterCursor:block style", zap.Any("cursor", cursor), zap.Any("node pos", positionOf(node)))
		// Only look line by line if the sequence is written in block style
		return positionOf(node).Line > cursor.Line
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
	logger.Debug("closestNodeInArray", zap.Any("node", node), zap.Any("pos", pos))
	for i, child := range node.Content {
		if i+1 >= len(node.Content) {
			logger.Debug("closestNodeInArray:on last child", zap.Any("child", child))
			return child
		}

		next := node.Content[i+1]
		logger.Debug("closestNodeInArray:next", zap.Any("child", child), zap.Any("next", next))
		if isAfterCursor(node.Style, next, pos) {
			logger.Debug("closestNodeInArray:found", zap.Any("child", child))
			return child
		}
	}
	return nil
}

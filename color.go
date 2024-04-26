package languageserver

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mazznoer/csscolorparser"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func decodeColorLiteral(raw string) protocol.Color {
	logger.Debug("decodeColorLiteral", zap.String("raw", raw))
	if hexLike(raw) {
		raw = "#" + raw
	}
	color, err := csscolorparser.Parse(raw)
	if err != nil {
		return protocol.Color{}
	}

	logger.Debug("decodeColorLiteral", zap.Any("color", color))

	return protocol.Color{
		Red:   color.R,
		Alpha: color.A,
		Blue:  color.B,
		Green: color.G,
	}
}

func encodeColorLiteral(color protocol.Color) string {
	logger.Debug("encodeColorLiteral", zap.Any("color", color))
	out := "#"
	out += encodeHexComponent(color.Red)
	out += encodeHexComponent(color.Blue)
	out += encodeHexComponent(color.Green)
	if color.Alpha != 1 {
		out += encodeHexComponent(color.Alpha)
	}
	logger.Debug("encodeColorLiteral", zap.String("out", out))
	return out
}

func encodeHexComponent(value float64) string {
	return fmt.Sprintf("%02s", strconv.FormatInt(int64(value*255), 16))
}

func hexLike(s string) bool {
	return regexp.MustCompile(`^[0-9a-fA-F]{3,8}$`).MatchString(s)
}

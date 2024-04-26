package languageserver

import (
	"math"
	"regexp"

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
		Red:   roundToThree(color.R),
		Alpha: roundToThree(color.A),
		Blue:  roundToThree(color.B),
		Green: roundToThree(color.G),
	}
}

func encodeColorLiteral(color protocol.Color) string {
	logger.Debug("encodeColorLiteral", zap.Any("color", color))
	out := csscolorparser.Color{
		R: roundToThree(color.Red),
		G: roundToThree(color.Green),
		B: roundToThree(color.Blue),
		A: roundToThree(color.Alpha),
	}.HexString()
	logger.Debug("encodeColorLiteral", zap.String("out", out))
	return out
}

func hexLike(s string) bool {
	return regexp.MustCompile(`^[0-9a-fA-F]{3,8}$`).MatchString(s)
}

func roundToThree(f float64) float64 {
    return math.Round(f*1_00) / 1_00
}

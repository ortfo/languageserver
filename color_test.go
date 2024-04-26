package languageserver

import (
	"math/rand"
	"testing"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func init() {
	if logger == nil {
		logger, _ = zap.NewDevelopmentConfig().Build()
	}
}

func TestColorEncoding(t *testing.T) {
	for i := 0; i < 20; i++ {
		color := randomColor()
		colorstring := encodeColorLiteral(color)
		decoded := decodeColorLiteral(colorstring)
		encoded := encodeColorLiteral(decoded)
		if colorstring != encoded {
			t.Errorf("Color %d: %s != %s (was decoded to %v)", i, colorstring, encoded, decoded)
		}
	}
}

func TestColorDecoding(t *testing.T) {
	// red
	if !compareColorStructs(protocol.Color{Red: 1, Green: 0, Blue: 0, Alpha: 1}, decodeColorLiteral("red")) {
		t.Errorf("Color 'red' != {1, 0, 0, 1} (was decoded to %v)", decodeColorLiteral("red"))
	}

	// blue
	if !compareColorStructs(protocol.Color{Red: 0, Green: 0, Blue: 1, Alpha: 1}, decodeColorLiteral("blue")) {
		t.Errorf("Color 'blue' != {0, 0, 1, 1} (was decoded to %v)", decodeColorLiteral("blue"))
	}

	// green
	if !compareColorStructs(protocol.Color{Red: 0, Green: 1, Blue: 0, Alpha: 1}, decodeColorLiteral("#00ff00")) {
		t.Errorf("Color '#00ff00' != {0, 1, 0, 1} (was decoded to %v)", decodeColorLiteral("#00ff00"))
	}


	for i := 0; i < 20; i++ {
		color := randomColor()
		originalLogger := *logger
		logger = logger.WithOptions(zap.Fields(zap.Int("test color number", i)))
		encoded := encodeColorLiteral(color)
		decoded := decodeColorLiteral(encoded)
		if !compareColorStructs(color, decoded) {
			t.Errorf("Color %d: %#v != %#v (was encoded to %s)", i, color, decoded, encoded)
		}
		logger = &originalLogger
	}
}

func compareColorStructs(a, b protocol.Color) bool {
	return a.Red == b.Red && a.Green == b.Green && a.Blue == b.Blue && a.Alpha == b.Alpha
}

func randomColor() protocol.Color {
	return protocol.Color{
		Red:   roundToThree(rand.Float64()),
		Green: roundToThree(rand.Float64()),
		Blue:  roundToThree(rand.Float64()),
		Alpha: roundToThree(rand.Float64()),
	}
}

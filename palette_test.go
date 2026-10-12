package tokens

import (
	"strings"
	"testing"
)

func TestOKLCHConvertsToTheSRGBChromeDraws(t *testing.T) {
	cases := map[OKLCH]uint32{
		{L: 0.165, C: 0.005, H: 264, Alpha: 1}:    0x0D0E11,
		{L: 0.62, C: 0.19, H: 250, Alpha: 1}:      0x0088F2,
		{L: 0.704, C: 0.191, H: 22.216, Alpha: 1}: 0xFF6467,
		{L: 0.78, C: 0.15, H: 70, Alpha: 1}:       0xF4A437,
		{L: 0.52, C: 0.13, H: 150, Alpha: 1}:      0x1D7D3E,
		{L: 1, C: 0, H: 0, Alpha: 1}:              0xFFFFFF,
	}
	for in, want := range cases {
		if got := in.SRGB().Hex(); got != want {
			t.Errorf("%+v = %06X, Chrome draws %06X", in, got, want)
		}
	}
}

func TestTheDarkSchemeIsTheWebPaletteWithSurfacesThatClimb(t *testing.T) {
	dark, err := Material(Dark)
	if err != nil {
		t.Fatal(err)
	}
	role := func(name string) RGB {
		for _, r := range dark.Scheme {
			if r.Name == name {
				return r.Color
			}
		}
		t.Fatalf("no %s role", name)
		return RGB{}
	}
	if role("background").Hex() != 0x0D0E11 || role("primary").Hex() != 0x0088F2 {
		t.Fatalf("background %06X, primary %06X: the web's dark background and blue", role("background").Hex(), role("primary").Hex())
	}
	levels := []string{"surfaceContainerLowest", "surfaceContainerLow", "surfaceContainer", "surfaceContainerHigh", "surfaceContainerHighest"}
	for i := 1; i < len(levels); i++ {
		if role(levels[i]).luminance() <= role(levels[i-1]).luminance() {
			t.Errorf("%s is not lighter than %s", levels[i], levels[i-1])
		}
	}
	for _, r := range dark.Status {
		if !strings.HasPrefix(r.Name, "on") {
			continue
		}
		base := strings.ToLower(r.Name[2:3]) + r.Name[3:]
		for _, s := range dark.Status {
			if s.Name == base && Contrast(r.Color, s.Color) < 4.5 {
				t.Errorf("%s on %s has contrast %.1f, under 4.5", r.Name, base, Contrast(r.Color, s.Color))
			}
		}
	}
}

func TestBothModesRenderForKotlinAndSwift(t *testing.T) {
	kotlin, err := Kotlin("one.aircast.map")
	if err != nil {
		t.Fatal(err)
	}
	swift, err := Swift()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"internal val AircastDark = darkColorScheme(", "internal val AircastLight = lightColorScheme(", "val warningContainer = Color(0xFF"} {
		if !strings.Contains(kotlin, want) {
			t.Errorf("Kotlin lacks %q", want)
		}
	}
	for _, want := range []string{"static let aircastDark = MaterialColors(", "static let aircastLight = MaterialColors(", "scrim: Color(hex: 0x000000)\n    )"} {
		if !strings.Contains(swift, want) {
			t.Errorf("Swift lacks %q", want)
		}
	}
}

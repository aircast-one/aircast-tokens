package tokens

import (
	"fmt"
	"maps"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Mode string

const (
	Light Mode = "light"
	Dark  Mode = "dark"
)

type OKLCH struct {
	L, C, H, Alpha float64
}

type RGB struct {
	R, G, B float64
}

var (
	blockPattern = regexp.MustCompile(`(?s)(:root|\.dark)\s*\{(.*?)\}`)
	colorPattern = regexp.MustCompile(`--([a-z0-9-]+):\s*oklch\(([0-9.]+)\s+([0-9.]+)\s+([0-9.]+)(?:\s*/\s*([0-9.]+)%)?\)`)
)

func Palette(mode Mode) (map[string]OKLCH, error) {
	selector := map[Mode]string{Light: ":root", Dark: ".dark"}[mode]
	blocks := blockPattern.FindAllStringSubmatch(string(CSS), -1)
	light := colors(blockBody(blocks, ":root"))
	if mode == Light {
		return light, nil
	}
	if blockBody(blocks, selector) == "" {
		return nil, fmt.Errorf("tokens.css has no %s block", selector)
	}
	return merge(light, colors(blockBody(blocks, selector))), nil
}

func blockBody(blocks [][]string, selector string) string {
	for _, block := range blocks {
		if block[1] == selector {
			return block[2]
		}
	}
	return ""
}

func colors(body string) map[string]OKLCH {
	found := map[string]OKLCH{}
	for _, m := range colorPattern.FindAllStringSubmatch(body, -1) {
		alpha := 1.0
		if m[5] != "" {
			alpha = number(m[5]) / 100
		}
		found[m[1]] = OKLCH{L: number(m[2]), C: number(m[3]), H: number(m[4]), Alpha: alpha}
	}
	return found
}

func merge(base, over map[string]OKLCH) map[string]OKLCH {
	out := maps.Clone(base)
	maps.Copy(out, over)
	return out
}

func number(s string) float64 {
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}

func (c OKLCH) SRGB() RGB {
	h := c.H * math.Pi / 180
	a, b := c.C*math.Cos(h), c.C*math.Sin(h)
	l := math.Pow(c.L+0.3963377774*a+0.2158037573*b, 3)
	m := math.Pow(c.L-0.1055613458*a-0.0638541728*b, 3)
	s := math.Pow(c.L-0.0894841775*a-1.2914855480*b, 3)
	return RGB{
		R: gamma(4.0767416621*l - 3.3077115913*m + 0.2309699292*s),
		G: gamma(-1.2684380046*l + 2.6097574011*m - 0.3413193965*s),
		B: gamma(-0.0041960863*l - 0.7034186147*m + 1.7076147010*s),
	}
}

func gamma(x float64) float64 {
	v := 12.92 * x
	if x > 0.0031308 {
		v = 1.055*math.Pow(x, 1/2.4) - 0.055
	}
	return math.Max(0, math.Min(1, v))
}

func (c RGB) Hex() uint32 {
	channel := func(v float64) uint32 { return uint32(math.Round(v * 255)) }
	return channel(c.R)<<16 | channel(c.G)<<8 | channel(c.B)
}

func (c RGB) luminance() float64 {
	linear := func(v float64) float64 {
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*linear(c.R) + 0.7152*linear(c.G) + 0.0722*linear(c.B)
}

func Contrast(a, b RGB) float64 {
	hi, lo := a.luminance(), b.luminance()
	if lo > hi {
		hi, lo = lo, hi
	}
	return (hi + 0.05) / (lo + 0.05)
}

func Over(top OKLCH, below RGB) RGB {
	fg := top.SRGB()
	blend := func(f, b float64) float64 { return f*top.Alpha + b*(1-top.Alpha) }
	return RGB{R: blend(fg.R, below.R), G: blend(fg.G, below.G), B: blend(fg.B, below.B)}
}

func Mix(a, b OKLCH, weightOfA float64) OKLCH {
	hue := a.H
	switch {
	case a.C < 0.02 && b.C >= 0.02:
		hue = b.H
	case b.C >= 0.02:
		delta := math.Mod(b.H-a.H+540, 360) - 180
		hue = math.Mod(a.H+delta*(1-weightOfA)+360, 360)
	}
	return OKLCH{
		L:     a.L*weightOfA + b.L*(1-weightOfA),
		C:     a.C*weightOfA + b.C*(1-weightOfA),
		H:     hue,
		Alpha: 1,
	}
}

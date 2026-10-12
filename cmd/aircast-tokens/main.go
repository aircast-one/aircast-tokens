package main

import (
	"flag"
	"fmt"
	"os"

	tokens "github.com/aircast-one/aircast-tokens"
)

func main() {
	format := flag.String("format", "css", "css, kotlin or swift")
	pkg := flag.String("package", "one.aircast.map", "Kotlin package for -format kotlin")
	flag.Parse()
	out, err := render(*format, *pkg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := os.Stdout.WriteString(out); err != nil {
		os.Exit(1)
	}
}

func render(format, pkg string) (string, error) {
	switch format {
	case "css":
		return string(tokens.CSS), nil
	case "kotlin":
		return tokens.Kotlin(pkg)
	case "swift":
		return tokens.Swift()
	}
	return "", fmt.Errorf("unknown -format %q: want css, kotlin or swift", format)
}

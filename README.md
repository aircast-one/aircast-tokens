# aircast-tokens

The colour and radius tokens shared by `aircast-flasher` and `aircastd`, as one Tailwind v4 stylesheet: `tokens.css`.

## Consuming

**aircast-flasher** (npm, pinned by commit):

```json
"@aircast/tokens": "github:aircast-one/aircast-tokens#<commit>"
```

```css
@import "@aircast/tokens/tokens.css";
```

**aircastd** (Go tool, pinned in `go.mod`, no Node needed):

```
go get -tool github.com/aircast-one/aircast-tokens/cmd/aircast-tokens@<commit>
go tool aircast-tokens > internal/web/assets/styles/tokens.css
```

Each app keeps its own fonts, type scale and `dark` variant selector; only the palette and radius scale live here.

**Native apps** (Android Compose and iOS SwiftUI in qgroundcontrol) take the same palette as Material 3 colour schemes:

```
go run github.com/aircast-one/aircast-tokens/cmd/aircast-tokens@<commit> -format kotlin -package one.aircast.map > AircastTokens.kt
go run github.com/aircast-one/aircast-tokens/cmd/aircast-tokens@<commit> -format swift > AircastTokens.swift
```

`material.go` maps the tokens onto Material roles once for both platforms: primary and the text colours map one to one, the surface containers step through background, sidebar, card and secondary, outlines are the web's input and border composited onto the background, and containers are the colour mixed into the background in OKLCH.

## Changing a token

Edit `tokens.css`, push, then bump the pin in both apps.

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

## Changing a token

Edit `tokens.css`, push, then bump the pin in both apps.

module github.com/rios0rios0/investmate

go 1.27.0

require (
	github.com/gocolly/colly/v2 v2.3.0
	github.com/olekukonko/tablewriter v1.1.4
	github.com/sirupsen/logrus v1.10.2
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/PuerkitoBio/goquery v1.13.0 // indirect
	github.com/andybalholm/cascadia v1.3.4 // indirect
	github.com/antchfx/htmlquery v1.3.6 // indirect
	github.com/antchfx/xmlquery v1.5.1 // indirect
	github.com/antchfx/xpath v1.3.8 // indirect
	github.com/bits-and-blooms/bitset v1.24.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.28 // indirect
	github.com/nlnwa/whatwg-url v0.6.2 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.3.0 // indirect
	github.com/olekukonko/ll v0.1.8 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

// gobwas/glob v1.0.0 is a breaking major: it removed the `Glob` interface in favour of a
// concrete `*Pattern` type. Every published colly release -- including v2.3.0, the latest --
// still declares `compiledGlob glob.Glob`, so selecting v1.0.0 makes colly fail to compile
// ("undefined: glob.Glob") and takes the whole module down with it. Excluded so that
// dependency automation cannot select it again until colly migrates upstream.
exclude github.com/gobwas/glob v1.0.0

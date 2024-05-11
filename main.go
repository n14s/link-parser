package main

import (
	"fmt"
	"io"
	"strings"
)

var htmlString = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>
`

func main() {
	r := strings.NewReader(htmlString)
	links, err := Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}

// Might move to a seperate package later

// Link represents a link in an HTML document
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and will return a slice of links
func Parse(r io.Reader) ([]Link, error) {
	return nil, nil
}

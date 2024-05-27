package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var htmlExample = `
<html>
<head>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
  <h1>Social stuffs</h1>
  <div>
    <a href="https://www.twitter.com/joncalhoun">
      Check me out on twitter
      <i class="fa fa-twitter" aria-hidden="true">watup</i>
    </a>
    <a href="https://github.com/gophercises">
      Gophercises is on <strong>Github</strong>!
    </a>
  </div>
</body>
</html>
`

const (
	urlTreibhaus = "https://treibhaus.at/programm"
)

func main() {
	htmlString, err := fetchHtml(urlTreibhaus)
	if err != nil {
		return
	}

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

// Event represents a public event
type Event struct {
	Venue       string
	Date        string
	Description string
}

func fetchHtml(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New("request did not respond 200")
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("failed to read response")
	}
	defer response.Body.Close()

	htmlString := string(bytes)

	return htmlString, nil
}

// Parse will take in an HTML document and will return a slice of links
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)

	var links []Link
	for _, n := range nodes {
		links = append(links, buildLink(n))
	}

	return links, nil
}

func buildLink(n *html.Node) Link {
	var link Link
	for _, a := range n.Attr {
		if a.Key == "href" {
			link.Href = a.Val
			break
		}
	}

	link.Text = extractText(n)

	return link
}

func ParseEvents(r io.Reader) ([]Event, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := treibhausEventNodes(doc)

	var events []Event
	for _, n := range nodes {
		events = append(events, buildTreibhausEvent(n))
	}

	return events, nil
}

func treibhausEventNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "id" && strings.Contains(a.Val, "event-") {
				return []*html.Node{n}
			}
		}
	}

	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, linkNodes(c)...)
	}
	return nodes
}

func buildTreibhausEvent(n *html.Node) Event {
	var event Event

	fmt.Println(n)

	return event
}

func buildPmkEvent(n *html.Node) Event {
	var event Event

	return event
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text = text + " " + extractText(c)
	}
	formattedText := strings.Join(strings.Fields(text), " ")
	return formattedText
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, linkNodes(c)...)
	}
	return nodes
}

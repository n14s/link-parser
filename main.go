package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	urlTreibhaus = "https://treibhaus.at/programm"
	urlPmk       = "https://www.pmk.or.at/termine"
)

func main() {
	htmlString, err := fetchHtml(urlPmk)
	if err != nil {
		return
	}

	r := strings.NewReader(htmlString)
	events, err := ParseEvents(r)
	if err != nil {
		panic(err)
	}

	for _, e := range events {
		fmt.Printf("%+v\n", e)
	}
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
	Title       string
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
	// nodes := treibhausEventNodes(doc)
	// var events []Event
	// for _, n := range nodes {
	// 	events = append(events, buildTreibhausEvent(n))
	// }

	nodes := pmkEventNodes(doc)
	var events []Event
	for _, n := range nodes {
		events = append(events, buildPmkEvent(n))
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
		nodes = append(nodes, treibhausEventNodes(c)...)
	}
	return nodes
}

func pmkEventNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "layout--pmktermin") {
				return []*html.Node{n}
			}
		}
	}

	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, pmkEventNodes(c)...)
	}
	return nodes
}

func buildTreibhausEvent(n *html.Node) Event {
	// Find title node
	var event Event

	// Get Date
	if n == nil {
		return Event{}
	}
	c1 := n.FirstChild
	if c1 == nil {
		return Event{}
	}
	c2 := c1.FirstChild
	if c2 == nil {
		return Event{}
	}
	c3 := c2.FirstChild
	if c3 == nil {
		return Event{}
	}
	c4 := c3.FirstChild
	if c4 == nil {
		return Event{}
	}
	c5 := c4.FirstChild
	if c5 == nil {
		return Event{}
	}
	c6 := c5.FirstChild
	if c6 == nil {
		return Event{}
	}

	c7 := ""
	for _, a := range c6.Attr {
		if a.Key == "content" {
			c7 = a.Val
		}
	}

	// Get Title. Starts halfway down the Date Tree
	t1 := c2.NextSibling
	if t1 == nil {
		return Event{}
	}
	t2 := t1.FirstChild
	if t2 == nil {
		return Event{}
	}
	t3 := t2.FirstChild
	if t3 == nil {
		return Event{}
	}
	t4 := t3.FirstChild
	if t4 == nil {
		return Event{}
	}
	t5 := t4.Data

	// Get Description. Starts halfway down the Title Tree
	d1 := t2.NextSibling
	if d1 == nil {
		return Event{}
	}
	d2 := d1.FirstChild
	if d2 == nil {
		return Event{}
	}
	d3 := d2.FirstChild
	if d3 == nil {
		return Event{}
	}
	d4 := d3.Data

	event.Date = c7
	event.Title = t5
	event.Description = d4

	return event
}

// Search recursively for a node with the attribute datetime. When found return the node
func pmkDateNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "time" {
		for _, a := range n.Attr {
			if a.Key == "datetime" {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		foundNode := pmkDateNode(c)
		if foundNode != nil {
			return foundNode
		}
	}

	return nil
}

// Search recursively for a node with the title class. When found return the node
func pmkTitleNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "field--name-field-titel") {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		foundNode := pmkTitleNode(c)
		if foundNode != nil {
			return foundNode
		}
	}

	return nil
}

// Search recursively for a node with the description class. When found return the node
func pmkDescriptionNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "field--type-text-with-summary") {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		foundNode := pmkDescriptionNode(c)
		if foundNode != nil {
			return foundNode
		}
	}

	return nil
}

func buildPmkEvent(n *html.Node) Event {
	var event Event

	// Get date
	dateNode := pmkDateNode(n)
	date := ""
	for _, a := range dateNode.Attr {
		if a.Key == "datetime" {
			date = a.Val
		}
	}

	// Get title
	titleNode := pmkTitleNode(n)
	title := extractText(titleNode)

	// Get description
	descriptionNode := pmkDescriptionNode(n)
	description := extractText(descriptionNode)

	event.Date = date
	event.Title = title
	event.Description = description

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

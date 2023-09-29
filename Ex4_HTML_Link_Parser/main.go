package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: go run main.go <html_file>")
		return
	}

	filename := args[0]
	fmt.Println(filename)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	links := findATags(doc)
	for _, link := range links {
		fmt.Println("Found link:", link.Href)
		fmt.Println("Link text:", link.Text)
		fmt.Println()
	}
}

func findATags(n *html.Node) []Link {
	var links []Link
	stack := []*html.Node{n}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.Type == html.ElementNode && node.Data == "a" {
			var link Link
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					link.Href = attr.Val
				}
			}
			link.Text = extractText(node)
			links = append(links, link)
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			stack = append(stack, c)
		}
	}

	return links
}

func extractText(n *html.Node) string {
	var textBuilder strings.Builder
	stack := []*html.Node{n}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if node.Type == html.TextNode {
			textBuilder.WriteString(strings.TrimSpace(node.Data))
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			stack = append(stack, c)
		}
	}

	return textBuilder.String()
}

package utils

import (
	"strings"

	"golang.org/x/net/html"
)

func traverse(node *html.Node, linksPtr *[]string, recursionCount int) {
	if recursionCount > 100 {
		return
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				*linksPtr = append(*linksPtr, attr.Val)
				break
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverse(child, linksPtr, recursionCount+1)
	}
}

func GetLinks(body string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	var links []string

	traverse(doc, &links, 0)
	return links, nil
}

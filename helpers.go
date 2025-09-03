package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// Inline elements that should not break onto a new line
var inlineElements = map[string]bool{
	"a": true, "span": true, "strong": true, "em": true, "b": true,
	"i": true, "u": true, "small": true, "abbr": true, "code": true,
	"label": true, "mark": true, "q": true, "cite": true,
}

// RenderHTMLIndented renders an html.Node to a well-indented string
func RenderHTMLIndented(n *html.Node) string {
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(&buf, c, 0, false)
	}
	return buf.String()
}

func renderNode(buf *bytes.Buffer, n *html.Node, indent int, inlineParent bool) {
	indentStr := strings.Repeat("  ", indent)

	switch n.Type {
	case html.ElementNode:
		isInline := inlineElements[n.Data]
		if !inlineParent && !isInline {
			buf.WriteString(indentStr)
		}
		buf.WriteString("<" + n.Data)
		for _, attr := range n.Attr {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, attr.Key, attr.Val))
		}
		buf.WriteString(">")

		// If element has children
		if n.FirstChild != nil {
			// Decide if children should be inline
			if !isInline {
				buf.WriteString("\n")
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				renderNode(buf, c, indent+1, isInline)
			}
			if !isInline {
				buf.WriteString(indentStr)
			}
			buf.WriteString("</" + n.Data + ">")
			if !isInline {
				buf.WriteString("\n")
			}
		} else {
			// No children
			buf.WriteString("</" + n.Data + ">")
			if !isInline {
				buf.WriteString("\n")
			}
		}
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			if !inlineParent {
				buf.WriteString(indentStr)
			}
			buf.WriteString(text)
			if !inlineParent {
				buf.WriteString("\n")
			}
		}
	case html.CommentNode:
		if !inlineParent {
			buf.WriteString(indentStr)
		}
		buf.WriteString(fmt.Sprintf("<!--%s-->", n.Data))
		if !inlineParent {
			buf.WriteString("\n")
		}
	case html.DoctypeNode:
		buf.WriteString(fmt.Sprintf("%s<!DOCTYPE %s>\n", indentStr, n.Data))
	}
}
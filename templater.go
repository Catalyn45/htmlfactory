package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"path/filepath"

	"golang.org/x/net/html"
)

func parseFile(fileName string) (*html.Node, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func templateFile(fileName string) {
	doc, err := parseFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "factory" {
			fmt.Println("Found factory")

			replaceTemplate(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	rendered := RenderHTMLIndented(doc)

    ext := filepath.Ext(fileName)
    base := strings.TrimSuffix(fileName, ext)
    newName := fmt.Sprintf("%s.compiled.%s", base, ext[1:])

	err = os.WriteFile(newName, []byte(rendered), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// helper function to extract text content
func replaceTemplate(n *html.Node) {
	variables := make(map[string]string)
	var source string

	for _, attr := range n.Attr {
		if attr.Key == "src" {
			fmt.Println("class attribute:", attr.Val)
			source = attr.Val
		} else {
			variables[attr.Key] = attr.Val
		}
	}

	if source == "" {
		return
	}

	file, err := os.Open(source)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer file.Close()

	docs, err := html.ParseFragment(file, n.Parent)
	if err != nil {
		fmt.Println(err.Error())
	}
	
	for _, doc := range docs {
		replaceVariables(doc, variables)
		n.Parent.InsertBefore(doc, n)
	}

	n.Parent.RemoveChild(n)

	return
}

var re = regexp.MustCompile(`\$\{([A-Za-z_]+[A-Za-z_0-9]*)\}`)

func replaceVariables(n *html.Node, variables map[string]string) {
	if n.Data != "" {
		matches := re.FindAllStringSubmatch(n.Data, -1)
		for _, match := range matches {
			// match[1] contains the variable name inside ${...}
			variableName := match[1]
			fmt.Println(variableName)

			variableValue, ok := variables[variableName]
			if ok {
				result := re.ReplaceAllString(n.Data, variableValue)
				n.Data = result
			}
		}
	}

	if len(n.Attr) > 0 {
		for i, attr := range n.Attr {
			for index, data := range []string {attr.Key, attr.Val} {
				matches := re.FindAllStringSubmatch(data, -1)
				for _, match := range matches {
					// match[1] contains the variable name inside ${...}
					variableName := match[1]
					fmt.Println(variableName)

					variableValue, ok := variables[variableName]
					if ok {
						result := re.ReplaceAllString(data, variableValue)

						if index == 0 {
							attr.Key = result
						} else {
							attr.Val = result
						}
					}
				}

				n.Attr[i] = attr
			}
		}
	}

	// Recurse to children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		replaceVariables(c, variables)
	}
}

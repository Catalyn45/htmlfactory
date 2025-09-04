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

func parseFragmentFile(fileName string, context *html.Node) ([]*html.Node, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	docs, err := html.ParseFragment(file, context)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

const fragmentTag = "fragment"

func walkHtml(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == fragmentTag {
		replaceTemplate(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkHtml(c)
	}
}

var queue StringQueue

func templateFile(fileName string, outdir *string) {
	fullPath, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ok := queue.push(fullPath)
	if !ok {
		fmt.Println("circular fragment import")
		return
	}
	defer queue.pop()

	doc, err := parseFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	walkHtml(doc)

	rendered := RenderHTMLIndented(doc)

	var newName string
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	abs1, err := filepath.Abs(cwd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	abs2, err := filepath.Abs(*outdir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if abs1 == abs2 {
		ext := filepath.Ext(fileName)
		base := strings.TrimSuffix(fileName, ext)
		newName = fmt.Sprintf("%s.compiled.%s", base, ext[1:])
	} else {
		newName = filepath.Join(*outdir, fileName)
	}

	err = os.WriteFile(newName, []byte(rendered), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func getVariables(n *html.Node) map[string]string {
	variables := make(map[string]string)

	for _, attr := range n.Attr {
		variables[attr.Key] = attr.Val
	}

	return variables
}

// helper function to extract text content
func replaceTemplate(n *html.Node) {
	variables := getVariables(n)

	source, ok := variables["src"]
	if !ok {
		return
	}

	fullPath, err := filepath.Abs(source)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ok = queue.push(fullPath)
	if !ok {
		fmt.Println("circular fragment import")
		return
	}
	defer queue.pop()

	docs, err := parseFragmentFile(source, n.Parent)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, doc := range docs {
		replaceVariables(doc, variables)
		walkHtml(doc)

		n.Parent.InsertBefore(doc, n)
	}

	n.Parent.RemoveChild(n)

	return
}

var re = regexp.MustCompile(`\$\{([A-Za-z_]+[A-Za-z_0-9]*)\}`)

func replaceVar(data string, variables map[string]string) string {
	matches := re.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		variableName := match[1]

		variableValue, ok := variables[variableName]
		if ok {
			 data = re.ReplaceAllString(data, variableValue)
		}
	}

	return data
}

func replaceVariables(n *html.Node, variables map[string]string) {
	if n.Type == html.ElementNode && n.Data == "script" {
		return
	}

	if n.Data != "" {
		n.Data = replaceVar(n.Data, variables)
	}

	for i, attr := range n.Attr {
		n.Attr[i].Key = replaceVar(attr.Key, variables)
		n.Attr[i].Val = replaceVar(attr.Val, variables)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		replaceVariables(c, variables)
	}
}

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"path/filepath"

	"sync"

	"golang.org/x/net/html"
)

type Templater struct {
	fileName string
	outdir string
	queue StringQueue
	wg *sync.WaitGroup
}

func parseFile(fileName string) (*html.Node, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
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

func (self *Templater) walkHtml(n *html.Node, content []*html.Node) {
	if n.Type == html.ElementNode && n.Data == fragmentTag {
		self.replaceTemplate(n)
		return
	}

	fragmentContent := false
	newAttrs := []html.Attribute{}
	if content != nil && n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "content" && attr.Val == "fragment" {
				fragmentContent = true
			} else {
				newAttrs = append(newAttrs, attr)
			}
		}
	}

	if fragmentContent {
		n.Attr = newAttrs

		for _, c := range content {
			newC := *c
			n.AppendChild(&newC)
		}
	}

	var nextSibling *html.Node
	for c := n.FirstChild; c != nil; c = nextSibling {
		nextSibling = c.NextSibling

		self.walkHtml(c, content)
	}
}

func (self *Templater) templateFile() {
	defer self.wg.Done()

	fullPath, err := filepath.Abs(self.fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ok := self.queue.push(fullPath)
	if !ok {
		fmt.Println("circular fragment import")
		return
	}
	defer self.queue.pop()

	doc, err := parseFile(self.fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	self.walkHtml(doc, nil)

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

	abs2, err := filepath.Abs(self.outdir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var newName string
	if abs1 == abs2 {
		ext := filepath.Ext(self.fileName)
		base := strings.TrimSuffix(self.fileName, ext)
		newName = fmt.Sprintf("%s.compiled.%s", base, ext[1:])
	} else {
		newName = filepath.Join(self.outdir, self.fileName)
	}

	file, err := os.OpenFile(newName, os.O_CREATE | os.O_WRONLY |  os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	html.Render(file, doc)
}

func getVariables(n *html.Node) map[string]string {
	variables := make(map[string]string)

	for _, attr := range n.Attr {
		variables[attr.Key] = attr.Val
	}

	return variables
}

// helper function to extract text content
func (self *Templater) replaceTemplate(n *html.Node) {
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

	ok = self.queue.push(fullPath)
	if !ok {
		fmt.Println("circular fragment import")
		return
	}
	defer self.queue.pop()

	docs, err := parseFragmentFile(source, n.Parent)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	orphanChilds := []*html.Node {}
	var nextSibling *html.Node
	for c := n.FirstChild; c != nil; c = nextSibling {
		nextSibling = c.NextSibling

		n.RemoveChild(c)

		orphanChilds = append(orphanChilds, c)
	}

	for _, doc := range docs {
		replaceVariables(doc, variables)
		self.walkHtml(doc, orphanChilds)

		n.Parent.InsertBefore(doc, n)
	}

	n.Parent.RemoveChild(n)
}

var re = regexp.MustCompile(`\$\{([A-Za-z_]+[A-Za-z_0-9]*)\}`)

func replaceVar(data string, variables map[string]string) string {
	data = re.ReplaceAllStringFunc(data, func (match string) string {
		variableName := re.FindStringSubmatch(match)[1]

		variableValue, ok := variables[variableName]
		if ok {
			return variableValue
		}

		return match
	})

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

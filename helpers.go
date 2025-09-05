package main

import (
	"fmt"
)

type StringQueue struct {
	container []string
}

func (self *StringQueue) push(item string) bool {
	if self.contains(item) {
		fmt.Println("importing fragment: ", item, " will cause circular dependency")
		return false
	}

	self.container = append(self.container, item)

	return true
}

func (self *StringQueue) pop() string {
	length := len(self.container)
	if length == 0 {
		return ""
	}

	lastIndex := length - 1
	item := self.container[lastIndex]
	self.container = self.container[:lastIndex]

	return item
}

func (self *StringQueue) contains(item string) bool {
	for _, internalItem := range self.container {
		if internalItem == item {
			return true
		}
	}

	return false
}
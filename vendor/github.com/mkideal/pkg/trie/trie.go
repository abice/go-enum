package trie

import (
//"unicode/utf8"
)

type Trie struct {
	root *Node
}

func NewTrie(words []string) *Trie {
	dict := new(Trie)
	dict.root = NewNode("", nil)

	if words != nil {
		for _, word := range words {
			dict.Add(word)
		}
	}
	return dict
}

func (this *Trie) isWordTail(node *Node) bool {
	tail, ok := node.Attr().(bool)
	return ok && tail
}

func (this *Trie) match(word string) (lastMatchNode *Node, deep int) {
	if word == "" {
		return this.root, 0
	}

	list := make([]string, 0, len(word))
	for _, r := range word {
		list = append(list, string(r))
	}

	lastMatchNode = this.root
	deep = 0
	wordsize := len(list)
	node := this.root.FirstChild()
	for node != nil && deep < wordsize {
		if node.Value().(string) == list[deep] {
			lastMatchNode = node
			node = node.FirstChild()
			deep++
		} else {
			node = node.Next()
		}
	}
	return
}

func (this *Trie) Match(word string) bool {
	node, _ := this.match(word)
	return this.isWordTail(node)
}

func (this *Trie) Add(word string) {
	node, deep := this.match(word)
	for i, r := range word {
		if i >= deep {
			node = node.AddChild(string(r))
		}
	}
	node.SetAttr(true)
}

func (this *Trie) AddWords(words []string) {
	for _, word := range words {
		this.Add(word)
	}
}

func (this *Trie) AutoComplete(word string) (add string, match bool, fullmatch bool) {
	node, deep := this.match(word)
	if deep != len(word) {
		return "", false, false
	}
	match = true
	for node.ChildSize() == 1 && !this.isWordTail(node) {
		node = node.FirstChild()
		add += node.Value().(string)
		if node.ChildSize() == 0 {
			fullmatch = true
		}
	}
	return
}

func (this *Trie) AutoCompleteList(word string) (addlist []string, match bool) {
	node, deep := this.match(word)
	if deep != len(word) {
		return []string{}, false
	}
	// TODO
	_ = node
	return
}

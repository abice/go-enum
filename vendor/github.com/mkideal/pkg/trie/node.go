package trie

type Element interface{}
type NodeAttr interface{}

type Node struct {
	e        Element
	parent   *Node
	next     *Node
	children []*Node
	attr     NodeAttr
}

func NewNode(e Element, parent *Node) *Node {
	node := new(Node)
	node.e = e
	node.parent = parent
	node.next = nil
	node.children = make([]*Node, 0, 2)
	return node
}

func (this *Node) Value() Element {
	return this.e
}

func (this *Node) Parent() *Node {
	return this.parent
}

func (this *Node) Next() *Node {
	return this.next
}

func (this *Node) Attr() NodeAttr {
	return this.attr
}

func (this *Node) SetAttr(attr NodeAttr) {
	this.attr = attr
}

func (this *Node) FirstChild() *Node {
	if len(this.children) == 0 {
		return nil
	}
	return this.children[0]
}

func (this *Node) Children() []*Node {
	return this.children
}

func (this *Node) ChildSize() int {
	return len(this.children)
}

func (this *Node) AddChildNode(child *Node) {
	childsize := len(this.children)
	if childsize > 0 {
		this.children[childsize-1].next = child
	}
	this.children = append(this.children, child)
	child.parent = this
}

func (this *Node) AddChild(e Element) *Node {
	node := NewNode(e, this)
	this.AddChildNode(node)
	return node
}

func (this *Node) RemoveChild(index int) {
	if index < 0 || index >= len(this.children) {
		return
	}
	for i, chi := range this.children {
		if i+1 == index {
			if index+1 == len(this.children) {
				chi.next = nil
			} else {
				chi.next = this.children[i+1]
			}
			break
		}
	}
	this.children = append(this.children[:index], this.children[index+1:]...)
}

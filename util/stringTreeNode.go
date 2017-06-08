package util

// StringTreeNode present all underlying data as strings. A StringTreeNode may
// have a value, and can be a node in a tree of other StringTreeNodes.
type StringTreeNode interface {

	// Value gets the wrapped string value.
	Value() string

	// Key gives a string, which should be unique in the set of keys of children
	// of the parent of this node, and describes this StringTreeNode or its
	// value.
	Key() string

	// Parent returns the parent StringTreeNode, or nil if there is no parent.
	Parent() StringTreeNode

	// Children returns the list of child StringTreeNodes of this
	// StringTreeNode.
	Children() []StringTreeNode

	// ChildrenByKey returns a map of child StringTreeNodes, keyed by the keys
	// of the children.
	ChildrenByKey() map[string]StringTreeNode

	// Path gives the uppermost parent StringTreeNode and a list of keys which
	// must be descended to reach this StringTreeNode in the tree.
	Path() (StringTreeNode, []string)

	// SetValue sets the value of this StringTreeNode.
	SetValue(string)

	// AddOrGetChild will create and return an empty child StringTreeNode with
	// the given key if no child with that key exists, or return the child with
	// that key.
	AddOrGetChild(childKey string) (child StringTreeNode)
}

// ObviousStringTreeNode is an obvious implementation of a StringTreeNode.
type ObviousStringTreeNode struct {
	value,
	key string
	parent   StringTreeNode
	children map[string]StringTreeNode
}

func (s *ObviousStringTreeNode) Value() string {
	return s.value
}

func (s *ObviousStringTreeNode) Key() string {
	return s.key
}

func (s *ObviousStringTreeNode) Parent() StringTreeNode {
	return s.parent
}

func (s *ObviousStringTreeNode) Children() []StringTreeNode {
	res := []StringTreeNode{}
	for _, v := range s.ChildrenByKey() {
		res = append(res, v)
	}
	return res
}

func (s *ObviousStringTreeNode) ChildrenByKey() map[string]StringTreeNode {
	return s.children
}

func (s *ObviousStringTreeNode) Path() (StringTreeNode, []string) {
	path := []string{s.Key()}
	var p StringTreeNode = s
	for ; p != nil; p = p.Parent() {
		path = append([]string{p.Key()}, path...)
	}
	return p, path
}

func (s *ObviousStringTreeNode) SetValue(v string) {
	s.value = v
}

func (s *ObviousStringTreeNode) AddOrGetChild(k string) StringTreeNode {
	if s.children == nil {
		s.children = map[string]StringTreeNode{}
	}
	if s.children[k] == nil {
		s.children[k] = &ObviousStringTreeNode{
			key:      k,
			value:    "",
			parent:   s,
			children: map[string]StringTreeNode{},
		}
	}
	return s.children[k]
}

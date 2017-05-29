package util

// StringTreeNode present all underlying data as strings. An StringTreeNode may
// have a value, and can be a node in a tree of other StringTreeNodes.
type StringTreeNode interface {

	// Value gets the wrapped string value.
	Value() string

	// Key gives a string, unique in the set of keys of children of the parent
	// of this node, which describes this StringTreeNode or its value, and which
	// can be used to identify this StringTreeNode within the set of its
	// siblings.
	Key() string

	// Parent returns the parent StringTreeNode, or nil if there is no parent.
	Parent() StringTreeNode

	// Children returns the list of child StringTreeNodes of this
	// StringTreeNode.
	Children() []StringTreeNode

	// ChildrenByKey returns a map of child StringTreeNodes, keyed by the keys
	// extracted from the children.
	ChildrenByKey() map[string]StringTreeNode

	// Path gives the uppermost parent StringTreeNode and a list of keys which
	// must be descended to reach this StringTreeNode in the tree.
	Path() ([]string, StringTreeNode)

	SetValue(string)

	SetChildValue(childKey, value string)

	AddChild(childKey string) (child StringTreeNode)
}

type StringTreeNode20170521 struct {
	value,
	key string
	parent   StringTreeNode
	children map[string]StringTreeNode
}

func (s *StringTreeNode20170521) Value() string {
	return s.value
}

func (s *StringTreeNode20170521) Key() string {
	return s.key
}

func (s *StringTreeNode20170521) Parent() StringTreeNode {
	return s.parent
}

func (s *StringTreeNode20170521) Children() []StringTreeNode {
	res := []StringTreeNode{}
	for _, v := range s.ChildrenByKey() {
		res = append(res, v)
	}
	return res
}

func (s *StringTreeNode20170521) ChildrenByKey() map[string]StringTreeNode {
	return s.children
}

func (s *StringTreeNode20170521) Path() ([]string, StringTreeNode) {
	path := []string{s.Key()}
	var p StringTreeNode = s
	for ; p != nil; p = p.Parent() {
		path = append([]string{p.Key()}, path...)
	}
	return path, p
}

func (s *StringTreeNode20170521) SetValue(v string) {
	s.value = v
}

func (s *StringTreeNode20170521) SetChildValue(k, v string) {
	if s.children == nil {
		s.children = map[string]StringTreeNode{}
	}
	if s.children[k] == nil {
		s.AddChild(k)
	}
	s.children[k].SetValue(v)
}

func (s *StringTreeNode20170521) AddChild(k string) StringTreeNode {
	if s.children == nil {
		s.children = map[string]StringTreeNode{}
	}
	if s.children[k] == nil {
		s.children[k] = &StringTreeNode20170521{
			key:      k,
			value:    "",
			parent:   s,
			children: map[string]StringTreeNode{},
		}
	}
	return s.children[k]
}

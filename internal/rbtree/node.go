package rbtree

type color bool

const (
	black, red color = false, true
)

type node struct {
	value               float64
	parent, left, right *node
	color               color
}

func (n *node) setLeft(l *node) {
	n.left = l
	l.parent = n
}

func (n *node) setRight(r *node) {
	n.right = r
	r.parent = n
}

func (n *node) safeColor() color {
	if n == nil {
		return black
	}
	return n.color
}

func (n *node) grandparent() *node {
	if n == nil || n.parent == nil {
		return nil
	}
	return n.parent.parent
}

func (n *node) sibling() *node {
	if n == nil || n.parent == nil {
		return nil
	}
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

func (n *node) grandparent_uncle() (*node, *node) {
	g := n.grandparent()
	if g == nil {
		return nil, nil
	}
	return g, n.parent.sibling()
}

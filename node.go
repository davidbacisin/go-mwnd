package mwnd

import (
	"fmt"
	"strings"
)

type color bool

const (
	black, red color = false, true
)

// Numeric describes a number-like type that may be stored as a sample in a moving window.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type node[T Numeric] struct {
	value               T
	parent, left, right *node[T]
	color               color
}

func (n *node[T]) setLeft(l *node[T]) {
	if n == nil {
		return
	}

	n.left = l

	if l != nil {
		l.parent = n
	}
}

func (n *node[T]) setRight(r *node[T]) {
	if n == nil {
		return
	}

	n.right = r

	if r != nil {
		r.parent = n
	}
}

func (n *node[T]) safeColor() color {
	if n == nil {
		return black
	}
	return n.color
}

func (n *node[T]) grandparent() *node[T] {
	if n == nil || n.parent == nil {
		return nil
	}
	return n.parent.parent
}

func (n *node[T]) sibling() *node[T] {
	if n == nil || n.parent == nil {
		return nil
	}
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

func (n *node[T]) grandparentAndUncle() (*node[T], *node[T]) {
	g := n.grandparent()
	if g == nil {
		return nil, nil
	}
	return g, n.parent.sibling()
}

func (n *node[T]) String() string {
	var sb strings.Builder
	printHelper(n, 0, &sb)
	return sb.String()
}

func printHelper[T Numeric](n *node[T], level int, sb *strings.Builder) {
	if n == nil {
		sb.WriteString("<empty>")
		return
	}

	if n.left != nil {
		printHelper(n.left, level+1, sb)
	}

	for i := 0; i < level; i++ {
		sb.WriteString("   ")
	}

	if n.color == black {
		sb.WriteString(fmt.Sprintf(" %v \n", n.value))
	} else {
		sb.WriteString(fmt.Sprintf("<%v>\n", n.value))
	}

	if n.right != nil {
		printHelper(n.right, level+1, sb)
	}
}

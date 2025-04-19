package rbtree

import (
	"cmp"
)

type tree[T cmp.Ordered] struct {
	// nodes is a ring buffer of all nodes, pre-allocated to the max capacity of the tree so that
	// memory allocations are minimized during normal operation.
	nodes          []node[T]
	root, min, max *node[T]

	// i represents the oldest node in the tree, which will be replaced
	// by the next inserted value
	i    int
	size int
}

func New[T cmp.Ordered](capacity int) *tree[T] {
	return &tree[T]{
		nodes: make([]node[T], capacity),
		i:     0,
		size:  0,
	}
}

func (t *tree[T]) nodeForInsert() *node[T] {
	next := &t.nodes[t.i]

	// If the node is already in the tree, remove it.
	if next.parent != nil {
		t.delete(next)
	}

	t.i = (t.i + 1) % cap(t.nodes)
	t.size = min(t.size+1, cap(t.nodes))

	// Inserted nodes start as red
	next.color = red
	return next
}

func (t *tree[T]) Size() int {
	return t.size
}

func (t *tree[T]) Insert(v T) {
	n := t.nodeForInsert()
	n.value = v

	if t.root == nil {
		t.root = n
		n.parent = nil
		t.rebalanceForInsert(n)
		return
	}

	p := t.root
	for {
		if v < p.value {
			if p.left == nil {
				p.setLeft(n)
				break
			} else {
				p = p.left
			}
		} else {
			if p.right == nil {
				p.setRight(n)
				break
			} else {
				p = p.right
			}
		}
	}
	t.rebalanceForInsert(n)
}

func (t *tree[T]) rebalanceForInsert(n *node[T]) {
	p := n.parent
	// Case 1
	if p == nil {
		n.color = black
		return
	}

	// Case 2
	if p.color == black {
		return
	}

	// Case 3
	g, u := n.grandparent_uncle()
	if u.safeColor() == red {
		p.color = black
		u.color = black
		g.color = red
		t.rebalanceForInsert(g)
		return
	}

	// Case 4
	if n == p.right && p == g.left {
		t.rotateLeft(p)
		n = n.left
	} else if n == p.left && p == g.right {
		t.rotateRight(p)
		n = n.right
	}

	// Case 5. Reset the parent and grandparent in case that case 4 rotated
	p = n.parent
	g = p.parent
	p.color = black
	g.color = red
	if n == p.left && p == g.left {
		t.rotateRight(g)
	} else if n == p.right && p == g.right {
		t.rotateLeft(g)
	}
}

func (t *tree[T]) replace(old, new *node[T]) {
	if old.parent == nil {
		t.root = new
		new.parent = nil
	} else if old == old.parent.left {
		old.parent.setLeft(new)
	} else {
		old.parent.setRight(new)
	}
}

func (t *tree[T]) swap(a, b *node[T]) {
	if a == b || a == nil || b == nil {
		return
	}

	if b.parent == a || b.parent == nil {
		// Swap to reduce number of conditions
		a, b = b, a
	}

	a.color, b.color = b.color, a.color
	aParent, aLeft, aRight := a.parent, a.left, a.right
	bParent, bLeft, bRight := b.parent, b.left, b.right

	var aWasLeft, bWasLeft bool
	if aParent != nil && aParent.left == a {
		aWasLeft = true
	}

	if bParent != nil && bParent.left == b {
		bWasLeft = true
	}

	a.setLeft(bLeft)
	a.setRight(bRight)
	b.setLeft(aLeft)
	b.setRight(aRight)

	if aParent == b {
		t.replace(b, a)
		if aWasLeft {
			a.setLeft(b)
		} else {
			a.setRight(b)
		}
		return
	}

	t.replace(a, b)
	if bWasLeft {
		bParent.setLeft(a)
	} else {
		bParent.setRight(a)
	}
}

func (t *tree[T]) rotateLeft(n *node[T]) {
	r := n.right
	t.replace(n, r)
	n.setRight(r.left)
	r.setLeft(n)
}

func (t *tree[T]) rotateRight(n *node[T]) {
	l := n.left
	t.replace(n, l)
	n.setLeft(l.right)
	l.setRight(n)
}

func (t *tree[T]) delete(n *node[T]) {
	if n == nil {
		return
	}

	orig := n
	t.size--

	if n.left != nil && n.right != nil {
		// Find the immediate predecessor
		pred := n.left
		for pred.right != nil {
			pred = pred.right
		}

		// Swap places with the in-order predecessor
		t.swap(n, pred)
	}

	// Invariant: n.left, n.right, or both are nil
	child := n.left
	if child == nil {
		child = n.right
	}

	if n.color == black {
		n.color = child.safeColor()
		t.rebalanceForDelete(n)
	}

	t.replace(n, child)

	// root should be black
	if n.parent == nil && child != nil {
		child.color = black
	}

	// Remove the node completely from the tree
	orig.parent = nil
	orig.left = nil
	orig.right = nil
}

func (t *tree[T]) rebalanceForDelete(n *node[T]) {
	p := n.parent
	// Case 1
	if p == nil {
		return
	}

	// Case 2
	s := n.sibling()
	if p.color == black &&
		s != nil &&
		s.safeColor() == black &&
		s.left.safeColor() == black &&
		s.right.safeColor() == black {
		s.color = red
		t.rebalanceForDelete(p)
		return
	}

	// Case 3
	if s.safeColor() == red {
		p.color = red
		s.color = black
		if n == p.left {
			t.rotateLeft(p)
		} else {
			t.rotateRight(p)
		}

		// Reassign p and s after rotation
		p = n.parent
		s = n.sibling()
	}

	// Case 4
	if p.safeColor() == red &&
		s != nil &&
		s.safeColor() == black &&
		s.left.safeColor() == black &&
		s.right.safeColor() == black {
		s.color = red
		p.color = black
		return
	}

	// Case 5
	if n == p.left &&
		s != nil &&
		s.safeColor() == black &&
		s.left.safeColor() == red &&
		s.right.safeColor() == black {
		s.color = red
		s.left.color = black
		t.rotateRight(s)
	} else if n == p.right &&
		s != nil &&
		s.safeColor() == black &&
		s.right.safeColor() == red &&
		s.left.safeColor() == black {
		s.color = red
		s.right.color = black
		t.rotateLeft(s)
	}

	// Case 6
	p = n.parent
	s = n.sibling()
	s.color = p.color
	p.color = black
	if n == p.left && s.right.safeColor() == red {
		s.right.color = black
		t.rotateLeft(p)
	} else if s.left.safeColor() == red {
		s.left.color = black
		t.rotateRight(p)
	}
}

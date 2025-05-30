package mwnd

// fixed aggregates a fixed number of values. Once the capacity is reached, each new value causes
// the oldest value to be evicted from the window.
type fixed[T Numeric] struct {
	// nodes is a ring buffer of all nodes, pre-allocated to the max capacity of the tree so that
	// memory allocations are minimized during normal operation.
	nodes          []node[T]
	root, min, max *node[T]

	// mean is the mean value of all values in the tree
	mean float64

	// m2 is the total sum of squared differences from the mean
	m2 float64

	// i represents the oldest node in the tree, which will be replaced
	// by the next inserted value
	i    int
	size int
}

// enforce compliance with interface
var _ Window[float64] = (*fixed[float64])(nil)

// Fixed initializes a moving window with the fixed capacity for values.
func Fixed[T Numeric](capacity int) *fixed[T] {
	return &fixed[T]{
		nodes: make([]node[T], capacity),
		i:     0,
		size:  0,
	}
}

func (t *fixed[T]) nodeForPut() *node[T] {
	next := &t.nodes[t.i]

	// If the node is already in the tree, remove it.
	if next.parent != nil || next == t.root {
		t.delete(next)
	}

	t.i = (t.i + 1) % cap(t.nodes)
	t.size = min(t.size+1, cap(t.nodes))

	// Inserted nodes start as red
	next.color = red
	return next
}

// Size returns the current number of values in the Window.
func (t *fixed[T]) Size() int {
	return t.size
}

// Min returns the lowest value currently in the Window.
// If the Window has no values, then it returns the zero value.
//
// Time complexity of O(1).
func (t *fixed[T]) Min() T {
	if t == nil || t.min == nil {
		var zero T
		return zero
	}

	return t.min.value
}

// Max returns the highest value currently in the Window.
// If the Window has no values, then it returns the zero value.
//
// Time complexity of O(1).
func (t *fixed[T]) Max() T {
	if t == nil || t.max == nil {
		var zero T
		return zero
	}

	return t.max.value
}

// Mean returns the arithmetic mean of all values currently in the Window.
// If the Window has no values, then it returns 0.0.
//
// Time complexity O(1).
func (t *fixed[T]) Mean() float64 {
	return t.mean
}

// Variance returns the population variance of all values currently in the Window.
// If the Window has no values, then it returns the zero value.
//
// Time complexity of O(1).
func (t *fixed[T]) Variance() float64 {
	if t.size == 0 {
		return 0
	}
	return t.m2 / float64(t.size)
}

// Put adds a new value to the Window. If the Window is at capacity, then the oldest value is
// evicted to be replaced by the new value.
//
// Time complexity of O(log n), where n is the number of values in the Window.
func (t *fixed[T]) Put(v T) {
	n := t.nodeForPut()
	n.value = v

	// Welford's algorithm for online variance, which is a numerically stable approach.
	delta := float64(v) - t.mean
	t.mean += delta / float64(t.size)
	delta2 := float64(v) - t.mean
	t.m2 += delta * delta2

	if t.root == nil {
		t.root = n
		t.min = n
		t.max = n
		n.parent = nil
		t.rebalanceForInsert(n)
		return
	}

	p := t.root
	for {
		if v < p.value {
			if p.left == nil {
				p.setLeft(n)

				if p == t.min {
					t.min = n
				}
				break
			}

			p = p.left
		} else {
			if p.right == nil {
				p.setRight(n)

				if p == t.max {
					t.max = n
				}
				break
			}

			p = p.right
		}
	}

	t.rebalanceForInsert(n)
}

func (t *fixed[T]) rebalanceForInsert(n *node[T]) {
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
	g, u := n.grandparentAndUncle()
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

func (t *fixed[T]) replace(old, n *node[T]) {
	if old.parent == nil {
		t.root = n
		if n != nil {
			n.parent = nil
		}
	} else if old == old.parent.left {
		old.parent.setLeft(n)
	} else {
		old.parent.setRight(n)
	}
}

func (t *fixed[T]) swap(a, b *node[T]) {
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

func (t *fixed[T]) rotateLeft(n *node[T]) {
	r := n.right
	t.replace(n, r)
	n.setRight(r.left)
	r.setLeft(n)
}

func (t *fixed[T]) rotateRight(n *node[T]) {
	l := n.left
	t.replace(n, l)
	n.setLeft(l.right)
	l.setRight(n)
}

func (t *fixed[T]) delete(n *node[T]) {
	if n == nil {
		return
	}

	t.size--

	// Adjust the mean and m2 for the removed value
	if t.size == 0 {
		t.mean = 0
		t.m2 = 0
	} else {
		delta2 := float64(n.value) - t.mean
		t.mean -= delta2 / float64(t.size)
		delta := float64(n.value) - t.mean
		t.m2 -= delta * delta2
	}

	if n.left != nil && n.right != nil {
		// Find the immediate predecessor
		pred := n.left
		for pred.right != nil {
			pred = pred.right
		}

		// Swap places with the in-order predecessor
		t.swap(n, pred)

		// Note that because this node had both left and right
		// children, it couldn't possibly be either min or max.
	}

	// Invariant: n.left, n.right, or both are nil
	child := n.left
	if child == nil {
		child = n.right
	}

	if n == t.min {
		if child == nil {
			t.min = n.parent
		} else {
			t.min = child
		}
	}

	if n == t.max {
		if child == nil {
			t.max = n.parent
		} else {
			t.max = child
		}
	}

	if n.color == black {
		n.color = child.safeColor()
		t.rebalanceForDelete(n)
	}

	t.replace(n, child)

	// Remove the node completely from the tree
	n.parent = nil
	n.left = nil
	n.right = nil
}

func (t *fixed[T]) rebalanceForDelete(n *node[T]) {
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

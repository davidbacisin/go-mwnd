package mwnd

import "testing"

func Test_node_safeColor(t *testing.T) {
	cases := []struct {
		name     string
		node     *node[int]
		expected color
	}{
		{
			name:     "nil node",
			node:     nil,
			expected: black,
		},
		{
			name:     "black node",
			node:     &node[int]{color: black},
			expected: black,
		},
		{
			name:     "red node",
			node:     &node[int]{color: red},
			expected: red,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			actual := c.node.safeColor()
			assertEqual(t, c.expected, actual)
		})
	}
}

func Test_node_relationships(t *testing.T) {
	root := &node[int]{value: 0}
	left := &node[int]{value: 1}
	leftLeft := &node[int]{value: 2}
	leftRight := &node[int]{value: 3}
	right := &node[int]{value: 4}
	rightLeft := &node[int]{value: 5}
	rightRight := &node[int]{value: 6}

	// Link the nodes
	root.setLeft(left)
	left.setLeft(leftLeft)
	left.setRight(leftRight)
	root.setRight(right)
	right.setLeft(rightLeft)
	right.setRight(rightRight)

	cases := []struct {
		name        string
		n           *node[int]
		grandparent *node[int]
		uncle       *node[int]
		sibling     *node[int]
	}{
		{
			name:        "nil",
			n:           nil,
			grandparent: nil,
			uncle:       nil,
			sibling:     nil,
		},
		{
			name:        "root",
			n:           root,
			grandparent: nil,
			uncle:       nil,
			sibling:     nil,
		},
		{
			name:        "left",
			n:           left,
			grandparent: nil,
			uncle:       nil,
			sibling:     right,
		},
		{
			name:        "right",
			n:           right,
			grandparent: nil,
			uncle:       nil,
			sibling:     left,
		},
		{
			name:        "left_left",
			n:           leftLeft,
			grandparent: root,
			uncle:       right,
			sibling:     leftRight,
		},
		{
			name:        "left_right",
			n:           leftRight,
			grandparent: root,
			uncle:       right,
			sibling:     leftLeft,
		},
		{
			name:        "right_left",
			n:           rightLeft,
			grandparent: root,
			uncle:       left,
			sibling:     rightRight,
		},
		{
			name:        "right_right",
			n:           rightRight,
			grandparent: root,
			uncle:       left,
			sibling:     rightLeft,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			grandparent, uncle := c.n.grandparentAndUncle()
			sibling := c.n.sibling()
			assertEqual(t, c.grandparent, grandparent)
			assertEqual(t, c.uncle, uncle)
			assertEqual(t, c.sibling, sibling)
		})
	}
}

func Test_node_String(t *testing.T) {
	root := &node[int]{value: 0}
	left := &node[int]{value: 1}
	leftLeft := &node[int]{value: 2}
	leftRight := &node[int]{value: 3}
	right := &node[int]{value: 4}
	rightLeft := &node[int]{value: 5}
	rightRight := &node[int]{value: 6}

	// Link the nodes
	root.setLeft(left)
	left.setLeft(leftLeft)
	left.setRight(leftRight)
	root.setRight(right)
	right.setLeft(rightLeft)
	right.setRight(rightRight)

	assertEqual(t, "       2 \n    1 \n       3 \n 0 \n       5 \n    4 \n       6 \n", root.String())
}

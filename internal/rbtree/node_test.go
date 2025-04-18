package rbtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, c.expected, actual)
		})
	}
}

func Test_node_relationships(t *testing.T) {
	root := &node[int]{value: 0}
	left := &node[int]{value: 1}
	left_left := &node[int]{value: 2}
	left_right := &node[int]{value: 3}
	right := &node[int]{value: 4}
	right_left := &node[int]{value: 5}
	right_right := &node[int]{value: 6}

	// Link the nodes
	root.setLeft(left)
	left.setLeft(left_left)
	left.setRight(left_right)
	root.setRight(right)
	right.setLeft(right_left)
	right.setRight(right_right)

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
			n:           left_left,
			grandparent: root,
			uncle:       right,
			sibling:     left_right,
		},
		{
			name:        "left_right",
			n:           left_right,
			grandparent: root,
			uncle:       right,
			sibling:     left_left,
		},
		{
			name:        "right_left",
			n:           right_left,
			grandparent: root,
			uncle:       left,
			sibling:     right_right,
		},
		{
			name:        "right_right",
			n:           right_right,
			grandparent: root,
			uncle:       left,
			sibling:     right_left,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			grandparent, uncle := c.n.grandparent_uncle()
			sibling := c.n.sibling()
			assert.Equal(t, c.grandparent, grandparent)
			assert.Equal(t, c.uncle, uncle)
			assert.Equal(t, c.sibling, sibling)
		})
	}
}

func Test_node_String(t *testing.T) {
	root := &node[int]{value: 0}
	left := &node[int]{value: 1}
	left_left := &node[int]{value: 2}
	left_right := &node[int]{value: 3}
	right := &node[int]{value: 4}
	right_left := &node[int]{value: 5}
	right_right := &node[int]{value: 6}

	// Link the nodes
	root.setLeft(left)
	left.setLeft(left_left)
	left.setRight(left_right)
	root.setRight(right)
	right.setLeft(right_left)
	right.setRight(right_right)

	assert.Equal(t, "       2 \n    1 \n       3 \n 0 \n       5 \n    4 \n       6 \n", root.String())
}

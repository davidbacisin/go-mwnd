package rbtree

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTree(values ...int) *tree[int] {
	tr := New[int](len(values))
	for _, v := range values {
		tr.Insert(v)
	}
	return tr
}

func assertRedBlackPropertiesNode[T cmp.Ordered](t *testing.T, n *node[T]) (blackCount int, ok bool) {
	if n == nil {
		return 0, true
	}

	var leftBlackCount, rightBlackCount int

	if n.left != nil {
		if !assert.Less(t, n.left.value, n.value, "left child should have a lesser value than parent") {
			return blackCount, false
		}

		leftCount, leftOk := assertRedBlackPropertiesNode(t, n.left)
		leftBlackCount = leftCount
		if !leftOk {
			return blackCount, false
		}
	}

	if n.right != nil {
		ok = ok && assert.GreaterOrEqual(t, n.right.value, n.value, "right child should have a greater or equal value to parent")

		rightCount, rightOk := assertRedBlackPropertiesNode(t, n.left)
		rightBlackCount = rightCount
		if !rightOk {
			return blackCount, false
		}
	}

	// Red-black properties
	assert.Equal(t, leftBlackCount, rightBlackCount, "should have equal number of black nodes to each leaf")

	if n.safeColor() == red {
		ok = ok && assert.Equal(t, black, n.left.safeColor(), "red node should have black left child")
		ok = ok && assert.Equal(t, black, n.right.safeColor(), "red node should have black right child")
	}

	if n.left.safeColor() == black {
		blackCount++
	}

	if n.right.safeColor() == black {
		blackCount++
	}

	return blackCount, ok
}

func assertRedBlackProperties[T cmp.Ordered](t *testing.T, tr *tree[T]) bool {
	_, ok := assertRedBlackPropertiesNode(t, tr.root)
	return ok
}

func Test_tree_Insert(t *testing.T) {
	t.Run("fully worked example", func(t *testing.T) {
		assert := assert.New(t)

		tr := New[int](11)
		assert.Equal(0, tr.Size())

		tr.Insert(1)
		assertRedBlackProperties(t, tr)
		assert.Equal(1, tr.root.value, "should insert root")

		tr.Insert(22)
		assertRedBlackProperties(t, tr)
		assert.Equal(22, tr.root.right.value, "should insert child")

		tr.Insert(27)
		assertRedBlackProperties(t, tr)
		assert.Equal(22, tr.root.value, "should rotate left")
		assert.Equal(1, tr.root.left.value, "should rotate left")
		assert.Equal(27, tr.root.right.value, "should rotate left")

		tr.Insert(15)
		assertRedBlackProperties(t, tr)
		assert.Equal(22, tr.root.value)
		assert.Equal(1, tr.root.left.value)
		assert.Equal(15, tr.root.left.right.value)

		tr.Insert(6)
		assertRedBlackProperties(t, tr)
		assert.Equal(22, tr.root.value)
		assert.Equal(6, tr.root.left.value, "should rotate right then left")
		assert.Equal(1, tr.root.left.left.value, "should rotate right then left")
		assert.Equal(15, tr.root.left.right.value, "should rotate right then left")
		assert.Equal(red, tr.root.left.right.color)

		tr.Insert(11)
		assertRedBlackProperties(t, tr)
		assert.Equal(22, tr.root.value)
		assert.Equal(6, tr.root.left.value)
		assert.Equal(15, tr.root.left.right.value)
		assert.Equal(black, tr.root.left.right.color, "should recolor 15")
		assert.Equal(11, tr.root.left.right.left.value)

		tr.Insert(17)
		assertRedBlackProperties(t, tr)

		tr.Insert(25)
		assertRedBlackProperties(t, tr)
		assert.Equal(22, tr.root.value)
		assert.Equal(27, tr.root.right.value)
		assert.Equal(25, tr.root.right.left.value)

		tr.Insert(13)
		assertRedBlackProperties(t, tr)
		assert.Equal(15, tr.root.value, "should rotate 15 up to root")
		assert.Equal(6, tr.root.left.value, "should rotate 15 up to root")
		assert.Equal(1, tr.root.left.left.value, "should rotate 15 up to root")
		assert.Equal(11, tr.root.left.right.value, "should rotate 15 up to root")
		assert.Equal(13, tr.root.left.right.right.value, "should rotate 15 up to root")
		assert.Equal(22, tr.root.right.value, "should rotate 15 up to root")
		assert.Equal(17, tr.root.right.left.value, "should rotate 15 up to root")
		assert.Equal(27, tr.root.right.right.value, "should rotate 15 up to root")

		tr.Insert(8)
		assertRedBlackProperties(t, tr)

		tr.Insert(1)
		assertRedBlackProperties(t, tr)
		assert.Equal(15, tr.root.value)
		assert.Equal(6, tr.root.left.value)
		assert.Equal(1, tr.root.left.left.value)
		assert.Equal(1, tr.root.left.left.right.value, "should insert duplicates to the right")

		assert.Equal(11, tr.Size(), "should reach its capacity")
	})

	t.Run("random tree", func(t *testing.T) {
		const size = 100
		tr := New[int](size)
		for range size {
			v := rand.Int()
			tr.Insert(v)
		}

		assertRedBlackProperties(t, tr)
	})
}

func Test_tree_swap(t *testing.T) {
	t.Run("nil and nil", func(t *testing.T) {
		tr := New[int](10)
		tr.swap(nil, nil)
		assert.Equal(t, 0, tr.Size())
	})

	t.Run("root with itself", func(t *testing.T) {
		tr := New[int](10)
		tr.root = &node[int]{value: 1}
		tr.swap(tr.root, tr.root)
		assert.Equal(t, 1, tr.root.value)
	})

	t.Run("root with nil", func(t *testing.T) {
		tr := New[int](10)
		tr.root = &node[int]{value: 1}
		tr.swap(tr.root, nil)
		assert.Equal(t, 1, tr.root.value)
	})

	t.Run("root with left", func(t *testing.T) {
		tr := New[int](10)
		n1 := &node[int]{value: 1}
		n2 := &node[int]{value: 2}
		n3 := &node[int]{value: 3}
		n4 := &node[int]{value: 4}
		n5 := &node[int]{value: 5}
		tr.root = n4
		n4.setLeft(n2)
		n4.setRight(n5)
		n2.setLeft(n1)
		n2.setRight(n3)

		tr.swap(tr.root, tr.root.left)
		assert.Equal(t, n2, tr.root)
		assert.Equal(t, n4, tr.root.left)
		assert.Equal(t, n1, tr.root.left.left)
		assert.Equal(t, n3, tr.root.left.right)
		assert.Equal(t, n5, tr.root.right)
	})

	t.Run("root with right", func(t *testing.T) {
		tr := New[int](10)
		n1 := &node[int]{value: 1}
		n2 := &node[int]{value: 2}
		n3 := &node[int]{value: 3}
		n4 := &node[int]{value: 4}
		n5 := &node[int]{value: 5}
		tr.root = n2
		n2.setLeft(n1)
		n2.setRight(n4)
		n4.setLeft(n3)
		n4.setRight(n5)

		tr.swap(tr.root, tr.root.right)
		assert.Equal(t, n4, tr.root)
		assert.Equal(t, n1, tr.root.left)
		assert.Equal(t, n2, tr.root.right)
		assert.Equal(t, n3, tr.root.right.left)
		assert.Equal(t, n5, tr.root.right.right)
	})

	t.Run("root with left left", func(t *testing.T) {
		tr := New[int](10)
		n1 := &node[int]{value: 1}
		n2 := &node[int]{value: 2}
		n3 := &node[int]{value: 3}
		n4 := &node[int]{value: 4}
		n5 := &node[int]{value: 5}
		n6 := &node[int]{value: 6}
		n7 := &node[int]{value: 7}
		tr.root = n6
		n6.setLeft(n4)
		n6.setRight(n7)
		n4.setLeft(n2)
		n4.setRight(n5)
		n2.setLeft(n1)
		n2.setRight(n3)

		tr.swap(tr.root, tr.root.left.left)
		assert.Equal(t, n2, tr.root)
		assert.Equal(t, n4, tr.root.left)
		assert.Equal(t, n7, tr.root.right)
		assert.Equal(t, n6, tr.root.left.left)
		assert.Equal(t, n5, tr.root.left.right)
		assert.Equal(t, n1, tr.root.left.left.left)
		assert.Equal(t, n3, tr.root.left.left.right)
	})

	t.Run("root with left right", func(t *testing.T) {
		tr := New[int](10)
		n1 := &node[int]{value: 1}
		n2 := &node[int]{value: 2}
		n3 := &node[int]{value: 3}
		n4 := &node[int]{value: 4}
		n5 := &node[int]{value: 5}
		n6 := &node[int]{value: 6}
		n7 := &node[int]{value: 7}
		tr.root = n6
		n6.setLeft(n4)
		n6.setRight(n7)
		n4.setLeft(n2)
		n4.setRight(n5)
		n2.setLeft(n1)
		n2.setRight(n3)

		tr.swap(tr.root, tr.root.left.right)
		assert.Equal(t, n5, tr.root)
		assert.Equal(t, n4, tr.root.left)
		assert.Equal(t, n7, tr.root.right)
		assert.Equal(t, n2, tr.root.left.left)
		assert.Equal(t, n6, tr.root.left.right)
		assert.Equal(t, n1, tr.root.left.left.left)
		assert.Equal(t, n3, tr.root.left.left.right)
	})

	t.Run("left with grandchild", func(t *testing.T) {
		tr := New[int](10)
		n1 := &node[int]{value: 1}
		n2 := &node[int]{value: 2}
		n3 := &node[int]{value: 3}
		n4 := &node[int]{value: 4}
		n5 := &node[int]{value: 5}
		n6 := &node[int]{value: 6}
		n7 := &node[int]{value: 7}
		tr.root = n6
		n6.setLeft(n4)
		n6.setRight(n7)
		n4.setLeft(n2)
		n4.setRight(n5)
		n2.setLeft(n1)
		n2.setRight(n3)

		tr.swap(n4, n1)
		assert.Equal(t, n6, tr.root)
		assert.Equal(t, n1, tr.root.left)
		assert.Equal(t, n7, tr.root.right)
		assert.Equal(t, n2, tr.root.left.left)
		assert.Equal(t, n5, tr.root.left.right)
		assert.Equal(t, n4, tr.root.left.left.left)
		assert.Equal(t, n3, tr.root.left.left.right)
	})

	t.Run("right with grandchild", func(t *testing.T) {
		tr := New[int](10)
		n1 := &node[int]{value: 1}
		n2 := &node[int]{value: 2}
		n3 := &node[int]{value: 3}
		n4 := &node[int]{value: 4}
		n5 := &node[int]{value: 5}
		n6 := &node[int]{value: 6}
		n7 := &node[int]{value: 7}
		tr.root = n2
		n2.setLeft(n1)
		n2.setRight(n4)
		n4.setLeft(n3)
		n4.setRight(n6)
		n6.setLeft(n5)
		n6.setRight(n7)

		tr.swap(n4, n7)
		assert.Equal(t, n2, tr.root)
		assert.Equal(t, n1, tr.root.left)
		assert.Equal(t, n7, tr.root.right)
		assert.Equal(t, n3, tr.root.right.left)
		assert.Equal(t, n6, tr.root.right.right)
		assert.Equal(t, n5, tr.root.right.right.left)
		assert.Equal(t, n4, tr.root.right.right.right)
	})
}

func Test_tree_delete(t *testing.T) {
	t.Run("remove leaf, no rotate", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(1, 22, 27, 15, 6, 11, 17, 25, 13, 8, 1)

		p := tr.root.left.left
		assert.Equal(1, p.value)
		n := p.right
		assert.Equal(1, n.value)
		tr.delete(n)
		assert.Equal(10, tr.Size())
		assertRedBlackProperties(t, tr)
		assert.Nil(n.parent)
		assert.Nil(n.left)
		assert.Nil(n.right)
		assert.Nil(p.left)
		assert.Nil(p.right)
	})

	t.Run("replace parent with child; case 4", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(1, 22, 27, 15, 6, 11, 17, 25, 13, 8, 1)

		p := tr.root.right
		assert.Equal(22, p.value)
		n := p.right
		assert.Equal(27, n.value)
		tr.delete(n)
		assert.Equal(10, tr.Size())
		assertRedBlackProperties(t, tr)
		assert.Nil(n.parent)
		assert.Nil(n.left)
		assert.Nil(n.right)
		assert.Equal(17, p.left.value)
		assert.Equal(25, p.right.value)
	})

	t.Run("remove parent with two children; cases 5 right and 6 left", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(1, 22, 27, 15, 6, 11, 17, 25, 13, 8, 1)

		p := tr.root
		assert.Equal(15, p.value)
		n := p.right
		assert.Equal(22, n.value)
		tr.delete(n)
		assert.Equal(10, tr.Size())
		assertRedBlackProperties(t, tr)
		assert.Nil(n.parent)
		assert.Nil(n.left)
		assert.Nil(n.right)
		assert.Equal(tr.root, p, "should keep 15 at root")
		assert.Equal(25, p.right.value)
		assert.Equal(17, p.right.left.value)
		assert.Equal(27, p.right.right.value)
	})

	t.Run("case 3 rotate left", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(5, 8, 1, 7, 9, 6)

		p := tr.root
		assert.Equal(5, p.value)
		n := p.left
		assert.Equal(1, n.value)
		tr.delete(n)
		assertRedBlackProperties(t, tr)
		assert.Nil(n.parent)
		assert.Nil(n.left)
		assert.Nil(n.right)
		assert.Equal(8, tr.root.value)
		assert.Equal(6, tr.root.left.value)
		assert.Equal(9, tr.root.right.value)
		assert.Equal(5, tr.root.left.left.value)
		assert.Equal(7, tr.root.left.right.value)
	})

	t.Run("case 3 rotate left", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(5, 8, 1, 7, 9, 6)

		p := tr.root
		assert.Equal(5, p.value)
		n := p.left
		assert.Equal(1, n.value)
		tr.delete(n)
		assertRedBlackProperties(t, tr)
		assert.Nil(n.parent)
		assert.Nil(n.left)
		assert.Nil(n.right)
		assert.Equal(8, tr.root.value)
		assert.Equal(6, tr.root.left.value)
		assert.Equal(9, tr.root.right.value)
		assert.Equal(5, tr.root.left.left.value)
		assert.Equal(7, tr.root.left.right.value)
	})

	t.Run("case 3 rotate right", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(5, 8, 2, 1, 3, 4)

		p := tr.root
		assert.Equal(5, p.value)
		n := p.right
		assert.Equal(8, n.value)
		tr.delete(n)
		assertRedBlackProperties(t, tr)
		assert.Nil(n.parent)
		assert.Nil(n.left)
		assert.Nil(n.right)
		assert.Equal(2, tr.root.value)
		assert.Equal(1, tr.root.left.value)
		assert.Equal(4, tr.root.right.value)
		assert.Equal(3, tr.root.right.left.value)
		assert.Equal(5, tr.root.right.right.value)
	})

	t.Run("case 2", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(5, 2, 8, 6)

		p := tr.root
		assert.Equal(5, p.value)
		n := p.left
		assert.Equal(2, n.value)

		// Delete the 6 to get the tree in the correct state
		n6 := p.right.left
		assert.Equal(6, n6.value)
		tr.delete(n6)

		assert.Equal(black, p.color)
		assert.Equal(black, n.color)
		assert.Equal(black, p.right.color)

		// Now it will trigger delete case 2
		tr.delete(n)
		assertRedBlackProperties(t, tr)
		assert.Nil(n.parent)
		assert.Nil(n.left)
		assert.Nil(n.right)
		assert.Equal(5, tr.root.value)
		assert.Nil(tr.root.left)
		assert.Equal(8, tr.root.right.value)
	})
}

func Test_tree_rollingWindowAtCapacity(t *testing.T) {
	t.Run("single node", func(t *testing.T) {
		assert := assert.New(t)
		tr := New[int](1)
		tr.Insert(1)
		assert.Equal(1, tr.Size())
		assert.Equal(1, tr.root.value)
		tr.Insert(2)
		assert.Equal(1, tr.Size())
		assert.Equal(2, tr.root.value, "should replace existing value")
	})

	t.Run("three nodes", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(1, 2, 3)
		assert.Equal(3, tr.Size())
		assert.Equal(2, tr.root.value)
		assert.Equal(1, tr.root.left.value)
		assert.Equal(3, tr.root.right.value)

		tr.Insert(4)
		assert.Equal(3, tr.Size(), "should replace oldest value")
		assert.Equal(3, tr.root.value)
		assert.Equal(2, tr.root.left.value)
		assert.Equal(4, tr.root.right.value)
	})

	t.Run("three nodes replace root", func(t *testing.T) {
		assert := assert.New(t)
		tr := makeTree(3, 1, 5)
		assert.Equal(3, tr.Size())
		assert.Equal(3, tr.root.value)
		assert.Equal(1, tr.root.left.value)
		assert.Equal(5, tr.root.right.value)

		tr.Insert(4)
		assert.Equal(3, tr.Size(), "should replace oldest value at root")
		assert.Equal(4, tr.root.value)
		assert.Nil(tr.root.parent)
		assert.Equal(1, tr.root.left.value)
		assert.Equal(5, tr.root.right.value)
	})
}

func Test_tree_minMax(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tr := New[int](1)
		assert.Equal(t, 0, tr.Min())
		assert.Equal(t, 0, tr.Max())
	})

	t.Run("single node", func(t *testing.T) {
		tr := New[int](1)
		tr.Insert(5)
		assert.Equal(t, 5, tr.Min())
		assert.Equal(t, 5, tr.Max())
		tr.Insert(6)
		assert.Equal(t, 6, tr.Min())
		assert.Equal(t, 6, tr.Max())
	})

	t.Run("three nodes", func(t *testing.T) {
		tr := makeTree(2, 1, 3)
		assert.Equal(t, 1, tr.Min())
		assert.Equal(t, 3, tr.Max())
	})

	t.Run("rolling three nodes", func(t *testing.T) {
		tr := makeTree(1, 2, 3)
		tr.Insert(4) // replaces 1
		assert.Equal(t, 2, tr.Min())
		assert.Equal(t, 4, tr.Max())
		tr.Insert(1) // replaces 2
		tr.Insert(2) // replaces 3
		tr.Insert(3) // replaces 4
		assert.Equal(t, 1, tr.Min())
		assert.Equal(t, 3, tr.Max())
	})

	t.Run("rolling 50 nodes random", func(t *testing.T) {
		const size = 50
		values := make([]int, 0, size)
		tr := New[int](size)
		for i := 0; i < 1000; i++ {
			v := rand.Int()
			if i >= size {
				k := i % size
				values[k] = v
			} else {
				values = append(values, v)
			}

			tr.Insert(v)
			expectedMin := slices.Min(values)
			expectedMax := slices.Max(values)
			if !assert.Equal(t, expectedMin, tr.Min(), "min should match") ||
				!assert.Equal(t, expectedMax, tr.Max(), "max should match") {
				break
			}
		}
	})
}

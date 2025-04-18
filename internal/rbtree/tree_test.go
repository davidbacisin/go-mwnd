package rbtree

import (
	"cmp"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

package mwnd

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func makeFixed(values ...int) *FixedWindow[int] {
	tr := Fixed[int](len(values))
	for _, v := range values {
		tr.Put(v)
	}
	return tr
}

func assertRedBlackPropertiesNode[T Numeric](t *testing.T, n *node[T]) (total int, blackCount int, ok bool) {
	if n == nil {
		return 0, 0, true
	}

	total, blackCount, ok = 1, 0, true
	if n.safeColor() == black {
		blackCount++
	}

	var (
		leftTotal, leftBlack, rightTotal, rightBlack int
		leftOk, rightOk                              bool
	)

	if n.left != nil {
		ok = ok && assertLessOrEqual(t, n.left.value, n.value, "left child should have a lesser or equal value to parent")

		leftTotal, leftBlack, leftOk = assertRedBlackPropertiesNode(t, n.left)
		total += leftTotal

		if !leftOk {
			return total, blackCount, false
		}
	}

	if n.right != nil {
		ok = ok && assertLessOrEqual(t, n.value, n.right.value, "right child should have a greater or equal value to parent")

		rightTotal, rightBlack, rightOk = assertRedBlackPropertiesNode(t, n.right)
		total += rightTotal

		if !rightOk {
			return total, blackCount, false
		}
	}

	// Subtree sizes
	if !assertEqual(t, leftTotal, n.nLeft, "incorrect node left subtree count") {
		return total, blackCount, false
	}

	if !assertEqual(t, rightTotal, n.nRight, "incorrect node right subtree count") {
		return total, blackCount, false
	}

	// Red-black properties
	if n.safeColor() == red {
		ok = ok && assertEqual(t, black, n.left.safeColor(), "red node should have black left child")
		ok = ok && assertEqual(t, black, n.right.safeColor(), "red node should have black right child")
	}

	assertEqual(t, leftBlack, rightBlack, "should have equal number of black nodes to each leaf")
	return total, leftBlack, ok
}

func assertRedBlackProperties[T Numeric](t *testing.T, tr *FixedWindow[T]) bool {
	_, _, ok := assertRedBlackPropertiesNode(t, tr.root)
	return ok
}

func Test_fixed_Insert(t *testing.T) {
	t.Run("fully worked example", func(t *testing.T) {
		tr := Fixed[int](11)
		assertEqual(t, 0, tr.Size())

		tr.Put(1)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 1, tr.root.value, "should insert root")

		tr.Put(22)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 22, tr.root.right.value, "should insert child")

		tr.Put(27)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 22, tr.root.value, "should rotate left")
		assertEqual(t, 1, tr.root.left.value, "should rotate left")
		assertEqual(t, 27, tr.root.right.value, "should rotate left")

		tr.Put(15)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 22, tr.root.value)
		assertEqual(t, 1, tr.root.left.value)
		assertEqual(t, 15, tr.root.left.right.value)

		tr.Put(6)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 22, tr.root.value)
		assertEqual(t, 6, tr.root.left.value, "should rotate right then left")
		assertEqual(t, 1, tr.root.left.left.value, "should rotate right then left")
		assertEqual(t, 15, tr.root.left.right.value, "should rotate right then left")
		assertEqual(t, red, tr.root.left.right.color)

		tr.Put(11)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 22, tr.root.value)
		assertEqual(t, 6, tr.root.left.value)
		assertEqual(t, 15, tr.root.left.right.value)
		assertEqual(t, black, tr.root.left.right.color, "should recolor 15")
		assertEqual(t, 11, tr.root.left.right.left.value)

		tr.Put(17)
		assertRedBlackProperties(t, tr)

		tr.Put(25)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 22, tr.root.value)
		assertEqual(t, 27, tr.root.right.value)
		assertEqual(t, 25, tr.root.right.left.value)

		tr.Put(13)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 15, tr.root.value, "should rotate 15 up to root")
		assertEqual(t, 6, tr.root.left.value, "should rotate 15 up to root")
		assertEqual(t, 1, tr.root.left.left.value, "should rotate 15 up to root")
		assertEqual(t, 11, tr.root.left.right.value, "should rotate 15 up to root")
		assertEqual(t, 13, tr.root.left.right.right.value, "should rotate 15 up to root")
		assertEqual(t, 22, tr.root.right.value, "should rotate 15 up to root")
		assertEqual(t, 17, tr.root.right.left.value, "should rotate 15 up to root")
		assertEqual(t, 27, tr.root.right.right.value, "should rotate 15 up to root")

		tr.Put(8)
		assertRedBlackProperties(t, tr)

		tr.Put(1)
		assertRedBlackProperties(t, tr)
		assertEqual(t, 15, tr.root.value)
		assertEqual(t, 6, tr.root.left.value)
		assertEqual(t, 1, tr.root.left.left.value)
		assertEqual(t, 1, tr.root.left.left.right.value, "should insert duplicates to the right")

		assertEqual(t, 11, tr.Size(), "should reach its capacity")
	})

	t.Run("many duplicates", func(t *testing.T) {
		const size = 100
		tr := Fixed[int](size)
		for range size {
			tr.Put(1)
			assertRedBlackProperties(t, tr)
		}
	})

	t.Run("random", func(t *testing.T) {
		const size = 100
		tr := Fixed[int](size)
		for range size {
			v := rand.Int()
			tr.Put(v)
		}

		assertRedBlackProperties(t, tr)
	})
}

func Test_fixed_swap(t *testing.T) {
	t.Run("nil and nil", func(t *testing.T) {
		tr := Fixed[int](10)
		tr.swap(nil, nil)
		assertEqual(t, 0, tr.Size())
	})

	t.Run("root with itself", func(t *testing.T) {
		tr := Fixed[int](10)
		tr.root = &node[int]{value: 1}
		tr.swap(tr.root, tr.root)
		assertEqual(t, 1, tr.root.value)
	})

	t.Run("root with nil", func(t *testing.T) {
		tr := Fixed[int](10)
		tr.root = &node[int]{value: 1}
		tr.swap(tr.root, nil)
		assertEqual(t, 1, tr.root.value)
	})

	t.Run("root with left", func(t *testing.T) {
		tr := Fixed[int](10)
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
		assertEqual(t, n2, tr.root)
		assertEqual(t, n4, tr.root.left)
		assertEqual(t, n1, tr.root.left.left)
		assertEqual(t, n3, tr.root.left.right)
		assertEqual(t, n5, tr.root.right)
	})

	t.Run("root with right", func(t *testing.T) {
		tr := Fixed[int](10)
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
		assertEqual(t, n4, tr.root)
		assertEqual(t, n1, tr.root.left)
		assertEqual(t, n2, tr.root.right)
		assertEqual(t, n3, tr.root.right.left)
		assertEqual(t, n5, tr.root.right.right)
	})

	t.Run("root with left left", func(t *testing.T) {
		tr := Fixed[int](10)
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
		assertEqual(t, n2, tr.root)
		assertEqual(t, n4, tr.root.left)
		assertEqual(t, n7, tr.root.right)
		assertEqual(t, n6, tr.root.left.left)
		assertEqual(t, n5, tr.root.left.right)
		assertEqual(t, n1, tr.root.left.left.left)
		assertEqual(t, n3, tr.root.left.left.right)
	})

	t.Run("root with left right", func(t *testing.T) {
		tr := Fixed[int](10)
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
		assertEqual(t, n5, tr.root)
		assertEqual(t, n4, tr.root.left)
		assertEqual(t, n7, tr.root.right)
		assertEqual(t, n2, tr.root.left.left)
		assertEqual(t, n6, tr.root.left.right)
		assertEqual(t, n1, tr.root.left.left.left)
		assertEqual(t, n3, tr.root.left.left.right)
	})

	t.Run("left with grandchild", func(t *testing.T) {
		tr := Fixed[int](10)
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
		assertEqual(t, n6, tr.root)
		assertEqual(t, n1, tr.root.left)
		assertEqual(t, n7, tr.root.right)
		assertEqual(t, n2, tr.root.left.left)
		assertEqual(t, n5, tr.root.left.right)
		assertEqual(t, n4, tr.root.left.left.left)
		assertEqual(t, n3, tr.root.left.left.right)
	})

	t.Run("right with grandchild", func(t *testing.T) {
		tr := Fixed[int](10)
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
		assertEqual(t, n2, tr.root)
		assertEqual(t, n1, tr.root.left)
		assertEqual(t, n7, tr.root.right)
		assertEqual(t, n3, tr.root.right.left)
		assertEqual(t, n6, tr.root.right.right)
		assertEqual(t, n5, tr.root.right.right.left)
		assertEqual(t, n4, tr.root.right.right.right)
	})
}

func Test_fixed_delete(t *testing.T) {
	t.Run("remove leaf, no rotate", func(t *testing.T) {
		tr := makeFixed(1, 22, 27, 15, 6, 11, 17, 25, 13, 8, 1)

		p := tr.root.left.left
		assertEqual(t, 1, p.value)
		n := p.right
		assertEqual(t, 1, n.value)
		tr.delete(n)
		assertEqual(t, 10, tr.Size())
		assertRedBlackProperties(t, tr)
		assertNil(t, n.parent)
		assertNil(t, n.left)
		assertNil(t, n.right)
		assertNil(t, p.left)
		assertNil(t, p.right)
	})

	t.Run("replace parent with child; case 4", func(t *testing.T) {
		tr := makeFixed(1, 22, 27, 15, 6, 11, 17, 25, 13, 8, 1)

		p := tr.root.right
		assertEqual(t, 22, p.value)
		n := p.right
		assertEqual(t, 27, n.value)
		tr.delete(n)
		assertEqual(t, 10, tr.Size())
		assertRedBlackProperties(t, tr)
		assertNil(t, n.parent)
		assertNil(t, n.left)
		assertNil(t, n.right)
		assertEqual(t, 17, p.left.value)
		assertEqual(t, 25, p.right.value)
	})

	t.Run("remove parent with two children; cases 5 right and 6 left", func(t *testing.T) {
		tr := makeFixed(1, 22, 27, 15, 6, 11, 17, 25, 13, 8, 1)

		p := tr.root
		assertEqual(t, 15, p.value)
		n := p.right
		assertEqual(t, 22, n.value)
		tr.delete(n)
		assertEqual(t, 10, tr.Size())
		assertRedBlackProperties(t, tr)
		assertNil(t, n.parent)
		assertNil(t, n.left)
		assertNil(t, n.right)
		assertEqual(t, tr.root, p, "should keep 15 at root")
		assertEqual(t, 25, p.right.value)
		assertEqual(t, 17, p.right.left.value)
		assertEqual(t, 27, p.right.right.value)
	})

	t.Run("case 3 rotate left", func(t *testing.T) {
		tr := makeFixed(5, 8, 1, 7, 9, 6)

		p := tr.root
		assertEqual(t, 5, p.value)
		n := p.left
		assertEqual(t, 1, n.value)
		tr.delete(n)
		assertRedBlackProperties(t, tr)
		assertNil(t, n.parent)
		assertNil(t, n.left)
		assertNil(t, n.right)
		assertEqual(t, 8, tr.root.value)
		assertEqual(t, 6, tr.root.left.value)
		assertEqual(t, 9, tr.root.right.value)
		assertEqual(t, 5, tr.root.left.left.value)
		assertEqual(t, 7, tr.root.left.right.value)
	})

	t.Run("case 3 rotate right", func(t *testing.T) {
		tr := makeFixed(5, 8, 2, 1, 3, 4)

		p := tr.root
		assertEqual(t, 5, p.value)
		n := p.right
		assertEqual(t, 8, n.value)
		tr.delete(n)
		assertRedBlackProperties(t, tr)
		assertNil(t, n.parent)
		assertNil(t, n.left)
		assertNil(t, n.right)
		assertEqual(t, 2, tr.root.value)
		assertEqual(t, 1, tr.root.left.value)
		assertEqual(t, 4, tr.root.right.value)
		assertEqual(t, 3, tr.root.right.left.value)
		assertEqual(t, 5, tr.root.right.right.value)
	})

	t.Run("case 2", func(t *testing.T) {
		tr := makeFixed(5, 2, 8, 6)

		p := tr.root
		assertEqual(t, 5, p.value)
		n := p.left
		assertEqual(t, 2, n.value)

		// Delete the 6 to get the tree in the correct state
		n6 := p.right.left
		assertEqual(t, 6, n6.value)
		tr.delete(n6)

		assertEqual(t, black, p.color)
		assertEqual(t, black, n.color)
		assertEqual(t, black, p.right.color)

		// Now it will trigger delete case 2
		tr.delete(n)
		assertRedBlackProperties(t, tr)
		assertNil(t, n.parent)
		assertNil(t, n.left)
		assertNil(t, n.right)
		assertEqual(t, 5, tr.root.value)
		assertNil(t, tr.root.left)
		assertEqual(t, 8, tr.root.right.value)
	})
}

func Test_fixed_rollingWindowAtCapacity(t *testing.T) {
	t.Run("single node", func(t *testing.T) {
		tr := Fixed[int](1)
		tr.Put(1)
		assertEqual(t, 1, tr.Size())
		assertEqual(t, 1, tr.root.value)
		tr.Put(2)
		assertEqual(t, 1, tr.Size())
		assertEqual(t, 2, tr.root.value, "should replace existing value")
	})

	t.Run("three nodes", func(t *testing.T) {
		tr := makeFixed(1, 2, 3)
		assertEqual(t, 3, tr.Size())
		assertEqual(t, 2, tr.root.value)
		assertEqual(t, 1, tr.root.left.value)
		assertEqual(t, 3, tr.root.right.value)

		tr.Put(4)
		assertEqual(t, 3, tr.Size(), "should replace oldest value")
		assertEqual(t, 3, tr.root.value)
		assertEqual(t, 2, tr.root.left.value)
		assertEqual(t, 4, tr.root.right.value)
	})

	t.Run("three nodes replace root", func(t *testing.T) {
		tr := makeFixed(3, 1, 5)
		assertEqual(t, 3, tr.Size())
		assertEqual(t, 3, tr.root.value)
		assertEqual(t, 1, tr.root.left.value)
		assertEqual(t, 5, tr.root.right.value)

		tr.Put(4)
		assertEqual(t, 3, tr.Size(), "should replace oldest value at root")
		assertEqual(t, 4, tr.root.value)
		assertNil(t, tr.root.parent)
		assertEqual(t, 1, tr.root.left.value)
		assertEqual(t, 5, tr.root.right.value)
	})

	t.Run("resets subtree counts for replaced node", func(t *testing.T) {
		// Can you tell that this was a particularly subtle bug?
		// The node that will be replaced (50784) did not correctly have its subtree counts
		// reset to zero before that node was reused for the new value (37314). This threw
		// off subtree counts up to the root, which in turn breaks quantile calculation.
		tr := makeFixed(36564, 50784, 30136, 31835, 44643, 2647, 63181, 13969, 43113, 33834)
		tr.i = 1
		tr.Put(37314)
		assertRedBlackProperties(t, tr)
	})

	t.Run("rolling 50 nodes random", func(t *testing.T) {
		const size = 50
		values := make([]int, 0, size)
		tr := Fixed[int](size)
		for i := range 1000 {
			v := rand.IntN(65536)
			if i >= size {
				k := i % size
				values[k] = v
			} else {
				values = append(values, v)
			}

			tr.Put(v)

			if !assertRedBlackProperties(t, tr) {
				t.Logf("failed at i=%d", i)
				break
			}
		}
	})
}

func Test_fixed_MinMax(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		tr := Fixed[int](1)
		assertEqual(t, 0, tr.Min())
		assertEqual(t, 0, tr.Max())
	})

	t.Run("single node", func(t *testing.T) {
		tr := Fixed[int](1)
		tr.Put(5)
		assertEqual(t, 5, tr.Min())
		assertEqual(t, 5, tr.Max())
		tr.Put(6)
		assertEqual(t, 6, tr.Min())
		assertEqual(t, 6, tr.Max())
	})

	t.Run("three nodes", func(t *testing.T) {
		tr := makeFixed(2, 1, 3)
		assertEqual(t, 1, tr.Min())
		assertEqual(t, 3, tr.Max())
	})

	t.Run("rolling three nodes", func(t *testing.T) {
		tr := makeFixed(1, 2, 3)
		tr.Put(4) // replaces 1
		assertEqual(t, 2, tr.Min())
		assertEqual(t, 4, tr.Max())
		tr.Put(1) // replaces 2
		tr.Put(2) // replaces 3
		tr.Put(3) // replaces 4
		assertEqual(t, 1, tr.Min())
		assertEqual(t, 3, tr.Max())
	})

	t.Run("rolling 50 nodes random", func(t *testing.T) {
		const size = 50
		values := make([]int, 0, size)
		tr := Fixed[int](size)
		for i := 0; i < 1000; i++ {
			v := rand.Int()
			if i >= size {
				k := i % size
				values[k] = v
			} else {
				values = append(values, v)
			}

			tr.Put(v)
			expectedMin := slices.Min(values)
			expectedMax := slices.Max(values)
			if !assertEqual(t, expectedMin, tr.Min(), "min should match") ||
				!assertEqual(t, expectedMax, tr.Max(), "max should match") {
				break
			}
		}
	})
}

func Test_fixed_MeanVariance(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		tr := Fixed[int](1)
		assertEqual(t, 0.0, tr.Mean())
		assertEqual(t, 0.0, tr.Variance())
	})

	t.Run("single node", func(t *testing.T) {
		tr := Fixed[int](1)
		tr.Put(5)
		assertEqual(t, 5.0, tr.Mean())
		assertEqual(t, 0.0, tr.Variance())
		tr.Put(6)
		assertEqual(t, 6.0, tr.Mean())
		assertEqual(t, 0.0, tr.Variance())
	})

	t.Run("three nodes", func(t *testing.T) {
		tr := makeFixed(2, 1, 3)
		assertEqual(t, 2.0, tr.Mean())
		assertEqual(t, 2.0/3.0, tr.Variance())
	})

	t.Run("rolling three nodes", func(t *testing.T) {
		tr := makeFixed(1, 2, 3)
		tr.Put(4) // replaces 1
		assertEqual(t, 3.0, tr.Mean())
		assertEqual(t, 2.0/3.0, tr.Variance())
		tr.Put(5) // replaces 2
		assertEqual(t, 4.0, tr.Mean())
		assertEqual(t, 2.0/3.0, tr.Variance())
		tr.Put(0) // replaces 3
		assertEqual(t, 3.0, tr.Mean())
		assertEqual(t, 14.0/3.0, tr.Variance())
		tr.Put(10) // replaces 4
		assertEqual(t, 5.0, tr.Mean())
		assertEqual(t, 50.0/3.0, tr.Variance())
	})

	t.Run("rolling 50 nodes random", func(t *testing.T) {
		const size = 50
		values := make([]int, 0, size)
		tr := Fixed[int](size)
		for i := 0; i < 1000; i++ {
			// Limit the random numbers to avoid overflow when summing
			v := rand.IntN(65536)
			if i >= size {
				k := i % size
				values[k] = v
			} else {
				values = append(values, v)
			}

			tr.Put(v)
			var sum float64
			for _, v := range values {
				sum += float64(v)
			}

			expectedMean := sum / float64(len(values))

			var tss float64
			for _, v := range values {
				delta := float64(v) - expectedMean
				tss += delta * delta
			}

			expectedVar := tss / float64(tr.size)

			// expectedVar can be rather large, so the allowed error delta is adjusted accordingly
			if !assertInDelta(t, expectedMean, tr.Mean(), 1e-6, "mean should be within error delta") ||
				!assertInDelta(t, expectedVar, tr.Variance(), expectedVar*1e-12, "variance should be within error delta") {
				break
			}
		}
	})
}

func Test_fixed_Quantile(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		tr := Fixed[int](1)
		assertEqual(t, 0, tr.Quantile(0.1))
		assertEqual(t, 0, tr.Quantile(0.5))
		assertEqual(t, 0, tr.Quantile(0.9))
	})

	t.Run("single node", func(t *testing.T) {
		tr := Fixed[int](1)
		tr.Put(5)
		assertEqual(t, 5, tr.Quantile(0.1))
		assertEqual(t, 5, tr.Quantile(0.5))
		assertEqual(t, 5, tr.Quantile(0.9))
		tr.Put(6)
		assertEqual(t, 6, tr.Quantile(0.1))
		assertEqual(t, 6, tr.Quantile(0.5))
		assertEqual(t, 6, tr.Quantile(0.9))
	})

	t.Run("two nodes", func(t *testing.T) {
		tr := makeFixed(2, 1)
		assertEqual(t, 1, tr.Quantile(0.1))
		assertEqual(t, 1, tr.Quantile(0.5))
		assertEqual(t, 2, tr.Quantile(0.9))
	})

	t.Run("rolling three nodes", func(t *testing.T) {
		tr := makeFixed(1, 2, 3)
		tr.Put(4) // replaces 1
		assertEqual(t, 2, tr.Quantile(0.1))
		assertEqual(t, 3, tr.Quantile(0.5))
		assertEqual(t, 4, tr.Quantile(0.9))
		tr.Put(5) // replaces 2
		assertEqual(t, 3, tr.Quantile(0.1))
		assertEqual(t, 4, tr.Quantile(0.5))
		assertEqual(t, 5, tr.Quantile(0.9))
		tr.Put(0) // replaces 3
		assertEqual(t, 0, tr.Quantile(0.1))
		assertEqual(t, 4, tr.Quantile(0.5))
		assertEqual(t, 5, tr.Quantile(0.9))
		tr.Put(10) // replaces 4
		assertEqual(t, 0, tr.Quantile(0.1))
		assertEqual(t, 5, tr.Quantile(0.5))
		assertEqual(t, 10, tr.Quantile(0.9))
	})

	t.Run("many duplicates", func(t *testing.T) {
		tr := makeFixed(2, 2, 2, 2, 2, 2, 2, 2, 2)
		assertEqual(t, 2, tr.Quantile(0.1))
		assertEqual(t, 2, tr.Quantile(0.5))
		assertEqual(t, 2, tr.Quantile(0.9))
		tr.Put(1)
		assertEqual(t, 1, tr.Quantile(0.1))
		assertEqual(t, 2, tr.Quantile(0.5))
		assertEqual(t, 2, tr.Quantile(0.9))
		tr.Put(3)
		assertEqual(t, 1, tr.Quantile(0.1))
		assertEqual(t, 2, tr.Quantile(0.5))
		assertEqual(t, 3, tr.Quantile(0.9))
	})

	t.Run("rolling 50 nodes random", func(t *testing.T) {
		const size = 50
		values := make([]int, 0, size)
		tr := Fixed[int](size)
		for i := 0; i < 1000; i++ {
			v := rand.IntN(65536)
			if i >= size {
				k := i % size
				values[k] = v
			} else {
				values = append(values, v)
			}

			tr.Put(v)

			// Make a copy of values so that slowQuantile's sort doesn't break the order
			// in which we replace old values.
			valuesCopy := slices.Clone(values)

			wantFirstDecile := slowQuantile(valuesCopy, 0.1)
			gotFirstDecile := tr.Quantile(0.1)
			ok := assertEqual(t, wantFirstDecile, gotFirstDecile, "unexpected first decile")

			wantMedian := slowQuantile(valuesCopy, 0.5)
			gotMedian := tr.Quantile(0.5)
			ok = ok && assertEqual(t, wantMedian, gotMedian, "unexpected median")

			wantNinthDecile := slowQuantile(valuesCopy, 0.9)
			gotNinthDecile := tr.Quantile(0.9)
			ok = ok && assertEqual(t, wantNinthDecile, gotNinthDecile, "unexpected ninth decile")

			if !ok {
				t.Logf("failed at i=%d", i)
				break
			}
		}
	})
}

// slowQuantile is a test helper function for computing an exact quantile in
// the most exact way possible. If two values are on the boundary of the quantile,
// it chooses the larger one (instead of averaging).
func slowQuantile(values []int, q float64) int {
	if len(values) == 0 {
		return 0
	}

	slices.Sort(values)
	fLen := float64(len(values))
	for i := 0; i < len(values); i++ {
		if float64(i+1)/fLen >= q {
			return values[i]
		}
	}

	return values[len(values)-1]
}

func Test_slowQuantile(t *testing.T) {
	v := []int{3, 6, 7, 8, 8, 10, 13, 15, 16, 20}
	assertEqual(t, 3, slowQuantile(v, 0.0))
	assertEqual(t, 7, slowQuantile(v, 0.25))
	assertEqual(t, 8, slowQuantile(v, 0.5))
	assertEqual(t, 15, slowQuantile(v, 0.75))
	assertEqual(t, 20, slowQuantile(v, 1.0))

	v = append(v, 21)
	assertEqual(t, 3, slowQuantile(v, 0.0))
	assertEqual(t, 7, slowQuantile(v, 0.25))
	assertEqual(t, 10, slowQuantile(v, 0.5))
	assertEqual(t, 16, slowQuantile(v, 0.75))
	assertEqual(t, 21, slowQuantile(v, 1.0))
}

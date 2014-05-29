package container

import (
	"github.com/clipperhouse/typewriter"
)

var containerTemplates = typewriter.TemplateSet{
	"SortedSet": &typewriter.Template{
		Text: `
		
// The primary type that represents a sorted set
// backed by a skiplist
type {{.Name}}SortedSet struct {
	less      func(a, b {{.Pointer}}{{.Name}}) bool
	head      []*sortedSet{{.Name}}Element
	length    int
	maxLevels int
	r         *rand.Rand
}

// the struct to hold elements of the skiplist
type sortedSet{{.Name}}Element struct {
	val  {{.Pointer}}{{.Name}}
	next []*sortedSet{{.Name}}Element
}

// Creates and returns a reference to an empty set.
func New{{.Name}}SortedSet(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) {{.Name}}SortedSet {
	return {{.Name}}SortedSet{
		less:      less,
		maxLevels: 64,
		head:      make([]*sortedSet{{.Name}}Element, 64),
		r:         rand.New(rand.NewSource(123123)),
	}
}

func newSortedSet{{.Name}}Element(v {{.Pointer}}{{.Name}}, levels int) *sortedSet{{.Name}}Element {
	return &sortedSet{{.Name}}Element{v, make([]*sortedSet{{.Name}}Element, levels)}
}

// Creates and returns a reference to a set from an existing slice
func New{{.Name}}SortedSetFromSlice(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, s []{{.Pointer}}{{.Name}}) {{.Name}}SortedSet {
	a := New{{.Name}}SortedSet(less)
	for _, item := range s {
		a.Add(item)
	}
	return a
}

func (ss {{.Name}}SortedSet) randomLevels() int {
	level := int(math.Log(1.0-ss.r.Float64()) / math.Log(0.5))
	if level >= ss.maxLevels {
		level = ss.maxLevels
	}
	if level == 0 {
		level++
	}
	return level
}

// Adds an item to the current set if it doesn't already exist in the set.
func (ss {{.Name}}SortedSet) Add(v {{.Pointer}}{{.Name}}) bool {
	var backPointer = make([]*sortedSet{{.Name}}Element, 64)
	// zeroing this causes the compiler to not allocate memory each time
	// for a 20-30% boost in speed
	for i := 0; i < 64; i++ {
		backPointer[i] = nil
	}
	for level := ss.maxLevels - 1; level >= 0; level-- {
		var e *sortedSet{{.Name}}Element = nil
		if level+1 == ss.maxLevels || backPointer[level+1] == nil {
			e = ss.head[level]
		} else {
			e = backPointer[level+1]
		}
		for e != nil {
			// if they are equal, overwrite?
			if ss.less(v, e.val) == ss.less(e.val, v) {
				return false
			}
			// if inspected val is greater than k, go back and down a level
			if ss.less(v, e.val) {
				break
			}
			backPointer[level] = e
			e = e.next[level]
		}
	}
	// create new element
	e := newSortedSet{{.Name}}Element(v, ss.randomLevels())

	// connect new element up with backPointer
	for level := 0; level < len(e.next); level++ {
		if backPointer[level] == nil {
			e.next[level] = ss.head[level]
			ss.head[level] = e
		} else {
			e.next[level] = backPointer[level].next[level]
			backPointer[level].next[level] = e
		}
	}

	ss.length++
	return true
}

// Determines if a given item is already in the set.
func (ss {{.Name}}SortedSet) Contains(v {{.Pointer}}{{.Name}}) bool {
	var backPointer = make([]*sortedSet{{.Name}}Element, 64)
	// zeroing this causes the compiler to not allocate memory each time
	// for a 20-30% boost in speed
	for i := 0; i < 64; i++ {
		backPointer[i] = nil
	}
	for level := ss.maxLevels - 1; level >= 0; level-- {
		var e *sortedSet{{.Name}}Element = nil
		if level+1 == ss.maxLevels || backPointer[level+1] == nil {
			e = ss.head[level]
		} else {
			e = backPointer[level+1]
		}
		for e != nil {
			// if they are equal, return val
			if ss.less(v, e.val) == ss.less(e.val, v) {
				return true
			}
			// if inspected val is greater than v, go back and down a level
			if ss.less(v, e.val) {
				break
			}
			backPointer[level] = e
			e = e.next[level]
		}
	}
	return false
}

// Determines if the given items are all in the set
func (ss {{.Name}}SortedSet) ContainsAll(i ...{{.Pointer}}{{.Name}}) bool {
	for _, elem := range i {
		if !ss.Contains(elem) {
			return false
		}
	}
	return true
}

// Determines if every item in the other set is in this set.
func (ss {{.Name}}SortedSet) IsSubset(other {{.Name}}SortedSet) bool {
	e := ss.head[0]
	for e != nil {
		if !other.Contains(e.val) {
			return false
		}
		e = e.next[0]
	}
	return true
}

// Determines if every item of this set is in the other set.
func (ss {{.Name}}SortedSet) IsSuperset(other {{.Name}}SortedSet) bool {
	return other.IsSubset(ss)
}

// Returns a new set with all items in both sets.
func (ss {{.Name}}SortedSet) Union(other {{.Name}}SortedSet) {{.Name}}SortedSet {
	unionedSet := New{{.Name}}SortedSet(ss.less)

	e := ss.head[0]
	for e != nil {
		unionedSet.Add(e.val)
		e = e.next[0]
	}
	e = other.head[0]
	for e != nil {
		unionedSet.Add(e.val)
		e = e.next[0]
	}
	return unionedSet
}

// Returns a new set with items that exist only in both sets.
func (ss {{.Name}}SortedSet) Intersect(other {{.Name}}SortedSet) {{.Name}}SortedSet {
	intersection := New{{.Name}}SortedSet(ss.less)
	// loop over smaller set
	if ss.Cardinality() < other.Cardinality() {
		e := ss.head[0]
		for e != nil {
			if other.Contains(e.val) {
				intersection.Add(e.val)
			}
			e = e.next[0]
		}
	} else {
		e := other.head[0]
		for e != nil {
			if ss.Contains(e.val) {
				intersection.Add(e.val)
			}
			e = e.next[0]
		}
	}
	return intersection
}

// Returns a new set with items in the current set but not in the other set
func (ss {{.Name}}SortedSet) Difference(other {{.Name}}SortedSet) {{.Name}}SortedSet {
	differencedSet := New{{.Name}}SortedSet(ss.less)
	e := ss.head[0]
	for e != nil {
		if !other.Contains(e.val) {
			differencedSet.Add(e.val)
		}
		e = e.next[0]
	}
	return differencedSet
}

// Returns a new set with items in the current set or the other set but not in both.
func (ss {{.Name}}SortedSet) SymmetricDifference(other {{.Name}}SortedSet) {{.Name}}SortedSet {
	aDiff := ss.Difference(other)
	bDiff := other.Difference(ss)
	return aDiff.Union(bDiff)
}

// Clears the entire set to be the empty set.
func (ss *{{.Name}}SortedSet) Clear() {
	*ss = {{.Name}}SortedSet{
		less:      ss.less,
		maxLevels: 64,
		head:      make([]*sortedSet{{.Name}}Element, 64),
		r:         rand.New(rand.NewSource(123123)),
	}
}

// Allows the removal of a single item in the set.
func (ss {{.Name}}SortedSet) Remove(v {{.Pointer}}{{.Name}}) {
	var backPointer = make([]*sortedSet{{.Name}}Element, 64)
	// zeroing this causes the compiler to not allocate memory each time
	// for a 20-30% boost in speed
	for i := 0; i < 64; i++ {
		backPointer[i] = nil
	}
	for level := ss.maxLevels - 1; level >= 0; level-- {
		var e *sortedSet{{.Name}}Element = nil
		if level+1 == ss.maxLevels || backPointer[level+1] == nil {
			e = ss.head[level]
		} else {
			e = backPointer[level+1]
		}
		for e != nil {
			// if they are equal, remove
			if level == 0 && ss.less(v, e.val) == ss.less(e.val, v) {
				for level := 0; level < len(e.next); level++ {
					if backPointer[level] == nil {
						ss.head[level] = e.next[level]
					} else {
						backPointer[level].next[level] = e.next[level]
					}
				}

				ss.length--
			}
			if ss.less(v, e.val) == ss.less(e.val, v) {
				break
			}
			// if inspected val is greater than k, go back and down a level
			if ss.less(v, e.val) {
				break
			}
			backPointer[level] = e
			e = e.next[level]
		}
	}
}

// Cardinality returns how many items are currently in the set.
func (ss {{.Name}}SortedSet) Cardinality() int {
	e := ss.head[0]
	ret := 0
	for e != nil {
		ret++
		e = e.next[0]
	}
	return ret
}

// Iter() returns a channel of type {{.Pointer}}{{.Name}} that you can range over.
func (ss {{.Name}}SortedSet) Iter() <-chan {{.Pointer}}{{.Name}} {
	ch := make(chan {{.Pointer}}{{.Name}})
	go func() {
		e := ss.head[0]
		for e != nil {
			ch <- e.val
			e = e.next[0]
		}
		close(ch)
	}()

	return ch
}

// Equal determines if two sets are equal to each other.
// If they both are the same size and have the same items they are considered equal.
// Order of items is not relevent for sets to be equal.
func (ss {{.Name}}SortedSet) Equal(other {{.Name}}SortedSet) bool {
	if ss.Cardinality() != other.Cardinality() {
		return false
	}
	e := ss.head[0]
	for e != nil {
		if !other.Contains(e.val) {
			return false
		}
		e = e.next[0]
	}
	return true
}

// Returns a clone of the set.
// Does NOT clone the underlying elements.
func (ss {{.Name}}SortedSet) Clone() {{.Name}}SortedSet {
	clonedSet := New{{.Name}}SortedSet(ss.less)
	e := ss.head[0]
	for e != nil {
		clonedSet.Add(e.val)
		e = e.next[0]
	}
	return clonedSet
}
`,
		RequiresComparable: true,
	},
}

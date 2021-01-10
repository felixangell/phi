package command_handler

import (
	"fmt"
	"sort"
)

type intHashSet map[int]bool

func (i *intHashSet) Contains(x int) bool {
	val, ok := (*i)[x]
	return ok && val
}
func (i *intHashSet) Delete(x int) {
	if i.Contains(x) {
		(*i)[x] = false
	}
}
func (i *intHashSet) Store(x int) {
	(*i)[x] = true
}

type intHashSetHashKey string

func (i *intHashSet) Hash() intHashSetHashKey {
	vals := make([]int, len(*i))

	idx := 0
	for key, _ := range *i {
		vals[idx] = key
		idx++
	}

	sort.Ints(vals)

	var res string
	for _, v := range vals {
		res += fmt.Sprintf("%d", v)
	}
	return intHashSetHashKey(res)
}

func newIntHashSet(vals ...int) intHashSet {
	hs := intHashSet{}
	for _, val := range vals {
		hs.Store(val)
	}
	return hs
}

package ds

type Rope struct {
	width, height uint
	left, right *Rope
	value []byte
}

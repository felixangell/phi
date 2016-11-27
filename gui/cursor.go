package gui

type Cursor struct {
	x, y   int
	rx, ry int
}

func (c *Cursor) move(x, y int) {
	c.moveRender(x, y, x, y)
}

// moves the cursors position, and the
// rendered coordinates by the given amount
func (c *Cursor) moveRender(x, y, rx, ry int) {
	c.x += x
	c.y += y

	c.rx += rx
	c.ry += ry
}

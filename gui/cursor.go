package gui

type Cursor struct {
	x, y   int
	rx, ry int
	dx, dy int
	moving bool
}

func (c *Cursor) gotoStart() {
	for c.x > 1 {
		c.move(-1, 0)
	}
}

func (c *Cursor) move(x, y int) {
	c.moveRender(x, y, x, y)
}

// moves the cursors position, and the
// rendered coordinates by the given amount
func (c *Cursor) moveRender(x, y, rx, ry int) {
	if x > 0 {
		c.dx = 1
	} else if x < 0 {
		c.dx = -1
	}

	if y > 0 {
		c.dy = 1
	} else if y < 0 {
		c.dy = -1
	}

	c.moving = true

	c.x += x
	c.y += y

	c.rx += rx
	c.ry += ry
}

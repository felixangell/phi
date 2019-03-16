package buff

import (
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
)

type Cursor struct {
	parent        *Buffer
	x, y          int
	rx, ry        int
	dx, dy        int
	moving        bool
	width, height int
}

func newCursor(parent *Buffer) *Cursor {
	return &Cursor{
		parent,
		0, 0,
		0, 0,
		0, 0,
		false,
		0, 0,
	}
}

func (c *Cursor) SetSize(w, h int) {
	c.width = w
	c.height = h
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

func (c *Cursor) Render(ctx *strife.Renderer, xOff, yOff int) {
	b := c.parent

	xPos := b.ex + (xOff + c.rx*lastCharW) - (b.cam.x * lastCharW)
	yPos := b.ey + (yOff + c.ry*c.height) - (b.cam.y * c.height)

	// NOTE: we dont have to scale the curor here because
	// it's based off the font size which has already been scaled.
	ctx.SetColor(strife.HexRGB(b.buffOpts.cursor))
	ctx.Rect(xPos, yPos, c.width, c.height, strife.Fill)

	if cfg.DebugMode {
		ctx.SetColor(strife.HexRGB(0xff00ff))
		ctx.Rect(xPos, yPos, c.width, c.height, strife.Line)
	}
}

package gui

func CloseBuffer(b *Buffer) bool {
	// no parent. this can happen sometimes
	// for example if we're trying to close
	// a buffer illegal, e.g. a palette has
	// a buffer
	if b.parent == nil {
		return false
	}

	b.parent.DeleteComponent(b.index)

	// we need to re-calculate the sizes of everything
	// since we've closed a buffer!

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	var bufferWidth int
	numComponents := b.parent.NumComponents()
	if numComponents > 0 {
		bufferWidth = b.parent.w / numComponents
	} else {
		bufferWidth = b.parent.w
	}

	// translate all the components accordingly.
	i := 0
	for _, p := range b.parent.components {
		if p == nil {
			continue
		}

		p.Resize(bufferWidth, b.parent.h)
		p.SetPosition(bufferWidth*i, 0)

		i = i + 1
	}

	return true
}

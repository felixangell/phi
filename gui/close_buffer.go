package gui

func CloseBuffer(b *Buffer) bool {
	// FIXME

	// no parent. this can happen sometimes
	// for example if we're trying to close
	// a buffer illegal, e.g. a palette has
	// a buffer
	if b.parent == nil {
		return false
	}

	// remove focus
	b.HasFocus = false
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

	// TODO track buffer visit history on a stack
	// when we remove a buffer we pop it from the stack
	// then we can focus on the top of the stack instead.
	// for now let's just focus on the most recently
	// added buffer

	// translate all the components accordingly.
	i := 0

	var lastBuffer *Buffer
	for _, p := range b.parent.components {
		if p == nil {
			continue
		}

		p.Resize(bufferWidth, b.parent.h)
		p.SetPosition(bufferWidth*i, 0)

		i = i + 1
		lastBuffer = p.(*Buffer)
	}

	if lastBuffer != nil {
		lastBuffer.HasFocus = true
	}

	return true
}

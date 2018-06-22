package piecetable

type Line struct {
	Buffer string
	parent *PieceTable
	mods   map[int]bool
	keys   []int
}

func NewLine(data string, parent *PieceTable) *Line {
	return &Line{
		data, parent, map[int]bool{}, []int{},
	}
}

func (l *Line) AppendNode(node *PieceNode) {
	nodeIndex := len(l.parent.nodes)
	l.mods[nodeIndex] = true
	l.keys = append(l.keys, nodeIndex)
	l.parent.nodes = append(l.parent.nodes, node)
}

func (l *Line) Len() int {
	return len(l.String())
}

func (l *Line) String() string {
	data := l.Buffer

	for _, keyName := range l.keys {
		thing, ok := l.mods[keyName]
		// ?
		if !ok || !thing {
			continue
		}

		mod := l.parent.nodes[keyName]

		if mod.Length >= 0 {

			// append!
			if mod.Start >= len(data) {
				data += mod.Data
				continue
			}

			fst, end := data[:mod.Start], data[mod.Start:]
			data = fst + mod.Data + end
		} else {
			data = data[:mod.Start-1] + data[mod.Start:]
		}
	}

	return data
}

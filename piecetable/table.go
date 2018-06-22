package piecetable

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type PieceTable struct {
	Lines    []*Line
	nodes    []*PieceNode
	redoList []*PieceNode
}

func MakePieceTable(data string) *PieceTable {
	readStrings := strings.Split(data, "\n")

	lines := make([]*Line, len(readStrings))
	table := &PieceTable{
		lines,
		[]*PieceNode{},
		[]*PieceNode{},
	}

	for idx, data := range readStrings {
		lines[idx] = NewLine(data, table)
	}

	return table
}

func (p *PieceTable) Redo() {
	if len(p.redoList) == 0 {
		return
	}

	action := p.redoList[len(p.redoList)-1]
	p.redoList = p.redoList[:len(p.redoList)-1]

	actionIndex := len(p.nodes)
	p.nodes = append(p.nodes, action)

	line := p.Lines[action.Index]
	line.mods[actionIndex] = true
}

func (p *PieceTable) Undo() {
	if len(p.nodes) == 0 {
		return
	}

	nodeIndex := len(p.nodes) - 1

	// get the value we pop
	change := p.nodes[nodeIndex]

	// remove the node index from
	// the mods (i.e. a dangling
	// pointer)
	line := p.Lines[change.Index]
	delete(line.mods, nodeIndex)

	// pop the most recent change
	p.nodes = p.nodes[:nodeIndex]

	// append it so we can redo it later if necessary
	p.redoList = append(p.redoList, change)
}

func (p *PieceTable) Delete(line int, idx int) {
	node := NewPiece("", line, idx)
	node.Length = -1
	p.Lines[line].AppendNode(node)
}

// TODO this builds the line and indexes it.
func (p *PieceTable) Index(line int, idx int) rune {
	r, _ := utf8.DecodeLastRuneInString(p.Lines[line].String()[idx:])
	return r
}

func (p *PieceTable) Insert(val string, line int, idx int) {
	node := NewPiece(val, line, idx)
	p.Lines[line].AppendNode(node)
}

func (p *PieceTable) Line(idx int) string {
	return p.Lines[idx].String()
}

func (p *PieceTable) String() string {
	var result string
	for idx, line := range p.Lines {
		if idx > 0 {
			result += string('\n')
		}
		result += fmt.Sprintf(line.String())
	}
	return result
}

func (p *PieceTable) Print() {
	for _, line := range p.Lines {
		fmt.Println(line.Buffer)
	}
}

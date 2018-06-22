package piecetable

type PieceNode struct {
	Index  int
	Start  int
	Length int
	Data   string
}

func NewPiece(data string, line int, start int) *PieceNode {
	return &PieceNode{
		line,
		start,
		len(data),
		data,
	}
}

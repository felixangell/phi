package gui

import (
	"github.com/felixangell/strife"
)

type debugPane struct {
	fpsGraph *lineGraph
}

type lineGraph struct {
	yAxisLabel string
	xAxisLabel string

	xValues []float64
	yValues []float64
}

func newLineGraph(xAxis, yAxis string) *lineGraph {
	return &lineGraph{
		xAxis,
		yAxis,
		[]float64{},
		[]float64{1, 2, 3, 4, 5, 6, 7},
	}
}

func (l *lineGraph) plot(x, y float64) {
	l.xValues = append(l.xValues, x)
	l.yValues = append(l.yValues, y)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (l *lineGraph) render(ctx *strife.Renderer, x, y int) {
	graphWidth, graphHeight := 256, 256

	ctx.SetColor(strife.RGBA(0, 0, 0, 255))
	ctx.Rect(x, y, graphWidth, graphHeight, strife.Fill)

	// render last ten values
	// xSnip := l.xValues[len(l.xValues)-10:]
	ySnip := l.yValues[max(len(l.yValues)-256, 0):]

	size := graphHeight / len(ySnip)

	for idx, yItem := range ySnip {
		ctx.SetColor(strife.RGB(255, 0, 255))
		ctx.Rect(x+(idx*size), y+(graphHeight-(int(yItem)*size)), 5, 5, strife.Fill)
	}
}

var pane = &debugPane{
	newLineGraph("time", "framerate"),
}

func renderDebugPane(ctx *strife.Renderer, x, y int) {
	ctx.SetColor(strife.HexRGB(0xff00ff))
	{
		pane.fpsGraph.render(ctx, x, y)
	}
}

package piecetable

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteSingle(t *testing.T) {
	text := `this is my testing
document i want to see how it fares
and all of that fun
stuff`
	output := ` my testing
document i want to see how it fares
and all of that fun
stuff`

	table := MakePieceTable(text)

	for i := 0; i < len("this is "); i++ {
		table.Delete(0, 1)
	}
	table.Undo()

	fmt.Println(table.String())

	assert.Equal(t, output, table.String())
}

func TestBadRedo(t *testing.T) {
	text := `this is my testing
document i want to see how it fares
and all of that fun
stuff`
	output := `this is my piece foo table testing
document i want to see how it fares
and all of that fun
stuff`

	table := MakePieceTable(text)
	table.Insert("piece table ", 0, 11)
	table.Undo()
	table.Redo()

	// no redo history.
	table.Redo()

	fmt.Println(table.String())

	assert.Equal(t, output, table.String())
}

func TestMultiStringInsertion(t *testing.T) {
	text := `this is my testing
document i want to see how it fares
and all of that fun
stuff`
	output := `this is my piece foo table testing
document i want to see how it fares
and all of that fun
stuff`

	table := MakePieceTable(text)
	table.Insert("piece table ", 0, 11)
	table.Insert("foo ", 0, 17)

	fmt.Println(table.String())

	assert.Equal(t, output, table.String())
}

func TestUndo(t *testing.T) {
	text := `this is my testing
document i want to see how it fares
and all of that fun
stuff`
	output := `this is my piece table testing
document i want to see how it fares
and all of that fun
stuff`

	table := MakePieceTable(text)
	table.Insert("piece table ", 0, 11)

	fmt.Println(table.String())
	assert.Equal(t, output, table.String())

	table.Undo()
	fmt.Println(table.String())
	assert.Equal(t, text, table.String())
}

func TestRedo(t *testing.T) {
	text := `this is my testing
document i want to see how it fares
and all of that fun
stuff`
	output := `this is my piece table testing
document i want to see how it fares
and all of that fun
stuff`

	table := MakePieceTable(text)
	table.Insert("piece table ", 0, 11)

	fmt.Println(table.String())
	assert.Equal(t, output, table.String())

	table.Undo()
	fmt.Println(table.String())
	assert.Equal(t, text, table.String())

	table.Redo()
	fmt.Println(table.String())
	assert.Equal(t, output, table.String())
}

func TestStringInsertion(t *testing.T) {
	text := `this is my testing
document i want to see how it fares
and all of that fun
stuff`
	output := `this is my piece table testing
document i want to see how it fares
and all of that fun
stuff`

	table := MakePieceTable(text)
	table.Insert("piece table ", 0, 11)

	fmt.Println(table.String())

	assert.Equal(t, output, table.String())
}

func TestPrintDocument(t *testing.T) {
	text := `this is my testing
document i want to see how it fares
and all of that fun
stuff`

	table := MakePieceTable(text)

	fmt.Println(table.String())

	assert.Equal(t, text, table.String(), "Un-modified piece table output doesn't match value expected")
}

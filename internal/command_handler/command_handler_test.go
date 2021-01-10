package command_handler

import (
	"github.com/felixangell/phi/internal/cfg"
	"github.com/felixangell/strife"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dummyConfig = &cfg.PhiEditorConfig{
	Commands: map[string]cfg.Command{
		"delete_line":  {Shortcut: "super+d"},
		"close_buffer": {Shortcut: "super+w"},
	},
}

func TestBuildsHashComboCorrectly(t *testing.T) {
	hash := parseHashCombo("super+d")

	hs := newIntHashSet(
		mapShortcutKeyword("super"), strife.KEY_D)
	expected := hs.Hash()

	assert.Equal(t, expected, hash)
}

func TestSimpleShortcutMapsCorrectly(t *testing.T) {
	ch := newCommandHandler(dummyConfig)
	victim, ok := ch.DeduceCommandName(true, false, strife.KEY_D)
	assert.True(t, ok)
	assert.Equal(t, "delete_line", victim)
}

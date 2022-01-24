package glob_test

import (
	"github.com/matdurand/go-import-checks/glob"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleWildcardShouldMatch(t *testing.T) {
	g, err := glob.NewGlob("a/*/c")
	assert.Nil(t, err)
	assert.True(t, g.Match("a/b/c"))
}

func TestSingleWildcardWithSameElementMultipleTimeShouldNotMatch(t *testing.T) {
	g, err := glob.NewGlob("a/*/c")
	assert.Nil(t, err)
	assert.False(t, g.Match("a/b/b/c"))
}

func TestPathDeeperThanGlobShouldMatch(t *testing.T) {
	g, err := glob.NewGlob("a/*/c")
	assert.Nil(t, err)
	assert.False(t, g.Match("a/b/c/c"))
}

func TestDoubleStarShouldMatchMultipleElements(t *testing.T) {
	g, err := glob.NewGlob("a/**/c")
	assert.Nil(t, err)
	assert.True(t, g.Match("a/b/b/b/c"))
}

func TestNegationShouldMatchWhenDifferent(t *testing.T) {
	g, err := glob.NewGlob("a/!z/c")
	assert.Nil(t, err)
	assert.True(t, g.Match("a/b/c"))
}

func TestNegationShouldntMatchWhenTheSame(t *testing.T) {
	g, err := glob.NewGlob("a/!z/c")
	assert.Nil(t, err)
	assert.False(t, g.Match("a/z/c"))
}
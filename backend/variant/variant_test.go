package variant

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateElo(t *testing.T)  {
	t.Run("w:500,l:100", func(t *testing.T) {
		newWin, newLose := calculateElo(500, 100)
		assert.Equal(t, 502, newWin)
		assert.Equal(t, 98, newLose)
	})
	t.Run("w:1000,l:1000", func(t *testing.T) {
		newWin, newLose := calculateElo(1000, 1000)
		assert.Equal(t, 1016, newWin)
		assert.Equal(t, 984, newLose)
	})
	t.Run("w:600,l:400", func(t *testing.T) {
		newWin, newLose := calculateElo(600, 400)
		assert.Equal(t, 607, newWin)
		assert.Equal(t, 393, newLose)
	})
	t.Run("w:853,l:1147", func(t *testing.T) {
		newWin, newLose := calculateElo(853 , 1147)
		assert.Equal(t, 880, newWin)
		assert.Equal(t, 1120, newLose)
	})
}

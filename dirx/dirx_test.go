package dirx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNestedDir(t *testing.T) {
	t.Run("create nested dir", func(t *testing.T) {
		err := CreateNestedDir("/tmp/test/dirx", 0755)
		assert.NoError(t, err)
	})
}

func TestCreateNestedDirFromFilepath(t *testing.T) {
	t.Run("create nested dir from filepath", func(t *testing.T) {
		err := CreateNestedDirFromFilepath("/tmp/test/dirx/file.txt", 0755)
		assert.NoError(t, err)
	})
}

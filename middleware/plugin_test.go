package middleware

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	s, err := Template("${tempdir}/logs/")
	assert.Equal(t, nil, err)
	assert.Equal(t, s, os.TempDir()+"/logs/")
}

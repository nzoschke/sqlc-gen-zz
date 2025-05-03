package gen_test

import (
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
)

func TestCamel(t *testing.T) {
	s := strcase.ToCamel("id")
	assert.Equal(t, "ID", s)
}

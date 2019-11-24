package postgres

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	someColumn := "name"
	someWord := "goddd"

	expected := fmt.Sprintf("WHERE (%s ILIKE '%s') LIMIT 5 OFFSET 10", someColumn, "%"+someWord+"%")

	underTest := NewParser()

	actual, err := underTest.Parse(fmt.Sprintf("name=match=*%s*", someWord), 5, 10)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

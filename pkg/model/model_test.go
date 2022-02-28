package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompanyGetContent(t *testing.T) {
	c := Company{
		Content: `{"location": "test"}`,
	}

	want := map[string]interface{}{
		"location": "test",
	}
	got, err := c.GetContent()

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, got["location"], "Location should not be nil")
	assert.Equal(t, want, got, "Result does not match expected.")
}

func TestCompanySetContent(t *testing.T) {
	c := Company{}
	content := map[string]interface{}{
		"location": "test",
	}

	want := Company{
		Content: `{"location":"test"}`,
	}

	err := c.SetContent(content)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, want, c)
}

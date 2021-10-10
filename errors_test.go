package venomoid

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorWrapper_Error(t *testing.T) {
	e := ErrorWrapper{
		InternalError: errors.New("internal error"),
		Label:         "error found",
	}
	assert.Equal(t, "error found, error: internal error", e.Error(), "unexpected formatted string")
}

package nsenter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	deferd, err := Set()
	require.NoError(t, err)
	defer func() { _ = deferd() }()


}

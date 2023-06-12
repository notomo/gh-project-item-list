package projectitem_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/notomo/gh-project-item-list/projectitem"
	"github.com/notomo/gh-project-item-list/projectitem/gqltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	gql, err := gqltest.New(
		t,

		// TODO: fill nodes
		gqltest.WithQueryOK("GetProjectItems", `
{
  "data": {
    "user": {
      "projectV2": {
        "items": {
          "nodes": []
        }
      }
    }
  }
}
`),
	)
	require.NoError(t, err)

	output := &bytes.Buffer{}
	ctx := context.Background()
	assert.NoError(t, projectitem.List(
		ctx,
		gql,
		"https://github.com/users/notomo/projects/1",
		".[]",
		10,
		output,
	))

	want := `[]
`
	got := output.String()
	assert.Equal(t, want, got)
}

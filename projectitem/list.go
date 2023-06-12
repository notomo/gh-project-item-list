package projectitem

import (
	"context"
	"encoding/json"
	"io"

	"github.com/cli/go-gh/pkg/api"
)

func List(
	ctx context.Context,
	gql api.GQLClient,
	projectUrl string,
	jqFilter string,
	limit int,
	writer io.Writer,
) error {
	projectDescriptor, err := GetProjectDescriptor(projectUrl)
	if err != nil {
		return err
	}

	filter, err := CreateFilter(jqFilter)
	if err != nil {
		return err
	}

	items, err := GetProjectItems(ctx, gql, *projectDescriptor, filter, limit)
	if err != nil {
		return err
	}

	encorder := json.NewEncoder(writer)
	encorder.SetIndent("", "  ")
	if err := encorder.Encode(items); err != nil {
		return err
	}

	return nil
}

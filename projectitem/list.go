package projectitem

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cli/go-gh/pkg/api"
)

func List(
	ctx context.Context,
	gql api.GQLClient,
	projectUrl string,
	jqFilter string,
	limit int,
	pageLimit int,
	writer io.Writer,
) error {
	projectDescriptor, err := GetProjectDescriptor(projectUrl)
	if err != nil {
		return fmt.Errorf("get project descriptor: %w", err)
	}

	filter, err := CreateFilter(jqFilter)
	if err != nil {
		return fmt.Errorf("create filter: %w", err)
	}

	items, err := GetProjectItems(ctx, gql, *projectDescriptor, filter, limit, pageLimit)
	if err != nil {
		return fmt.Errorf("get project items: %w", err)
	}

	encorder := json.NewEncoder(writer)
	encorder.SetIndent("", "  ")
	if err := encorder.Encode(items); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	return nil
}

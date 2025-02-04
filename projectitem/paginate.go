package projectitem

import (
	"context"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type PageInfo struct {
	EndCursor   string
	HasNextPage bool
}

type PageReceiver interface {
	Page() ([]any, PageInfo)
}

const limitPerRequest = 100

func Paginate(
	ctx context.Context,
	gql *api.GraphQLClient,
	query string,
	pageReceiver PageReceiver,
	variables map[string]any,
	each func(context.Context, []any) (bool, error),
) error {
	var cursor *graphql.String
	for {
		vars := map[string]any{
			"limit": graphql.Int(limitPerRequest),
			"after": cursor,
		}
		for k, v := range variables {
			vars[k] = v
		}

		if err := gql.DoWithContext(ctx, query, vars, pageReceiver); err != nil {
			return err
		}

		page, pageInfo := pageReceiver.Page()
		continued, err := each(ctx, page)
		if err != nil {
			return err
		}
		if !pageInfo.HasNextPage || !continued {
			break
		}

		cursor = graphql.NewString(graphql.String(pageInfo.EndCursor))
	}
	return nil
}

func Collect(
	ctx context.Context,
	gql *api.GraphQLClient,
	query string,
	pageReceiver PageReceiver,
	variables map[string]any,
	filter Filter,
	limit int,
	pageLimit int,
) ([]any, error) {
	pageCountDown := pageLimit
	collected := []any{}
	if err := Paginate(
		ctx,
		gql,
		query,
		pageReceiver,
		variables,
		func(ctx context.Context, page []any) (bool, error) {
			filtered, err := filter(ctx, page)
			if err != nil {
				return false, err
			}
			collected = append(collected, filtered...)

			if done := len(collected) >= limit; done {
				collected = collected[:limit]
				return false, nil
			}

			pageCountDown--

			return pageCountDown > 0, nil
		},
	); err != nil {
		return nil, err
	}
	return collected, nil
}

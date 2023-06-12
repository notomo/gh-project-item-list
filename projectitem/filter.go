package projectitem

import (
	"context"

	"github.com/itchyny/gojq"
)

type Filter func(context.Context, []any) ([]any, error)

func CreateFilter(jqFilter string) (Filter, error) {
	query, err := gojq.Parse(jqFilter)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, nodes []any) ([]any, error) {
		result := []any{}
		iter := query.RunWithContext(ctx, nodes)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				return nil, err
			}
			result = append(result, v)
		}
		return result, nil
	}, nil
}

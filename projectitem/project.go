package projectitem

import (
	"context"
	"strings"

	_ "embed"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

//go:embed query/get_project_items.graphql
var queryForUser string

type GetProjectItemsQuery struct {
	ProjectV2 struct {
		Items struct {
			Nodes    []any
			PageInfo PageInfo
		}
	}
}

type PageForUser struct {
	User GetProjectItemsQuery
}

func (p *PageForUser) Page() ([]any, PageInfo) {
	return p.User.ProjectV2.Items.Nodes, p.User.ProjectV2.Items.PageInfo
}

type PageForOrganization struct {
	Organization GetProjectItemsQuery
}

func (p *PageForOrganization) Page() ([]any, PageInfo) {
	return p.Organization.ProjectV2.Items.Nodes, p.Organization.ProjectV2.Items.PageInfo
}

func GetProjectItems(
	ctx context.Context,
	gql *api.GraphQLClient,
	descriptor ProjectDescriptor,
	filter Filter,
	limit int,
	pageLimit int,
) ([]any, error) {
	variables := map[string]any{
		"owner":         graphql.String(descriptor.Owner),
		"projectNumber": graphql.Int(descriptor.Number),
	}

	var query string
	var pageReceiver PageReceiver
	if descriptor.OwnerIsOrganization {
		query = strings.Replace(queryForUser, "user(login: $owner)", "organization(login: $owner)", 1)
		pageReceiver = &PageForOrganization{}
	} else {
		query = queryForUser
		pageReceiver = &PageForUser{}
	}

	return Collect(
		ctx,
		gql,
		query,
		pageReceiver,
		variables,
		filter,
		limit,
		pageLimit,
	)
}

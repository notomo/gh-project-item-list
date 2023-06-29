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

		gqltest.WithQueryOK("GetProjectItems", `
{
  "data": {
    "user": {
      "projectV2": {
        "items": {
          "nodes": [
            {
              "content": {
                "__typename": "Issue",
                "author": {
                  "login": "notomo"
                },
                "createdAt": "2023-01-01T00:00:00Z",
                "id": "I_kwDOGjOmFs5px3kM",
                "state": "OPEN",
                "title": "issue1",
                "updatedAt": "2023-01-01T00:00:01Z",
                "url": "https://github.com/notomo/example/issues/1"
              },
              "fieldValues": {
                "nodes": [
                  {
                    "__typename": "ProjectV2ItemFieldUserValue",
                    "field": {
                      "name": "Assignees"
                    },
                    "users": {
                      "nodes": [
                        {
                          "login": "notomo"
                        }
                      ]
                    }
                  },
                  {
                    "__typename": "ProjectV2ItemFieldRepositoryValue",
                    "field": {
                      "name": "Repository"
                    },
                    "repository": {
                      "nameWithOwner": "notomo/example"
                    }
                  },
                  {
                    "__typename": "ProjectV2ItemFieldLabelValue",
                    "field": {
                      "name": "Labels"
                    },
                    "labels": {
                      "nodes": [
                        {
                          "name": "label1"
                        }
                      ]
                    }
                  },
                  {
                    "__typename": "ProjectV2ItemFieldTextValue",
                    "field": {
                      "name": "Title"
                    },
                    "text": "issue1"
                  },
                  {
                    "__typename": "ProjectV2ItemFieldSingleSelectValue",
                    "field": {
                      "name": "Status"
                    },
                    "id": "PVTFSV_lQHOARqWjM1IA84B5GM-zgU4oNQ",
                    "name": "Todo"
                  }
                ]
              },
              "id": "PVTI_1111111111111111111-",
              "isArchived": false,
              "type": "ISSUE"
            },
            {
              "content": {
                "__typename": "Issue",
                "author": {
                  "login": "notomo"
                },
                "createdAt": "2023-02-02T00:00:00Z",
                "id": "I_kwDOGjOmFs5px3kM",
                "state": "OPEN",
                "title": "issue2",
                "updatedAt": "2023-02-02T00:00:01Z",
                "url": "https://github.com/notomo/example/issues/2"
              },
              "fieldValues": {
                "nodes": [
                  {
                    "__typename": "ProjectV2ItemFieldRepositoryValue",
                    "field": {
                      "name": "Repository"
                    },
                    "repository": {
                      "nameWithOwner": "notomo/example"
                    }
                  },
                  {
                    "__typename": "ProjectV2ItemFieldTextValue",
                    "field": {
                      "name": "Title"
                    },
                    "text": "issue2"
                  }
                ]
              },
              "id": "PVTI_2222222222222222222-",
              "isArchived": true,
              "type": "ISSUE"
            }
          ]
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
		`.[] | select(.fieldValues.nodes|any(.field.name == "Status" and .name == "Todo"))`,
		10,
		0,
		output,
	))

	want := `[
  {
    "content": {
      "__typename": "Issue",
      "author": {
        "login": "notomo"
      },
      "createdAt": "2023-01-01T00:00:00Z",
      "id": "I_kwDOGjOmFs5px3kM",
      "state": "OPEN",
      "title": "issue1",
      "updatedAt": "2023-01-01T00:00:01Z",
      "url": "https://github.com/notomo/example/issues/1"
    },
    "fieldValues": {
      "nodes": [
        {
          "__typename": "ProjectV2ItemFieldUserValue",
          "field": {
            "name": "Assignees"
          },
          "users": {
            "nodes": [
              {
                "login": "notomo"
              }
            ]
          }
        },
        {
          "__typename": "ProjectV2ItemFieldRepositoryValue",
          "field": {
            "name": "Repository"
          },
          "repository": {
            "nameWithOwner": "notomo/example"
          }
        },
        {
          "__typename": "ProjectV2ItemFieldLabelValue",
          "field": {
            "name": "Labels"
          },
          "labels": {
            "nodes": [
              {
                "name": "label1"
              }
            ]
          }
        },
        {
          "__typename": "ProjectV2ItemFieldTextValue",
          "field": {
            "name": "Title"
          },
          "text": "issue1"
        },
        {
          "__typename": "ProjectV2ItemFieldSingleSelectValue",
          "field": {
            "name": "Status"
          },
          "id": "PVTFSV_lQHOARqWjM1IA84B5GM-zgU4oNQ",
          "name": "Todo"
        }
      ]
    },
    "id": "PVTI_1111111111111111111-",
    "isArchived": false,
    "type": "ISSUE"
  }
]
`
	got := output.String()
	assert.Equal(t, want, got)
}

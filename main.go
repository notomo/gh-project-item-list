package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/henvic/httpretty"
	"github.com/notomo/gh-project-item-list/projectitem"
	"github.com/notomo/httpwriter"
	"github.com/urfave/cli/v2"
)

const (
	paramProjectUrl = "project-url"
	paramJq         = "jq"
	paramLimit      = "limit"
	paramPageLimit  = "page-limit"
	paramLog        = "log"
)

func main() {
	app := &cli.App{
		Name: "gh-project-item-list",
		Action: func(c *cli.Context) error {
			opts := api.ClientOptions{}
			logDirPath := c.String(paramLog)
			if logDirPath != "" {
				opts.Transport = &httpwriter.Transport{
					TransportFactory: func(writer io.Writer) http.RoundTripper {
						logger := &httpretty.Logger{
							Time:            true,
							TLS:             false,
							RequestHeader:   true,
							RequestBody:     true,
							ResponseHeader:  true,
							ResponseBody:    true,
							MaxResponseBody: 1000000,
							Formatters:      []httpretty.Formatter{&httpretty.JSONFormatter{}},
						}
						logger.SetOutput(writer)
						return logger.RoundTripper(nil)
					},
					GetWriter: httpwriter.MustDirectoryWriter(
						&httpwriter.Directory{Path: logDirPath},
					),
				}
			}
			gql, err := api.NewGraphQLClient(opts)
			if err != nil {
				return fmt.Errorf("create gql client: %w", err)
			}
			return projectitem.List(
				c.Context,
				gql,
				c.String(paramProjectUrl),
				c.String(paramJq),
				c.Int(paramLimit),
				c.Int(paramPageLimit),
				os.Stdout,
			)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     paramProjectUrl,
				Value:    "",
				Required: true,
				Usage:    "project url",
			},
			&cli.StringFlag{
				Name:     paramJq,
				Value:    ".[]",
				Required: false,
				Usage:    "jq query to filter project item nodes in one page",
			},
			&cli.IntFlag{
				Name:     paramLimit,
				Value:    10,
				Required: false,
				Usage:    "limit",
			},
			&cli.IntFlag{
				Name:     paramPageLimit,
				Value:    0,
				Required: false,
				Usage:    "page limit",
			},
			&cli.StringFlag{
				Name:     paramLog,
				Value:    "",
				Required: false,
				Usage:    "log directory path",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

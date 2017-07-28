package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()

	app.Name = "beachfront"
	app.Usage = "access the Beachfront services"

	var feedName, sceneId string

	app.Commands = []cli.Command{
		{
			Name:    "feedinfo",
			Aliases: []string{"fi"},
			Usage:   "get info about a feed or scene",

			Subcommands: []cli.Command{
				{
					Name:    "feedinfo",
					Aliases: []string{"fi"},
					Usage:   "get info about a feed or scene",

					Subcommands: []cli.Command{
						{
							Name:    "feedinfo",
							Aliases: []string{"fi"},
							Usage:   "get info about a feed or scene",

							Flags: []cli.Flag{
								cli.StringFlag{
									Name:        "feed,f",
									Usage:       "name of the feed (required)",
									Destination: &feedName,
								},
								cli.StringFlag{
									Name:        "scene,s",
									Usage:       "id of the scene in the feed (optional)",
									Destination: &sceneId,
								},
							},
							Action: func(c *cli.Context) error {
								switch {
								case c.IsSet("feed") && c.IsSet("scene"):
									return runFeedInfoScene(feedName, sceneId)
								case c.IsSet("feed"):
									return runFeedInfoFeed(feedName)
								default:
									return cli.NewExitError("feed info: either '--feed' or '-feed --scene' is required", 2)
								}
							},
						},
					},
				},
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "feed,f",
					Usage:       "name of the feed (required)",
					Destination: &feedName,
				},
				cli.StringFlag{
					Name:        "scene,s",
					Usage:       "id of the scene in the feed (optional)",
					Destination: &sceneId,
				},
			},
		},
		{
			Name:    "feeddownload",
			Aliases: []string{"fd"},
			Usage:   "download file(s) for a scene",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "feed,f",
					Usage:       "name of the feed (required)",
					Destination: &feedName,
				},
				cli.StringFlag{
					Name:        "scene,s",
					Usage:       "id of the scene in the feed (optional)",
					Destination: &sceneId,
				},
			},
			Action: func(c *cli.Context) error {
				switch {
				case c.IsSet("feed") && c.IsSet("scene"):
					return runFeedDownload(feedName, sceneId)
				default:
					return cli.NewExitError("feed download: --feed and --scene are both required", 2)
				}
			},
		},
		{
			Name:    "jobinfo",
			Aliases: []string{"ji"},
			Usage:   "get information about a job",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "feed,f",
					Usage:       "name of the feed (required)",
					Destination: &feedName,
				},
				cli.StringFlag{
					Name:        "scene,s",
					Usage:       "id of the scene in the feed (optional)",
					Destination: &sceneId,
				},
			},
			Action: func(c *cli.Context) error {
				switch {
				case c.IsSet("feed") && c.IsSet("scene"):
					return runFeedDownload(feedName, sceneId)
				default:
					return cli.NewExitError("feed download: --feed and --scene are both required", 2)
				}
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}

	app.Run(os.Args)
}

//---------------------------------------------------------------------

func runFeedDownload(feedName string, sceneId string) error {
	return nil
}

func runFeedInfoScene(feedName string, sceneId string) error {
	return nil
}

func runFeedInfoFeed(feedName string) error {
	return nil
}

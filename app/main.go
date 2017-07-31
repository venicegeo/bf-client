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

	app.Commands = []cli.Command{
		{
			Name:    "catalog",
			Aliases: []string{"cat"},
			Usage:   "access catalog (imagery feed) services",

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "info,i",
					Usage: "get information about a scene, a catalog, or all catalogs",
				},
				cli.StringFlag{
					Name:  "download,d",
					Usage: "downoad a scene from a catalog",
				},
			},
			Action: func(c *cli.Context) error {
				info := c.IsSet("info")
				download := c.IsSet("download")

				switch {
				case info && !download:
					arg0, err := getZeroOrOneArg("catalog info", c)
					if err != nil {
						return err
					}
					return runCatalogInfo(arg0)
				case !info && download:
					arg, err := getOneArg("catalog download", c)
					if err != nil {
						return err
					}
					return runCatalogDownload(arg)
				default:
					return cli.NewExitError("catalog: exactly one of --info and --download is required", 2)
				}
			},
		},
		{
			Name: "job",
			//Aliases: []string{""},
			Usage: "access job services",

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "info,i",
					Usage: "get information about a job or all jobs",
				},
				cli.BoolFlag{
					Name:  "submit,s",
					Usage: "submit a new job for execution",
				},
				cli.BoolFlag{
					Name:  "delete,d",
					Usage: "delete a job",
				},
			},
			Action: func(c *cli.Context) error {
				info := c.IsSet("info")
				submit := c.IsSet("submit")
				delete := c.IsSet("delete")

				switch {
				case info && !submit && !delete:
					arg, err := getZeroOrOneArg("job info", c)
					if err != nil {
						return err
					}
					return runJobInfo(arg)
				case !info && submit && !delete:
					return runJobSubmit()
				case !info && !submit && delete:
					arg, err := getOneArg("job deletion", c)
					if err != nil {
						return err
					}
					return runJobDelete(arg)
				default:
					return cli.NewExitError("job: exactly one of --info and --submit is required", 2)
				}
			},
		},
		{
			Name:    "coastline",
			Aliases: []string{"coast"},
			Usage:   "access coastline data",

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "download,d",
					Usage: "download the geojson coastline file",
				},
			},
			Action: func(c *cli.Context) error {
				download := c.IsSet("download")

				switch {
				case download:
					arg, err := getOneArg("coastline download", c)
					if err != nil {
						return err
					}
					return runCoastlineDownload(arg)
				default:
					return cli.NewExitError("coastline: --download is required", 2)
				}
			},
		},
		{
			Name:    "algorithm",
			Aliases: []string{"alg"},
			Usage:   "access the algorithm services",

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "info,i",
					Usage: "get information about an algorithm or all algorithms",
				},
			},
			Action: func(c *cli.Context) error {
				info := c.IsSet("info")

				arg, err := getZeroOrOneArg("algorithm info", c)
				if err != nil {
					return err
				}

				switch {
				case info:
					return runAlgorithmInfo(arg)
				default:
					return cli.NewExitError("algorithm: --info is required", 2)
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

func getZeroOrOneArg(area string, c *cli.Context) (string, error) {
	switch c.NArg() {
	case 0:
		return "", nil
	case 1:
		return c.Args().Get(0), nil
	default:
		return "", cli.NewExitError(area+": exactly zero or one argument is required", 2)
	}
}

func getOneArg(area string, c *cli.Context) (string, error) {
	switch c.NArg() {
	case 1:
		return c.Args().Get(0), nil
	default:
		return "", cli.NewExitError(area+": exactly one argument is required", 2)
	}
}

//---------------------------------------------------------------------

func runCatalogInfo(id string) error {
	return nil
}

func runCatalogDownload(id string) error {
	return nil
}

func runJobInfo(id string) error {
	return nil
}

func runJobSubmit() error {
	return cli.NewExitError("job: --submit not yet supported", 2)
}

func runJobDelete(id string) error {
	return cli.NewExitError("job: --delete not yet supported", 2)
}

func runCoastlineDownload(id string) error {
	return nil
}

func runAlgorithmInfo(id string) error {
	return nil
}

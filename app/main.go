/* Copyright 2017, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/venicegeo/bf-client/client"

	"gopkg.in/urfave/cli.v1"
)

func main() {

	catalogCommand := cli.Command{
		Name:    "catalog",
		Aliases: []string{"cat"},
		Usage:   "access catalog (imagery feed) services",

		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "info,i",
				Usage: "get information about a scene, a catalog, or all catalogs",
			},
			cli.BoolFlag{
				Name:  "download,d",
				Usage: "downoad a scene from a catalog",
			},
		},
		Action: func(c *cli.Context) error {
			info := c.IsSet("info")
			download := c.IsSet("download")

			switch {
			case info && !download:
				arg, err := getZeroOrOneArg("catalog info", c)
				if err != nil {
					return err
				}
				switch {
				case arg == "":
					return runCatalogInfoForCatalogs()
				case strings.Contains(arg, ":"):
					return runCatalogInfoForScene(arg)
				default:
					return runCatalogInfoForCatalog(arg)
				}
			case !info && download:
				arg, err := getOneArg("catalog download", c)
				if err != nil {
					return err
				}
				return runCatalogSceneDownload(arg)
			default:
				return cli.NewExitError("catalog: exactly one of --info and --download is required", 2)
			}
		},
	}

	jobCommand := cli.Command{
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
				if arg == "" {
					return runJobInfoForJobs()
				} else {
					return runJobInfoForJob(arg)
				}
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
	}

	coastlineCommand := cli.Command{
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
	}

	algorithmCommand := cli.Command{
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

			switch {
			case info:
				arg, err := getZeroOrOneArg("algorithm info", c)
				if err != nil {
					return err
				}
				if arg == "" {
					return runAlgorithmInfoForAll()
				} else {
					return runAlgorithmInfoForOne(arg)
				}
			default:
				return cli.NewExitError("algorithm: --info is required", 2)
			}
		},
	}

	app := cli.NewApp()
	app.Name = "beachfront"
	app.Usage = "access the Beachfront services"

	app.Commands = []cli.Command{
		catalogCommand,
		jobCommand,
		coastlineCommand,
		algorithmCommand,
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
	log.Print(c.Args())
	switch c.NArg() {
	case 1:
		return c.Args().Get(0), nil
	default:
		return "", cli.NewExitError(area+": exactly one argument is required", 2)
	}
}

func newBFClient() (*client.BFClient, error) {

	c, err := client.NewBFClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newCatalogClient() (*client.CatalogClient, error) {

	c, err := client.NewCatalogClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newJobClient() (*client.JobClient, error) {

	c, err := client.NewJobClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newCoastlineClient() (*client.CoastlineClient, error) {

	c, err := client.NewCoastlineClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newAlgorithmClient() (*client.AlgorithmClient, error) {

	c, err := client.NewAlgorithmClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

//---------------------------------------------------------------------

func runCatalogInfoForCatalogs() error {
	c, err := newCatalogClient()
	if err != nil {
		return err
	}
	s, err := c.GetCatalogInfoForCatalogs()
	if err != nil {
		return err
	}

	fmt.Print(s)

	return nil
}

func runCatalogInfoForScene(id string) error {
	c, err := newCatalogClient()
	if err != nil {
		return err
	}
	s, err := c.GetCatalogInfoForScene(id)
	if err != nil {
		return err
	}

	fmt.Print(s)

	return nil
}

func runCatalogInfoForCatalog(id string) error {
	c, err := newCatalogClient()
	if err != nil {
		return err
	}
	s, err := c.GetCatalogInfoForCatalog(id)
	if err != nil {
		return err
	}

	fmt.Print(s)

	return nil
}

func runCatalogSceneDownload(id string) error {
	c, err := newCatalogClient()
	if err != nil {
		return err
	}
	info, err := c.DoCatalogSceneDownload(id)
	if err != nil {
		return err
	}

	for k, v := range info {
		fmt.Printf("%s: %d bytes\n", k, v)
	}

	return nil
}

func runJobInfoForJobs() error {
	c, err := newJobClient()
	if err != nil {
		return err
	}
	s, err := c.GetJobInfoForJobs()
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

func runJobInfoForJob(id string) error {
	c, err := newJobClient()
	if err != nil {
		return err
	}
	s, err := c.GetJobInfoForJob(id)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

func runJobSubmit() error {
	return cli.NewExitError("job: --submit not yet supported", 2)
}

func runJobDelete(id string) error {
	return cli.NewExitError("job: --delete not yet supported", 2)
}

func runCoastlineDownload(id string) error {
	c, err := newCoastlineClient()
	if err != nil {
		return err
	}
	s, err := c.DoCoastlineDownload(id)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

func runAlgorithmInfoForAll() error {
	c, err := newAlgorithmClient()
	if err != nil {
		return err
	}
	s, err := c.GetAlgorithmInfoForAll()
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

func runAlgorithmInfoForOne(id string) error {
	c, err := newAlgorithmClient()
	if err != nil {
		return err
	}
	s, err := c.GetAlgorithmInfoForOne(id)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

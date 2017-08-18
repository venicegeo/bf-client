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

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
)

const catalogServer = "bf-ia-broker"

const catalogTimeout = time.Duration(2 * time.Minute)

type CatalogClient struct {
	url      string // "https://<bf-ia-broker>.<int.geointservices.io>"
	keyParam string // "PL_API_KEY=<123abc>"
}

func NewCatalogClient() (*CatalogClient, error) {

	fields, err := ReadBeachfrontrcFields([]string{"domain", "planet_key"})
	if err != nil {
		return nil, err
	}

	keyParam := "PL_API_KEY=" + fields["planet_key"]

	return &CatalogClient{
		url:      "https://" + catalogServer + "." + fields["domain"],
		keyParam: keyParam,
	}, nil
}

//---------------------------------------------------------------------

type Catalog struct {
	Type     string
	Features []*CatalogFeature
}

func (c *Catalog) String() string {
	s := ""
	for _, v := range c.Features {
		s += v.String() + "\n"
	}
	return s
}

type CatalogFeature struct {
	Type       string
	Geometry   *GeometryInfo
	Properties *PropertiesInfo
	Id         string
	Bbox       [4]float64
}

func (c *CatalogFeature) String() string {
	return fmt.Sprintf("[catalog-feature %s]", c.Id)
}

type Point []float64
type PointList []Point
type Any interface{}
type AnyList []Any

type GeometryInfo struct {
	Type        string
	Coordinates []AnyList // either [][]Point or [][]PointList
}

type PropertiesInfo struct {
	AcquiredDate string
	Bands        map[string]string
	CloudCover   float64
	FileFormat   string
	resolution   float64
	sensorName   string
}

type SceneInfo struct {
	Properties PropertiesInfo
}

//---------------------------------------------------------------------

func toPoint(v interface{}) (Point, bool) {
	point, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	x, ok := point[0].(float64)
	if !ok {
		return nil, false
	}
	y, ok := point[1].(float64)
	if !ok {
		return nil, false
	}
	return Point{x, y}, true
}

func toPointList(v interface{}) ([]Point, bool) {
	list, ok := v.([]interface{})
	if !ok {
		return nil, false
	}

	cl := []Point{}
	for _, v := range list {
		point, ok := toPoint(v)
		if !ok {
			return nil, false
		}
		cl = append(cl, point)
	}
	return cl, true
}

//---------------------------------------------------------------------

// landsat:LC80480102017209LGN00 -> (landsat, LC80480102017209LGN00)
func splitId(id string) (string, string, error) {

	ids := strings.Split(id, ":")
	if len(ids) != 2 || ids[0] == "" || ids[1] == "" {
		return "", "", fmt.Errorf("scene id format must be '<catalogname>:<sceneid>'")
	}

	return ids[0], ids[1], nil
}

//---------------------------------------------------------------------

func (c *CatalogClient) GetInfoForCatalogs() error {

	log.Printf("Catalog.GetInfoForCatalogs")

	return cli.NewExitError("catalog: --info for catalogs not yet supported", 2)
}

func (c *CatalogClient) GetInfoForScene(id string) (string, error) {

	log.Printf("Catalog.GetInfoForScene")

	sensor, scene, err := splitId(id)
	if err != nil {
		return "", err
	}

	path := "/planet/" + sensor + "/" + scene

	params := c.keyParam

	url := fmt.Sprintf("%s%s?%s", c.url, path, params)

	jsn, err := doHttpGetJSON(url, catalogTimeout, 200)

	obj := &Catalog{}
	err = json.Unmarshal([]byte(jsn), obj)
	if err != nil {
		return "", err
	}

	return obj.String(), nil
}

func (c *CatalogClient) GetInfoForCatalog(id string) (string, error) {

	log.Printf("Catalog.GetInfoForCatalog")

	path := "/planet/discover/" + id

	params := c.keyParam

	url := fmt.Sprintf("%s%s?%s", c.url, path, params)

	jsn, err := doHttpGetJSON(url, catalogTimeout, 200)
	if err != nil {
		return "", err
	}

	obj := &Catalog{}
	err = json.Unmarshal([]byte(jsn), obj)
	if err != nil {
		return "", err
	}

	for _, f := range obj.Features {
		newCoordinates := make([]AnyList, len(f.Geometry.Coordinates))
		for i, anyArray := range f.Geometry.Coordinates {
			newArray := make(AnyList, len(anyArray))
			for j, any := range anyArray {
				point, ok := toPoint(any)
				if ok {
					fmt.Printf("%f %f\n", point[0], point[1])
					newArray[j] = point
				} else {
					pointList, ok := toPointList(any)
					if ok {
						fmt.Printf("** %#v\n", pointList)
						newArray[j] = pointList
					} else {
						fmt.Printf("=== %T === %v ===", any, any)
						panic(17)
					}
				}
			}
			newCoordinates[i] = newArray
		}
		f.Geometry.Coordinates = newCoordinates
	}
	os.Exit(9)
	return obj.String(), nil
}

// returns map: file name -> file size
func (c *CatalogClient) DoCatalogSceneDownload(id string) (map[string]int, error) {

	log.Printf("Catalog.DoSceneDownload")

	body, err := c.GetInfoForScene(id)

	info := SceneInfo{}
	err = json.Unmarshal([]byte(body), &info)
	if err != nil {
		return nil, err
	}

	fileInfo := map[string]int{}

	i := 1
	for bandName, value := range info.Properties.Bands {

		status, byts, err := doHttpGetBytes(value, catalogTimeout)
		if err != nil {
			return nil, err
		}
		if status != 200 {
			return nil, fmt.Errorf("HTTP download failed with status %d", status)
		}

		idx := strings.LastIndex(value, "/")
		if idx == -1 {
			return nil, fmt.Errorf("unable to parse URL path: %s", value)
		}
		filename := value[idx+1 : len(value)]

		// TODO: fix path root
		err = ioutil.WriteFile("./"+filename, byts, 0600)
		if err != nil {
			return nil, err
		}

		fileInfo[filename] = len(byts)

		log.Printf("%d/%d: %s\n", i, len(info.Properties.Bands), bandName)

		i++
	}

	return fileInfo, nil
}

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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func doHttpGetJSON(
	url string,
	timeout time.Duration,
) (int, string, error) {

	client := &http.Client{Timeout: timeout}

	log.Printf("URL: %s %s", "GET", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}

	status := resp.StatusCode

	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0, "", err
	}

	return status, string(responseBody), nil
}

func doHttpGetJSONWithAuth(
	url string,
	timeout time.Duration,
	auth string,
) (int, string, error) {

	auth64 := ""
	if auth != "" {
		if auth[len(auth)-1:len(auth)] != ":" {
			auth = auth + ":"
		}
		auth64 = base64.StdEncoding.EncodeToString([]byte(auth))
	}

	client := &http.Client{Timeout: timeout}

	log.Printf("URL: %s %s", "GET", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	if auth64 != "" {
		req.Header.Set("Authorization", "Basic "+auth64)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}

	status := resp.StatusCode

	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0, "", err
	}

	return status, string(responseBody), nil
}

func doHttpGetBytes(
	url string,
	timeout time.Duration,
) (int, []byte, error) {

	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	status := resp.StatusCode

	byts, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0, nil, err
	}

	return status, byts, nil
}

// if any one field fails, the whole thing fails
func ReadBeachfrontrcFields(fields []string) (map[string]string, error) {
	user := os.Getenv("HOME")
	file, err := os.Open(user + "/.beachfrontrc")
	if err != nil {
		return nil, err
	}

	byts, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	obj := &map[string]string{}
	err = json.Unmarshal(byts, obj)
	if err != nil {
		return nil, err
	}
	//log.Printf("OBJ: %v", obj)
	results := map[string]string{}

	for _, field := range fields {
		value, ok := (*obj)[field]
		if !ok || value == "" {
			return nil, fmt.Errorf("Missing item in .beachfrontrc: '%s'", field)
		}
		results[field] = value
	}

	return results, nil
}

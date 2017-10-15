// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// InvokeFunction a function
func InvokeFunction(gateway string, name string, query []string, bytesIn *[]byte, contentType string) (*[]byte, error) {
	var resBytes []byte

	gateway = strings.TrimRight(gateway, "/")

	funcURL, err := url.Parse(gateway + "/function/" + name)
	if err != nil {
		log.Fatal(err)
	}

	if query != nil {
		q := u.Query()
		for _, valueStr := range query {
			value := strings.Split(valueStr, "=")
			if len(value) != 2 {
				return nil, fmt.Errorf("Wrong query format, should key=value %s", valueStr)
			}
			q.Add(value[0], value[1])
		}
		funcURL.RawQuery = q.Encode()
	}

	reader := bytes.NewReader(*bytesIn)
	res, err := http.Post(funcURL.String(), contentType, reader)
	if err != nil {
		fmt.Println()
		fmt.Println(err)
		return nil, fmt.Errorf("cannot connect to OpenFaaS on URL: %s", gateway)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case 200:
		var readErr error
		resBytes, readErr = ioutil.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("cannot read result from OpenFaaS on URL: %s %s", gateway, readErr)
		}

	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return nil, fmt.Errorf("Server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}

	return &resBytes, nil
}

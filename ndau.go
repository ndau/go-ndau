package ndau

//  ----- ---- --- -- -
//  Copyright 2022 Oneiro NA, Inc. All Rights Reserved.
//
//  Licensed under the Apache License 2.0 (the "License").  You may not use
//  this file except in compliance with the License.  You can obtain a copy
//  in the file LICENSE in the source distribution or at
//  https://www.apache.org/licenses/LICENSE-2.0.txt
//  - -- --- ---- -----

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	uuid "github.com/satori/uuid"

	logger "github.com/ndau/go-logger"
)

type NdauConfig struct {
	Network string
	NodeAPI string
}

type Ndau struct {
	Config *NdauConfig
	Client HttpClient
	Log    logger.Logger
}

// New creates a new Ndau client
func New(client HttpClient, config *NdauConfig, loggers ...logger.Logger) (*Ndau, error) {
	// Attach an optional logger
	var log logger.Logger
	if len(loggers) > 0 {
		log = loggers[0]
	} else {
		log = &logger.NoopLogger{}
	}

	log.Info("New go-ndau")

	return &Ndau{
		Config: config,
		Client: client,
		Log:    log,
	}, nil
}

// GetData is a general-purpose query helper
func (n *Ndau) GetData(api string, params interface{}) ([]byte, error) {
	ctx := context.Background()
	trackingNumber := uuid.NewV4().String()
	return n.GetDataWithContext(context.WithValue(ctx, "tracking_number", trackingNumber), api, params)
}

// GetDataWithContext is a general-purpose query helper
func (n *Ndau) GetDataWithContext(ctx context.Context, api string, params interface{}) ([]byte, error) {
	return n.DoWithContext(ctx, http.MethodGet, api, params)
}

// GetData is a general-purpose query helper
func (n *Ndau) PostData(api string, params interface{}) ([]byte, error) {
	ctx := context.Background()
	trackingNumber := uuid.NewV4().String()
	return n.PostDataWithContext(context.WithValue(ctx, "tracking_number", trackingNumber), api, params)
}

// PostDataWithContexts is a general-purpose query helper
func (n *Ndau) PostDataWithContext(ctx context.Context, api string, params interface{}) ([]byte, error) {
	return n.DoWithContext(ctx, http.MethodPost, api, params)
}

// DoWithContext is a general-purpose http POST/GET helper
func (n *Ndau) DoWithContext(ctx context.Context, method, api string, params interface{}) ([]byte, error) {
	trackingNumber := ctx.Value("tracking_number")

	// Create the HTTP request object
	var err error
	var req *http.Request
	endpoint := n.Config.NodeAPI + api

	// Parse query params
	queryParams := url.Values{}
	switch params.(type) {
	case map[string]interface{}:
		for k, v := range params.(map[string]interface{}) {
			// n.Log.Infof("%s | Debug %v: type = %v, val = %v", trackingNumber, k, fmt.Sprintf("%T", v), v)
			switch v.(type) {
			case int:
				queryParams.Add(k, fmt.Sprintf("%v", v))
			case float64:
				queryParams.Add(k, fmt.Sprintf("%v", v))
			case bool:
				queryParams.Add(k, fmt.Sprintf("%v", v))
			case string:
				queryParams.Add(k, fmt.Sprintf("%v", v))
				// Other cases will be added later, when needed
			default:
				queryParams.Add(k, fmt.Sprintf("%v", v))
			}
		}

		if len(queryParams) > 0 {
			switch strings.ToUpper(method) {
			case http.MethodGet:
				endpoint = endpoint + "?" + queryParams.Encode()
				req, err = http.NewRequest(method, endpoint, nil)
			case http.MethodPost:
				req, err = http.NewRequest(method, endpoint, strings.NewReader(queryParams.Encode()))
			default:
			}
		}
	case []interface{}:
		jsonData, err := json.Marshal(params)
		if err != nil {
			n.Log.Errorf("%s | Failed to unmarshall post data. Error %+v", trackingNumber, err)
			return nil, err
		}
		req, err = http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	default:
		req, err = http.NewRequest(method, endpoint, nil)
	}

	if err != nil {
		n.Log.Errorf("%s | Failed to new a request object. Error %+v", trackingNumber, err)
		return nil, err
	}

	//Make the request and parse the response
	resp, err := n.Client.Do(req)
	if err != nil {
		n.Log.Errorf("%s | Failed to send request to %s. Error %+v", trackingNumber, endpoint, err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		n.Log.Errorf("%s | Failed to read from %s. Status code %d, Msg = ", trackingNumber, endpoint, resp.StatusCode, resp.Status)
		return nil, fmt.Errorf(resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		n.Log.Errorf("%s | Failed to read from %s. Error %+v", trackingNumber, endpoint, err)
		return nil, err
	}

	return body, nil
}

// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/mainflux/mainflux/things"
	thingsapi "github.com/mainflux/mainflux/things/api/things/http"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/api"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	contentType  = "application/json"
	token        = "token"
	invalidToken = "invalid"
	email        = "user@example.com"
	validJson    = "{\n	\"name\": \"eu-agents\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
	invalidJson  = "{"
	wrongID      = 0
	maxNameSize  = 1024
	channelsNum  = 3
	limit        = 10
)

var (
	agentGroup = fleet.AgentGroup{
		ID:          "",
		MFOwnerID:   "",
		Name:        types.Identifier{},
		Description: "",
		MFChannelID: "",
		Tags:        nil,
		Created:     time.Time{},
	}
	metadata    = map[string]interface{}{"meta": "data"}
	invalidName = strings.Repeat("m", maxNameSize+1)
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
	body        io.Reader
}

type clientServer struct {
	service fleet.AgentGroupService
	server  *httptest.Server
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}
	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}
	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}
	return tr.client.Do(req)
}

func generateChannels() map[string]things.Channel {
	channels := make(map[string]things.Channel, channelsNum)
	for i := 0; i < channelsNum; i++ {
		id := strconv.Itoa(i + 1)
		channels[id] = things.Channel{
			ID:       id,
			Owner:    email,
			Metadata: metadata,
		}
	}
	return channels
}

func newThingsService(auth mainflux.AuthServiceClient) things.Service {
	return flmocks.NewThingsService(map[string]things.Thing{}, generateChannels(), auth)
}

func newThingsServer(svc things.Service) *httptest.Server {
	mux := thingsapi.MakeHandler(mocktracer.New(), svc)
	return httptest.NewServer(mux)
}

func newService(auth mainflux.AuthServiceClient, url string) fleet.Service {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()
	var logger *zap.Logger
	config := mfsdk.Config{
		BaseURL: url,
	}

	mfsdk := mfsdk.NewSDK(config)
	return fleet.NewFleetService(logger, auth, agentRepo, agentGroupRepo, nil, mfsdk)
}

func newServer(svc fleet.Service) *httptest.Server {
	mux := api.MakeHandler(mocktracer.New(), "fleet", svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestCreateAgentGroup(t *testing.T) {
	cli := newClientServer(t)
	defer cli.server.Close()

	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"add a valid agent group": {
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusCreated,
			location:    "/agent_groups",
		},
		"add a valid agent group with invalid token": {
			req:         validJson,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/agent_groups",
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client:      cli.server.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/agent_groups", cli.server.URL),
			contentType: tc.contentType,
			token:       tc.auth,
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
	}

}

func TestViewAgentGroup(t *testing.T) {
	cli := newClientServer(t)

	ag, err := createAgentGroup(t, "ue-agent-group", &cli)
	require.Nil(t, err, "unexpected error: %s", err)

	cases := map[string]struct {
		id          string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"view a existing agent group": {
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
		},
		"view a non-existing agent group": {
			id:          "9bb1b244-a199-93c2-aa03-28067b431e2c",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"view a agent group with invalid content type": {
			id:          ag.ID,
			contentType: "application",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"view a agent group with empty content type": {
			id:          ag.ID,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"view a agent group with a invalid token": {
			id:          ag.ID,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
		},
		"view a agent group with a empty token": {
			id:          ag.ID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client:      cli.server.Client(),
			method:      http.MethodGet,
			url:         fmt.Sprintf("%s/agent_groups/%s", cli.server.URL, tc.id),
			contentType: tc.contentType,
			token:       tc.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected erro %s", desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
	}

}

func TestListAgentGroup(t *testing.T) {
	cli := newClientServer(t)

	var data []agentGroupRes
	for i := 0; i < limit; i++ {
		ag, err := createAgentGroup(t, fmt.Sprintf("ue-agent-group-%d", i), &cli)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		data = append(data, agentGroupRes{
			ID:             ag.ID,
			Name:           ag.Name.String(),
			Description:    ag.Description,
			Tags:           ag.Tags,
			TsCreated:      ag.Created,
			MatchingAgents: nil,
		})
	}

	cases := map[string]struct {
		auth   string
		status int
		url    string
		res    []agentGroupRes
		total  uint64
	}{
		"retrieve a list of agent groups": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group with empty token": {
			auth:   "",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid token": {
			auth:   invalidToken,
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid dir": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=name&dir=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with negative offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", -1, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with negative limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, -5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with zero limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 0),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group without offset": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?limit=%d", limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group without limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d", 1),
			res:    data[1:limit],
			total:  limit - 1,
		},
		"get a list of agent group with limit greater than max": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 110),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with default URL": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s", ""),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group with invalid number of params": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=4&limit=4&limit=5&offset=5"),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=e&limit=5"),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=5&limit=e"),
			res:    nil,
			total:  0,
		},
		"get a list of agent group filtering with invalid name": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&name=%s", 0, 5, invalidName),
			res:    nil,
			total:  0,
		},
		"get a list of agent group sorted with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong&dir=desc", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group sorted with invalid direction": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=name&dir=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client:      cli.server.Client(),
			method:      http.MethodGet,
			url:         fmt.Sprintf(fmt.Sprintf("%s/agent_groups%s", cli.server.URL, tc.url)),
			contentType: contentType,
			token:       tc.auth,
		}
		res, err := req.make()
		require.Nil(t, err, "%s: unexpected error: %s", desc, err)
		var body agentGroupsPageRes
		json.NewDecoder(res.Body).Decode(&body)
		total := uint64(len(body.AgentGroups))
		assert.Equal(t, res.StatusCode, tc.status, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		assert.Equal(t, total, tc.total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.total, total))
	}

}

func createAgentGroup(t *testing.T, name string, cli *clientServer) (fleet.AgentGroup, error) {
	agCopy := agentGroup
	validName, err := types.NewIdentifier(name)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	agCopy.Name = validName
	ag, err := cli.service.CreateAgentGroup(context.Background(), token, agCopy)
	if err != nil {
		return fleet.AgentGroup{}, err
	}
	return ag, nil
}

func newClientServer(t *testing.T) clientServer {
	t.Helper()
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)
	fleetServer := newServer(fleetService)

	return clientServer{
		service: fleetService,
		server:  fleetServer,
	}
}

type agentGroupRes struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description,omitempty"`
	Tags           types.Tags     `json:"tags"`
	TsCreated      time.Time      `json:"ts_created,omitempty"`
	MatchingAgents types.Metadata `json:"matching_agents,omitempty"`
	created        bool
}

type agentGroupsPageRes struct {
	Total       uint64          `json:"total"`
	Offset      uint64          `json:"offset"`
	Limit       uint64          `json:"limit"`
	AgentGroups []agentGroupRes `json:"agentGroups"`
}

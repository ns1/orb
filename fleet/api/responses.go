/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/types"
	"net/http"
)

var (
	_ types.Response = (*agentGroupRes)(nil)
	_ types.Response = (*agentRes)(nil)
)

type agentGroupRes struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Tags    types.Tags `json:"tags"`
	created bool
}

func (s agentGroupRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s agentGroupRes) Headers() map[string]string {
	return map[string]string{}
}

func (s agentGroupRes) Empty() bool {
	return false
}

type agentRes struct {
	ID        string `json:"id"`
	Key       string `json:"key,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	Name      string `json:"name"`
	State     string `json:"state"`
	created   bool
}

func (s agentRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s agentRes) Headers() map[string]string {
	return map[string]string{}
}

func (s agentRes) Empty() bool {
	return false
}

type viewAgentRes struct {
	ID           string                 `json:"id"`
	ChannelID    string                 `json:"channel_id,omitempty"`
	Owner        string                 `json:"-"`
	Name         string                 `json:"name,omitempty"`
	State        string                 `json:"state"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
}

func (res viewAgentRes) Code() int {
	return http.StatusOK
}

func (res viewAgentRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewAgentRes) Empty() bool {
	return false
}

type agentsPageRes struct {
	pageRes
	Agents []viewAgentRes `json:"agents"`
}

func (res agentsPageRes) Code() int {
	return http.StatusOK
}

func (res agentsPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res agentsPageRes) Empty() bool {
	return false
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
	Order  string `json:"order"`
	Dir    string `json:"direction"`
}

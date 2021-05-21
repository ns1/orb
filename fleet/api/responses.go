/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/types"
	"net/http"
)

var (
	_ types.Response = (*selectorRes)(nil)
	_ types.Response = (*agentRes)(nil)
)

type selectorRes struct {
	Name    string `json:"name"`
	created bool
}

func (s selectorRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s selectorRes) Headers() map[string]string {
	return map[string]string{}
}

func (s selectorRes) Empty() bool {
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

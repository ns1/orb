/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otel

import (
	"net/http"
)

type totalAgents struct {
	Total uint64 `json:"total"`
}

type agentBackendTapsRes struct {
	Name             string      `json:"name"`
	InputType        string      `json:"input_type"`
	ConfigPredefined []string    `json:"config_predefined"`
	TotalAgents      totalAgents `json:"agents"`
}

func (res agentBackendTapsRes) Code() int {
	return http.StatusOK
}

func (res agentBackendTapsRes) Headers() map[string]string {
	return map[string]string{}
}

func (res agentBackendTapsRes) Empty() bool {
	return true
}

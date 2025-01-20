/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/orb-community/orb/agent/policies"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func (p *pktvisorBackend) ApplyPolicy(data policies.PolicyData, updatePolicy bool) error {

	if updatePolicy {
		// To update a policy it's necessary first remove it and then apply a new version
		if err := p.RemovePolicy(data); err != nil {
			p.logger.Warn("policy failed to remove", zap.String("policy_id", data.ID), zap.String("policy_name", data.Name), zap.Error(err))
		}
	}

	p.logger.Debug("pktvisor policy apply", zap.String("policy_id", data.ID), zap.Any("data", data.Data))

	fullPolicy := map[string]interface{}{
		"version": "1.0",
		"visor": map[string]interface{}{
			"policies": map[string]interface{}{
				data.Name: data.Data,
			},
		},
	}

	policyYaml, err := yaml.Marshal(fullPolicy)
	if err != nil {
		p.logger.Warn("yaml policy marshal failure", zap.String("policy_id", data.ID), zap.Any("policy", fullPolicy))
		return err
	}

	var resp map[string]interface{}
	err = p.request("policies", &resp, http.MethodPost, bytes.NewBuffer(policyYaml), "application/x-yaml", ApplyPolicyTimeout)
	if err != nil {
		p.logger.Warn("yaml policy application failure", zap.String("policy_id", data.ID), zap.ByteString("policy", policyYaml))
		return err
	}

	return nil
}

func (p *pktvisorBackend) RemovePolicy(data policies.PolicyData) error {
	p.logger.Debug("pktvisor policy remove", zap.String("policy_id", data.ID))
	var resp interface{}
	var name string
	// Since we use Name for removing policies not IDs, if there is a change, we need to remove the previous name of the policy
	if data.PreviousPolicyData != nil && data.PreviousPolicyData.Name != data.Name {
		name = data.PreviousPolicyData.Name
	} else {
		name = data.Name
	}
	if err := p.request(fmt.Sprintf("policies/%s", name), &resp, http.MethodDelete, http.NoBody, "application/json", RemovePolicyTimeout); err != nil {
		if strings.Contains(err.Error(), "404") {
			p.logger.Warn("ignoring error from removing a policy not found", zap.String("policy_id", data.ID), zap.String("policy_name", name))
			return nil
		}
		return err
	}
	return nil
}

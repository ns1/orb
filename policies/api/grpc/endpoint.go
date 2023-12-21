// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/orb-community/orb/policies"
	"gopkg.in/yaml.v3"
)

func retrievePolicyEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(accessByIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		policy, err := svc.ViewPolicyByIDInternal(ctx, req.PolicyID, req.OwnerID)
		if err != nil {
			return policyRes{}, err
		}
		data, format, err := extractData(policy)
		if err != nil {
			return policyRes{}, err
		}

		return policyRes{
			id:      policy.ID,
			name:    policy.Name.String(),
			format:  format,
			backend: policy.Backend,
			version: policy.Version,
			data:    data,
		}, nil
	}
}

func extractData(policy policies.Policy) (data []byte, format string, err error) {
	if policy.Format == "yaml" {
		data, err = yaml.Marshal(policy.Policy)
		format = "yaml"
	} else {
		data, err = json.Marshal(policy.Policy)
		format = "json"
	}
	if err != nil {
		return nil, "", err
	}
	return data, format, nil
}

func retrievePoliciesByGroupsEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(accessByGroupIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		plist, err := svc.ListPoliciesByGroupIDInternal(ctx, req.GroupIDs, req.OwnerID)
		if err != nil {
			return policyInDSListRes{}, err
		}
		policies := make([]policyInDSRes, len(plist))
		for i, policy := range plist {
			data, format, err := extractData(policy.Policy)
			if err != nil {
				return policyInDSListRes{}, err
			}
			policies[i] = policyInDSRes{
				id:           policy.ID,
				name:         policy.Name.String(),
				backend:      policy.Backend,
				version:      policy.Version,
				format:       format,
				data:         data,
				datasetID:    policy.DatasetID,
				agentGroupID: policy.AgentGroupID,
			}
		}

		return policyInDSListRes{policies: policies}, nil
	}
}

func retrieveDatasetEnpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(accessDatasetByIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		dataset, err := svc.ViewDatasetByIDInternal(ctx, req.ownerID, req.datasetID)
		if err != nil {
			return nil, err
		}
		return datasetRes{
			id:           dataset.ID,
			agentGroupID: dataset.AgentGroupID,
			policyID:     dataset.PolicyID,
			sinkIDs:      *dataset.SinkIDs,
		}, nil
	}
}

func retrieveDatasetsByGroupsEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(accessByGroupIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		dsList, err := svc.ListDatasetsByGroupIDInternal(ctx, req.GroupIDs, req.OwnerID)
		if err != nil {
			return datasetListRes{}, err
		}
		datasets := make([]datasetRes, len(dsList))
		for i, ds := range dsList {
			datasets[i] = datasetRes{
				id:           ds.ID,
				agentGroupID: ds.AgentGroupID,
				sinkIDs:      *ds.SinkIDs,
				policyID:     ds.PolicyID,
			}
		}

		return datasetListRes{datasets: datasets}, nil
	}
}

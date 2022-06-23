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
	"github.com/ns1labs/orb/policies"

	"github.com/go-kit/kit/endpoint"
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
		data, err := json.Marshal(policy.Policy)
		if err != nil {
			return policyRes{}, err
		}
		return policyRes{
			id:      policy.ID,
			name:    policy.Name.String(),
			backend: policy.Backend,
			version: policy.Version,
			data:    data,
		}, nil
	}
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
			data, err := json.Marshal(policy.Policy.Policy)
			if err != nil {
				return policyInDSListRes{}, err
			}
			policies[i] = policyInDSRes{
				id:           policy.ID,
				name:         policy.Name.String(),
				backend:      policy.Backend,
				version:      policy.Version,
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
			sinkIDs:      dataset.SinkIDs,
		}, nil
	}
}

// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
)

type addReq struct {
	Name       string         `json:"name"`
	Backend    string         `json:"backend"`
	Policy     types.Metadata `json:"policy,omitempty"`
	Format     string         `json:"format,omitempty"`
	PolicyData string         `json:"policy_data,omitempty"`
	token      string
}

func (req addReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.Name == "" {
		return errors.ErrMalformedEntity
	}
	if req.Backend == "" {
		return errors.ErrMalformedEntity
	}

	if req.Policy == nil {
		// passing policy data blob in the specified format
		if req.Format == "" || req.PolicyData == "" {
			return errors.ErrMalformedEntity
		}
	} else {
		// policy is in json, verified by the back ends later
		if req.Format != "" || req.PolicyData != "" {
			return errors.ErrMalformedEntity
		}
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

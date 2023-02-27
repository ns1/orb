/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package orb

import (
	"errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/policies/backend"
)

var _ backend.Backend = (*orbBackend)(nil)

type orbBackend struct {
}

func (p orbBackend) Validate(policy types.Metadata) error {
	if version, ok := policy["version"]; ok {
		if version != CurrentSchemaVersion {
			return errors.New("unsupported version")
		}
	} else {
		return errors.New("missing version")
	}
	if _, ok := policy["orb"]; !ok {
		return errors.New("malformed policy")
	}
	// todo finish validation
	return nil
}

func (p orbBackend) ConvertFromFormat(format string, policy string) (types.Metadata, error) {
	return nil, errors.New("unsupported format")
}

func (p orbBackend) SupportsFormat(format string) bool {
	return false
}

func Register() bool {
	backend.Register("orb", &orbBackend{})
	return true
}

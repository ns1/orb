// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

var (
	logger, _   = zap.NewDevelopment()
	description = "An example prometheus sink"
)

func TestSinkSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	skID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	conflictNameID, err := types.NewIdentifier("my-sink-conflict")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		ID:          skID.String(),
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	invalidOwnerSink := sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		ID:          skID.String(),
		Created:     time.Now(),
		MFOwnerID:   "",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkMalformedOwnerID := sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		ID:          skID.String(),
		Created:     time.Now(),
		MFOwnerID:   "123",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkCopy := sink
	sinkCopy.Name = conflictNameID
	_, err = sinkRepo.Save(context.Background(), sinkCopy)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		sink sinks.Sink
		err  error
	}{
		"create a new sink": {
			sink: sink,
			err:  nil,
		},
		"create a sink that already exist": {
			sink: sinkCopy,
			err:  errors.ErrConflict,
		},
		"create a sink with invalid ownerID": {
			sink: invalidOwnerSink,
			err:  errors.ErrMalformedEntity,
		},
		"create a sink with a malformed ownerID": {
			sink: sinkMalformedOwnerID,
			err:  errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := sinkRepo.Save(context.Background(), tc.sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}

}

func TestSinkUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalideOwnerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalideID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	sink.ID = sinkID

	nameConflict, err := types.NewIdentifier("my-sink-conflict")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	sinkConflictName := sinks.Sink{
		Name:        nameConflict,
		Description: &description,
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	_, err = sinkRepo.Save(context.Background(), sinkConflictName)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		sink sinks.Sink
		err  error
	}{
		"update a existing sink": {
			sink: sink,
			err:  nil,
		},
		"update a non-existing sink with a existing user": {
			sink: sinks.Sink{
				ID:        invalideID.String(),
				MFOwnerID: oID.String(),
			},
			err: sinks.ErrNotFound,
		},
		"update a existing sink with a non-existing user": {
			sink: sinks.Sink{
				ID:        sinkID,
				MFOwnerID: invalideOwnerID.String(),
			},
			err: sinks.ErrNotFound,
		},
		"update a non-existing sink with a non-existing user": {
			sink: sinks.Sink{
				ID:        invalideID.String(),
				MFOwnerID: invalideOwnerID.String(),
			},
			err: sinks.ErrNotFound,
		},
		"update a sink with malformed ownerID": {
			sink: sinks.Sink{
				ID:        sinkID,
				MFOwnerID: "123",
			},
			err: errors.ErrMalformedEntity,
		},
		"update a existing sink with conflict name": {
			sink: sinks.Sink{
				ID:        sinkID,
				Name:      nameConflict,
				Backend:   "prometheus",
				MFOwnerID: oID.String(),
				Config:    map[string]interface{}{"remote_host": "data", "username": "dbuser"},
				Tags:      map[string]string{"cloud": "aws"},
			},
			err: errors.ErrConflict,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := sinkRepo.Update(context.Background(), tc.sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}

}

func TestSinkRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		sinkID string
		nameID string
		err    error
	}{
		"retrive existing sink by sinkID": {
			sinkID: sinkID,
			nameID: sink.Name.String(),
			err:    nil,
		},
		"retrive non-existing sink by sinkID": {
			sinkID: "",
			nameID: sink.Name.String(),
			err:    errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := sinkRepo.RetrieveById(context.Background(), tc.sinkID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}

}

func TestMultiSinkRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	wrongoID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {

		nameID, err := types.NewIdentifier(fmt.Sprintf("my-sink-%d", i))
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		sink := sinks.Sink{
			Name:        nameID,
			Description: &description,
			Backend:     "prometheus",
			Created:     time.Now(),
			MFOwnerID:   oID.String(),
			Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
			Tags:        map[string]string{"cloud": "aws"},
		}

		_, err = sinkRepo.Save(context.Background(), sink)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner        string
		pageMetadata sinks.PageMetadata
		size         uint64
	}{
		"retrieve all sinks with existing owner": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
			},
			size: n,
		},
		"retrieve subset of sinks with existing owner": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: n / 2,
				Limit:  n,
				Total:  n,
			},
			size: n / 2,
		},
		"retrieve sinks with no-existing owner": {
			owner: wrongoID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
			},
			size: 0,
		},
		"retrieve sinks with no-existing name": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Name:   "wrong",
				Total:  0,
			},
			size: 0,
		},
		"retrieve sinks sorted by name ascendent": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "asc",
			},
			size: n,
		},
		"retrieve sinks sorted by name descendent": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "desc",
			},
			size: n,
		},
		"retrieve sinks filtered by tags": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Tags:   map[string]string{"cloud": "aws"},
			},
			size: n,
		},
		"retrieve sinks filtered by metadata": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset:   0,
				Limit:    n,
				Total:    0,
				Metadata: map[string]interface{}{"username": "dbuser", "remote_host": "my.prometheus-host.com"},
			},
			size: 0,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			page, err := sinkRepo.RetrieveAllByOwnerID(context.Background(), tc.owner, tc.pageMetadata)
			require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
			size := uint64(len(page.Sinks))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d", desc, tc.size, size))
			assert.Equal(t, tc.pageMetadata.Total, page.Total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.pageMetadata.Total, page.Total))

			if size > 0 {
				testSortSinks(t, tc.pageMetadata, page.Sinks)
			}
		})
	}
}

func TestSinkRemoval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkName, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        sinkName,
		Description: &description,
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	sink.ID = sinkID

	cases := map[string]struct {
		sink sinks.Sink
		err  error
	}{
		"delete existing sink": {
			sink: sink,
			err:  nil,
		},
		"delete non-existent sink": {
			sink: sink,
			err:  nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := sinkRepo.Remove(context.Background(), tc.sink.MFOwnerID, tc.sink.ID)
			require.Nil(t, err, fmt.Sprintf("%s: failed to remove sink due to: %s", desc, err))

			_, err = sinkRepo.RetrieveById(context.Background(), tc.sink.ID)
			require.True(t, errors.Contains(err, sinks.ErrNotFound), fmt.Sprintf("%s: expected %s got %s", desc, sinks.ErrNotFound, err))
		})
	}
}

func TestSinkRetrieveInternal(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	sink.ID = sinkID

	cases := map[string]struct {
		sinkID  string
		ownerID string
		nameID  string
		err     error
	}{
		"retrieve existing sink by sinkID and ownerID": {
			sinkID:  sinkID,
			ownerID: sink.MFOwnerID,
			nameID:  sink.Name.String(),
			err:     nil,
		},
		"retrieve sink with empty sinkID and ownerID": {
			sinkID:  "",
			ownerID: "",
			nameID:  sink.Name.String(),
			err:     errors.ErrSelectEntity,
		},
		"retrieve non-existing sink by sinkID and ownerID": {
			sinkID:  "invalid",
			ownerID: "invalid",
			nameID:  sink.Name.String(),
			err:     errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := sinkRepo.RetrieveByOwnerAndId(context.Background(), tc.ownerID, tc.sinkID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}

}

func TestUpdateSinkState(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	fakeSinkID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	fakeOwnerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	sink.ID = sinkID

	cases := map[string]struct {
		sinkID  string
		ownerID string
		state   sinks.State
		msg     string
		err     error
	}{
		"update sink state with no error": {
			sinkID:  sinkID,
			ownerID: sink.MFOwnerID,
			state:   sinks.State(1),
			msg:     "",
			err:     nil,
		},
		"update sink state with error": {
			sinkID:  sinkID,
			ownerID: sink.MFOwnerID,
			state:   sinks.State(2),
			msg:     "failed",
			err:     nil,
		},
		"update sink state of a non-existent sink": {
			sinkID:  fakeSinkID.String(),
			ownerID: sink.MFOwnerID,
			state:   sinks.State(2),
			msg:     "failed",
			err:     sinks.ErrUpdateEntity,
		},
		"update sink state of a non-existent owner": {
			sinkID:  sinkID,
			ownerID: fakeOwnerID.String(),
			state:   sinks.State(2),
			msg:     "failed",
			err:     sinks.ErrUpdateEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "test", desc)
			err := sinkRepo.UpdateSinkState(context.Background(), tc.sinkID, tc.msg, tc.ownerID, tc.state)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
			// only validate success scenarios
			if tc.err == nil {
				got, err := sinkRepo.RetrieveById(ctx, sinkID)
				assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
				assert.Equal(t, tc.state, got.State, fmt.Sprintf("%s: expected state %d got %d", desc, tc.state, got.State))
				assert.Equal(t, tc.msg, got.Error, fmt.Sprintf("%s: expected msg %s got %s", desc, tc.msg, got.Error))
			}
		})
	}

}

func testSortSinks(t *testing.T, pm sinks.PageMetadata, sks []sinks.Sink) {
	t.Helper()
	switch pm.Order {
	case "name":
		current := sks[0]
		for _, res := range sks {
			if pm.Dir == "asc" {
				assert.GreaterOrEqual(t, res.Name.String(), current.Name.String())
			}
			if pm.Dir == "desc" {
				assert.GreaterOrEqual(t, current.Name.String(), res.Name.String())
			}
			current = res
		}
	default:
		break
	}
}

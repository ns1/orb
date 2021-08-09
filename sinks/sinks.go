/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks/backend"
	"time"
)

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = errors.New("failed to scan metadata in db")

	// ErrSelectEntity indicates error while reading entity from database
	ErrSelectEntity = errors.New("select entity from db error")

	// ErrEntityConnected indicates error while checking connection in database
	ErrEntityConnected = errors.New("check connection in database error")

	// ErrUpdateEntity indicates error while updating a entity
	ErrUpdateEntity = errors.New("failed to update entity")

	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	ErrRemoveEntity = errors.New("failed to remove entity")
)

type Sink struct {
	ID          string
	Name        types.Identifier
	MFOwnerID   string
	Description string
	Backend     string
	Config      types.Metadata
	Tags        types.Tags
	Status      string
	Error       string
	Created     time.Time
}

// Page contains page related metadata as well as list of sinks that
// belong to this page
type Page struct {
	PageMetadata
	Sinks []Sink
}

// SinkService Sink CRUD interface
type SinkService interface {
	// CreateSink creates new data sink
	CreateSink(ctx context.Context, token string, s Sink) (Sink, error)
	// UpdateSink by id
	UpdateSink(ctx context.Context, token string, s Sink) error
	// ListSinks retrieves data about sinks
	ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error)
	// ListBackends retrieves a list of available backends
	ListBackends(ctx context.Context, token string) ([]string, error)
	// ViewBackend retrieves a backend by the name
	ViewBackend(ctx context.Context, token string, key string) (backend.Backend, error)
	// ViewSink retrieves a sink by id
	ViewSink(ctx context.Context, token string, key string) (Sink, error)
	// DeleteSink delete a existing sink by id
	DeleteSink(ctx context.Context, token string, key string) error
	// ValidateSink validates a sink
	ValidateSink(ctx context.Context, token string, s Sink) (Sink, error)
}

type SinkRepository interface {
	// Save persists the Sink. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, sink Sink) (string, error)
	// Update performs an update to the existing sink, A non-nil error is
	// returned to indicate operation failure
	Update(ctx context.Context, sink Sink) error
	// RetrieveAll retrieves Sinks
	RetrieveAll(ctx context.Context, owner string, pm PageMetadata) (Page, error)
	// RetrieveById retrieves a Sink by Id
	RetrieveById(ctx context.Context, key string) (Sink, error)
	// Remove a existing Sink by id
	Remove(ctx context.Context, owner string, id string) error
	// RetrieveToValidate validates a sink
	RetrieveToValidate(ctx context.Context, name, owner string) error
}

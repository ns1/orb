/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/authentication_type/basicauth"
	"github.com/orb-community/orb/sinks/backend"
	"go.uber.org/zap"
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

	ErrInvalidBackend = errors.New("No available backends")
)

const (
	Unknown State = iota
	Active
	Error
	Idle
	Warning
	Provisioning
	ProvisioningError
)

type State int

var stateMap = [...]string{
	"unknown",
	"active",
	"error",
	"idle",
	"warning",
	"provisioning",
	"provisioning_error",
}

const MetadataLabelOtel = "opentelemetry"

type Filter struct {
	StateFilter string
}

var stateRevMap = map[string]State{
	"unknown":            Unknown,
	"active":             Active,
	"error":              Error,
	"idle":               Idle,
	"warning":            Warning,
	"provisioning":       Provisioning,
	"provisioning_error": ProvisioningError,
}

func (s State) String() string {
	return stateMap[s]
}

func (s *State) Scan(value interface{}) error {
	asString, ok := value.(string)
	if !ok {
		asBytes, ok := value.([]byte)
		if !ok {
			return errors.New("Scan source is not []byte")
		}
		asString = string(asBytes)
	}
	*s = stateRevMap[asString]
	return nil
}
func (s State) Value() (driver.Value, error) { return s.String(), nil }

func NewStateFromString(state string) State {
	return stateRevMap[state]
}

func NewConfigBackends(e backend.Backend, a authentication_type.AuthenticationType) Configuration {
	return Configuration{
		Exporter:       e,
		Authentication: a,
	}
}

type Configuration struct {
	Exporter       backend.Backend                        `json:"exporter" ,yaml:"exporter"`
	Authentication authentication_type.AuthenticationType `json:"authentication" ,yaml:"authentication"`
}

type Sink struct {
	ID          string
	Name        types.Identifier
	MFOwnerID   string
	Description *string
	Backend     string
	Config      types.Metadata
	Format      string
	ConfigData  string
	Tags        types.Tags
	State       State
	Error       string
	Created     time.Time
}

func (s *Sink) GetAuthenticationTypeName() string {
	authMeta := s.Config.GetSubMetadata("authentication")
	// Defaults to basicauth
	if authMeta == nil {
		return basicauth.AuthType
	}
	return authMeta["type"].(string)
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
	UpdateSink(ctx context.Context, token string, s Sink) (Sink, error)
	// UpdateSinkInternal by id
	UpdateSinkInternal(ctx context.Context, s Sink) (Sink, error)
	// ListSinks retrieves data about sinks
	ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error)
	// ListSinksInternal retrieves data from sinks filtered by SinksFilter for Services like Maestro, to build DeploymentEntries
	ListSinksInternal(ctx context.Context, filter Filter) (Page, error)
	// ListBackends retrieves a list of available backends
	ListBackends(ctx context.Context, token string) ([]string, error)
	// ViewBackend retrieves a backend by the name
	ViewBackend(ctx context.Context, token string, key string) (backend.Backend, error)
	// ListAuthenticationTypes retrieves a list of available AuthenticationTypes
	ListAuthenticationTypes(ctx context.Context, token string) ([]authentication_type.AuthenticationTypeConfig, error)
	// ViewAuthenticationType retrieves a AuthenticationType by the name
	ViewAuthenticationType(ctx context.Context, token string, key string) (authentication_type.AuthenticationTypeConfig, error)
	// ViewSink retrieves a sink by id, for View, does not send password
	ViewSink(ctx context.Context, token string, key string) (Sink, error)
	// ViewSinkInternal retrieves a sink by id, via GRPC, sends password
	ViewSinkInternal(ctx context.Context, ownerID string, key string) (Sink, error)
	// DeleteSink delete a existing sink by id
	DeleteSink(ctx context.Context, token string, key string) error
	// ValidateSink validate a sink configuration without saving
	ValidateSink(ctx context.Context, token string, sink Sink) (Sink, error)
	// ChangeSinkStateInternal change the sink internal state from new/idle/active
	ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state State) error
	// GetLogger gets service logger to log within gokit's packages
	GetLogger() *zap.Logger
}

type SinkRepository interface {
	// Save persists the Sink. Successful operation is indicated by non-nil error response.
	Save(ctx context.Context, sink Sink) (string, error)
	// Update performs an update to the existing sink, A non-nil error is
	// returned to indicate operation failure
	Update(ctx context.Context, sink Sink) error
	// RetrieveAllByOwnerID retrieves Sinks by OwnerID
	RetrieveAllByOwnerID(ctx context.Context, owner string, pm PageMetadata) (Page, error)
	// SearchAllSinks search Sinks for internal usage like services
	SearchAllSinks(ctx context.Context, filter Filter) ([]Sink, error)
	// RetrieveById retrieves a Sink by ID
	RetrieveById(ctx context.Context, key string) (Sink, error)
	// RetrieveByOwnerAndId retrieves a By OwnerId And SinkId
	RetrieveByOwnerAndId(ctx context.Context, ownerID string, key string) (Sink, error)
	// Remove an existing Sink by id
	Remove(ctx context.Context, owner string, key string) error
	// UpdateSinkState updates sink state like active, idle, new, unknown
	UpdateSinkState(ctx context.Context, sinkID string, msg string, ownerID string, state State) error
	// GetVersion for migrate service
	GetVersion(ctx context.Context) (string, error)
	// UpsertVersion for migrate service
	UpdateVersion(ctx context.Context, version string) error
}

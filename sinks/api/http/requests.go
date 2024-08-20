// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/authentication_type/basicauth"
	"github.com/orb-community/orb/sinks/backend"
	"gopkg.in/yaml.v3"
)

const (
	maxLimitSize = 100
	maxNameSize  = 1024
	nameOrder    = "name"
	idOrder      = "id"
	ascDir       = "asc"
	descDir      = "desc"
)

type addReq struct {
	Name        string         `json:"name,omitempty"`
	Backend     string         `json:"backend,omitempty"`
	Config      types.Metadata `json:"config,omitempty"`
	Format      string         `json:"format,omitempty"`
	ConfigData  string         `json:"config_data,omitempty"`
	Description string         `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags,omitempty"`
	token       string
}

func GetConfigurationAndMetadataFromMeta(backendName string, config types.Metadata) (configSvc *sinks.Configuration, exporter types.Metadata, authentication types.Metadata, err error) {

	if !backend.HaveBackend(backendName) {
		return nil, nil, nil, errors.Wrap(errors.ErrInvalidBackend, errors.New("invalid backend"))
	}

	if config == nil {
		return nil, nil, nil, errors.Wrap(errors.ErrConfigFieldNotFound, errors.New("backend must not be nil"))
	}

	configSvc = &sinks.Configuration{
		Exporter: backend.GetBackend(backendName),
	}
	exporter = config.GetSubMetadata("exporter")
	if exporter == nil {
		return nil, nil, nil, errors.Wrap(errors.ErrExporterFieldNotFound, errors.New("exporter field must not be nil"))
	}
	err = configSvc.Exporter.ValidateConfiguration(exporter)
	if err != nil {
		return
	}

	authentication = config.GetSubMetadata(authentication_type.AuthenticationKey)
	if authentication == nil {
		return nil, nil, nil, errors.Wrap(errors.ErrAuthFieldNotFound, errors.New("authentication field must not be nil"))
	}
	authtype, ok := authentication["type"]
	if !ok {
		authtype = basicauth.AuthType
	}
	switch authtype.(type) {
	case string:
		break
	default:
		return
	}
	authTypeSvc, ok := authentication_type.GetAuthType(authtype.(string))
	if !ok {
		err = errors.Wrap(errors.ErrAuthInvalidType, errors.New("invalid required field authentication type"))
		return
	}
	configSvc.Authentication = authTypeSvc
	err = configSvc.Authentication.ValidateConfiguration("object", authentication)
	return
}

func GetConfigurationAndMetadataFromYaml(backendName string, config string) (configSvc *sinks.Configuration, exporter types.Metadata, authentication types.Metadata, err error) {

	if !backend.HaveBackend(backendName) {
		return nil, nil, nil, errors.Wrap(errors.ErrInvalidBackend, errors.New("invalid backend"))
	}

	configSvc = &sinks.Configuration{
		Exporter: backend.GetBackend(backendName),
	}
	var configStr types.Metadata
	err = yaml.Unmarshal([]byte(config), &configStr)
	if err != nil {
		return
	}
	exporter = configStr.GetSubMetadata("exporter")
	if exporter == nil {
		return nil, nil, nil, errors.New("malformed entity specification. exporter field is expected on configuration field")
	}
	err = configSvc.Exporter.ValidateConfiguration(exporter)
	if err != nil {
		return
	}

	authentication = configStr.GetSubMetadata(authentication_type.AuthenticationKey)
	if authentication == nil {
		return nil, nil, nil, errors.New("malformed entity specification. authentication fields are expected on configuration field")
	}
	authtype, ok := authentication["type"]
	if !ok {
		authtype = basicauth.AuthType
	}
	switch authtype.(type) {
	case string:
		break
	default:
		err = errors.ErrAuthInvalidType
		return
	}
	authTypeSvc, ok := authentication_type.GetAuthType(authtype.(string))
	if !ok {
		err = errors.Wrap(errors.ErrAuthInvalidType, errors.New("invalid required field authentication type"))
		return
	}
	configSvc.Authentication = authTypeSvc
	err = configSvc.Authentication.ValidateConfiguration("object", authentication)
	return
}

func (req addReq) validate() (err error) {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.Name == "" {
		return errors.Wrap(errors.ErrEntityNameNotFound, errors.New("name not found"))
	}

	_, err = types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrConflict, errors.New("identifier duplicated"))
	}
	return nil
}

type updateSinkReq struct {
	Name        string         `json:"name,omitempty"`
	Config      types.Metadata `json:"config,omitempty"`
	Backend     string         `json:"backend,omitempty"`
	Format      string         `json:"format,omitempty"`
	ConfigData  string         `json:"config_data,omitempty"`
	Description *string        `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags,omitempty"`
	id          string
	token       string
}

func (req updateSinkReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return errors.ErrMalformedEntity
	}

	if req.Description == nil && req.Name == "" && req.ConfigData == "" && len(req.Config) == 0 && req.Tags == nil {
		return errors.ErrMalformedEntity
	}

	return nil
}

type viewResourceReq struct {
	token string
	id    string
}

func (req viewResourceReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.id == "" {
		return errors.ErrMalformedEntity
	}
	return nil
}

type listResourcesReq struct {
	token        string
	pageMetadata sinks.PageMetadata
}

func (req *listResourcesReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.pageMetadata.Limit == 0 {
		req.pageMetadata.Limit = defLimit
	}

	if req.pageMetadata.Limit > maxLimitSize {
		return errors.ErrMalformedEntity
	}

	if len(req.pageMetadata.Name) > maxNameSize {
		return errors.ErrMalformedEntity
	}

	if req.pageMetadata.Order != "" &&
		req.pageMetadata.Order != nameOrder && req.pageMetadata.Order != idOrder {
		return errors.ErrMalformedEntity
	}

	if req.pageMetadata.Dir != "" &&
		req.pageMetadata.Dir != ascDir && req.pageMetadata.Dir != descDir {
		return errors.ErrMalformedEntity
	}

	return nil
}

type listBackendsReq struct {
	token string
}

func (req *listBackendsReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	return nil
}

type listAuthTypesReq struct {
	token string
}

func (req *listAuthTypesReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	return nil
}

type deleteSinkReq struct {
	token string
	id    string
}

func (req deleteSinkReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return errors.ErrMalformedEntity
	}

	return nil
}

type validateReq struct {
	Name        string         `json:"name,omitempty"`
	Backend     string         `json:"backend,omitempty"`
	Config      types.Metadata `json:"config,omitempty"`
	Description string         `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags,omitempty"`
	token       string
}

func (req validateReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.Name == "" {
		return errors.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/opentracing/opentracing-go"
	"github.com/orb-community/orb/buildinfo"
	"github.com/orb-community/orb/fleet"
	"github.com/orb-community/orb/fleet/backend"
	"github.com/orb-community/orb/internal/httputil"
	"github.com/orb-community/orb/pkg/db"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType = "application/json"
	offsetKey   = "offset"
	limitKey    = "limit"
	nameKey     = "name"
	orderKey    = "order"
	dirKey      = "dir"
	metadataKey = "metadata"
	tagsKey     = "tags"
	defOffset   = 0
	defLimit    = 10
)

func MakeHandler(tracer opentracing.Tracer, svcName string, svc fleet.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/agent_groups", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_agent_group")(addAgentGroupEndpoint(svc)),
		decodeAddAgentGroup,
		types.EncodeResponse,
		opts...))
	r.Get("/agent_groups", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_agent_group")(listAgentGroupsEndpoint(svc)),
		decodeList,
		types.EncodeResponse,
		opts...))
	r.Get("/agent_groups/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_group")(viewAgentGroupEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Put("/agent_groups/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "edit_agent_group")(editAgentGroupEndpoint(svc)),
		decodeAgentGroupUpdate,
		types.EncodeResponse,
		opts...))
	r.Delete("/agent_groups/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_agent_group")(removeAgentGroupEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Post("/agent_groups/validate", kithttp.NewServer(
		kitot.TraceServer(tracer, "validate_agent_group")(validateAgentGroupEndpoint(svc)),
		decodeValidateAgentGroup,
		types.EncodeResponse,
		opts...))

	r.Post("/agents", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_agent")(addAgentEndpoint(svc)),
		decodeAddAgent,
		types.EncodeResponse,
		opts...))
	r.Get("/agents", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_agents")(listAgentsEndpoint(svc)),
		decodeList,
		types.EncodeResponse,
		opts...))
	r.Get("/agents/backends", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_backends")(listAgentBackendsEndpoint(svc)),
		decodeListBackends,
		types.EncodeResponse,
		opts...))
	r.Post("/agents/:id/rpc/reset", kithttp.NewServer(
		kitot.TraceServer(tracer, "reset_agent")(resetAgentEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Get("/agents/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "edit_agent")(viewAgentEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Get("/agents/:id/matching_groups", kithttp.NewServer(
		kitot.TraceServer(tracer, "matching_groups")(viewAgentMatchingGroups(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Put("/agents/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "edit_agent")(editAgentEndpoint(svc)),
		decodeAgentUpdate,
		types.EncodeResponse,
		opts...))
	r.Post("/agents/validate", kithttp.NewServer(
		kitot.TraceServer(tracer, "validate_agent")(validateAgentEndpoint(svc)),
		decodeAddAgent,
		types.EncodeResponse,
		opts...))
	r.Delete("/agents/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_agent")(removeAgentEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))

	bks := backend.GetList()
	if len(bks) > 0 {
		for _, v := range bks {
			backend.GetBackend(v).MakeHandler(tracer, opts, r)
		}
	}

	r.GetFunc("/version", buildinfo.Version(svcName))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeAddAgentGroup(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}

	req := addAgentGroupReq{token: parseJwt(r)}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeView(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewResourceReq{
		token: parseJwt(r),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeAgentGroupUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}

	req := updateAgentGroupReq{
		token: parseJwt(r),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeAddAgent(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}

	req := addAgentReq{token: parseJwt(r)}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeAgentUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}
	req := updateAgentReq{
		token: parseJwt(r),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeListBackends(_ context.Context, r *http.Request) (interface{}, error) {
	req := listAgentBackendsReq{token: parseJwt(r)}
	return req, nil
}

func decodeList(_ context.Context, r *http.Request) (interface{}, error) {
	o, err := httputil.ReadUintQuery(r, offsetKey, defOffset)
	if err != nil {
		return nil, err
	}

	l, err := httputil.ReadUintQuery(r, limitKey, defLimit)
	if err != nil {
		return nil, err
	}

	n, err := httputil.ReadStringQuery(r, nameKey, "")
	if err != nil {
		return nil, err
	}

	or, err := httputil.ReadStringQuery(r, orderKey, "")
	if err != nil {
		return nil, err
	}

	d, err := httputil.ReadStringQuery(r, dirKey, "")
	if err != nil {
		return nil, err
	}

	m, err := httputil.ReadMetadataQuery(r, metadataKey, nil)
	if err != nil {
		return nil, err
	}

	t, err := httputil.ReadTagQuery(r, tagsKey, nil)
	if err != nil {
		return nil, err
	}

	req := listResourcesReq{
		token: parseJwt(r),
		pageMetadata: fleet.PageMetadata{
			Offset:   o,
			Limit:    l,
			Name:     n,
			Order:    or,
			Dir:      d,
			Metadata: m,
			Tags:     t,
		},
	}

	return req, nil
}

func decodeValidateAgentGroup(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}

	req := addAgentGroupReq{token: parseJwt(r)}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errorVal := err.(type) {
	case errors.Error:
		w.Header().Set("Content-Type", types.ContentType)
		switch {
		case errors.Contains(errorVal, errors.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusUnauthorized)

		case errors.Contains(errorVal, errors.ErrInvalidQueryParams):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, errors.ErrUnsupportedContentType):
			w.WriteHeader(http.StatusUnsupportedMediaType)

		case errors.Contains(errorVal, errors.ErrMalformedEntity):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, errors.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Contains(errorVal, errors.ErrConflict):
			w.WriteHeader(http.StatusConflict)

		case errors.Contains(errorVal, db.ErrScanMetadata):
			w.WriteHeader(http.StatusUnprocessableEntity)

		case errors.Contains(errorVal, fleet.ErrCreateAgentGroup):
			w.WriteHeader(http.StatusBadRequest)

		case errors.Contains(errorVal, io.ErrUnexpectedEOF),
			errors.Contains(errorVal, io.EOF):
			w.WriteHeader(http.StatusBadRequest)

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if errorVal.Msg() != "" {
			if err := json.NewEncoder(w).Encode(types.ErrorRes{Err: errorVal.Msg()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseJwt(r *http.Request) (token string) {
	if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		token = r.Header.Get("Authorization")[7:]
	}
	return
}

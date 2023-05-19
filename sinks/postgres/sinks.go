// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/orb-community/orb/pkg/db"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"go.uber.org/zap"
	"strings"
	"time"
)

var _ sinks.SinkRepository = (*sinksRepository)(nil)

type sinksRepository struct {
	db     Database
	logger *zap.Logger
}

func (s sinksRepository) UpdateVersion(ctx context.Context, incomingVersion string) error {
	q := `UPDATE current_version SET version = :version, last_updated = :currenttime`
	params := map[string]interface{}{
		"version":     incomingVersion,
		"currenttime": time.Now(),
	}
	res, err := s.db.NamedExecContext(ctx, q, params)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(sinks.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return errors.Wrap(errors.ErrConflict, err)
			}
		}
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	if count == 0 {
		return sinks.ErrNotFound
	}
	return nil
}

func (s sinksRepository) GetVersion(ctx context.Context) (string, error) {
	q := `SELECT version FROM current_version`
	params := map[string]interface{}{}
	rows, err := s.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		version := ""
		err := rows.Scan(&version)
		if err != nil {
			return "", err
		}
		return version, nil
	}
	return "", err
}

func (s sinksRepository) SearchAllSinks(ctx context.Context, filter sinks.Filter) ([]sinks.Sink, error) {
	q := `SELECT id, name, mf_owner_id, description, tags, state, coalesce(error, '') as error, backend, metadata, ts_created FROM sinks`
	params := map[string]interface{}{}
	if (filter != sinks.Filter{} && filter.StateFilter != "") {
		q += `WHERE state == :state`
		params["state"] = filter.StateFilter
	}

	rows, err := s.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return nil, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			s.logger.Error("error closing rows", zap.Error(err))
		}
	}(rows)

	items := make([]sinks.Sink, 0)
	for rows.Next() {
		dbSink := dbSink{}
		if err := rows.StructScan(&dbSink); err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}

		sink, err := toSink(dbSink)
		if err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}
		// metadataFilters will apply only after Fetching in metadata, due to struct
		filterFunc := func(key string, value interface{}) bool {
			if key == sinks.MetadataLabelOtel {
				if value.(string) == filter.OpenTelemetry {
					return true
				}
			}
			return false
		}
		if sink.Config.IsApplicable(filterFunc) {
			items = append(items, sink)
		}
	}

	return items, err
}

func (s sinksRepository) Save(ctx context.Context, sink sinks.Sink) (string, error) {
	q := `INSERT INTO sinks (name, mf_owner_id, metadata, config_data, format, description, backend, tags, state, error)         
			  VALUES (:name, :mf_owner_id, :metadata, :config_data, :format, :description, :backend, :tags, :state, :error) RETURNING id`

	if !sink.Name.IsValid() || sink.MFOwnerID == "" {
		return "", errors.ErrMalformedEntity
	}

	dba, err := toDBSink(sink)
	if err != nil {
		return "", errors.Wrap(db.ErrSaveDB, err)
	}

	row, err := s.db.NamedQueryContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return "", errors.Wrap(errors.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return "", errors.Wrap(errors.ErrConflict, err)
			}
		}
		return "", errors.Wrap(db.ErrSaveDB, err)
	}

	defer row.Close()
	row.Next()
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil

}

func (s sinksRepository) Update(ctx context.Context, sink sinks.Sink) error {
	q := `UPDATE sinks 
			SET description = :description, 
			    tags = :tags, 
			    metadata = :metadata,  
			    config_data = :config_data, 
			    format = :format, 
			    name = :name 
			WHERE mf_owner_id = :mf_owner_id 
			  AND id = :id;`

	sinkDB, err := toDBSink(sink)
	if err != nil {
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	res, err := s.db.NamedExecContext(ctx, q, sinkDB)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(sinks.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return errors.Wrap(errors.ErrConflict, err)
			}
		}
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	if count == 0 {
		return sinks.ErrNotFound
	}
	return nil
}

func (s sinksRepository) RetrieveAllByOwnerID(ctx context.Context, owner string, pm sinks.PageMetadata) (sinks.Page, error) {
	name, nameQuery := getNameQuery(pm.Name)
	orderQuery := getOrderQuery(pm.Order)
	dirQuery := getDirQuery(pm.Dir)
	metadata, metadataQuery, err := getMetadataQuery(pm.Metadata)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	tags, tagsQuery, err := getTagsQuery(pm.Tags)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(`SELECT id, name, mf_owner_id, description, tags, state, coalesce(error, '') as error, backend, metadata, ts_created
								FROM sinks 
								WHERE mf_owner_id = :mf_owner_id %s%s%s 
								ORDER BY %s %s LIMIT :limit OFFSET :offset;`,
		tagsQuery, metadataQuery, nameQuery, orderQuery, dirQuery)
	params := map[string]interface{}{
		"mf_owner_id": owner,
		"limit":       pm.Limit,
		"offset":      pm.Offset,
		"name":        name,
		"metadata":    metadata,
		"tags":        tags,
	}
	rows, err := s.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []sinks.Sink
	for rows.Next() {
		d := dbSink{MFOwnerID: owner}
		if err := rows.StructScan(&d); err != nil {
			return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
		}

		sink, err := toSink(d)
		if err != nil {
			return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
		}

		items = append(items, sink)
	}

	count := fmt.Sprintf(`SELECT COUNT(*) FROM sinks WHERE mf_owner_id = :mf_owner_id %s%s%s`, tagsQuery, metadataQuery, nameQuery)

	total, err := total(ctx, s.db, count, params)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	page := sinks.Page{
		Sinks: items,
		PageMetadata: sinks.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
			Order:  pm.Order,
			Dir:    pm.Dir,
		},
	}

	return page, nil
}

func (s sinksRepository) RetrieveById(ctx context.Context, id string) (sinks.Sink, error) {

	q := `SELECT id, name, mf_owner_id, description, tags, backend, metadata, format, config_data, ts_created, state, coalesce(error, '') as error
			FROM sinks where id = $1`

	dba := dbSink{}

	if err := s.db.QueryRowxContext(ctx, q, id).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
			return sinks.Sink{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return sinks.Sink{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toSink(dba)
}

func (s sinksRepository) RetrieveByOwnerAndId(ctx context.Context, ownerID string, id string) (sinks.Sink, error) {

	q := `SELECT id, name, mf_owner_id, description, tags, backend, metadata, format, config_data, ts_created, state, coalesce(error, '') as error
			FROM sinks where id = $1 and mf_owner_id = $2`

	if ownerID == "" || id == "" {
		return sinks.Sink{}, errors.ErrSelectEntity
	}

	dba := dbSink{}

	if err := s.db.QueryRowxContext(ctx, q, id, ownerID).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
			return sinks.Sink{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return sinks.Sink{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toSink(dba)
}

func (s sinksRepository) Remove(ctx context.Context, owner, id string) error {
	dbsk := dbSink{
		ID:        id,
		MFOwnerID: owner,
	}

	q := `DELETE FROM sinks WHERE id = :id AND mf_owner_id = :mf_owner_id;`
	if _, err := s.db.NamedExecContext(ctx, q, dbsk); err != nil {
		return errors.Wrap(sinks.ErrRemoveEntity, err)
	}

	return nil
}

func (s sinksRepository) UpdateSinkState(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) error {
	dbsk := dbSink{
		ID:        sinkID,
		MFOwnerID: ownerID,
		State:     state,
		Error:     msg,
	}

	q := "update sinks set state = :state, error = :error where mf_owner_id = :mf_owner_id and id = :id"

	res, err := s.db.NamedExecContext(ctx, q, dbsk)
	if err != nil {
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	if count == 0 {
		return sinks.ErrUpdateEntity
	}

	return nil
}

type dbSink struct {
	ID          string           `db:"id"`
	Name        types.Identifier `db:"name"`
	MFOwnerID   string           `db:"mf_owner_id"`
	Metadata    db.Metadata      `db:"metadata"`
	ConfigData  string           `db:"config_data"`
	Format      string           `db:"format"`
	Backend     string           `db:"backend"`
	Description string           `db:"description"`
	Created     time.Time        `db:"ts_created"`
	Tags        db.Tags          `db:"tags"`
	State       sinks.State      `db:"state"`
	Error       string           `db:"error"`
}

func toDBSink(sink sinks.Sink) (dbSink, error) {

	var uID uuid.UUID
	err := uID.Scan(sink.MFOwnerID)
	if err != nil {
		return dbSink{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	var description string
	if sink.Description == nil {
		description = ""
	} else {
		description = *sink.Description
	}

	return dbSink{
		ID:          sink.ID,
		Name:        sink.Name,
		MFOwnerID:   uID.String(),
		Metadata:    db.Metadata(sink.Config),
		ConfigData:  sink.ConfigData,
		Format:      sink.Format,
		Backend:     sink.Backend,
		Description: description,
		Created:     sink.Created,
		Tags:        db.Tags(sink.Tags),
		State:       sink.State,
		Error:       sink.Error,
	}, nil

}

func toSink(dba dbSink) (sinks.Sink, error) {
	sink := sinks.Sink{
		ID:          dba.ID,
		Name:        dba.Name,
		MFOwnerID:   dba.MFOwnerID,
		Backend:     dba.Backend,
		Description: &dba.Description,
		State:       dba.State,
		Error:       dba.Error,
		Config:      types.Metadata(dba.Metadata),
		Created:     dba.Created,
		Tags:        types.Tags(dba.Tags),
	}
	return sink, nil
}

func getNameQuery(name string) (string, string) {
	if name == "" {
		return "", ""
	}
	name = fmt.Sprintf(`%%%s%%`, strings.ToLower(name))
	nameQuey := ` AND LOWER(name) LIKE :name`
	return name, nameQuey
}

func getOrderQuery(order string) string {
	switch order {
	case "name":
		return "name"
	default:
		return "id"
	}
}

func getDirQuery(dir string) string {
	switch dir {
	case "asc":
		return "ASC"
	default:
		return "DESC"
	}
}

func getMetadataQuery(m types.Metadata) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		mq = ` AND metadata @> :metadata`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
}

func getTagsQuery(m types.Tags) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		// todo add in orb tags
		mq = ` AND tags @> :tags`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
}

func total(ctx context.Context, db Database, query string, params interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func NewSinksRepository(db Database, logger *zap.Logger) sinks.SinkRepository {
	return &sinksRepository{db: db, logger: logger}
}

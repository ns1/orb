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
	"github.com/lib/pq"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"strings"
	"time"
)

var _ fleet.AgentRepository = (*agentRepository)(nil)

type agentRepository struct {
	db     Database
	logger *zap.Logger
}

func (r agentRepository) RetrieveMatchingAgents(ctx context.Context, ownerID string, tags types.Tags) (types.Metadata, error) {
	t, tmq, err := getOrbOrAgentTagsQuery(tags)
	if err != nil {
		return types.Metadata{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(
		`select
			json_build_object('total', coalesce(total,0), 'online', coalesce(online,0)) AS matching_agents
		from
			(select
				sum(case when mf_thing_id is not null then 1 else 0 end) as total,
				sum(case when state = 'online' then 1 else 0 end) as online
			from agents WHERE mf_owner_id = :mf_owner_id %s) as agent_groups`, tmq)

	params := map[string]interface{}{
		"tags":        t,
		"mf_owner_id": ownerID,
	}

	rows, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return types.Metadata{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	dbma := dbMatchingAgent{}
	if rows.Next() {
		if err := rows.StructScan(&dbma); err != nil {
			return types.Metadata{}, errors.Wrap(errors.ErrSelectEntity, err)
		}
	}

	return types.Metadata(dbma.MatchingAgents), nil
}

func (r agentRepository) RetrieveAllByAgentGroupID(ctx context.Context, owner string, agentGroupID string, onlinishOnly bool) ([]fleet.Agent, error) {

	q := `SELECT agent_mf_thing_id AS mf_thing_id, agent_mf_channel_id AS mf_channel_id FROM agent_group_membership 
			WHERE mf_owner_id = :mf_owner_id AND agent_groups_id = :group_id`

	if onlinishOnly {
		q = q + ` AND (agent_state = 'online' OR agent_state = 'stale')`
	}

	if agentGroupID == "" || owner == "" {
		return nil, errors.ErrMalformedEntity
	}

	params := map[string]interface{}{
		"mf_owner_id": owner,
		"group_id":    agentGroupID,
	}

	rows, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return nil, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []fleet.Agent
	for rows.Next() {
		dbth := dbAgent{MFOwnerID: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}

		th, err := toAgent(dbth)
		if err != nil {
			return nil, errors.Wrap(errors.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	return items, nil
}

func (r agentRepository) RetrieveAll(ctx context.Context, owner string, pm fleet.PageMetadata) (fleet.Page, error) {
	nq, name := getNameQuery(pm.Name)
	oq := getOrderQuery(pm.Order)
	dq := getDirQuery(pm.Dir)
	m, mq, err := getMetadataQuery(pm.Metadata)
	if err != nil {
		return fleet.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	t, tmq, err := getOrbOrAgentTagsQuery(pm.Tags)
	if err != nil {
		return fleet.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(`SELECT mf_thing_id, name, mf_owner_id, mf_channel_id, ts_created, orb_tags, agent_tags, agent_metadata, state, last_hb_data, ts_last_hb
			FROM agents WHERE mf_owner_id = :mf_owner_id %s%s%s ORDER BY %s %s LIMIT :limit OFFSET :offset;`, tmq, mq, nq, oq, dq)
	params := map[string]interface{}{
		"mf_owner_id": owner,
		"limit":       pm.Limit,
		"offset":      pm.Offset,
		"name":        name,
		"metadata":    m,
		"tags":        t,
	}

	rows, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return fleet.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []fleet.Agent
	for rows.Next() {
		dbth := dbAgent{MFOwnerID: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return fleet.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
		}

		th, err := toAgent(dbth)
		if err != nil {
			return fleet.Page{}, errors.Wrap(errors.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM agents WHERE mf_owner_id = :mf_owner_id %s%s%s;`, nq, tmq, mq)

	total, err := total(ctx, r.db, cq, params)
	if err != nil {
		return fleet.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	page := fleet.Page{
		Agents: items,
		PageMetadata: fleet.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
			Order:  pm.Order,
			Dir:    pm.Dir,
		},
	}

	return page, nil
}

func (r agentRepository) UpdateDataByIDWithChannel(ctx context.Context, agent fleet.Agent) error {

	q := `UPDATE agents SET (orb_tags, agent_tags, agent_metadata)         
			= (:orb_tags, :agent_tags, :agent_metadata) 
			WHERE mf_thing_id = :mf_thing_id AND mf_channel_id = :mf_channel_id;`

	if agent.MFThingID == "" || agent.MFChannelID == "" {
		return errors.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
	if err != nil {
		return errors.Wrap(errors.ErrUpdateEntity, err)
	}

	res, err := r.db.NamedExecContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(errors.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return errors.Wrap(errors.ErrConflict, err)
			}
		}
		return errors.Wrap(db.ErrUpdateDB, err)
	}

	cnt, errdb := res.RowsAffected()
	if errdb != nil {
		return errors.Wrap(errors.ErrUpdateEntity, errdb)
	}

	if cnt == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r agentRepository) UpdateHeartbeatByIDWithChannel(ctx context.Context, agent fleet.Agent) error {

	q := `UPDATE agents SET (last_hb_data, ts_last_hb, state)         
			= (:last_hb_data, now(), :state) 
			WHERE mf_thing_id = :mf_thing_id AND mf_channel_id = :mf_channel_id;`

	if agent.MFThingID == "" || agent.MFChannelID == "" {
		return errors.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
	if err != nil {
		return errors.Wrap(errors.ErrUpdateEntity, err)
	}
	res, err := r.db.NamedExecContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(errors.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return errors.Wrap(errors.ErrConflict, err)
			}
		}
		return errors.Wrap(db.ErrUpdateDB, err)
	}

	cnt, errdb := res.RowsAffected()
	if errdb != nil {
		return errors.Wrap(errors.ErrUpdateEntity, errdb)
	}

	if cnt == 0 {
		return errors.ErrNotFound
	}

	return nil

}

func (r agentRepository) RetrieveByIDWithChannel(ctx context.Context, thingID string, channelID string) (fleet.Agent, error) {

	q := `SELECT mf_thing_id, name, mf_owner_id, mf_channel_id, ts_created, orb_tags, agent_tags, agent_metadata, state, last_hb_data, ts_last_hb FROM agents WHERE mf_thing_id = $1 AND mf_channel_id = $2;`

	dba := dbAgent{}

	if err := r.db.QueryRowxContext(ctx, q, thingID, channelID).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
			return fleet.Agent{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return fleet.Agent{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toAgent(dba)
}

func (r agentRepository) Save(ctx context.Context, agent fleet.Agent) error {

	q := `INSERT INTO agents (name, mf_thing_id, mf_owner_id, mf_channel_id, orb_tags, agent_tags, agent_metadata, state)         
			  VALUES (:name, :mf_thing_id, :mf_owner_id, :mf_channel_id, :orb_tags, :agent_tags, :agent_metadata, :state)`

	if !agent.Name.IsValid() || agent.MFOwnerID == "" || agent.MFThingID == "" || agent.MFChannelID == "" {
		return errors.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
	if err != nil {
		return errors.Wrap(db.ErrSaveDB, err)
	}

	_, err = r.db.NamedExecContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(errors.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return errors.Wrap(errors.ErrConflict, err)
			}
		}
		return errors.Wrap(db.ErrSaveDB, err)
	}

	return nil

}

func (r agentRepository) UpdateAgentByID(ctx context.Context, ownerID string, agent fleet.Agent) error {
	q := `UPDATE agents SET (name, orb_tags)         
			= (:name, :orb_tags) 
			WHERE mf_thing_id = :mf_thing_id AND mf_owner_id = :mf_owner_id;`

	agent.MFOwnerID = ownerID
	if agent.MFThingID == "" || agent.MFOwnerID == "" {
		return errors.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
	if err != nil {
		return errors.Wrap(errors.ErrUpdateEntity, err)
	}

	res, err := r.db.NamedExecContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(errors.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return errors.Wrap(errors.ErrConflict, err)
			}
		}
		return errors.Wrap(db.ErrUpdateDB, err)
	}

	cnt, errdb := res.RowsAffected()
	if errdb != nil {
		return errors.Wrap(errors.ErrUpdateEntity, errdb)
	}

	if cnt == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r agentRepository) RetrieveByID(ctx context.Context, ownerID string, thingID string) (fleet.Agent, error) {
	q := `SELECT 
			mf_thing_id, 
			name, 
			mf_owner_id, 
			mf_channel_id, 
			ts_created, 
			orb_tags, 
			agent_tags, 
			agent_metadata, 
			state, 
			last_hb_data, 
			ts_last_hb 
		FROM agents 
		WHERE 
			mf_thing_id = $1 
			AND mf_owner_id = $2;`

	dba := dbAgent{}

	if err := r.db.QueryRowxContext(ctx, q, thingID, ownerID).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
			return fleet.Agent{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return fleet.Agent{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toAgent(dba)
}

func (r agentRepository) Delete(ctx context.Context, ownerID string, thingID string) error {
	params := map[string]interface{}{
		"id":    thingID,
		"owner": ownerID,
	}

	q := `DELETE FROM agents WHERE mf_thing_id = :id AND mf_owner_id = :owner;`

	if _, err := r.db.NamedQueryContext(ctx, q, params); err != nil {
		return errors.Wrap(fleet.ErrRemoveEntity, err)
	}

	return nil
}

type dbAgent struct {
	Name          types.Identifier `db:"name"`
	MFOwnerID     string           `db:"mf_owner_id"`
	MFThingID     string           `db:"mf_thing_id"`
	MFChannelID   string           `db:"mf_channel_id"`
	OrbTags       db.Tags          `db:"orb_tags"`
	AgentTags     db.Tags          `db:"agent_tags"`
	AgentMetadata db.Metadata      `db:"agent_metadata"`
	State         fleet.State      `db:"state"`
	Created       time.Time        `db:"ts_created"`
	LastHBData    db.Metadata      `db:"last_hb_data"`
	LastHB        sql.NullTime     `db:"ts_last_hb"`
}

type dbMatchingAgent struct {
	MatchingAgents db.Metadata `db:"matching_agents"`
}

func toDBAgent(agent fleet.Agent) (dbAgent, error) {

	a := dbAgent{
		Name:          agent.Name,
		MFOwnerID:     agent.MFOwnerID,
		MFThingID:     agent.MFThingID,
		MFChannelID:   agent.MFChannelID,
		OrbTags:       db.Tags(agent.OrbTags),
		AgentTags:     db.Tags(agent.AgentTags),
		AgentMetadata: db.Metadata(agent.AgentMetadata),
		State:         agent.State,
		Created:       agent.Created,
		LastHBData:    db.Metadata(agent.LastHBData),
	}

	if !agent.LastHB.IsZero() {
		a.LastHB = sql.NullTime{
			Time:  agent.LastHB,
			Valid: true,
		}
	}
	return a, nil
}

func toAgent(dba dbAgent) (fleet.Agent, error) {

	agent := fleet.Agent{
		Name:          dba.Name,
		MFOwnerID:     dba.MFOwnerID,
		MFThingID:     dba.MFThingID,
		MFChannelID:   dba.MFChannelID,
		Created:       dba.Created,
		OrbTags:       types.Tags(dba.OrbTags),
		AgentTags:     types.Tags(dba.AgentTags),
		AgentMetadata: types.Metadata(dba.AgentMetadata),
		State:         dba.State,
		LastHBData:    types.Metadata(dba.LastHBData),
		LastHB:        dba.LastHB.Time,
	}

	return agent, nil

}
func getNameQuery(name string) (string, string) {
	if name == "" {
		return "", ""
	}
	name = fmt.Sprintf(`%%%s%%`, strings.ToLower(name))
	nq := ` AND LOWER(name) LIKE :name`
	return nq, name
}

func getOrderQuery(order string) string {
	switch order {
	case "name":
		return "name"
	default:
		return "mf_thing_id"
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
		mq = ` AND agent_metadata @> :metadata`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
}

func getOrbOrAgentTagsQuery(m types.Tags) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		// todo add in orb tags
		mq = ` AND (agent_tags @> :tags OR orb_tags @> :tags)`

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

func NewAgentRepository(db Database, logger *zap.Logger) fleet.AgentRepository {
	return &agentRepository{db: db, logger: logger}
}

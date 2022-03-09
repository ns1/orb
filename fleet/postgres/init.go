// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"fmt"
	"github.com/ns1labs/orb/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required for SQL access
	migrate "github.com/rubenv/sql-migrate"
)

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg config.PostgresConfig) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.DB, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "fleet_1",
				Up: []string{
					`CREATE TYPE agent_state AS ENUM ('new', 'online', 'offline', 'stale', 'removed', 'upgrade_required');`,
					`CREATE TABLE IF NOT EXISTS agents (

						mf_thing_id        UUID NOT NULL,
						name        	   TEXT NOT NULL,
						mf_owner_id        UUID NOT NULL,

						mf_channel_id      UUID NOT NULL,

                        ts_created         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

						orb_tags           JSONB NOT NULL DEFAULT '{}',

						agent_tags         JSONB NOT NULL DEFAULT '{}',
						agent_metadata     JSONB NOT NULL DEFAULT '{}',

						state              agent_state NOT NULL DEFAULT 'new',

						last_hb_data       JSONB NOT NULL DEFAULT '{}',
                        ts_last_hb         TIMESTAMPTZ DEFAULT NULL,
						PRIMARY KEY(mf_thing_id),
						UNIQUE(name, mf_owner_id)
					)`,
					`CREATE INDEX ON agents (mf_owner_id)`,
					`CREATE INDEX ON agents USING gin (orb_tags)`,
					`CREATE INDEX ON agents USING gin (agent_tags)`,
					`CREATE TABLE IF NOT EXISTS agent_groups (
						id			       UUID NOT NULL DEFAULT gen_random_uuid(),
						name        	   TEXT NOT NULL,
						description        TEXT NOT NULL,
						mf_owner_id        UUID NOT NULL,

						mf_channel_id      UUID NOT NULL,
	
						tags			   JSONB NOT NULL DEFAULT '{}',
                        ts_created         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
						PRIMARY KEY (name, mf_owner_id),
					    UNIQUE(id)
					)`,
					`CREATE INDEX ON agent_groups (mf_owner_id)`,
					`CREATE INDEX ON agent_groups USING gin (tags)`,
					`CREATE VIEW agent_group_membership(agent_groups_id, agent_groups_name, agent_mf_thing_id, agent_mf_channel_id, group_mf_channel_id, mf_owner_id, agent_state) as
					SELECT agent_groups.id,
						   agent_groups.name,
						   agents.mf_thing_id,
						   agents.mf_channel_id,
						   agent_groups.mf_channel_id,
						   agent_groups.mf_owner_id,
						   agents.state
					FROM agents,
						 agent_groups
					WHERE agent_groups.mf_owner_id = agents.mf_owner_id
					  AND (agent_groups.tags <@ agents.agent_tags OR agent_groups.tags <@ agents.orb_tags)`,
				},
				Down: []string{
					"DROP TABLE agents",
					"DROP TABLE agent_groups",
					"DROP VIEW agent_group_membership",
				},
			}, {
				Id: "fleet_2",
				Up: []string{
					`CREATE or REPLACE VIEW agent_group_membership(agent_groups_id, agent_groups_name, agent_mf_thing_id, agent_mf_channel_id, group_mf_channel_id, mf_owner_id, agent_state) as
					SELECT agent_groups.id,
						   agent_groups.name,
						   agents.mf_thing_id,
						   agents.mf_channel_id,
						   agent_groups.mf_channel_id,
						   agent_groups.mf_owner_id,
						   agents.state
					FROM agents,
						 agent_groups
					WHERE agent_groups.mf_owner_id = agents.mf_owner_id
					  AND (agent_groups.tags <@ coalesce(agents.agent_tags || agents.orb_tags, agents.agent_tags, agents.orb_tags))`,
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)

	return err
}

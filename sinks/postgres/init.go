// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"fmt"
	"github.com/orb-community/orb/pkg/config"

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
				Id: "sinks_1",
				Up: []string{
					`CREATE TYPE sinks_state as enum ('unknown', 'active', 'error', 'idle');`,
					`CREATE TABLE IF NOT EXISTS sinks (
						id			   UUID NOT NULL DEFAULT gen_random_uuid(),
						name           TEXT NOT NULL,
						mf_owner_id    UUID NOT NULL,
						description    TEXT NOT NULL,
						tags           JSONB NOT NULL DEFAULT '{}',

						state          sinks_state NOT NULL DEFAULT 'unknown',

						error          TEXT,
						backend        TEXT NOT NULL,
						metadata       JSONB NOT NULL DEFAULT '{}',
						ts_created     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
						PRIMARY KEY (name, mf_owner_id),
						UNIQUE(id)
					)`,
					`CREATE INDEX ON sinks (mf_owner_id)`,
				},
				Down: []string{
					"DROP TABLE sinks",
				},
			},
			{
				Id: "sinks_2",
				Up: []string{
					`ALTER TABLE sinks ADD COLUMN format TEXT ;`,
					`ALTER TABLE sinks ADD COLUMN config_data TEXT ;`,
				},
				Down: []string{
					`ALTER TABLE sinks DROP COLUMN format;`,
					`ALTER TABLE sinks DROP COLUMN config_data;`,
				},
			},
			{
				Id: "sinks_3",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS current_version (
						id			     UUID NOT NULL DEFAULT gen_random_uuid(),
						version          TEXT NOT NULL,
    					last_updated     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
						PRIMARY KEY (id, version)
					)`,
					`INSERT INTO current_version  (id, version, last_updated) VALUES (DEFAULT, '0.25.0', DEFAULT);`,
				},
				Down: []string{
					"DROP TABLE current_version",
				},
			},
			{
				Id: "sinks_4",
				Up: []string{
					`ALTER TYPE public.sinks_state ADD VALUE IF NOT EXISTS 'warning';`,
					`ALTER TYPE public.sinks_state ADD VALUE IF NOT EXISTS 'provisioning';`,
					`ALTER TYPE public.sinks_state ADD VALUE IF NOT EXISTS 'provisioning_error';`,
				},
				Down: []string{
					`ALTER TYPE public.sinks_state DROP VALUE IF EXISTS 'warning';`,
					`ALTER TYPE public.sinks_state DROP VALUE IF EXISTS 'provisioning';`,
					`ALTER TYPE public.sinks_state DROP VALUE IF EXISTS 'provisioning_error';`,
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)

	return err
}

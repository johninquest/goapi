// internal/database/migrations.go
package database

import "database/sql"

const createPoliciesTable = `
CREATE TABLE IF NOT EXISTS policies (
	id TEXT PRIMARY KEY,
	policy_type TEXT NOT NULL,
	policy_number TEXT NOT NULL,
	insurance_provider TEXT NOT NULL,
	policy_comment TEXT,
	start_date DATETIME NOT NULL,
	end_date DATETIME NOT NULL,
	automatic_renewal BOOLEAN NOT NULL,
	created_by TEXT NOT NULL,
	premium REAL NOT NULL,
	payment_frequency INTEGER NOT NULL,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);`

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(createPoliciesTable)
	return err
}

package storage

import (
	"database/sql"
	"fmt"

	"github.com/authelia/authelia/internal/utils"
)

func (p *SQLProvider) upgradeCreateTableStatements(tx *sql.Tx, statements map[string]string, existingTables []string) error {
	for table, statement := range statements {
		if !utils.IsStringInSlice(table, existingTables) {
			_, err := tx.Exec(fmt.Sprintf(statement, table))
			if err != nil {
				return fmt.Errorf("Unable to create table %s: %v", table, err)
			}
		}
	}

	return nil
}

func (p *SQLProvider) upgradeRunMultipleStatements(tx *sql.Tx, statements []string) error {
	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

// upgradeFinalize sets the schema version and logs a message, as well as any other future finalization tasks.
func (p *SQLProvider) upgradeFinalize(tx *sql.Tx, version SchemaVersion) error {
	_, err := tx.Exec(p.sqlConfigSetValue, "schema", "version", version.ToString())
	if err != nil {
		return err
	}

	p.log.Debugf("%s%d", storageSchemaUpgradeMessage, version)

	return nil
}

// upgradeSchemaToVersion001 upgrades the schema to version 1.
func (p *SQLProvider) upgradeSchemaToVersion001(tx *sql.Tx, tables []string) error {
	version := SchemaVersion(1)

	err := p.upgradeCreateTableStatements(tx, p.sqlUpgradesCreateTableStatements[version], tables)
	if err != nil {
		return err
	}

	// Skip mysql create index statements. It doesn't support CREATE INDEX IF NOT EXIST. May be able to work around this with an Index struct.
	if p.name != "mysql" {
		err = p.upgradeRunMultipleStatements(tx, p.sqlUpgradesCreateTableIndexesStatements[1])
		if err != nil {
			return fmt.Errorf("Unable to create index: %v", err)
		}
	}

	err = p.upgradeFinalize(tx, version)
	if err != nil {
		return err
	}

	return nil
}

// upgradeSchemaToVersion002 upgrades the schema to faux version 2.
func (p *SQLProvider) upgradeSchemaToVersion002(tx *sql.Tx, tables []string) error {
	err := p.upgradeFinalize(tx, SchemaVersion(2))
	if err != nil {
		return err
	}

	p.log.Tracef("tables are %v", tables)

	return nil
}

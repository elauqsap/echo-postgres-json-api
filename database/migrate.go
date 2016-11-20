package database

var schema = []string{
	`DROP SCHEMA IF EXISTS app CASCADE`,
	`CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION appbot`,
	`ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO appbot`,
	`CREATE TABLE IF NOT EXISTS app.users (
				id SERIAL PRIMARY KEY,
				login varchar(100) NOT NULL UNIQUE CHECK (login <> ''),
				password bytea NOT NULL,
				jwt text
		)`,
}

// SetSchema migrates the schema over to the database
func (s *Store) SetSchema() error {
	for _, cmd := range schema {
		if _, err := s.Exec(cmd); err != nil {
			return err
		}
	}
	return nil
}

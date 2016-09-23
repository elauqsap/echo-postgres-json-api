package database

const (
	// enums for peforming migrations in an order
	mSCHEMA = iota
	mTYPE
	mTABLE
	mINDEX
)

var enums = []int{mSCHEMA, mTYPE, mTABLE, mINDEX}

// migrations to perform
var migrations = map[int]map[int]string{
	// schema commands in the order to be performed
	mSCHEMA: map[int]string{
		0: `DROP SCHEMA IF EXISTS app CASCADE`,
		1: `CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION appbot`,
	},
	// type commands in the order to be performed
	mTYPE: map[int]string{
		0: `CREATE TYPE role AS ENUM ('user','manager','admin')`,
	},
	// table commands in the order to be performed
	mTABLE: map[int]string{
		0: `CREATE TABLE IF NOT EXISTS app.users (
					id SERIAL PRIMARY KEY,
					first varchar(100) NOT NULL,
					last varchar(100) NOT NULL,
					role role NOT NULL DEFAULT user,
			)`,
	},
	// index commands in the order to be performed
	mINDEX: map[int]string{
		0: `CREATE INDEX role_idx ON app.users (role)`,
	},
}

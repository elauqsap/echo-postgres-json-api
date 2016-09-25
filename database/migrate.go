package database

import (
	"fmt"
	"sort"
)

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
		0: `CREATE TYPE app.roles AS ENUM ('user','manager','admin')`,
	},
	// table commands in the order to be performed
	mTABLE: map[int]string{
		0: `CREATE TABLE IF NOT EXISTS app.users (
					id SERIAL PRIMARY KEY,
					first varchar(100) NOT NULL,
					last varchar(100) NOT NULL,
					role app.roles NOT NULL DEFAULT 'user',
					api_key char(32) NOT NULL UNIQUE
			)`,
	},
	// index commands in the order to be performed
	mINDEX: map[int]string{
		0: `CREATE INDEX role_idx ON app.users (role)`,
	},
}

// Migrate performs the migrations on the databse
func (s *Store) Migrate(verbose bool) error {
	sort.Ints(enums)
	for key := range enums {
		var order []int
		for index := range migrations[key] {
			order = append(order, index)
		}
		var msgs []string
		sort.Ints(order)
		for _, index := range order {
			if err := s.ExecTransact(migrations[key][index]); err != nil {
				msgs = append(msgs, fmt.Sprintf("[-] %s => %s\n", err.Error(), migrations[key][index]))
			}
		}
		if verbose {
			switch key {
			case mSCHEMA:
				printMsgs("Schema", msgs)
				break
			case mTYPE:
				printMsgs("Type", msgs)
				break
			case mTABLE:
				printMsgs("Table", msgs)
				break
			case mINDEX:
				printMsgs("Index", msgs)
				break
			default:
				break
			}
		}
	}
	return nil
}

// print any error msgs during migration
func printMsgs(key string, msgs []string) {
	fmt.Printf("[+] %s\t migrations completed\n", key)
	if len(msgs) > 0 {
		for _, msg := range msgs {
			fmt.Printf("\t %s\n", msg)
		}
	}
}

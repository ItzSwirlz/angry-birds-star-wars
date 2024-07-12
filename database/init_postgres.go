package database

import "github.com/ItzSwirlz/angry-birds-star-wars/globals"

func initPostgres() {
	var err error

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS rankings (
		owner_pid integer,
		unique_id bigint,
		category integer,
		score integer,
		order_by integer,
		update_mode integer,
		groups bytea,
		param bigint,
		updated_at bigint,
		PRIMARY KEY (unique_id, owner_pid, category)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS common_datas (
		unique_id bigserial,
		owner_pid integer,
		common_data bytea,
		updated_at bigint,
		PRIMARY KEY (unique_id, owner_pid)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres tables created")
}

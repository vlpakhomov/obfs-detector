package model

import "time"

type BlockedIPAddress struct {
	Address string `db:"address"`
	Verdict string `db:"verdict"`
	//rawdata tcp
	CreatedAt time.Time `db:"created_at"`
}

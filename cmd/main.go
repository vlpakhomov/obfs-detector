package main

import (
	"fmt"
	"github.com/florianl/go-nfqueue"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"obfs-detector/config"
	"obfs-detector/internal/app"
	"os"
	"time"
)

func main() {
	cfg := config.Load("./config/config.json")

	zerolog.DurationFieldUnit = time.Millisecond
	logger := zerolog.New(os.Stdout)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		logger.Fatal().Err(err).Msg("connect to postgres fail")
	}
	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			logger.Fatal().Err(err).Msg("close postgres fail")
		}
	}(db)

	if err = db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("ping postgres fail")
	}

	config := nfqueue.Config{
		NfQueue:      100,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 15 * time.Millisecond,
	}

	nfq, err := nfqueue.Open(&config)
	if err != nil {
		logger.Fatal().Err(err).Msg("open nfqueue socket fail")
	}
	defer func(nfq *nfqueue.Nfqueue) {
		if err = nfq.Close(); err != nil {
			logger.Fatal().Err(err).Msg("close connection to the nfqueue socket fail")
		}
	}(nfq)

	obfsDetector := app.New(&logger, db, nfq)
	obfsDetector.Start()
}

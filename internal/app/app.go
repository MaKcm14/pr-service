package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/MaKcm14/pr-service/internal/config/cfg"
	"github.com/MaKcm14/pr-service/internal/controller/chttp"
	"github.com/MaKcm14/pr-service/internal/repo/postgres"
	"github.com/MaKcm14/pr-service/internal/services/usecase"
)

// Service defines the main service's structure with all dependencies in it.
type Service struct {
	log     *slog.Logger
	logFile *os.File
	contr   chttp.HttpController
}

func NewService() Service {
	date := strings.Split(time.Now().String()[:19], " ")

	logFile, err := os.Create(fmt.Sprintf("../../logs/pull-request-service-main-logs_%s___%s.txt",
		date[0], strings.Join(strings.Split(date[1], ":"), "-")))

	if err != nil {
		panic(fmt.Sprintf("error of creating the main-log-file: %v", err))
	}

	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	config := cfg.Config{}
	if err := config.Configure(log,
		cfg.ConfigSocket("SERVICE_SOCKET"),
		cfg.ConfigDBSocket("DB_SOCKET")); err != nil {
		logFile.Close()
		panic(fmt.Sprintf("error while configuring the service: %s", err))
	}

	contr, err := configureLayers(log, config)
	if err != nil {
		logFile.Close()
		panic(fmt.Sprintf("error while configuring the service: %s", err))
	}

	return Service{
		log:     log,
		logFile: logFile,
		contr:   contr,
	}
}

func configureLayers(log *slog.Logger, config cfg.Config) (chttp.HttpController, error) {
	log.Info("configuring the DB")

	repo, err := postgres.New(log, config.DBSocket)
	if err != nil {
		return chttp.HttpController{}, fmt.Errorf("error while configuring the service: %s", err)
	}

	contr := chttp.New(
		log,
		config.Socket,
		usecase.NewUseCase(log, repo, repo, repo),
	)

	return contr, nil
}

func (s *Service) Start() error {
	defer s.log.Info("STOP THE PULL-REQUEST SERVICE")
	defer s.close()

	s.log.Info("start the pull-request service")
	if err := s.contr.Run(); err != nil {
		s.log.Error(fmt.Sprintf("error of starting the pull-request service: %s", err))
		return err
	}

	return nil
}

func (s *Service) close() {
	s.logFile.Close()
}

package cfg

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
)

// Config defines the service's configuration data-object.
type Config struct {
	Socket string
	DSN    string
}

func (c *Config) Configure(log *slog.Logger, opts ...ConfigOpt) error {
	const op = "cfg.new-config"

	log.Info("starting configuring the service")

	if err := godotenv.Load("/project/pr-service/configs/.env"); err != nil {
		retErr := fmt.Errorf("%w: %w", ErrEnvFile, err)
		log.Error(fmt.Sprintf(
			"error of the %s: %s", op, retErr,
		))
		return retErr
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			log.Error(fmt.Sprintf(
				"error of the %s: %s", op, err,
			))
			return err
		}
	}

	return nil
}

package cfg

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
)

// Config defines the service's configuration data-object.
type Config struct {
	Socket string
}

func NewConfig(log *slog.Logger, opts ...ConfigOpt) (Config, error) {
	const op = "cfg.new-config"

	if err := godotenv.Load("../../configs/.env"); err != nil {
		retErr := fmt.Errorf("%w: %w", ErrEnvFile, err)
		log.Error(fmt.Sprintf(
			"error of the %s: %s", op, retErr,
		))
		return Config{}, retErr
	}

	conf := Config{}
	for _, opt := range opts {
		if err := opt(&conf); err != nil {
			log.Error(fmt.Sprintf(
				"error of the %s: %s", op, err,
			))
			return Config{}, err
		}
	}

	return conf, nil
}

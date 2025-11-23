package cfg

import (
	"fmt"
	"os"
)

// ConfigOpt defines the function for configuration the option
type ConfigOpt func(conf *Config) error

// ConfigSocket defines the logic of the socket's configuration.
func ConfigSocket(key string) ConfigOpt {
	return func(conf *Config) error {
		val := os.Getenv(key)

		if len(val) == 0 {
			return fmt.Errorf("%w: is empty", ErrSocketConfig)
		}
		conf.Socket = val

		return nil
	}
}

func ConfigDBSocket(key string) ConfigOpt {
	return func(conf *Config) error {
		val := os.Getenv(key)

		if len(val) == 0 {
			return fmt.Errorf("%w: is empty", ErrSocketConfig)
		}
		conf.DBSocket = val

		return nil
	}
}

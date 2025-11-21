package cfg

import "errors"

var (
	ErrSocketConfig = errors.New("cfg: error of the socket's configuration")
	ErrEnvFile      = errors.New("cfg: error of the .env file")
)

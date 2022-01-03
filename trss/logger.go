package transmissionrss

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

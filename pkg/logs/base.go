package logs

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	mu         sync.Mutex
	logLevel   = log.LevelDebug
	keyMasks   = []string{"password", "pass", "pwd", "secret", "token", "key", "otp"}
	valueMasks = []string{"password", "pass", "pwd", "secret", "token", "key", "otp"}
)

func newHelper(logger log.Logger) *log.Helper {
	f := log.NewFilter(
		logger,
		log.FilterLevel(logLevel),
		log.FilterKey(keyMasks...),
		log.FilterValue(valueMasks...),
	)
	return log.NewHelper(f)
}

func SetLogLevel(level log.Level) {
	mu.Lock()
	defer mu.Unlock()
	logLevel = level
}

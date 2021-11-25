package di

import (
	"fmt"
	"io"

	"go.uber.org/zap"
)

// var terminationGracePeriod = 3

// LogInitFatal logs initialization error and then forces to stop application.
func LogInitFatal(name string, err error) {
	msg := fmt.Sprintf("failed to init %s", name)
	GetLogger().Fatal(msg, zap.Error(err))
}


type closerFunc struct {
	key    string
	closer io.Closer
}

var closerFuncs = make([]*closerFunc, 0)

// RegisterCloser registers the given closer.
func RegisterCloser(key string, closer io.Closer) {
	l := GetLogger().Named("closer")

	for _, c := range closerFuncs {
		if key == c.key {
			l.Fatal(fmt.Sprintf("duplicate closer key: %s", key))
		}
	}

	closerFuncs = append(closerFuncs, &closerFunc{
		key:    key,
		closer: closer,
	})

	l.Info(fmt.Sprintf("%s closer is registered", key))
}

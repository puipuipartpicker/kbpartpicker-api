package retry

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// provide errors for retry package.
var (
	ErrBackOffTimeout  = errors.New("backoff timeout")
	ErrReachedMaxRetry = errors.New("reached max retry")
)

// Retrier is an interface for retry.
type Retrier interface {
	Do(context.Context, func(ctx context.Context) error) error
}

type retry struct {

	// maxRetryCnt is the maximum number of retries
	maxRetryCnt int

	// initialDuration is the amount of time to backoff after the first failure.
	initialDuration time.Duration

	// backoffTimeout is the backoff timeout.
	// If backoff time is exceeded this value, Returns an error.
	backoffTimeout time.Duration

	// factor is applied to the backoff after each retry.
	factor float64

	// jitter is the factor with which backoffs are randomized.
	jitter float64

	// errorFunc is the function to report an error.
	errorFunc func(error)
}

// New returns Retrier implementation and error.
func New(opts ...Option) (Retrier, error) {
	r := new(retry)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(r); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	rand.Seed(time.Now().UnixNano())

	return r, nil
}

func (r *retry) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	limit := time.NewTimer(r.backoffTimeout)
	defer limit.Stop()

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	durf := float64(r.initialDuration)

	var err error
	for cnt := 0; cnt < r.maxRetryCnt; cnt++ {
		if err = fn(ctx); err != nil {
			r.errorFunc(err)

			timer.Reset(time.Duration(durf))
			select {
			case <-limit.C:
				return fmt.Errorf("%s: %w", err.Error(), ErrBackOffTimeout)
			case <-ctx.Done():
				return fmt.Errorf("%s: %w", err.Error(), ctx.Err())
			case <-timer.C:
				// NOTE: https://chromium.googlesource.com/external/github.com/grpc/grpc-go/+/refs/heads/v1.10.x/backoff.go
				durf *= r.factor
				durf *= 1 + r.jitter*(rand.Float64()*2-1)

				continue
			}
		}

		return nil
	}

	return fmt.Errorf("%s: %w", err.Error(), ErrReachedMaxRetry)
}

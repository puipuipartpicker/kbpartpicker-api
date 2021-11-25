package retry

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Option configures "retry" structure.
type Option func(*retry) error

var defaultOptions = []Option{
	WithInitialDelay("10ms"),
	WithBackoffFactor(1.6),
	WithJitter(0.2),
	WithBackoffTimeout("3m"),
	WithRetryCnt(5),
	WithErrorFunc(nopErrorFunc),
}

var (
	nopErrorFunc     = func(error) {}
	errInvalidNumber = errors.New("invalid number")
)

// WithRetryCnt returns Option that sets the value to the retry.maxRetryCnt.
func WithRetryCnt(n int) Option {
	return func(r *retry) error {
		if n <= 0 {
			return fmt.Errorf("failed to set retry.maxRetryCnt: %d", n)
		}

		r.maxRetryCnt = n

		return nil
	}
}

// WithInitialDelay returns Option that sets the value to the retry.initialDuration.
func WithInitialDelay(str string) Option {
	return func(r *retry) error {
		d, err := parseDuration(str)
		if err != nil {
			return fmt.Errorf("failed to set retry.initialDuration: %s: %w", str, err)
		}

		r.initialDuration = d

		return nil
	}
}

// WithInitialDelayMilliSeconds returns Option that sets the value to the retry.initialDuration.
func WithInitialDelayMilliSeconds(str string) Option {
	return func(r *retry) error {
		n, err := parseInt(str)
		if err != nil {
			return fmt.Errorf("failed to set retry.initialDuration: %s: %w", str, err)
		}

		r.initialDuration = time.Duration(n) * time.Millisecond

		return nil
	}
}

// WithBackoffTimeout returns Option that sets the value to the retry.backoffTimeout.
func WithBackoffTimeout(str string) Option {
	return func(r *retry) error {
		d, err := parseDuration(str)
		if err != nil {
			return fmt.Errorf("failed to set retry.backoffTimeout: %s: %w", str, err)
		}

		r.backoffTimeout = d

		return nil
	}
}

// WithBackoffTimeoutSeconds returns Option that sets the value to the retry.backoffTimeout.
func WithBackoffTimeoutSeconds(str string) Option {
	return func(r *retry) error {
		n, err := parseInt(str)
		if err != nil {
			return fmt.Errorf("failed to set retry.backoffTimeout: %s: %w", str, err)
		}

		r.backoffTimeout = time.Duration(n) * time.Second

		return nil
	}
}

// WithBackoffFactor returns Option that sets the value to the retry.factor.
func WithBackoffFactor(f float64) Option {
	return func(r *retry) error {
		if f <= 0 {
			return fmt.Errorf("failed to set retry.factor: %v", f)
		}

		r.factor = f

		return nil
	}
}

// WithJitter returns Option that sets the value to the retry.jitter.
func WithJitter(f float64) Option {
	return func(r *retry) error {
		if f <= 0 {
			return fmt.Errorf("failed to set retry.jitter: %v", f)
		}

		r.jitter = f

		return nil
	}
}

// WithErrorFunc returns Option that sets the value to the retry.errorFunc.
func WithErrorFunc(fn func(err error)) Option {
	return func(r *retry) error {
		if fn == nil {
			return errors.New("failed to set retry.errorFunc")
		}

		r.errorFunc = fn

		return nil
	}
}

func parseInt(str string) (int, error) {
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("failed to convert: %w", err)
	}

	if n <= 0 {
		return 0, errInvalidNumber
	}

	return n, nil
}

func parseDuration(str string) (time.Duration, error) {
	d, err := time.ParseDuration(str)
	if err != nil {
		return 0, fmt.Errorf("failed to parse: %w", err)
	}

	if d <= 0 {
		return 0, errInvalidNumber
	}

	return d, nil
}

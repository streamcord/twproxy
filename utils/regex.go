package utils

import (
	"errors"
	"regexp"
	"time"
)

// CompileWithTimeout asynchronously compiles a regular expression pattern, and times out if the operation takes longer
// than 10ms to complete.
func CompileWithTimeout(r SpyglassNotificationRegexPattern) (*regexp.Regexp, error) {
	outCh := make(chan *regexp.Regexp, 1)
	errCh := make(chan error, 1)

	go func() {
		p, err := regexp.Compile(r.Pattern)

		if err != nil {
			errCh <- err
		} else {
			outCh <- p
		}

		close(outCh)
		close(errCh)
	}()

	select {
	case p := <-outCh:
		// Successful result
		return p, nil
	case err := <-errCh:
		// Error from regexp library
		return nil, err
	case <-time.After(10 * time.Millisecond):
		// Operation timed out
		err := errors.New("compile execution time exceeded >10ms")
		return nil, err
	}
}

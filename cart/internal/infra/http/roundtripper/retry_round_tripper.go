package roundtripper

import (
	"fmt"
	"net/http"
)

// RetryRoundTripper - http.RoundTripper (middleware) с поддержкой повторных попыток.
type RetryRoundTripper struct {
	rt                http.RoundTripper
	triggeredStatuses map[int]struct{}
	maxRetries        int
}

// NewRetryRoundTripper создает http.RoundTripper с поддержкой повторных попыток.
// Повторная попытка выполняется, если код ответа содержится в triggerStatuses.
func NewRetryRoundTripper(rt http.RoundTripper, triggerStatuses []int, maxRetries int) http.RoundTripper {
	l := &RetryRoundTripper{
		rt:                rt,
		triggeredStatuses: make(map[int]struct{}, len(triggerStatuses)),
		maxRetries:        maxRetries,
	}

	for _, status := range triggerStatuses {
		l.triggeredStatuses[status] = struct{}{}
	}

	return l
}

func (l *RetryRoundTripper) isTriggeredStatus(statusCode int) bool {
	_, ok := l.triggeredStatuses[statusCode]
	return ok
}

func copyRequest(r *http.Request) (*http.Request, error) {
	rCopy := r
	if r.Body != nil && r.GetBody != nil {
		bodyCopy, err := r.GetBody()
		if err != nil {
			return nil, fmt.Errorf("r.GetBody: %w", err)
		}

		rCopy = r.Clone(r.Context())
		rCopy.Body = bodyCopy
	}

	return rCopy, nil
}

func (l *RetryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rCopy, errCopy := copyRequest(r)
	if errCopy != nil {
		return nil, fmt.Errorf("copyRequest: %w", errCopy)
	}

	resp, err := l.rt.RoundTrip(rCopy)
	if err != nil {
		return resp, fmt.Errorf("l.rt.RoundTrip: %w", err)
	}

	curRetry := 1
	for l.isTriggeredStatus(resp.StatusCode) && curRetry <= l.maxRetries {
		resp.Body.Close()

		rCopy, errCopy = copyRequest(r)
		if errCopy != nil {
			return nil, fmt.Errorf("copyRequest: %w", errCopy)
		}

		resp, err = l.rt.RoundTrip(rCopy)
		if err != nil {
			return resp, fmt.Errorf("l.rt.RoundTrip: %w", err)
		}

		curRetry++
	}

	return resp, err
}

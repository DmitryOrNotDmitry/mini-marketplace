package roundtripper

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	args := m.Called(r)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestRetryRoundTripper_RoundTripOnlyErrors(t *testing.T) {
	mockRoundTripper := new(MockRoundTripper)

	// Мокнутый ответ для 420
	resp420 := &http.Response{
		StatusCode: 420,
		Body:       io.NopCloser(bytes.NewBufferString(``)),
	}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(resp420, nil)

	retryRT := NewRetryRoundTripper(mockRoundTripper, []int{420}, 3)

	respGot, err := retryRT.RoundTrip(&http.Request{})
	if err != nil {
		t.Fatal(err)
	}

	if respGot.StatusCode != resp420.StatusCode {
		t.Errorf("expected 420, got %d", respGot.StatusCode)
	}
}

func TestRetryRoundTripper_RoundTripSuccessOnSecondCall(t *testing.T) {
	mockRoundTripper := new(MockRoundTripper)

	resp420 := &http.Response{
		StatusCode: 420,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}
	resp200 := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("ok")),
	}

	mockRoundTripper.
		On("RoundTrip", mock.Anything).
		Return(resp420, nil).Once()
	mockRoundTripper.
		On("RoundTrip", mock.Anything).
		Return(resp200, nil).Once()

	retryRT := NewRetryRoundTripper(mockRoundTripper, []int{420}, 3)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	resp, err := retryRT.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRetryRoundTripper_RoundTripSuccessOnSecondCallButZeroRetries(t *testing.T) {
	mockRoundTripper := new(MockRoundTripper)

	resp420 := &http.Response{
		StatusCode: 420,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}
	resp200 := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("ok")),
	}

	mockRoundTripper.
		On("RoundTrip", mock.Anything).
		Return(resp420, nil).Once()
	mockRoundTripper.
		On("RoundTrip", mock.Anything).
		Return(resp200, nil).Once()

	retryRT := NewRetryRoundTripper(mockRoundTripper, []int{420}, 0)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	resp, err := retryRT.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 420 {
		t.Errorf("expected 420, got %d", resp.StatusCode)
	}
}

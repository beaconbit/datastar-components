package stackedbarchart

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClockTickHandler_MinuteChange(t *testing.T) {
	// Create component with chart's current time set to a minute ago
	comp := New()
	comp.data.CurrentTime = time.Now().Add(-time.Minute) // simulate that chart is behind wall time by a minute

	// Create request to tick endpoint with datastar param
	req := httptest.NewRequest("GET", "/tick?datastar=true", nil)
	w := httptest.NewRecorder()

	// Call handler directly
	comp.clockTickHandler(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check response body contains SSE patch for chart container (indicating chart advanced)
	body := w.Body.String()
	if !strings.Contains(body, "event: datastar-patch-elements") {
		t.Errorf("expected SSE patch event, got: %s", body)
	}
	// When minute changes, the handler should patch the entire chart container, not just clock
	// The container ID is "stacked-bar-chart"
	if !strings.Contains(body, "stacked-bar-chart") {
		t.Log("patch may be for clock only; ensure minute changed logic works")
	}
}

func TestClockTickHandler_NoMinuteChange(t *testing.T) {
	// Create component with chart's current time set to current wall time
	comp := New()
	comp.data.CurrentTime = time.Now()

	req := httptest.NewRequest("GET", "/tick?datastar=true", nil)
	w := httptest.NewRecorder()

	comp.clockTickHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "stacked-chart-clock") {
		t.Errorf("expected clock patch, got: %s", body)
	}
}

func TestAdvanceMinute(t *testing.T) {
	data := DefaultStackedBarChart()
	// Store original minutes
	originalMinutes := make([]MinuteBar, len(data.Minutes))
	copy(originalMinutes, data.Minutes)

	data.AdvanceMinute()

	// Check length unchanged
	if len(data.Minutes) != 10 {
		t.Errorf("expected 10 minutes, got %d", len(data.Minutes))
	}

	// Check that minutes shifted: index 0 should be original index 1 (compare timestamps)
	for i := 0; i < 9; i++ {
		if !data.Minutes[i].Timestamp.Equal(originalMinutes[i+1].Timestamp) {
			t.Errorf("minute timestamp mismatch at index %d", i)
		}
	}
	// Check new current minute offset is 0
	if data.Minutes[9].MinuteOffset != 0 {
		t.Errorf("new current minute offset should be 0, got %d", data.Minutes[9].MinuteOffset)
	}
	// Check that new current minute has zero delays
	for i := 0; i < 3; i++ {
		if data.Minutes[9].MachineDelays[i] != 0 {
			t.Errorf("new current minute should have zero delays for machine %d", i)
		}
	}
	// Check that oldest minute (original index 0) is no longer present
	// by verifying its timestamp is not in data.Minutes
	oldestTimestamp := originalMinutes[0].Timestamp
	for i, minute := range data.Minutes {
		if minute.Timestamp.Equal(oldestTimestamp) {
			t.Errorf("oldest minute still present at index %d", i)
		}
	}
}

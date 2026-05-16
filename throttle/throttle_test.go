package throttle

import (
	"testing"
	"time"
)

func TestThrottler_DisabledWhenZero(t *testing.T) {
	th := New(0)
	if !th.Disabled() {
		t.Fatal("expected Throttler to be disabled when maxPerSec=0")
	}
}

func TestThrottler_DisabledWhenNegative(t *testing.T) {
	th := New(-5)
	if !th.Disabled() {
		t.Fatal("expected Throttler to be disabled when maxPerSec<0")
	}
}

func TestThrottler_EnabledWhenPositive(t *testing.T) {
	th := New(10)
	if th.Disabled() {
		t.Fatal("expected Throttler to be enabled when maxPerSec>0")
	}
}

func TestThrottler_AllowDoesNotSleepUnderLimit(t *testing.T) {
	slept := time.Duration(0)
	th := newWithSleep(100, func(d time.Duration) { slept += d })

	// First call should not sleep (count=1, expected=windowStart+0).
	th.Allow()

	if slept > 0 {
		t.Errorf("expected no sleep on first Allow, got %v", slept)
	}
}

func TestThrottler_AllowDoesNotSleepWhenDisabled(t *testing.T) {
	slept := time.Duration(0)
	th := newWithSleep(0, func(d time.Duration) { slept += d })

	for i := 0; i < 1000; i++ {
		th.Allow()
	}

	if slept != 0 {
		t.Errorf("disabled Throttler should never sleep, got %v", slept)
	}
}

func TestThrottler_SleepsWhenOverLimit(t *testing.T) {
	slept := time.Duration(0)
	th := newWithSleep(2, func(d time.Duration) { slept += d })

	// Force the window to look full by pre-setting count.
	th.count = 2
	th.windowStart = time.Now()

	// This call is the 3rd in the window — over the limit of 2.
	th.Allow()

	if slept == 0 {
		t.Error("expected Throttler to sleep when over limit, but it did not")
	}
}

func TestThrottler_ResetsWindowAfterOneSecond(t *testing.T) {
	slept := time.Duration(0)
	th := newWithSleep(1, func(d time.Duration) { slept += d })

	// Simulate a window that started more than one second ago.
	th.windowStart = time.Now().Add(-2 * time.Second)
	th.count = 999 // Would normally trigger sleep.

	th.Allow()

	// Window should have reset; count is now 1 which is within limit.
	if th.count != 1 {
		t.Errorf("expected count=1 after window reset, got %d", th.count)
	}
}

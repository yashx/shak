package test

import (
	"testing"
	"time"

	"github.com/yashx/shak/rule"
)

var (
	epoch  = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	past   = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	future = time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	midday = time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
)

// Before

func TestBeforeValid(t *testing.T) {
	r := rule.Before(epoch)
	if err := r(past); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestBeforeError(t *testing.T) {
	r := rule.Before(epoch)
	err := r(future)
	if err == nil {
		t.Errorf("Error expected for time after bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestBeforeEqualError(t *testing.T) {
	r := rule.Before(epoch)
	if err := r(epoch); err == nil {
		t.Errorf("Error expected for time equal to bound (strict before)")
	}
}

// After

func TestAfterValid(t *testing.T) {
	r := rule.After(epoch)
	if err := r(future); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestAfterError(t *testing.T) {
	r := rule.After(epoch)
	err := r(past)
	if err == nil {
		t.Errorf("Error expected for time before bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestAfterEqualError(t *testing.T) {
	r := rule.After(epoch)
	if err := r(epoch); err == nil {
		t.Errorf("Error expected for time equal to bound (strict after)")
	}
}

// Between

func TestBetweenValid(t *testing.T) {
	r := rule.Between(past, future)
	if err := r(epoch); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestBetweenAtMin(t *testing.T) {
	r := rule.Between(past, future)
	if err := r(past); err != nil {
		t.Errorf("No errors expected at lower bound, got: %v", err)
	}
}

func TestBetweenAtMax(t *testing.T) {
	r := rule.Between(past, future)
	if err := r(future); err != nil {
		t.Errorf("No errors expected at upper bound, got: %v", err)
	}
}

func TestBetweenBeforeMinError(t *testing.T) {
	r := rule.Between(epoch, future)
	err := r(past)
	if err == nil {
		t.Errorf("Error expected for time before range")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestBetweenAfterMaxError(t *testing.T) {
	r := rule.Between(past, epoch)
	err := r(future)
	if err == nil {
		t.Errorf("Error expected for time after range")
	} else {
		t.Logf("got error: %v", err)
	}
}

// InThePast

func TestInThePastValid(t *testing.T) {
	r := rule.InThePast()
	if err := r(past); err != nil {
		t.Errorf("No errors expected for past time, got: %v", err)
	}
}

func TestInThePastError(t *testing.T) {
	r := rule.InThePast()
	err := r(future)
	if err == nil {
		t.Errorf("Error expected for future time")
	} else {
		t.Logf("got error: %v", err)
	}
}

// InThePastOrPresent

func TestInThePastOrPresentValid(t *testing.T) {
	r := rule.InThePastOrPresent()
	if err := r(past); err != nil {
		t.Errorf("No errors expected for past time, got: %v", err)
	}
}

func TestInThePastOrPresentNowValid(t *testing.T) {
	r := rule.InThePastOrPresent()
	if err := r(time.Now()); err != nil {
		t.Errorf("No errors expected for current time, got: %v", err)
	}
}

func TestInThePastOrPresentError(t *testing.T) {
	r := rule.InThePastOrPresent()
	err := r(future)
	if err == nil {
		t.Errorf("Error expected for future time")
	} else {
		t.Logf("got error: %v", err)
	}
}

// InTheFuture

func TestInTheFutureValid(t *testing.T) {
	r := rule.InTheFuture()
	if err := r(future); err != nil {
		t.Errorf("No errors expected for future time, got: %v", err)
	}
}

func TestInTheFutureError(t *testing.T) {
	r := rule.InTheFuture()
	err := r(past)
	if err == nil {
		t.Errorf("Error expected for past time")
	} else {
		t.Logf("got error: %v", err)
	}
}

// InTheFutureOrPresent

func TestInTheFutureOrPresentValid(t *testing.T) {
	r := rule.InTheFutureOrPresent()
	if err := r(future); err != nil {
		t.Errorf("No errors expected for future time, got: %v", err)
	}
}

func TestInTheFutureOrPresentError(t *testing.T) {
	r := rule.InTheFutureOrPresent()
	err := r(past)
	if err == nil {
		t.Errorf("Error expected for past time")
	} else {
		t.Logf("got error: %v", err)
	}
}

// OnDate

func TestOnDateValid(t *testing.T) {
	r := rule.OnDate(epoch)
	if err := r(midday); err != nil {
		t.Errorf("No errors expected for same date different time, got: %v", err)
	}
}

func TestOnDateError(t *testing.T) {
	r := rule.OnDate(epoch)
	err := r(past)
	if err == nil {
		t.Errorf("Error expected for different date")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestOnDateDifferentTimezone(t *testing.T) {
	// same instant in UTC and UTC+9 falls on different calendar dates
	loc := time.FixedZone("UTC+9", 9*60*60)
	ref := time.Date(2000, 1, 2, 0, 0, 0, 0, loc)         // 2000-01-02 in UTC+9
	value := time.Date(2000, 1, 1, 20, 0, 0, 0, time.UTC) // still 2000-01-01 in UTC, but 2000-01-02 in UTC+9
	r := rule.OnDate(ref)
	if err := r(value); err != nil {
		t.Errorf("No errors expected when value falls on ref date in ref timezone, got: %v", err)
	}
}

func TestOnDateTimezoneError(t *testing.T) {
	loc := time.FixedZone("UTC+9", 9*60*60)
	ref := time.Date(2000, 1, 1, 0, 0, 0, 0, loc)         // 2000-01-01 in UTC+9
	value := time.Date(2000, 1, 1, 20, 0, 0, 0, time.UTC) // 2000-01-02 in UTC+9
	r := rule.OnDate(ref)
	err := r(value)
	if err == nil {
		t.Errorf("Error expected when value falls on different date in ref timezone")
	} else {
		t.Logf("got error: %v", err)
	}
}

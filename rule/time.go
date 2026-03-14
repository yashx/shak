package rule

import (
	"time"

	"github.com/yashx/shak/types"
)

// Before validates that the time is before t.
func Before(t time.Time) types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		if value.Before(t) {
			return nil
		}
		return types.NewValidationError("time %v must be before %v", value, t)
	}
}

// After validates that the time is after t.
func After(t time.Time) types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		if value.After(t) {
			return nil
		}
		return types.NewValidationError("time %v must be after %v", value, t)
	}
}

// Between validates that the time is between min and max inclusive.
func Between(min, max time.Time) types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		if !value.Before(min) && !value.After(max) {
			return nil
		}
		return types.NewValidationError("time %v must be between %v and %v", value, min, max)
	}
}

// InThePast validates that the time is before the current time.
func InThePast() types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		if value.Before(time.Now()) {
			return nil
		}
		return types.NewValidationError("time %v must be in the past", value)
	}
}

// InThePastOrPresent validates that the time is before or equal to the current time.
func InThePastOrPresent() types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		if !value.After(time.Now()) {
			return nil
		}
		return types.NewValidationError("time %v must be in the past or present", value)
	}
}

// InTheFuture validates that the time is after the current time.
func InTheFuture() types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		if value.After(time.Now()) {
			return nil
		}
		return types.NewValidationError("time %v must be in the future", value)
	}
}

// InTheFutureOrPresent validates that the time is after or equal to the current time.
func InTheFutureOrPresent() types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		if !value.Before(time.Now()) {
			return nil
		}
		return types.NewValidationError("time %v must be in the future or present", value)
	}
}

// OnDate validates that the time falls on the same calendar date as t, using t's location.
func OnDate(t time.Time) types.Rule[time.Time] {
	return func(value time.Time) *types.ValidationError {
		v := value.In(t.Location())
		if v.Year() == t.Year() && v.Month() == t.Month() && v.Day() == t.Day() {
			return nil
		}
		return types.NewValidationError("time %v is not on date %v", value.In(t.Location()).Format("2006-01-02"), t.Format("2006-01-02"))
	}
}

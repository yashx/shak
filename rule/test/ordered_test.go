package test

import (
	"testing"

	"github.com/yashx/shak/rule"
)

// Min

func TestMinIntValid(t *testing.T) {
	r := rule.Min[int](5)
	if err := r(10); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMinIntEqual(t *testing.T) {
	r := rule.Min[int](5)
	if err := r(5); err != nil {
		t.Errorf("No errors expected when value equals min, got: %v", err)
	}
}

func TestMinIntError(t *testing.T) {
	r := rule.Min[int](5)
	err := r(4)
	if err == nil {
		t.Errorf("Error expected for value below min")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMinFloat64Valid(t *testing.T) {
	r := rule.Min[float64](1.5)
	if err := r(2.0); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMinFloat64Error(t *testing.T) {
	r := rule.Min[float64](1.5)
	err := r(1.0)
	if err == nil {
		t.Errorf("Error expected for float below min")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMinStringValid(t *testing.T) {
	r := rule.Min[string]("b")
	if err := r("c"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMinStringError(t *testing.T) {
	r := rule.Min[string]("m")
	err := r("a")
	if err == nil {
		t.Errorf("Error expected for string less than min")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMinCustomType(t *testing.T) {
	type Celsius float64
	r := rule.Min[Celsius](-273.15)
	if err := r(Celsius(0)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// Max

func TestMaxIntValid(t *testing.T) {
	r := rule.Max[int](100)
	if err := r(50); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMaxIntEqual(t *testing.T) {
	r := rule.Max[int](100)
	if err := r(100); err != nil {
		t.Errorf("No errors expected when value equals max, got: %v", err)
	}
}

func TestMaxIntError(t *testing.T) {
	r := rule.Max[int](100)
	err := r(101)
	if err == nil {
		t.Errorf("Error expected for value above max")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMaxFloat64Valid(t *testing.T) {
	r := rule.Max[float64](9.99)
	if err := r(5.5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMaxFloat64Error(t *testing.T) {
	r := rule.Max[float64](9.99)
	err := r(10.0)
	if err == nil {
		t.Errorf("Error expected for float above max")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestMaxStringValid(t *testing.T) {
	r := rule.Max[string]("z")
	if err := r("a"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMaxCustomType(t *testing.T) {
	type Percentage float64
	r := rule.Max[Percentage](100.0)
	if err := r(Percentage(99.9)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestMaxCustomTypeError(t *testing.T) {
	type Percentage float64
	r := rule.Max[Percentage](100.0)
	err := r(Percentage(100.1))
	if err == nil {
		t.Errorf("Error expected for value above max")
	} else {
		t.Logf("got error: %v", err)
	}
}

// Range

func TestRangeIntValid(t *testing.T) {
	r := rule.Range[int](1, 10)
	if err := r(5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRangeIntAtMin(t *testing.T) {
	r := rule.Range[int](1, 10)
	if err := r(1); err != nil {
		t.Errorf("No errors expected at lower bound, got: %v", err)
	}
}

func TestRangeIntAtMax(t *testing.T) {
	r := rule.Range[int](1, 10)
	if err := r(10); err != nil {
		t.Errorf("No errors expected at upper bound, got: %v", err)
	}
}

func TestRangeIntBelowMin(t *testing.T) {
	r := rule.Range[int](1, 10)
	err := r(0)
	if err == nil {
		t.Errorf("Error expected for value below range")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestRangeIntAboveMax(t *testing.T) {
	r := rule.Range[int](1, 10)
	err := r(11)
	if err == nil {
		t.Errorf("Error expected for value above range")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestRangeFloat64Valid(t *testing.T) {
	r := rule.Range[float64](0.0, 1.0)
	if err := r(0.5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestRangeFloat64Error(t *testing.T) {
	r := rule.Range[float64](0.0, 1.0)
	err := r(1.1)
	if err == nil {
		t.Errorf("Error expected for float above range")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestRangeStringValid(t *testing.T) {
	r := rule.Range[string]("a", "z")
	if err := r("m"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// LessThan

func TestLessThanIntValid(t *testing.T) {
	r := rule.LessThan[int](10)
	if err := r(9); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLessThanIntEqualError(t *testing.T) {
	r := rule.LessThan[int](10)
	err := r(10)
	if err == nil {
		t.Errorf("Error expected when value equals the bound (strict less than)")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestLessThanIntError(t *testing.T) {
	r := rule.LessThan[int](10)
	err := r(11)
	if err == nil {
		t.Errorf("Error expected for value not less than bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestLessThanFloat64Valid(t *testing.T) {
	r := rule.LessThan[float64](1.0)
	if err := r(0.99); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestLessThanCustomType(t *testing.T) {
	type Weight float64
	r := rule.LessThan[Weight](100.0)
	if err := r(Weight(99.9)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// LessThanEqual

func TestLessThanEqualIntValid(t *testing.T) {
	r := rule.LessThanEqual[int](10)
	if err := r(10); err != nil {
		t.Errorf("No errors expected at bound, got: %v", err)
	}
}

func TestLessThanEqualIntBelow(t *testing.T) {
	r := rule.LessThanEqual[int](10)
	if err := r(5); err != nil {
		t.Errorf("No errors expected below bound, got: %v", err)
	}
}

func TestLessThanEqualIntError(t *testing.T) {
	r := rule.LessThanEqual[int](10)
	err := r(11)
	if err == nil {
		t.Errorf("Error expected for value above bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestLessThanEqualFloat64(t *testing.T) {
	r := rule.LessThanEqual[float64](3.14)
	if err := r(3.14); err != nil {
		t.Errorf("No errors expected at exact bound, got: %v", err)
	}
}

func TestLessThanEqualString(t *testing.T) {
	r := rule.LessThanEqual[string]("m")
	if err := r("a"); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

// GreaterThan

func TestGreaterThanIntValid(t *testing.T) {
	r := rule.GreaterThan[int](5)
	if err := r(6); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestGreaterThanIntEqualError(t *testing.T) {
	r := rule.GreaterThan[int](5)
	err := r(5)
	if err == nil {
		t.Errorf("Error expected when value equals bound (strict greater than)")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestGreaterThanIntError(t *testing.T) {
	r := rule.GreaterThan[int](5)
	err := r(4)
	if err == nil {
		t.Errorf("Error expected for value not greater than bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestGreaterThanFloat64Valid(t *testing.T) {
	r := rule.GreaterThan[float64](0.0)
	if err := r(0.001); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestGreaterThanCustomType(t *testing.T) {
	type Age int
	r := rule.GreaterThan[Age](17)
	if err := r(Age(18)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestGreaterThanCustomTypeError(t *testing.T) {
	type Age int
	r := rule.GreaterThan[Age](17)
	err := r(Age(17))
	if err == nil {
		t.Errorf("Error expected for non-strictly-greater value")
	} else {
		t.Logf("got error: %v", err)
	}
}

// GreaterThanEqual

func TestGreaterThanEqualIntValid(t *testing.T) {
	r := rule.GreaterThanEqual[int](5)
	if err := r(5); err != nil {
		t.Errorf("No errors expected at bound, got: %v", err)
	}
}

func TestGreaterThanEqualIntAbove(t *testing.T) {
	r := rule.GreaterThanEqual[int](5)
	if err := r(10); err != nil {
		t.Errorf("No errors expected above bound, got: %v", err)
	}
}

func TestGreaterThanEqualIntError(t *testing.T) {
	r := rule.GreaterThanEqual[int](5)
	err := r(4)
	if err == nil {
		t.Errorf("Error expected for value below bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestGreaterThanEqualFloat64(t *testing.T) {
	r := rule.GreaterThanEqual[float64](0.0)
	if err := r(0.0); err != nil {
		t.Errorf("No errors expected at zero bound, got: %v", err)
	}
}

func TestGreaterThanEqualString(t *testing.T) {
	r := rule.GreaterThanEqual[string]("b")
	err := r("a")
	if err == nil {
		t.Errorf("Error expected for string less than bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

// Positive

func TestPositiveInt(t *testing.T) {
	r := rule.Positive[int]()
	if err := r(5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestPositiveIntError(t *testing.T) {
	r := rule.Positive[int]()
	err := r(-1)
	if err == nil {
		t.Errorf("Error expected for negative int")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestPositiveZeroError(t *testing.T) {
	r := rule.Positive[int]()
	err := r(0)
	if err == nil {
		t.Errorf("Error expected for zero (not positive)")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestPositiveFloat64(t *testing.T) {
	r := rule.Positive[float64]()
	if err := r(0.001); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestPositiveFloat64Error(t *testing.T) {
	r := rule.Positive[float64]()
	err := r(-0.5)
	if err == nil {
		t.Errorf("Error expected for negative float64")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestPositiveFloat64ZeroError(t *testing.T) {
	r := rule.Positive[float64]()
	err := r(0.0)
	if err == nil {
		t.Errorf("Error expected for zero float64 (not positive)")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestPositiveCustomIntType(t *testing.T) {
	type MyInt int
	r := rule.Positive[MyInt]()
	if err := r(MyInt(3)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestPositiveCustomIntTypeError(t *testing.T) {
	type MyInt int
	r := rule.Positive[MyInt]()
	err := r(MyInt(-3))
	if err == nil {
		t.Errorf("Error expected for negative custom type")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestPositiveCustomFloatType(t *testing.T) {
	type Celsius float64
	r := rule.Positive[Celsius]()
	if err := r(Celsius(37.0)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestPositiveCustomFloatTypeError(t *testing.T) {
	type Celsius float64
	r := rule.Positive[Celsius]()
	err := r(Celsius(-10.0))
	if err == nil {
		t.Errorf("Error expected for negative custom float type")
	} else {
		t.Logf("got error: %v", err)
	}
}

// PositiveOrZero

func TestPositiveOrZeroInt(t *testing.T) {
	r := rule.PositiveOrZero[int]()
	if err := r(0); err != nil {
		t.Errorf("No errors expected for zero, got: %v", err)
	}
}

func TestPositiveOrZeroIntPositive(t *testing.T) {
	r := rule.PositiveOrZero[int]()
	if err := r(10); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestPositiveOrZeroIntError(t *testing.T) {
	r := rule.PositiveOrZero[int]()
	err := r(-1)
	if err == nil {
		t.Errorf("Error expected for negative value")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestPositiveOrZeroFloat64(t *testing.T) {
	r := rule.PositiveOrZero[float64]()
	if err := r(0.0); err != nil {
		t.Errorf("No errors expected for zero float64, got: %v", err)
	}
}

func TestPositiveOrZeroFloat64Positive(t *testing.T) {
	r := rule.PositiveOrZero[float64]()
	if err := r(3.14); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestPositiveOrZeroFloat64Error(t *testing.T) {
	r := rule.PositiveOrZero[float64]()
	err := r(-0.001)
	if err == nil {
		t.Errorf("Error expected for negative float64")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestPositiveOrZeroCustomType(t *testing.T) {
	type Score float32
	r := rule.PositiveOrZero[Score]()
	if err := r(Score(0)); err != nil {
		t.Errorf("No errors expected for zero custom float type, got: %v", err)
	}
}

// Negative

func TestNegativeInt(t *testing.T) {
	r := rule.Negative[int]()
	if err := r(-5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNegativeIntError(t *testing.T) {
	r := rule.Negative[int]()
	err := r(1)
	if err == nil {
		t.Errorf("Error expected for positive int")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNegativeZeroError(t *testing.T) {
	r := rule.Negative[int]()
	err := r(0)
	if err == nil {
		t.Errorf("Error expected for zero (not negative)")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNegativeFloat64(t *testing.T) {
	r := rule.Negative[float64]()
	if err := r(-0.001); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNegativeFloat64Error(t *testing.T) {
	r := rule.Negative[float64]()
	err := r(1.5)
	if err == nil {
		t.Errorf("Error expected for positive float64")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNegativeCustomIntType(t *testing.T) {
	type Delta int64
	r := rule.Negative[Delta]()
	if err := r(Delta(-99)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNegativeCustomFloatType(t *testing.T) {
	type Offset float64
	r := rule.Negative[Offset]()
	if err := r(Offset(-3.14)); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNegativeCustomFloatTypeError(t *testing.T) {
	type Offset float64
	r := rule.Negative[Offset]()
	err := r(Offset(1.0))
	if err == nil {
		t.Errorf("Error expected for positive custom float type")
	} else {
		t.Logf("got error: %v", err)
	}
}

// NegativeOrZero

func TestNegativeOrZeroInt(t *testing.T) {
	r := rule.NegativeOrZero[int]()
	if err := r(0); err != nil {
		t.Errorf("No errors expected for zero, got: %v", err)
	}
}

func TestNegativeOrZeroIntNegative(t *testing.T) {
	r := rule.NegativeOrZero[int]()
	if err := r(-10); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNegativeOrZeroIntError(t *testing.T) {
	r := rule.NegativeOrZero[int]()
	err := r(1)
	if err == nil {
		t.Errorf("Error expected for positive value")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNegativeOrZeroFloat64(t *testing.T) {
	r := rule.NegativeOrZero[float64]()
	if err := r(0.0); err != nil {
		t.Errorf("No errors expected for zero float64, got: %v", err)
	}
}

func TestNegativeOrZeroFloat64Negative(t *testing.T) {
	r := rule.NegativeOrZero[float64]()
	if err := r(-2.5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNegativeOrZeroFloat64Error(t *testing.T) {
	r := rule.NegativeOrZero[float64]()
	err := r(0.001)
	if err == nil {
		t.Errorf("Error expected for positive float64")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNegativeOrZeroCustomType(t *testing.T) {
	type Adjustment float32
	r := rule.NegativeOrZero[Adjustment]()
	err := r(Adjustment(5))
	if err == nil {
		t.Errorf("Error expected for positive custom type")
	} else {
		t.Logf("got error: %v", err)
	}
}

// NotInRange

func TestNotInRangeValid(t *testing.T) {
	r := rule.NotInRange[int](1, 10)
	if err := r(20); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotInRangeBelowMin(t *testing.T) {
	r := rule.NotInRange[int](1, 10)
	if err := r(0); err != nil {
		t.Errorf("No errors expected for value below range, got: %v", err)
	}
}

func TestNotInRangeAtMinError(t *testing.T) {
	r := rule.NotInRange[int](1, 10)
	err := r(1)
	if err == nil {
		t.Errorf("Error expected for value at lower bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotInRangeAtMaxError(t *testing.T) {
	r := rule.NotInRange[int](1, 10)
	err := r(10)
	if err == nil {
		t.Errorf("Error expected for value at upper bound")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotInRangeInsideError(t *testing.T) {
	r := rule.NotInRange[int](1, 10)
	err := r(5)
	if err == nil {
		t.Errorf("Error expected for value inside range")
	} else {
		t.Logf("got error: %v", err)
	}
}

func TestNotInRangeFloat64Valid(t *testing.T) {
	r := rule.NotInRange[float64](0.0, 1.0)
	if err := r(1.5); err != nil {
		t.Errorf("No errors expected, got: %v", err)
	}
}

func TestNotInRangeCustomType(t *testing.T) {
	type Port int
	r := rule.NotInRange[Port](1024, 49151)
	if err := r(Port(80)); err != nil {
		t.Errorf("No errors expected for well-known port, got: %v", err)
	}
}

func TestNotInRangeCustomTypeError(t *testing.T) {
	type Port int
	r := rule.NotInRange[Port](1024, 49151)
	err := r(Port(8080))
	if err == nil {
		t.Errorf("Error expected for registered port inside range")
	} else {
		t.Logf("got error: %v", err)
	}
}

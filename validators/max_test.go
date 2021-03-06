package validators_test

import (
	"errors"
	"github.com/typerandom/validator/core"
	. "github.com/typerandom/validator/validators"
	"testing"
)

func TestThatMaxValidatorFailsForInvalidOptions(t *testing.T) {
	dummy := 100

	ctx := core.NewTestContext(dummy)
	err := MaxValidator(ctx, []interface{}{})

	if err == nil {
		t.Fatal(errors.New("Expected error, didn't get any."))
	}

	if err.Error() != "arguments.singleRequired" {
		t.Fatal(errors.New("Expected single argument required error."))
	}

	err = MaxValidator(ctx, []interface{}{"abc"})

	if err == nil {
		t.Fatal(errors.New("Expected error, didn't get any."))
	}

	if err.Error() != "arguments.invalidType" {
		t.Fatal(errors.New("Expected invalid arguments error."))
	}

	err = MaxValidator(ctx, []interface{}{"123", "123"})

	if err == nil {
		t.Fatal(errors.New("Expected error, didn't get any."))
	}

	if err.Error() != "arguments.singleRequired" {
		t.Fatal(errors.New("Expected single argument required error."))
	}
}

func testThatMaxValidatorFailsForValueOverLimit(t *testing.T, limit float64, dummy interface{}, expectedErr string) {
	ctx := core.NewTestContext(dummy)
	opts := []interface{}{limit}

	err := MaxValidator(ctx, opts)

	if err == nil {
		t.Fatal(errors.New("Expected error, didn't get any."))
	}

	if err.Error() != expectedErr {
		t.Fatalf("Expected cannot be more than error, got %s.", err)
	}
}

func testThatMaxValidatorSucceedsForValueOnLimit(t *testing.T, limit float64, dummy interface{}) {
	ctx := core.NewTestContext(dummy)
	opts := []interface{}{limit}

	if err := MaxValidator(ctx, opts); err != nil {
		t.Fatalf("Didn't expect error, but got one (%s).", err)
	}
}

func testThatMaxValidatorSucceedsForValueUnderLimit(t *testing.T, limit float64, dummy interface{}) {
	ctx := core.NewTestContext(dummy)
	opts := []interface{}{limit}

	if err := MaxValidator(ctx, opts); err != nil {
		t.Fatalf("Didn't expect error, but got one (%s).", err)
	}
}

func TestThatMaxValidatorFailsForIntValueOverLimit(t *testing.T) {
	testThatMaxValidatorFailsForValueOverLimit(t, 5, 6, "max.cannotBeGreaterThan")
}

func TestThatMaxValidatorSucceedsForIntValueOnLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueOnLimit(t, 5, 5)
}

func TestThatMaxValidatorSucceedsForIntValueUnderLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueUnderLimit(t, 5, 4)
}

func TestThatMaxValidatorFailsForFloatValueOverLimit(t *testing.T) {
	testThatMaxValidatorFailsForValueOverLimit(t, 5.5, 5.6, "max.cannotBeGreaterThan")
}

func TestThatMaxValidatorSucceedsForFloatValueOnLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueOnLimit(t, 5.5, 5.5)
}

func TestThatMaxValidatorSucceedsForFloatValueUnderLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueUnderLimit(t, 5.5, 5.4)
}

func TestThatMaxValidatorFailsForStringValueOverLimit(t *testing.T) {
	testThatMaxValidatorFailsForValueOverLimit(t, 5, "123456", "max.cannotBeLongerThan")
}

func TestThatMaxValidatorSucceedsForStringValueOnLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueOnLimit(t, 5, "12345")
}

func TestThatMaxValidatorSucceedsForStringValueUnderLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueUnderLimit(t, 5, "1234")
}

func TestThatMaxValidatorFailsForSliceLengthOverLimit(t *testing.T) {
	testThatMaxValidatorFailsForValueOverLimit(t, 5, []string{"1", "2", "3", "4", "5", "6"}, "max.cannotContainMoreItemsThan")
}

func TestThatMaxValidatorSucceedsForSliceLengthOnLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueOnLimit(t, 5, []string{"1", "2", "3", "4", "5"})
}

func TestThatMaxValidatorSucceedsForSliceLengthUnderLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueUnderLimit(t, 5, []string{"1", "2", "3", "4"})
}

func TestThatMaxValidatorFailsForMapLengthOverLimit(t *testing.T) {
	testThatMaxValidatorFailsForValueOverLimit(t, 5, map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6"}, "max.cannotContainMoreKeysThan")
}

func TestThatMaxValidatorSucceedsForMapLengthOnLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueOnLimit(t, 5, map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "5": "5"})
}

func TestThatMaxValidatorSucceedsForMapLengthUnderLimit(t *testing.T) {
	testThatMaxValidatorSucceedsForValueUnderLimit(t, 5, map[string]string{"1": "1", "2": "2", "3": "3", "4": "4"})
}

func TestThatMaxValidatorFailsForUnsupportedType(t *testing.T) {
	type Dummy struct{}
	testThatMaxValidatorFailsForValueOverLimit(t, 5, &Dummy{}, "type.unsupported")
}

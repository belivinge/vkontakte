package forms

import (
	"testing"
)

func TestErrors_Add(t *testing.T) {
	errs := make(errors)

	errs.Add("field1", "This is an error message")

	if len(errs["field1"]) != 1 {
		t.Errorf("Expected 1 error message for field 'field1', got %d", len(errs["field1"]))
	}
	if errs["field1"][0] != "This is an error message" {
		t.Errorf("Expected error message 'This is an error message', got '%s'", errs["field1"][0])
	}
}

func TestErrors_Get(t *testing.T) {
	errs := make(errors)

	errs.Add("field1", "This is an error message")

	message := errs.Get("field1")
	if message != "This is an error message" {
		t.Errorf("Expected error message 'This is an error message', got '%s'", message)
	}

	message = errs.Get("field2")
	if message != "" {
		t.Errorf("Expected empty string for field 'field2', got '%s'", message)
	}
}

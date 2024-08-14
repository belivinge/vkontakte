package forms

import (
	"net/url"
	"testing"
)

func TestForm_Required(t *testing.T) {
	data := url.Values{
		"field1": []string{""},
		"field2": []string{"non-empty"},
	}

	form := New(data)
	form.Required("field1", "field2")

	if len(form.Errors["field1"]) == 0 {
		t.Error("Expected error for field1 but got none")
	}
	if len(form.Errors["field2"]) > 0 {
		t.Error("Expected no error for field2 but got some")
	}
}

func TestForm_MaxLength(t *testing.T) {
	data := url.Values{
		"field1": []string{"short"},
		"field2": []string{"this string is definitely too long"},
	}

	form := New(data)
	form.MaxLength("field1", 10)
	form.MaxLength("field2", 10)

	if len(form.Errors["field2"]) == 0 {
		t.Error("Expected error for field2 but got none")
	}
	if len(form.Errors["field1"]) > 0 {
		t.Error("Expected no error for field1 but got some")
	}
}

func TestForm_PermittedValues(t *testing.T) {
	data := url.Values{
		"field1": []string{"valid"},
		"field2": []string{"invalid"},
	}

	form := New(data)
	form.PermittedValues("field1", "valid", "also-valid")
	form.PermittedValues("field2", "valid", "also-valid")

	if len(form.Errors["field2"]) == 0 {
		t.Error("Expected error for field2 but got none")
	}
	if len(form.Errors["field1"]) > 0 {
		t.Error("Expected no error for field1 but got some")
	}
}

func TestForm_MinLength(t *testing.T) {
	data := url.Values{
		"field1": []string{"short"},
		"field2": []string{"this is long enough"},
	}

	form := New(data)
	form.MinLength("field1", 10)
	form.MinLength("field2", 10)

	if len(form.Errors["field1"]) == 0 {
		t.Error("Expected error for field1 but got none")
	}
	if len(form.Errors["field2"]) > 0 {
		t.Error("Expected no error for field2 but got some")
	}
}

func TestForm_MatchesPattern(t *testing.T) {
	validEmail := "test@example.com"
	invalidEmail := "not-an-email"

	data := url.Values{
		"email1": []string{validEmail},
		"email2": []string{invalidEmail},
	}

	form := New(data)
	emailPattern := EmailRX
	form.MatchesPattern("email1", emailPattern)
	form.MatchesPattern("email2", emailPattern)

	if len(form.Errors["email2"]) == 0 {
		t.Error("Expected error for email2 but got none")
	}
	if len(form.Errors["email1"]) > 0 {
		t.Error("Expected no error for email1 but got some")
	}
}

func TestForm_Valid(t *testing.T) {
	dataWithErrors := url.Values{
		"field1": []string{""},
	}
	dataWithoutErrors := url.Values{
		"field2": []string{"valid"},
	}

	formWithErrors := New(dataWithErrors)
	formWithErrors.Required("field1")

	formWithoutErrors := New(dataWithoutErrors)

	if formWithErrors.Valid() {
		t.Error("Expected formWithErrors to be invalid but it was valid")
	}
	if !formWithoutErrors.Valid() {
		t.Error("Expected formWithoutErrors to be valid but it was invalid")
	}
}

func TestForm_RequiredAtLeastOne(t *testing.T) {
	data := url.Values{
		"field1": []string{""},
		"field2": []string{""},
		"field3": []string{"value"},
	}

	form := New(data)
	form.RequiredAtLeastOne("field1", "field2", "field3")

	if len(form.Errors["field1"]) > 0 {
		t.Error("Expected no error for field1 but got some")
	}
	if len(form.Errors["field2"]) > 0 {
		t.Error("Expected no error for field2 but got some")
	}
	if len(form.Errors["field3"]) > 0 {
		t.Error("Expected no error for field3 but got some")
	}
}

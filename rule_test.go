package formspec

import (
	"strings"
	"testing"
)

func TestSingleFilter(t *testing.T) {
	// name=" " -> Filter(Trim) -> name="" -> ValidationError

	f := newDummyform()
	f.Set("name", " ")

	aFormspec := New()
	aFormspec.Rule("name", RuleRequired()).Filter(func(value string) string {
		return strings.Trim(value, " ")
	})

	ok, _ := aFormspec.Validate(f)

	if ok {
		t.Errorf("expected validation error")
	}
}

func TestMultipleFilter(t *testing.T) {
	// Test multiple filter
	// name=" \n" -> Filter(Trim) -> name="" -> ValidationError

	f := newDummyform()
	f.Set("name", `
`)

	aFormspec := New()
	aFormspec.Rule("name", RuleRequired()).Filter(func(value string) string {
		return strings.Trim(value, " ")
	}).Filter(func(value string) string {
		return strings.Trim(value, "\n")
	})

	ok, _ := aFormspec.Validate(f)

	if ok {
		t.Errorf("expected validation error")
	}
}

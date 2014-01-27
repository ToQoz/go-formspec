package formspec_test

import (
	"fmt"
	"github.com/ToQoz/go-formspec"
)

type exampleForm struct {
	form map[string]string
}

func (e *exampleForm) Set(key, value string) {
	e.form[key] = value
}

func (e *exampleForm) FormValue(value string) string {
	return ""
}

func ExampleFormspec_basic() {
	aFormspec := formspec.New()
	aFormspec.Rule("name", formspec.RuleRequired())
	aFormspec.Rule("age", formspec.RuleRequired()).Message("must be integer. ok?").AllowBlank()
	aFormspec.Rule("nick", formspec.RuleRequired()).FullMessage("Please enter your cool nick.")

	f := &exampleForm{}
	// f.Set("name", "ToQoz")
	f.Set("age", "invalid int")
	// f.Set("age", "22")
	// f.Set("nick", "Toqoz")
	errs, ok := aFormspec.Validate(f)
	fmt.Printf("%q, %v\n", errs, ok)
}

func ExampleFormspec_getValidationErrorDetail() {
	aFormspec := formspec.New()
	aFormspec.Rule("name", formspec.RuleRequired())
	aFormspec.Rule("age", formspec.RuleRequired()).Message("must be integer. ok?").AllowBlank()
	aFormspec.Rule("nick", formspec.RuleRequired()).FullMessage("Please enter your cool nick.")

	f := &exampleForm{}

	// f.Set("name", "ToQoz")
	f.Set("age", "invalid int")
	// f.Set("age", "22")
	// f.Set("nick", "Toqoz")

	errs, ok := aFormspec.Validate(f)

	if !ok {
		for _, err := range errs {
			verr := err.(*formspec.Error)
			fmt.Printf("Validation error in %s. Message is %s.\n", verr.Field, verr.Message)
		}
	}
}
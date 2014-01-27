package formspec

import (
	"errors"
	"testing"
)

func TestBasic(t *testing.T) {
	f := newDummyform()

	aFormspec := New()
	aFormspec.Rule("name", RuleRequired())
	aFormspec.Rule("nick", RuleRequired()).FullMessage("Please enter your cool nick.")
	aFormspec.Rule("goodpoint", RuleRequired()).Message("must not be blank.")

	// Test
	//   when name is not given
	//     formspec returns error `name is required.`
	//     formspec returns error `Please enter your cool nick.`
	//     formspec returns error `goodpoing must not be blank.`

	errs, ok := aFormspec.Validate(f)

	// check ok is false
	if ok {
		t.Errorf("expected validation error")
		return
	}

	// check all errors expected occur
	func() {
		for _, err := range errs {
			if err.Error() == "name is required." {
				return
			}
		}

		t.Error("name is required.")
	}()

	func() {
		for _, err := range errs {
			if err.Error() == "Please enter your cool nick." {
				return
			}
		}

		t.Error("expected error `Please enter your cool nick.` is not got")
	}()

	func() {
		for _, err := range errs {
			if err.Error() == "goodpoint must not be blank." {
				return
			}
		}

		t.Error("expected error `goodpoint must not be blank` is not got.")
	}()
}

func TestClone(t *testing.T) {
	signInFormspec := New()
	signInFormspec.Rule("password", RuleRequired()).FullMessage("name is required")

	signUpFormspec := signInFormspec.Clone()
	signUpFormspec.Rule("password_confirmation", RuleRequired())
	signUpFormspec.Rule("password_confirmation", func(value string, f Form) error {
		if value != f.FormValue("password") {
			return errors.New("must be same as password")
		}

		return nil
	}).AllowBlank()

	// Test signInFormspec
	//   when password is given
	//     formspec should not return error
	f := newDummyform()
	f.Set("password", "hoge")
	_, ok := signInFormspec.Validate(f)

	if !ok {
		t.Errorf("validation error is not expected, but got it.")
	}

	// Test signUpFormspec
	//   when password is given and password_confirmation isn't given
	//     formspec should return 2 error
	//     its' message for 'password_confirmation' should be 'password_confirmation is required'
	f = newDummyform()
	f.Set("password", "hoge")

	if _, ok := signUpFormspec.Validate(f); ok {
		t.Errorf("validation error is expected, but not got it.")
	}

	// Test signUpFormspec
	//   when password is given and password_confirmation not same as password is given
	//     formspec should return 1 error
	//     its' message for field 'password_confirmation' must be 'password_confirmation is required'
	f = newDummyform()
	f.Set("password", "hoge")
	f.Set("password_confirmation", "hoge_different")

	if _, ok := signUpFormspec.Validate(f); ok {
		t.Errorf("validation error is expected, but not got it.")
	}

	// Test signUpFormspec
	//   when password is given and password_confirmation same as password is given
	//     formspec should not return error
	f = newDummyform()
	f.Set("password", "hoge")
	f.Set("password_confirmation", "hoge")

	if _, ok := signUpFormspec.Validate(f); !ok {
		t.Errorf("validation error is not expected, but got it.")
	}
}

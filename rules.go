package formspec

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	// Default rormats

	// regexp for number format. You change override.
	RuleFormatNumber = regexp.MustCompile(`\A[+-]?\d+\.?\d*\z`)
	// regexp for integer format. You change override.
	RuleFormatInt = regexp.MustCompile(`\A[+-]?\d+\z`)

	// Default messages

	RuleMessageRequired    = "is required."
	RuleMessageMaxLen      = "is too long. Max is %d character."
	RuleMessageMinLen      = "is too short. Min is %d character."
	RuleInvalidMessage     = "is invalid."
	RuleMessageNumber      = "must be number."
	RuleMessageInt         = "must be integer."
	RuleMessageLessThan    = "must be less than %d"
	RuleMessageGreaterThan = "must be greater than %d"
)

// funcs that return RuleFunc
// They must have prefix `Rule`.

func RuleRequired() RuleFunc {
	return func(value string, _ Form) error {
		if value == "" {
			return errors.New(RuleMessageRequired)
		}

		return nil
	}
}

func RuleMaxLen(maxLen int) RuleFunc {
	return func(value string, _ Form) error {
		if len(value) > maxLen {
			return fmt.Errorf(RuleMessageMaxLen, maxLen)
		}

		return nil
	}
}

func RuleMinLen(minLen int) RuleFunc {
	return func(value string, _ Form) error {
		if len(value) < minLen {
			return fmt.Errorf(RuleMessageMinLen, minLen)
		}

		return nil
	}
}

func RuleFormat(r *regexp.Regexp) RuleFunc {
	return func(value string, _ Form) error {
		if !r.MatchString(value) {
			return errors.New(RuleInvalidMessage)
		}

		return nil
	}
}

func RuleNumber() RuleFunc {
	return func(value string, _ Form) error {
		if !RuleFormatNumber.MatchString(value) {
			return errors.New(RuleMessageNumber)
		}

		return nil
	}
}

func RuleInt() RuleFunc {
	return func(value string, _ Form) error {
		if !RuleFormatInt.MatchString(value) {
			return errors.New(RuleMessageInt)
		}

		return nil
	}
}

func RuleFloatLessThan(a float64) RuleFunc {
	return func(value string, f Form) error {
		err := RuleNumber()(value, f)

		if err != nil {
			return err
		}

		i, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return errors.New(RuleMessageNumber)
		}

		if !(i < a) {
			return fmt.Errorf(RuleMessageLessThan, a)
		}

		return nil
	}
}

func RuleFloatGreaterThan(a float64) RuleFunc {
	return func(value string, f Form) error {
		err := RuleNumber()(value, f)

		if err != nil {
			return err
		}

		i, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return errors.New(RuleMessageNumber)
		}

		if !(i > a) {
			return fmt.Errorf(RuleMessageGreaterThan, a)
		}

		return nil
	}
}

func RuleIntLessThan(a float64) RuleFunc {
	return func(value string, f Form) error {
		err := RuleInt()(value, f)

		if err != nil {
			return err
		}

		i, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return errors.New(RuleMessageInt)
		}

		if !(i < a) {
			return fmt.Errorf(RuleMessageLessThan, a)
		}

		return nil
	}
}

func RuleIntGreaterThan(a int) RuleFunc {
	return func(value string, f Form) error {
		err := RuleInt()(value, f)

		if err != nil {
			return err
		}

		i, err := strconv.Atoi(value)

		if err != nil {
			return errors.New(RuleMessageInt)
		}

		if !(i > a) {
			return fmt.Errorf(RuleMessageGreaterThan, a)
		}

		return nil
	}
}

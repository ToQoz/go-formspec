package formspec

import (
	"regexp"
	"testing"
)

type dummyForm struct {
	form map[string]string
}

func (d *dummyForm) Set(key, value string) *dummyForm {
	d.form[key] = value
	return d
}

func (d *dummyForm) FormValue(key string) string {
	return d.form[key]
}

func newDummyform() *dummyForm {
	d := &dummyForm{}
	d.form = map[string]string{}
	return d
}

type ruleTestExample struct {
	input    string
	expected bool
}

// -----------------------------------------------------------------------------
// Test formspec.RuleRequired
// -----------------------------------------------------------------------------

func TestRuleRequired_WhenValueIsGiven_ItReturns__NO__Error(t *testing.T) {
	s := New()
	s.Rule("name", RuleRequired())

	if r := s.Validate(newDummyform().Set("name", "toqoz")); !r.Ok {
		t.Error("Test RuleRequired(): When value is given, returns no error. But got error.")
	}
}

func TestRuleRequired_WhenValueIsBlankOrNotGiven_ItReturnsError(t *testing.T) {
	s := New()
	s.Rule("name", RuleRequired())

	// value is blank
	if r := s.Validate(newDummyform().Set("name", "")); r.Ok {
		t.Error("Test RuleRequired(): When value is blank, returns error. But got no error.")
	}

	// value is not given
	if r := s.Validate(newDummyform()); r.Ok {
		t.Error("Test RuleRequired(): When value is not given, returns error. But got no error.")
		return
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleFormat
// -----------------------------------------------------------------------------

func TestRuleFormat(t *testing.T) {
	s := New()
	s.Rule("name", RuleFormat(regexp.MustCompile(`\Atoqoz403\+?[0-9]*@gmail.com\z`)))

	examples := []ruleTestExample{
		{"toqoz403@toqoz.net", false}, {"toqoz403@gmail.com@toqoz.net", false}, {"toqoz403@toqoz.net@gmail.com", false},
		{"toqoz403@gmail.com", true}, {"toqoz403+200@gmail.com", true},
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("name", example.input)

		if r := s.Validate(f); r.Ok != example.expected {
			t.Errorf("Test RuleFormat(regexp.MustCompile(`"+`\Atoqoz403\+?[0-9]*@gmail.com\z`+")): When `%s` is given, expected result is (_, %v), but got (_, %v).", example.input, example.expected, r.Ok)
		}
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleMaxLen
// -----------------------------------------------------------------------------

func TestRuleMaxLen(t *testing.T) {
	s := New()
	s.Rule("name", RuleMaxLen(10))

	examples := []ruleTestExample{
		{"aaaaaaaaa", true},    // len() => 9
		{"aaaaaaaaaa", true},   // len() => 10
		{"aaaaaaaaaaa", false}, // len() => 11
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("name", example.input)

		if r := s.Validate(f); r.Ok != example.expected {
			t.Errorf("Test RuleMaxLen(10): When `%s(len=%d)` is given, expected result is (_, %v). But got (_, %v).", example.input, len(example.input), example.expected, r.Ok)
		}
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleMinLen
// -----------------------------------------------------------------------------

func TestRuleMinLen(t *testing.T) {
	s := New()
	s.Rule("name", RuleMinLen(10))

	examples := []ruleTestExample{
		{"aaaaaaaaa", false},  // len() => 9
		{"aaaaaaaaaa", true},  // len() => 10
		{"aaaaaaaaaaa", true}, // len() => 11
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("name", example.input)

		if r := s.Validate(f); r.Ok != example.expected {
			t.Errorf("Test RuleMinLen(10): When `%s(len=%d)` is given, expected result is (_, %v). But got (_, %v).", example.input, len(example.input), example.expected, r.Ok)
		}
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleInt
// -----------------------------------------------------------------------------

func TestRuleInt(t *testing.T) {
	s := New()
	s.Rule("age", RuleInt())

	examples := []ruleTestExample{
		{"12", true}, {"0", true}, {"1000", true}, {"01", true},
		{"12.00", false}, {"0ab", false}, {"xxx", false}, {"x.x", false},
	}

	for _, test := range examples {
		f := newDummyform()
		f.Set("age", test.input)

		if r := s.Validate(f); r.Ok != test.expected {
			t.Errorf("Test RuleInt: When `%s` is given, expected result is (_, %v). But got (_, %v).", test.input, test.expected, r.Ok)
		}
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleNumber
// -----------------------------------------------------------------------------

func TestRuleNumber(t *testing.T) {
	s := New()
	s.Rule("age", RuleNumber())

	examples := []ruleTestExample{
		{"12", true}, {"0", true}, {"1000", true}, {"01", true}, {"12.521", true},
		{"12.0.1", false}, {"12.a0", false}, {"0ab", false}, {"xx", false}, {"x.x", false},
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("age", example.input)

		if r := s.Validate(f); r.Ok != example.expected {
			t.Errorf("Test RuleFloat: When value is `%s`, it returns (_, %v). But got (_, %v).", example.input, example.expected, r.Ok)
		}
	}
}

var ()

// -----------------------------------------------------------------------------
// Test formspec.RuleIntGreaterThan
// -----------------------------------------------------------------------------

func TestRuleIntGreaterThan(t *testing.T) {
	s := New()
	s.Rule("age", RuleIntGreaterThan(10))

	examples := []ruleTestExample{
		// (can't parse to int) -> false
		{"12x", false}, {"12.x", false}, {"12.0", false}, {"x.0", false},
		// (less than) -> false
		{"9", false},
		// (eq to) -> false
		{"10", false},
		// (greater than) -> true
		{"11", true}, {"12", true},
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("age", example.input)

		r := s.Validate(f)

		if example.expected != r.Ok {
			t.Errorf("Test RuleIntGreaterThan(10): When `%s` is given, expected result is (_, %v). But got = (_, %v)", example.input, example.expected, r.Ok)
		}
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleIntLessThan
// -----------------------------------------------------------------------------

func TestRuleIntLessThan(t *testing.T) {
	s := New()
	s.Rule("age", RuleIntLessThan(10))

	examples := []ruleTestExample{
		// (can't parse to int) -> false
		{"12x", false}, {"12.x", false}, {"12.0", false}, {"x.0", false},
		// (less than) -> false
		{"9", true},
		// (eq to) -> false
		{"10", false},
		// (greater than) -> true
		{"11", false}, {"12", false},
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("age", example.input)

		r := s.Validate(f)

		if example.expected != r.Ok {
			t.Errorf("Test RuleIntLessThan(10): When `%s` is given, expected result is (_, %v). But got = (_, %v)", example.input, example.expected, r.Ok)
		}
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleFloatGreaterThan
// -----------------------------------------------------------------------------

func TestRuleFloatGreaterThan(t *testing.T) {
	s := New()
	s.Rule("age", RuleFloatGreaterThan(10.0))

	examples := []ruleTestExample{
		// (can't parse to float) -> false
		{"12x", false}, {"12.x", false}, {"12.s0", false}, {"x.0", false},
		// (less than) -> false
		{"9.0", false},
		// (eq to) -> false
		{"10.0", false}, {"10", false},
		// (greater than) -> true
		{"11.0", true}, {"12.0", true}, {"13", true},
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("age", example.input)

		r := s.Validate(f)

		if example.expected != r.Ok {
			t.Errorf("Test RuleFloatGreaterThan(10.0): When `%s` is given, expected result is (_, %v). But got = (_, %v)", example.input, example.expected, r.Ok)
		}
	}
}

// -----------------------------------------------------------------------------
// Test formspec.RuleFloatLessThan
// -----------------------------------------------------------------------------

func TestRuleFloatLessThan(t *testing.T) {
	s := New()
	s.Rule("age", RuleFloatLessThan(10.0))

	examples := []ruleTestExample{
		// (can't parse to int) -> false
		{"12x", false}, {"12.x", false}, {"12.s0", false}, {"x.0", false},
		// (less than) -> false
		{"9.0", true}, {"9.99", true},
		// (eq to) -> false
		{"10.0", false}, {"10", false},
		// (greater than) -> true
		{"11.0", false}, {"12.0", false}, {"13", false},
	}

	for _, example := range examples {
		f := newDummyform()
		f.Set("age", example.input)

		r := s.Validate(f)

		if example.expected != r.Ok {
			t.Errorf("Test RuleFloatLessThan(10.0): When `%s` is given, expected result is (_, %v). But got = (_, %v)", example.input, example.expected, r.Ok)
		}
	}
}

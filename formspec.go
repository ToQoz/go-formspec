package formspec

import (
	"errors"
	"fmt"
)

// ----------------------------------------------------------------------------
// Error
// ----------------------------------------------------------------------------

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (ve *Error) Error() string {
	return ve.Message
}

// ----------------------------------------------------------------------------
// Formspec
// ----------------------------------------------------------------------------

type Form interface {
	FormValue(string) string
}

func New() *Formspec {
	return &Formspec{}
}

type Formspec struct {
	Rules []*Rule
}

func (f *Formspec) Rule(field string, ruleFunc RuleFunc) *Rule {
	rule := &Rule{Field: field, RuleFunc: ruleFunc}
	f.Rules = append(f.Rules, rule)
	return rule
}

func (fspec *Formspec) Validate(f Form) (errors []error, ok bool) {
	ok = true

	for _, rule := range fspec.Rules {
		err := rule.Call(f)

		if err != nil {
			ok = false
			errors = append(errors, &Error{Field: rule.Field, Message: err.Error()})
		}
	}

	return
}

func (f *Formspec) Clone() *Formspec {
	clone := &Formspec{}

	for _, rule := range f.Rules {
		clone.Rules = append(clone.Rules, rule.clone())
	}

	return clone
}

// ----------------------------------------------------------------------------
// Rule
// ----------------------------------------------------------------------------

type FilterFunc func(string) string
type RuleFunc func(value string, f Form) error

type Rule struct {
	Field       string
	RuleFunc    RuleFunc
	FilterFuncs []FilterFunc
	allowBlank  bool

	// This is used prior to Rule.message.
	fullMessage string
	// This is used prior to error message that is returned from Rule.RuleFunc.
	message string
}

func (v *Rule) AllowBlank() *Rule {
	v.allowBlank = true
	return v
}

// If you override error message. Use following funcs `FullMessage()/Message()`.

// FullMessage sets Rule.fullMessage.
func (v *Rule) FullMessage(m string) *Rule {
	v.fullMessage = m
	return v
}

// Message sets Rule.message.
func (v *Rule) Message(m string) *Rule {
	v.message = m
	return v
}

func (r *Rule) Filter(filterFunc FilterFunc) *Rule {
	r.FilterFuncs = append(r.FilterFuncs, filterFunc)
	return r
}

func (r *Rule) Call(f Form) error {
	v := f.FormValue(r.Field)

	// Filter value
	if len(r.FilterFuncs) > 0 {
		for _, filterFunc := range r.FilterFuncs {
			v = filterFunc(v)
		}
	}

	// If rule.allowblank is true, all rule returns no error when value is blank.
	if v == "" && r.allowBlank {
		return nil
	}

	err := r.RuleFunc(v, f)

	if err != nil {
		if r.fullMessage != "" {
			return errors.New(r.fullMessage)
		}

		if r.message != "" {
			return fmt.Errorf("%s %s", r.Field, r.message)
		}

		return fmt.Errorf("%s %s", r.Field, err.Error())
	}

	return nil
}

func (v *Rule) clone() *Rule {
	return &Rule{
		Field:       v.Field,
		RuleFunc:    v.RuleFunc,
		FilterFuncs: v.FilterFuncs,
		allowBlank:  v.allowBlank,
		message:     v.message,
		fullMessage: v.fullMessage,
	}
}

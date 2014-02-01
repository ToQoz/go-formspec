package formspec

import (
	"errors"
	"fmt"
)

type Result struct {
	Ok     bool
	Errors []*Error `json:"errors"`
}

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
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

func (f *Formspec) Validate(form Form) *Result {
	r := &Result{Ok: true}

	for _, rule := range f.Rules {
		err := rule.Call(form)

		if err != nil {
			r.Ok = false
			r.Errors = append(r.Errors, &Error{Field: rule.Field, Message: err.Error()})
		}
	}

	return r
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

func (r *Rule) AllowBlank() *Rule {
	r.allowBlank = true
	return r
}

// If you override error message. Use following funcs `FullMessage()/Message()`.

// FullMessage sets Rule.fullMessage.
func (r *Rule) FullMessage(m string) *Rule {
	r.fullMessage = m
	return r
}

// Message sets Rule.message.
func (r *Rule) Message(m string) *Rule {
	r.message = m
	return r
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

func (r *Rule) clone() *Rule {
	return &Rule{
		Field:       r.Field,
		RuleFunc:    r.RuleFunc,
		FilterFuncs: r.FilterFuncs,
		allowBlank:  r.allowBlank,
		message:     r.message,
		fullMessage: r.fullMessage,
	}
}

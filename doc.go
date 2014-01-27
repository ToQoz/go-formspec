// Copyright 2014 Takatoshi Matsumoto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package github.com/ToQoz/formspec validates a form. So it will expresses spec for form. This is generally used in http.Handler.

Usage

	package main

	import (
		"encoding/json"
		"fmt"
		"github.com/ToQoz/go-formspec"
		"log"
		"net/http"
	)

	type ErrorJson struct {
		Errors []*formspec.Error `json:"errors"`
	}

	func main() {
		s := formspec.New()
		s.Rule("name", formspec.RuleRequired())
		s.Rule("age", formspec.RuleInt()).Message("must be integer. ok?").AllowBlank()
		s.Rule("nick", formspec.RuleRequired()).FullMessage("Please enter your cool nickname.")

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			errs, ok := s.Validate(r)

			if !ok {
				w.WriteHeader(403)

				for _, err := range errs {
					verr := err.(*formspec.Error)
					w.Write([]byte(fmt.Sprintf("Validation error in %s. %s\n", verr.Field, verr.Message)))
				}

				return
			}

			w.Write([]byte("ok"))
		})

		http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
			errs, ok := s.Validate(r)

			if !ok {
				var verrs []*formspec.Error

				for _, err := range errs {
					verrs = append(verrs, err.(*formspec.Error))
				}

				if j, err := json.Marshal(&validationErrorJson{Errors: verrs}); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json; charset=utf8")
				w.WriteHeader(403)
				w.Write(j)

				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf8")
			w.Write([]byte(`{"message": "ok"}`))
		})

		log.Fatal(http.ListenAndServe(":8888", nil))
	}
*/

package formspec

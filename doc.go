// Copyright 2014 Takatoshi Matsumoto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package github.com/ToQoz/formspec validates a form. So it will expresses spec for form. This is generally used in http.Handler.

Simple usage in http.Handler.

	package main

	import (
		"fmt"
		"github.com/ToQoz/go-formspec"
		"log"
		"net/http"
	)

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

		log.Fatal(http.ListenAndServe(":8888", nil))
	}
*/
package formspec

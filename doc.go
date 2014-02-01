// Copyright 2014 Takatoshi Matsumoto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package github.com/ToQoz/formspec validates a form. So it will expresses **spec** for form. This is generally used in http.Handler, but you can use *formspec.Result as a retuen value of your validation func in models.

ExampleApp: https://github.com/ToQoz/go-formspec/tree/master/_example

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
			vr := s.Validate(r)

			if !vr.Ok {
				w.WriteHeader(403)

				for _, verr := range vr.Errors {
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

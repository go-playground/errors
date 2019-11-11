Package errors
============
![Project status](https://img.shields.io/badge/version-5.0.1-green.svg)
[![Build Status](https://travis-ci.org/go-playground/errors.svg?branch=master)](https://travis-ci.org/go-playground/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-playground/errors)](https://goreportcard.com/report/github.com/go-playground/errors)
[![GoDoc](https://godoc.org/github.com/go-playground/errors?status.svg)](https://godoc.org/github.com/go-playground/errors)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Package errors is an errors wrapping package to help propagate and chain errors as well as attach
stack traces, tags(additional information) and even a Type classification system to categorize errors into types eg. Permanent vs Transient.


Common Questions

Why another package?
There are two main reasons.
- I think that the programs generating the original error(s) should be responsible for handling them, even though this package allows access to the original error, and that the callers are mainly interested in:
  - If the error is Transient or Permanent for retries.
  - Additional details for logging.

- IMO most of the existing packages either don't take the error handling far enough, too far or down right unfriendly to use/consume. 

Features
--------
- [x] works with go-playground/log, the Tags will be added as Field Key Values and Types will be concatenated as well when using `WithError`
- [x] helpers to extract and classify error types using `RegisterHelper(...)`, many already existing such as ioerrors, neterrors, awserrors...
- [x] built in helpers only need to be imported, eg. `_ github.com/go-playground/errors/v5/helpers/neterrors` allowing libraries to register their own helpers not needing the caller to do or guess what needs to be imported.

Installation
------------

Use go get.

	go get -u github.com/go-playground/errors/v5
    
Usage
-----
```go
package main

import (
	"fmt"
	"io"

	"github.com/go-playground/errors/v5"
)

func main() {
	err := level1("testing error")
	fmt.Println(err)
	if errors.HasType(err, "Permanent") {
		// os.Exit(1)
		fmt.Println("it is a permanent error")
	}

	// root error
	cause := errors.Cause(err)
	fmt.Println("CAUSE:", cause)

	// can even still inspect the internal error
	fmt.Println(errors.Cause(err) == io.EOF) // will extract the cause for you
	fmt.Println(errors.Cause(cause) == io.EOF)

	// and still in a switch
	switch errors.Cause(err) {
	case io.EOF:
		fmt.Println("EOF error")
	default:
		fmt.Println("unknown error")
	}
}

func level1(value string) error {
	if err := level2(value); err != nil {
		return errors.Wrap(err, "level2 call failed")
	}
	return nil
}

func level2(value string) error {
	err := io.EOF
	return errors.Wrap(err, "failed to do something").AddTypes("Permanent").AddTags(errors.T("value", value))
}
```

or using stack only

```go
package main

import (
	"fmt"

	"github.com/go-playground/errors/v5"
)

func main() {
	// maybe you just want to grab a stack trace and process on your own like go-playground/log
	// uses it to produce a stack trace log message
	frame := errors.Stack()
	fmt.Printf("Function: %s File: %s Line: %d\n", frame.Function(), frame.File(), frame.Line())

	// and still have access to the underlying runtime.Frame
	fmt.Printf("%+v\n", frame.Frame)
}
```

Package Versioning
----------
Using Go modules and proper semantic version releases (as always).

License
------
Distributed under MIT License, please see license file in code for more details.

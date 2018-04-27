Package errors
============
![Project status](https://img.shields.io/badge/version-3.2.1-green.svg)
[![Build Status](https://semaphoreci.com/api/v1/joeybloggs/errors/branches/master/badge.svg)](https://semaphoreci.com/joeybloggs/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-playground/errors)](https://goreportcard.com/report/github.com/go-playground/errors)
[![GoDoc](https://godoc.org/github.com/go-playground/errors?status.svg)](https://godoc.org/github.com/go-playground/errors)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Package errors is an errors wrapping package to help propagate and chain errors as well as attach
stack traces, tags(additional information) and even a Type classification system to categorize errors into types eg. Permanent vs Transient.


Common Questions

Why another package?
Because IMO most of the existing packages either don't take the error handling far enough, too far or down right unfriendly to use/consume. 

Features
--------
- [x] works with go-playground/log, the Tags will be added as Field Key Values and Types will be concatenated as well when using `WithError`
- [x] helpers to extract and classify error types using `RegisterHelper(...)`, many already existing such as ioerrors, neterrors, awserrors...

Installation
------------

Use go get.

	go get -u github.com/go-playground/errors
    
Usage
-----
```go
package main

import (
	"fmt"
	"io"

	"github.com/go-playground/errors"
)

func main() {
	err := level1("testing error")
	fmt.Println(err)
	if errors.HasType(err, "Permanent") {
		// os.Exit(1)
	}

	// root error
	cause := errors.Cause(err)
	fmt.Println(cause)

	// can even still inspect the internal error
	fmt.Println(errors.Cause(err) == io.EOF) // will extract the cause for you
	fmt.Println(errors.Cause(cause) == io.EOF)
}

func level1(value string) error {
	if err := level2(value); err != nil {
		return errors.Wrap(err, "level2 call failed")
	}
	return nil
}

func level2(value string) error {
	err := fmt.Errorf("this is an %s", "error")
	return errors.Wrap(err, "failed to do something").AddTypes("Permanent").AddTags(errors.T("value", value))
}
```

or using stack only

```go
package main

import (
	"fmt"

	"strings"

	"github.com/go-playground/errors"
)

func main() {
	// maybe you just want to grab a stack trace and process on your own like go-playground/log
	// uses it to produce a stack trace log message
	frame := errors.Stack()
	name := fmt.Sprintf("%n", frame)
	file := fmt.Sprintf("%+s", frame)
	line := fmt.Sprintf("%d", frame)
	parts := strings.Split(file, "\n\t")
	if len(parts) > 1 {
		file = parts[1]
	}

	fmt.Printf("Name: %s File: %s Line: %s\n", name, file, line)
}
```

Package Versioning
----------
I'm jumping on the vendoring bandwagon, you should vendor this package as I will not
be creating different version with gopkg.in like allot of my other libraries.

Why? because my time is spread pretty thin maintaining all of the libraries I have + LIFE,
it is so freeing not to worry about it and will help me keep pouring out bigger and better
things for you the community.

License
------
Distributed under MIT License, please see license file in code for more details.

// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: reserved

//nolint:golint
package reserved

import (
	"fmt"

	"v.io/v23/context"
	"v.io/v23/verror"
)

var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

//////////////////////////////////////////////////
// Error definitions

var (

	// ErrGlobMaxRecursionReached indicates that the Glob request exceeded the
	// max recursion level.
	ErrGlobMaxRecursionReached = verror.NewIDAction("v.io/v23/rpc/reserved.GlobMaxRecursionReached", verror.NoRetry)
	// ErrGlobMatchesOmitted indicates that some of the Glob results might
	// have been omitted due to access restrictions.
	ErrGlobMatchesOmitted = verror.NewIDAction("v.io/v23/rpc/reserved.GlobMatchesOmitted", verror.NoRetry)
	// ErrGlobNotImplemented indicates that Glob is not implemented by the
	// object.
	ErrGlobNotImplemented = verror.NewIDAction("v.io/v23/rpc/reserved.GlobNotImplemented", verror.NoRetry)
)

// ErrorfErrGlobMaxRecursionReached calls ErrGlobMaxRecursionReached.Errorf with the supplied arguments.
func ErrorfErrGlobMaxRecursionReached(ctx *context.T, format string) error {
	return ErrGlobMaxRecursionReached.Errorf(ctx, format)
}

// MessageErrGlobMaxRecursionReached calls ErrGlobMaxRecursionReached.Message with the supplied arguments.
func MessageErrGlobMaxRecursionReached(ctx *context.T, message string) error {
	return ErrGlobMaxRecursionReached.Message(ctx, message)
}

// ParamsErrGlobMaxRecursionReached extracts the expected parameters from the error's ParameterList.
func ParamsErrGlobMaxRecursionReached(argumentError error) (verrorComponent string, verrorOperation string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	return
}

// ErrorfErrGlobMatchesOmitted calls ErrGlobMatchesOmitted.Errorf with the supplied arguments.
func ErrorfErrGlobMatchesOmitted(ctx *context.T, format string) error {
	return ErrGlobMatchesOmitted.Errorf(ctx, format)
}

// MessageErrGlobMatchesOmitted calls ErrGlobMatchesOmitted.Message with the supplied arguments.
func MessageErrGlobMatchesOmitted(ctx *context.T, message string) error {
	return ErrGlobMatchesOmitted.Message(ctx, message)
}

// ParamsErrGlobMatchesOmitted extracts the expected parameters from the error's ParameterList.
func ParamsErrGlobMatchesOmitted(argumentError error) (verrorComponent string, verrorOperation string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	return
}

// ErrorfErrGlobNotImplemented calls ErrGlobNotImplemented.Errorf with the supplied arguments.
func ErrorfErrGlobNotImplemented(ctx *context.T, format string) error {
	return ErrGlobNotImplemented.Errorf(ctx, format)
}

// MessageErrGlobNotImplemented calls ErrGlobNotImplemented.Message with the supplied arguments.
func MessageErrGlobNotImplemented(ctx *context.T, message string) error {
	return ErrGlobNotImplemented.Message(ctx, message)
}

// ParamsErrGlobNotImplemented extracts the expected parameters from the error's ParameterList.
func ParamsErrGlobNotImplemented(argumentError error) (verrorComponent string, verrorOperation string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	return
}

type paramListIterator struct {
	err      error
	idx, max int
	params   []interface{}
}

func (pl *paramListIterator) next() (interface{}, error) {
	if pl.err != nil {
		return nil, pl.err
	}
	if pl.idx+1 > pl.max {
		pl.err = fmt.Errorf("too few parameters: have %v", pl.max)
		return nil, pl.err
	}
	pl.idx++
	return pl.params[pl.idx-1], nil
}

func (pl *paramListIterator) preamble() (component, operation string, err error) {
	var tmp interface{}
	if tmp, err = pl.next(); err != nil {
		return
	}
	var ok bool
	if component, ok = tmp.(string); !ok {
		return "", "", fmt.Errorf("ParamList[0]: component name is not a string: %T", tmp)
	}
	if tmp, err = pl.next(); err != nil {
		return
	}
	if operation, ok = tmp.(string); !ok {
		return "", "", fmt.Errorf("ParamList[1]: operation name is not a string: %T", tmp)
	}
	return
}

var initializeVDLCalled bool

// initializeVDL performs vdl initialization.  It is safe to call multiple times.
// If you have an init ordering issue, just insert the following line verbatim
// into your source files in this package, right after the "package foo" clause:
//
//    var _ = initializeVDL()
//
// The purpose of this function is to ensure that vdl initialization occurs in
// the right order, and very early in the init sequence.  In particular, vdl
// registration and package variable initialization needs to occur before
// functions like vdl.TypeOf will work properly.
//
// This function returns a dummy value, so that it can be used to initialize the
// first var in the file, to take advantage of Go's defined init order.
func initializeVDL() struct{} {
	if initializeVDLCalled {
		return struct{}{}
	}
	initializeVDLCalled = true

	return struct{}{}
}

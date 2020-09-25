// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: vine

//nolint:golint
package vine

import (
	"fmt"

	v23 "v.io/v23"
	"v.io/v23/context"
	"v.io/v23/rpc"
	"v.io/v23/vdl"
	"v.io/v23/verror"
)

var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

//////////////////////////////////////////////////
// Type definitions

// PeerKey is a key that represents a connection from a Dialer tag to an Acceptor tag.
type PeerKey struct {
	Dialer   string
	Acceptor string
}

func (PeerKey) VDLReflect(struct {
	Name string `vdl:"v.io/x/ref/runtime/protocols/vine.PeerKey"`
}) {
}

func (x PeerKey) VDLIsZero() bool { //nolint:gocyclo
	return x == PeerKey{}
}

func (x PeerKey) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	if x.Dialer != "" {
		if err := enc.NextFieldValueString(0, vdl.StringType, x.Dialer); err != nil {
			return err
		}
	}
	if x.Acceptor != "" {
		if err := enc.NextFieldValueString(1, vdl.StringType, x.Acceptor); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *PeerKey) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = PeerKey{}
	if err := dec.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct1 {
			index = vdlTypeStruct1.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			switch value, err := dec.ReadValueString(); {
			case err != nil:
				return err
			default:
				x.Dialer = value
			}
		case 1:
			switch value, err := dec.ReadValueString(); {
			case err != nil:
				return err
			default:
				x.Acceptor = value
			}
		}
	}
}

// PeerBehavior specifies characteristics of a connection.
type PeerBehavior struct {
	// Reachable specifies whether the outgoing or incoming connection can be
	// dialed or accepted.
	// TODO(suharshs): Make this a user defined error which vine will return instead of a bool.
	Reachable bool
	// Discoverable specifies whether the Dialer can advertise a discovery packet
	// to the Acceptor. This is useful for emulating neighborhoods.
	// TODO(suharshs): Discoverable should always be bidirectional. It is unrealistic for
	// A to discover B, but not vice versa.
	Discoverable bool
}

func (PeerBehavior) VDLReflect(struct {
	Name string `vdl:"v.io/x/ref/runtime/protocols/vine.PeerBehavior"`
}) {
}

func (x PeerBehavior) VDLIsZero() bool { //nolint:gocyclo
	return x == PeerBehavior{}
}

func (x PeerBehavior) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct2); err != nil {
		return err
	}
	if x.Reachable {
		if err := enc.NextFieldValueBool(0, vdl.BoolType, x.Reachable); err != nil {
			return err
		}
	}
	if x.Discoverable {
		if err := enc.NextFieldValueBool(1, vdl.BoolType, x.Discoverable); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *PeerBehavior) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = PeerBehavior{}
	if err := dec.StartValue(vdlTypeStruct2); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct2 {
			index = vdlTypeStruct2.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			switch value, err := dec.ReadValueBool(); {
			case err != nil:
				return err
			default:
				x.Reachable = value
			}
		case 1:
			switch value, err := dec.ReadValueBool(); {
			case err != nil:
				return err
			default:
				x.Discoverable = value
			}
		}
	}
}

//////////////////////////////////////////////////
// Error definitions

var (
	ErrInvalidAddress       = verror.NewIDAction("v.io/x/ref/runtime/protocols/vine.InvalidAddress", verror.NoRetry)
	ErrAddressNotReachable  = verror.NewIDAction("v.io/x/ref/runtime/protocols/vine.AddressNotReachable", verror.NoRetry)
	ErrNoRegisteredProtocol = verror.NewIDAction("v.io/x/ref/runtime/protocols/vine.NoRegisteredProtocol", verror.NoRetry)
	ErrCantAcceptFromTag    = verror.NewIDAction("v.io/x/ref/runtime/protocols/vine.CantAcceptFromTag", verror.NoRetry)
)

// ErrorfErrInvalidAddress calls ErrInvalidAddress.Errorf with the supplied arguments.
func ErrorfErrInvalidAddress(ctx *context.T, format string, address string) error {
	return ErrInvalidAddress.Errorf(ctx, format, address)
}

// MessageErrInvalidAddress calls ErrInvalidAddress.Message with the supplied arguments.
func MessageErrInvalidAddress(ctx *context.T, message string, address string) error {
	return ErrInvalidAddress.Message(ctx, message, address)
}

// ParamsErrInvalidAddress extracts the expected parameters from the error's ParameterList.
func ParamsErrInvalidAddress(argumentError error) (verrorComponent string, verrorOperation string, address string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	var (
		tmp interface{}
		ok  bool
	)
	tmp, returnErr = iter.next()
	if address, ok = tmp.(string); !ok {
		if returnErr != nil {
			return
		}
		returnErr = fmt.Errorf("parameter list contains the wrong type for return value address, has %T and not string", tmp)
		return
	}

	return
}

// ErrorfErrAddressNotReachable calls ErrAddressNotReachable.Errorf with the supplied arguments.
func ErrorfErrAddressNotReachable(ctx *context.T, format string, address string) error {
	return ErrAddressNotReachable.Errorf(ctx, format, address)
}

// MessageErrAddressNotReachable calls ErrAddressNotReachable.Message with the supplied arguments.
func MessageErrAddressNotReachable(ctx *context.T, message string, address string) error {
	return ErrAddressNotReachable.Message(ctx, message, address)
}

// ParamsErrAddressNotReachable extracts the expected parameters from the error's ParameterList.
func ParamsErrAddressNotReachable(argumentError error) (verrorComponent string, verrorOperation string, address string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	var (
		tmp interface{}
		ok  bool
	)
	tmp, returnErr = iter.next()
	if address, ok = tmp.(string); !ok {
		if returnErr != nil {
			return
		}
		returnErr = fmt.Errorf("parameter list contains the wrong type for return value address, has %T and not string", tmp)
		return
	}

	return
}

// ErrorfErrNoRegisteredProtocol calls ErrNoRegisteredProtocol.Errorf with the supplied arguments.
func ErrorfErrNoRegisteredProtocol(ctx *context.T, format string, protocol string) error {
	return ErrNoRegisteredProtocol.Errorf(ctx, format, protocol)
}

// MessageErrNoRegisteredProtocol calls ErrNoRegisteredProtocol.Message with the supplied arguments.
func MessageErrNoRegisteredProtocol(ctx *context.T, message string, protocol string) error {
	return ErrNoRegisteredProtocol.Message(ctx, message, protocol)
}

// ParamsErrNoRegisteredProtocol extracts the expected parameters from the error's ParameterList.
func ParamsErrNoRegisteredProtocol(argumentError error) (verrorComponent string, verrorOperation string, protocol string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	var (
		tmp interface{}
		ok  bool
	)
	tmp, returnErr = iter.next()
	if protocol, ok = tmp.(string); !ok {
		if returnErr != nil {
			return
		}
		returnErr = fmt.Errorf("parameter list contains the wrong type for return value protocol, has %T and not string", tmp)
		return
	}

	return
}

// ErrorfErrCantAcceptFromTag calls ErrCantAcceptFromTag.Errorf with the supplied arguments.
func ErrorfErrCantAcceptFromTag(ctx *context.T, format string, tag string) error {
	return ErrCantAcceptFromTag.Errorf(ctx, format, tag)
}

// MessageErrCantAcceptFromTag calls ErrCantAcceptFromTag.Message with the supplied arguments.
func MessageErrCantAcceptFromTag(ctx *context.T, message string, tag string) error {
	return ErrCantAcceptFromTag.Message(ctx, message, tag)
}

// ParamsErrCantAcceptFromTag extracts the expected parameters from the error's ParameterList.
func ParamsErrCantAcceptFromTag(argumentError error) (verrorComponent string, verrorOperation string, tag string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	var (
		tmp interface{}
		ok  bool
	)
	tmp, returnErr = iter.next()
	if tag, ok = tmp.(string); !ok {
		if returnErr != nil {
			return
		}
		returnErr = fmt.Errorf("parameter list contains the wrong type for return value tag, has %T and not string", tmp)
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

//////////////////////////////////////////////////
// Interface definitions

// VineClientMethods is the client interface
// containing Vine methods.
//
// Vine is the interface to a vine service that can dynamically change the network
// behavior of connection's on the vine service's process.
type VineClientMethods interface {
	// SetBehaviors sets the policy that the accepting vine service's process
	// will use on connections.
	// behaviors is a map from server tag to the desired connection behavior.
	// For example,
	//   client.SetBehaviors(map[PeerKey]PeerBehavior{PeerKey{"foo", "bar"}, PeerBehavior{Reachable: false}})
	// will cause all vine protocol dial calls from "foo" to "bar" to fail.
	SetBehaviors(_ *context.T, behaviors map[PeerKey]PeerBehavior, _ ...rpc.CallOpt) error
}

// VineClientStub embeds VineClientMethods and is a
// placeholder for additional management operations.
type VineClientStub interface {
	VineClientMethods
}

// VineClient returns a client stub for Vine.
func VineClient(name string) VineClientStub {
	return implVineClientStub{name}
}

type implVineClientStub struct {
	name string
}

func (c implVineClientStub) SetBehaviors(ctx *context.T, i0 map[PeerKey]PeerBehavior, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "SetBehaviors", []interface{}{i0}, nil, opts...)
	return
}

// VineServerMethods is the interface a server writer
// implements for Vine.
//
// Vine is the interface to a vine service that can dynamically change the network
// behavior of connection's on the vine service's process.
type VineServerMethods interface {
	// SetBehaviors sets the policy that the accepting vine service's process
	// will use on connections.
	// behaviors is a map from server tag to the desired connection behavior.
	// For example,
	//   client.SetBehaviors(map[PeerKey]PeerBehavior{PeerKey{"foo", "bar"}, PeerBehavior{Reachable: false}})
	// will cause all vine protocol dial calls from "foo" to "bar" to fail.
	SetBehaviors(_ *context.T, _ rpc.ServerCall, behaviors map[PeerKey]PeerBehavior) error
}

// VineServerStubMethods is the server interface containing
// Vine methods, as expected by rpc.Server.
// There is no difference between this interface and VineServerMethods
// since there are no streaming methods.
type VineServerStubMethods VineServerMethods

// VineServerStub adds universal methods to VineServerStubMethods.
type VineServerStub interface {
	VineServerStubMethods
	// DescribeInterfaces the Vine interfaces.
	Describe__() []rpc.InterfaceDesc
}

// VineServer returns a server stub for Vine.
// It converts an implementation of VineServerMethods into
// an object that may be used by rpc.Server.
func VineServer(impl VineServerMethods) VineServerStub {
	stub := implVineServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implVineServerStub struct {
	impl VineServerMethods
	gs   *rpc.GlobState
}

func (s implVineServerStub) SetBehaviors(ctx *context.T, call rpc.ServerCall, i0 map[PeerKey]PeerBehavior) error {
	return s.impl.SetBehaviors(ctx, call, i0)
}

func (s implVineServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implVineServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{VineDesc}
}

// VineDesc describes the Vine interface.
var VineDesc rpc.InterfaceDesc = descVine

// descVine hides the desc to keep godoc clean.
var descVine = rpc.InterfaceDesc{
	Name:    "Vine",
	PkgPath: "v.io/x/ref/runtime/protocols/vine",
	Doc:     "// Vine is the interface to a vine service that can dynamically change the network\n// behavior of connection's on the vine service's process.",
	Methods: []rpc.MethodDesc{
		{
			Name: "SetBehaviors",
			Doc:  "// SetBehaviors sets the policy that the accepting vine service's process\n// will use on connections.\n// behaviors is a map from server tag to the desired connection behavior.\n// For example,\n//   client.SetBehaviors(map[PeerKey]PeerBehavior{PeerKey{\"foo\", \"bar\"}, PeerBehavior{Reachable: false}})\n// will cause all vine protocol dial calls from \"foo\" to \"bar\" to fail.",
			InArgs: []rpc.ArgDesc{
				{Name: "behaviors", Doc: ``}, // map[PeerKey]PeerBehavior
			},
		},
	},
}

// Hold type definitions in package-level variables, for better performance.
//nolint:unused
var (
	vdlTypeStruct1 *vdl.Type
	vdlTypeStruct2 *vdl.Type
)

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

	// Register types.
	vdl.Register((*PeerKey)(nil))
	vdl.Register((*PeerBehavior)(nil))

	// Initialize type definitions.
	vdlTypeStruct1 = vdl.TypeOf((*PeerKey)(nil)).Elem()
	vdlTypeStruct2 = vdl.TypeOf((*PeerBehavior)(nil)).Elem()

	return struct{}{}
}

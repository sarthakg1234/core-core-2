// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: naming

//nolint:golint
package naming

import (
	"fmt"

	"v.io/v23/vdl"
	vdltime "v.io/v23/vdlroot/time"
	"v.io/v23/verror"
)

var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

//////////////////////////////////////////////////
// Type definitions

// MountFlag is a bit mask of options to the mount call.
type MountFlag uint32

func (MountFlag) VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.MountFlag"`
}) {
}

func (x MountFlag) VDLIsZero() bool { //nolint:gocyclo
	return x == 0
}

func (x MountFlag) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.WriteValueUint(vdlTypeUint321, uint64(x)); err != nil {
		return err
	}
	return nil
}

func (x *MountFlag) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	switch value, err := dec.ReadValueUint(32); {
	case err != nil:
		return err
	default:
		*x = MountFlag(value)
	}
	return nil
}

// MountedServer represents a server mounted on a specific name.
type MountedServer struct {
	// Server is the OA that's mounted.
	Server string
	// Deadline before the mount entry expires.
	Deadline vdltime.Deadline
}

func (MountedServer) VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.MountedServer"`
}) {
}

func (x MountedServer) VDLIsZero() bool { //nolint:gocyclo
	if x.Server != "" {
		return false
	}
	if !x.Deadline.Time.IsZero() {
		return false
	}
	return true
}

func (x MountedServer) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct2); err != nil {
		return err
	}
	if x.Server != "" {
		if err := enc.NextFieldValueString(0, vdl.StringType, x.Server); err != nil {
			return err
		}
	}
	if !x.Deadline.Time.IsZero() {
		if err := enc.NextField(1); err != nil {
			return err
		}
		var wire vdltime.WireDeadline
		if err := vdltime.WireDeadlineFromNative(&wire, x.Deadline); err != nil {
			return err
		}
		if err := wire.VDLWrite(enc); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *MountedServer) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = MountedServer{}
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
			switch value, err := dec.ReadValueString(); {
			case err != nil:
				return err
			default:
				x.Server = value
			}
		case 1:
			var wire vdltime.WireDeadline
			if err := wire.VDLRead(dec); err != nil {
				return err
			}
			if err := vdltime.WireDeadlineToNative(wire, &x.Deadline); err != nil {
				return err
			}
		}
	}
}

// MountEntry represents a given name mounted in the mounttable.
type MountEntry struct {
	// Name is the mounted name.
	Name string
	// Servers (if present) specifies the mounted names.
	Servers []MountedServer
	// ServesMountTable is true if the servers represent mount tables.
	ServesMountTable bool
	// IsLeaf is true if this entry represents a leaf object.
	IsLeaf bool
}

func (MountEntry) VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.MountEntry"`
}) {
}

func (x MountEntry) VDLIsZero() bool { //nolint:gocyclo
	if x.Name != "" {
		return false
	}
	if len(x.Servers) != 0 {
		return false
	}
	if x.ServesMountTable {
		return false
	}
	if x.IsLeaf {
		return false
	}
	return true
}

func (x MountEntry) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct4); err != nil {
		return err
	}
	if x.Name != "" {
		if err := enc.NextFieldValueString(0, vdl.StringType, x.Name); err != nil {
			return err
		}
	}
	if len(x.Servers) != 0 {
		if err := enc.NextField(1); err != nil {
			return err
		}
		if err := vdlWriteAnonList1(enc, x.Servers); err != nil {
			return err
		}
	}
	if x.ServesMountTable {
		if err := enc.NextFieldValueBool(2, vdl.BoolType, x.ServesMountTable); err != nil {
			return err
		}
	}
	if x.IsLeaf {
		if err := enc.NextFieldValueBool(3, vdl.BoolType, x.IsLeaf); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func vdlWriteAnonList1(enc vdl.Encoder, x []MountedServer) error {
	if err := enc.StartValue(vdlTypeList5); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for _, elem := range x {
		if err := enc.NextEntry(false); err != nil {
			return err
		}
		if err := elem.VDLWrite(enc); err != nil {
			return err
		}
	}
	if err := enc.NextEntry(true); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *MountEntry) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = MountEntry{}
	if err := dec.StartValue(vdlTypeStruct4); err != nil {
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
		if decType != vdlTypeStruct4 {
			index = vdlTypeStruct4.FieldIndexByName(decType.Field(index).Name)
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
				x.Name = value
			}
		case 1:
			if err := vdlReadAnonList1(dec, &x.Servers); err != nil {
				return err
			}
		case 2:
			switch value, err := dec.ReadValueBool(); {
			case err != nil:
				return err
			default:
				x.ServesMountTable = value
			}
		case 3:
			switch value, err := dec.ReadValueBool(); {
			case err != nil:
				return err
			default:
				x.IsLeaf = value
			}
		}
	}
}

func vdlReadAnonList1(dec vdl.Decoder, x *[]MountedServer) error {
	if err := dec.StartValue(vdlTypeList5); err != nil {
		return err
	}
	if len := dec.LenHint(); len > 0 {
		*x = make([]MountedServer, 0, len)
	} else {
		*x = nil
	}
	for {
		switch done, err := dec.NextEntry(); {
		case err != nil:
			return err
		case done:
			return dec.FinishValue()
		default:
			var elem MountedServer
			if err := elem.VDLRead(dec); err != nil {
				return err
			}
			*x = append(*x, elem)
		}
	}
}

// GlobError is returned by namespace.Glob to indicate a subtree of the namespace
// that could not be traversed.
type GlobError struct {
	// Root of the subtree.
	Name string
	// The error that occurred fulfilling the request.
	Error error
}

func (GlobError) VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.GlobError"`
}) {
}

func (x GlobError) VDLIsZero() bool { //nolint:gocyclo
	return x == GlobError{}
}

func (x GlobError) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct6); err != nil {
		return err
	}
	if x.Name != "" {
		if err := enc.NextFieldValueString(0, vdl.StringType, x.Name); err != nil {
			return err
		}
	}
	if x.Error != nil {
		if err := enc.NextField(1); err != nil {
			return err
		}
		if err := verror.VDLWrite(enc, x.Error); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *GlobError) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = GlobError{}
	if err := dec.StartValue(vdlTypeStruct6); err != nil {
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
		if decType != vdlTypeStruct6 {
			index = vdlTypeStruct6.FieldIndexByName(decType.Field(index).Name)
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
				x.Name = value
			}
		case 1:
			if err := verror.VDLRead(dec, &x.Error); err != nil {
				return err
			}
		}
	}
}

type (
	// GlobReply represents any single field of the GlobReply union type.
	//
	// GlobReply is the data type returned by Glob__.
	GlobReply interface {
		// Index returns the field index.
		Index() int
		// Interface returns the field value as an interface.
		Interface() interface{}
		// Name returns the field name.
		Name() string
		// VDLReflect describes the GlobReply union type.
		VDLReflect(vdlGlobReplyReflect)
		VDLIsZero() bool
		VDLWrite(vdl.Encoder) error
	}
	// GlobReplyEntry represents field Entry of the GlobReply union type.
	GlobReplyEntry struct{ Value MountEntry }
	// GlobReplyError represents field Error of the GlobReply union type.
	GlobReplyError struct{ Value GlobError }
	// vdlGlobReplyReflect describes the GlobReply union type.
	vdlGlobReplyReflect struct {
		Name  string `vdl:"v.io/v23/naming.GlobReply"`
		Type  GlobReply
		Union struct {
			Entry GlobReplyEntry
			Error GlobReplyError
		}
	}
)

func (x GlobReplyEntry) Index() int                     { return 0 }
func (x GlobReplyEntry) Interface() interface{}         { return x.Value }
func (x GlobReplyEntry) Name() string                   { return "Entry" }
func (x GlobReplyEntry) VDLReflect(vdlGlobReplyReflect) {}

func (x GlobReplyError) Index() int                     { return 1 }
func (x GlobReplyError) Interface() interface{}         { return x.Value }
func (x GlobReplyError) Name() string                   { return "Error" }
func (x GlobReplyError) VDLReflect(vdlGlobReplyReflect) {}

func (x GlobReplyEntry) VDLIsZero() bool { //nolint:gocyclo
	return x.Value.VDLIsZero()
}

func (x GlobReplyError) VDLIsZero() bool {
	return false
}

func (x GlobReplyEntry) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeUnion7); err != nil {
		return err
	}
	if err := enc.NextField(0); err != nil {
		return err
	}
	if err := x.Value.VDLWrite(enc); err != nil {
		return err
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x GlobReplyError) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeUnion7); err != nil {
		return err
	}
	if err := enc.NextField(1); err != nil {
		return err
	}
	if err := x.Value.VDLWrite(enc); err != nil {
		return err
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func VDLReadGlobReply(dec vdl.Decoder, x *GlobReply) error { //nolint:gocyclo
	if err := dec.StartValue(vdlTypeUnion7); err != nil {
		return err
	}
	decType := dec.Type()
	index, err := dec.NextField()
	switch {
	case err != nil:
		return err
	case index == -1:
		return fmt.Errorf("missing field in union %T, from %v", x, decType)
	}
	if decType != vdlTypeUnion7 {
		name := decType.Field(index).Name
		index = vdlTypeUnion7.FieldIndexByName(name)
		if index == -1 {
			return fmt.Errorf("field %q not in union %T, from %v", name, x, decType)
		}
	}
	switch index {
	case 0:
		var field GlobReplyEntry
		if err := field.Value.VDLRead(dec); err != nil {
			return err
		}
		*x = field
	case 1:
		var field GlobReplyError
		if err := field.Value.VDLRead(dec); err != nil {
			return err
		}
		*x = field
	}
	switch index, err := dec.NextField(); {
	case err != nil:
		return err
	case index != -1:
		return fmt.Errorf("extra field %d in union %T, from %v", index, x, dec.Type())
	}
	return dec.FinishValue()
}

type (
	// GlobChildrenReply represents any single field of the GlobChildrenReply union type.
	//
	// GlobChildrenReply is the data type returned by GlobChildren__.
	GlobChildrenReply interface {
		// Index returns the field index.
		Index() int
		// Interface returns the field value as an interface.
		Interface() interface{}
		// Name returns the field name.
		Name() string
		// VDLReflect describes the GlobChildrenReply union type.
		VDLReflect(vdlGlobChildrenReplyReflect)
		VDLIsZero() bool
		VDLWrite(vdl.Encoder) error
	}
	// GlobChildrenReplyName represents field Name of the GlobChildrenReply union type.
	GlobChildrenReplyName struct{ Value string }
	// GlobChildrenReplyError represents field Error of the GlobChildrenReply union type.
	GlobChildrenReplyError struct{ Value GlobError }
	// vdlGlobChildrenReplyReflect describes the GlobChildrenReply union type.
	vdlGlobChildrenReplyReflect struct {
		Name  string `vdl:"v.io/v23/naming.GlobChildrenReply"`
		Type  GlobChildrenReply
		Union struct {
			Name  GlobChildrenReplyName
			Error GlobChildrenReplyError
		}
	}
)

func (x GlobChildrenReplyName) Index() int                             { return 0 }
func (x GlobChildrenReplyName) Interface() interface{}                 { return x.Value }
func (x GlobChildrenReplyName) Name() string                           { return "Name" }
func (x GlobChildrenReplyName) VDLReflect(vdlGlobChildrenReplyReflect) {}

func (x GlobChildrenReplyError) Index() int                             { return 1 }
func (x GlobChildrenReplyError) Interface() interface{}                 { return x.Value }
func (x GlobChildrenReplyError) Name() string                           { return "Error" }
func (x GlobChildrenReplyError) VDLReflect(vdlGlobChildrenReplyReflect) {}

func (x GlobChildrenReplyName) VDLIsZero() bool { //nolint:gocyclo
	return x.Value == ""
}

func (x GlobChildrenReplyError) VDLIsZero() bool {
	return false
}

func (x GlobChildrenReplyName) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeUnion8); err != nil {
		return err
	}
	if err := enc.NextFieldValueString(0, vdl.StringType, x.Value); err != nil {
		return err
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x GlobChildrenReplyError) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeUnion8); err != nil {
		return err
	}
	if err := enc.NextField(1); err != nil {
		return err
	}
	if err := x.Value.VDLWrite(enc); err != nil {
		return err
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func VDLReadGlobChildrenReply(dec vdl.Decoder, x *GlobChildrenReply) error { //nolint:gocyclo
	if err := dec.StartValue(vdlTypeUnion8); err != nil {
		return err
	}
	decType := dec.Type()
	index, err := dec.NextField()
	switch {
	case err != nil:
		return err
	case index == -1:
		return fmt.Errorf("missing field in union %T, from %v", x, decType)
	}
	if decType != vdlTypeUnion8 {
		name := decType.Field(index).Name
		index = vdlTypeUnion8.FieldIndexByName(name)
		if index == -1 {
			return fmt.Errorf("field %q not in union %T, from %v", name, x, decType)
		}
	}
	switch index {
	case 0:
		var field GlobChildrenReplyName
		switch value, err := dec.ReadValueString(); {
		case err != nil:
			return err
		default:
			field.Value = value
		}
		*x = field
	case 1:
		var field GlobChildrenReplyError
		if err := field.Value.VDLRead(dec); err != nil {
			return err
		}
		*x = field
	}
	switch index, err := dec.NextField(); {
	case err != nil:
		return err
	case index != -1:
		return fmt.Errorf("extra field %d in union %T, from %v", index, x, dec.Type())
	}
	return dec.FinishValue()
}

//////////////////////////////////////////////////
// Const definitions

const Replace = MountFlag(1) // Replace means the mount should replace what is currently at the mount point
const MT = MountFlag(2)      // MT means that the target server is a mount table.
const Leaf = MountFlag(4)    // Leaf means that the target server is a leaf.

// Hold type definitions in package-level variables, for better performance.
//nolint:unused
var (
	vdlTypeUint321 *vdl.Type
	vdlTypeStruct2 *vdl.Type
	vdlTypeStruct3 *vdl.Type
	vdlTypeStruct4 *vdl.Type
	vdlTypeList5   *vdl.Type
	vdlTypeStruct6 *vdl.Type
	vdlTypeUnion7  *vdl.Type
	vdlTypeUnion8  *vdl.Type
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
	vdl.Register((*MountFlag)(nil))
	vdl.Register((*MountedServer)(nil))
	vdl.Register((*MountEntry)(nil))
	vdl.Register((*GlobError)(nil))
	vdl.Register((*GlobReply)(nil))
	vdl.Register((*GlobChildrenReply)(nil))

	// Initialize type definitions.
	vdlTypeUint321 = vdl.TypeOf((*MountFlag)(nil))
	vdlTypeStruct2 = vdl.TypeOf((*MountedServer)(nil)).Elem()
	vdlTypeStruct3 = vdl.TypeOf((*vdltime.WireDeadline)(nil)).Elem()
	vdlTypeStruct4 = vdl.TypeOf((*MountEntry)(nil)).Elem()
	vdlTypeList5 = vdl.TypeOf((*[]MountedServer)(nil))
	vdlTypeStruct6 = vdl.TypeOf((*GlobError)(nil)).Elem()
	vdlTypeUnion7 = vdl.TypeOf((*GlobReply)(nil))
	vdlTypeUnion8 = vdl.TypeOf((*GlobChildrenReply)(nil))

	return struct{}{}
}

// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: access

// Package access defines types and interfaces for dynamic access control.
// Examples: "allow app to read this photo", "prevent user from modifying this
// file".
//
// Target Developers
//
// Developers creating functionality to share data or services between
// multiple users/devices/apps.
//
// Overview
//
// Vanadium objects provide GetPermissions and SetPermissions methods.  An
// AccessList contains the set of blessings that grant principals access to the
// object. All methods on objects can have "tags" on them and the AccessList
// used for the method is selected based on that tag (from a Permissions).
//
// An object can have multiple names, so GetPermissions and SetPermissions can
// be invoked on any of these names, but the object itself has a single
// AccessList.
//
// SetPermissions completely replaces the Permissions. To perform an atomic
// read-modify-write of the AccessList, use the version parameter.
//
// Conventions
//
// Service implementors should follow the conventions below to be consistent
// with other parts of Vanadium and with each other.
//
// All methods that create an object (e.g. Put, Mount, Link) should take an
// optional AccessList parameter.  If the AccessList is not specified, the new
// object, O, copies its AccessList from the parent.  Subsequent changes to the
// parent AccessList are not automatically propagated to O.  Instead, a client
// library must make recursive AccessList changes.
//
// Resolve access is required on all components of a name, except the last one,
// in order to access the object referenced by that name.  For example, for
// principal P to access the name "a/b/c", P must have resolve access to "a"
// and "a/b".
//
// The Resolve tag means that a principal can traverse that component of the
// name to access the child.  It does not give the principal permission to list
// the children via Glob or a similar method.  For example, a server might have
// an object named "home" with a child for each user of the system.  If these
// users were allowed to list the contents of "home", they could discover the
// other users of the system.  That could be a privacy violation.  Without
// Resolve, every user of the system would need read access to "home" to access
// "home/<user>".  If the user called Glob("home/*"), it would then be up to
// the server to filter out the names that the user could not access.  That
// could be a very expensive operation if there were a lot of children of
// "home".  Resolve protects these servers against potential denial of service
// attacks on these large, shared directories.
//
// Blessings allow for sweeping access changes. In particular, a blessing is
// useful for controlling access to objects that are always accessed together.
// For example, a document may have embedded images and comments, each with a
// unique name. When accessing a document, the server would generate a blessing
// that the client would use to fetch the images and comments; the images and
// comments would have this blessed identity in their AccessLists. Changes to
// the document's AccessLists are therefore "propagated" to the images and
// comments.
//
// In the future, we may add some sort of "groups" mechanism to provide an
// alternative way to express access control policies.
//
// Some services will want a concept of implicit access control. They are free
// to implement this as appropriate for their service. However, GetPermissions
// should respond with the correct Permissions. For example, a corporate file
// server would allow all employees to create their own directory and have full
// control within that directory. Employees should not be allowed to modify
// other employee directories. In other words, within the directory "home",
// employee E should be allowed to modify only "home/E". The file server doesn't
// know the list of all employees a priori, so it uses an
// implementation-specific rule to map employee identities to their home
// directory.
//
// Examples
//
//   client := access.ObjectClient(name)
//   for {
//     perms, version, err := client.GetPermissions()
//     if err != nil {
//       return err
//     }
//     perms[newTag] = AccessList{In: []security.BlessingPattern{newPattern}}
//     // Use the same version with the modified perms to ensure that no other
//     // client has modified the perms since GetPermissions returned.
//     if err := client.SetPermissions(perms, version); err != nil {
//       if errors.Is(err, verror.ErrBadVersion) {
//         // Another client replaced the Permissions after our GetPermissions
//         // returned. Try again.
//         continue
//       }
//       return err
//     }
//   }
//nolint:golint
package access

import (
	"v.io/v23/context"
	"v.io/v23/i18n"
	"v.io/v23/security"
	"v.io/v23/uniqueid"
	"v.io/v23/vdl"
	"v.io/v23/verror"
)

var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

//////////////////////////////////////////////////
// Type definitions

// AccessList represents a set of blessings that should be granted access.
//
// See also: https://vanadium.github.io/glossary.html#access-list
type AccessList struct {
	// In denotes the set of blessings (represented as BlessingPatterns) that
	// should be granted access, unless blacklisted by an entry in NotIn.
	//
	// For example:
	//   In: {"alice:family"}
	// grants access to a principal that presents at least one of
	// "alice:family", "alice:family:friend", "alice:family:friend:spouse" etc.
	// as a blessing.
	In []security.BlessingPattern
	// NotIn denotes the set of blessings (and their delegates) that
	// have been explicitly blacklisted from the In set.
	//
	// For example:
	//   In: {"alice:friend"}, NotIn: {"alice:friend:bob"}
	// grants access to principals that present "alice:friend",
	// "alice:friend:carol" etc. but NOT to a principal that presents
	// "alice:friend:bob" or "alice:friend:bob:spouse" etc.
	NotIn []string
}

func (AccessList) VDLReflect(struct {
	Name string `vdl:"v.io/v23/security/access.AccessList"`
}) {
}

func (x AccessList) VDLIsZero() bool { //nolint:gocyclo
	if len(x.In) != 0 {
		return false
	}
	if len(x.NotIn) != 0 {
		return false
	}
	return true
}

func (x AccessList) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	if len(x.In) != 0 {
		if err := enc.NextField(0); err != nil {
			return err
		}
		if err := vdlWriteAnonList1(enc, x.In); err != nil {
			return err
		}
	}
	if len(x.NotIn) != 0 {
		if err := enc.NextField(1); err != nil {
			return err
		}
		if err := vdlWriteAnonList2(enc, x.NotIn); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func vdlWriteAnonList1(enc vdl.Encoder, x []security.BlessingPattern) error {
	if err := enc.StartValue(vdlTypeList2); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for _, elem := range x {
		if err := enc.NextEntryValueString(vdlTypeString4, string(elem)); err != nil {
			return err
		}
	}
	if err := enc.NextEntry(true); err != nil {
		return err
	}
	return enc.FinishValue()
}

func vdlWriteAnonList2(enc vdl.Encoder, x []string) error {
	if err := enc.StartValue(vdlTypeList3); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for _, elem := range x {
		if err := enc.NextEntryValueString(vdl.StringType, elem); err != nil {
			return err
		}
	}
	if err := enc.NextEntry(true); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *AccessList) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = AccessList{}
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
			if err := vdlReadAnonList1(dec, &x.In); err != nil {
				return err
			}
		case 1:
			if err := vdlReadAnonList2(dec, &x.NotIn); err != nil {
				return err
			}
		}
	}
}

func vdlReadAnonList1(dec vdl.Decoder, x *[]security.BlessingPattern) error {
	if err := dec.StartValue(vdlTypeList2); err != nil {
		return err
	}
	if len := dec.LenHint(); len > 0 {
		*x = make([]security.BlessingPattern, 0, len)
	} else {
		*x = nil
	}
	for {
		switch done, elem, err := dec.NextEntryValueString(); {
		case err != nil:
			return err
		case done:
			return dec.FinishValue()
		default:
			*x = append(*x, security.BlessingPattern(elem))
		}
	}
}

func vdlReadAnonList2(dec vdl.Decoder, x *[]string) error {
	if err := dec.StartValue(vdlTypeList3); err != nil {
		return err
	}
	if len := dec.LenHint(); len > 0 {
		*x = make([]string, 0, len)
	} else {
		*x = nil
	}
	for {
		switch done, elem, err := dec.NextEntryValueString(); {
		case err != nil:
			return err
		case done:
			return dec.FinishValue()
		default:
			*x = append(*x, elem)
		}
	}
}

// Permissions maps string tags to access lists specifying the blessings
// required to invoke methods with that tag.
//
// These tags are meant to add a layer of interposition between the set of
// users (blessings, specifically) and the set of methods, much like "Roles" do
// in Role Based Access Control.
// (http://en.wikipedia.org/wiki/Role-based_access_control)
type Permissions map[string]AccessList

func (Permissions) VDLReflect(struct {
	Name string `vdl:"v.io/v23/security/access.Permissions"`
}) {
}

func (x Permissions) VDLIsZero() bool { //nolint:gocyclo
	return len(x) == 0
}

func (x Permissions) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeMap5); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for key, elem := range x {
		if err := enc.NextEntryValueString(vdl.StringType, key); err != nil {
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

func (x *Permissions) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	if err := dec.StartValue(vdlTypeMap5); err != nil {
		return err
	}
	var tmpMap Permissions
	if len := dec.LenHint(); len > 0 {
		tmpMap = make(Permissions, len)
	}
	for {
		switch done, key, err := dec.NextEntryValueString(); {
		case err != nil:
			return err
		case done:
			*x = tmpMap
			return dec.FinishValue()
		default:
			var elem AccessList
			if err := elem.VDLRead(dec); err != nil {
				return err
			}
			if tmpMap == nil {
				tmpMap = make(Permissions)
			}
			tmpMap[key] = elem
		}
	}
}

// Tag is used to associate methods with an AccessList in a Permissions.
//
// While services can define their own tag type and values, many
// services should be able to use the type and values defined in
// this package.
type Tag string

func (Tag) VDLReflect(struct {
	Name string `vdl:"v.io/v23/security/access.Tag"`
}) {
}

func (x Tag) VDLIsZero() bool { //nolint:gocyclo
	return x == ""
}

func (x Tag) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.WriteValueString(vdlTypeString6, string(x)); err != nil {
		return err
	}
	return nil
}

func (x *Tag) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	switch value, err := dec.ReadValueString(); {
	case err != nil:
		return err
	default:
		*x = Tag(value)
	}
	return nil
}

//////////////////////////////////////////////////
// Const definitions

const Admin = Tag("Admin")     // Operations that require privileged access for object administration.
const Debug = Tag("Debug")     // Operations that return debugging information (e.g., logs, statistics etc.) about the object.
const Read = Tag("Read")       // Operations that do not mutate the state of the object.
const Write = Tag("Write")     // Operations that mutate the state of the object.
const Resolve = Tag("Resolve") // Operations involving namespace navigation.
// AccessTagCaveat represents a caveat that validates iff the method being invoked has
// at least one of the tags listed in the caveat.
var AccessTagCaveat = security.CaveatDescriptor{
	Id: uniqueid.Id{
		239,
		205,
		227,
		117,
		20,
		22,
		199,
		59,
		24,
		156,
		232,
		156,
		204,
		147,
		128,
		0,
	},
	ParamType: vdl.TypeOf((*[]Tag)(nil)),
}

//////////////////////////////////////////////////
// Error definitions

var (

	// The AccessList is too big.  Use groups to represent large sets of principals.
	ErrTooBig                    = verror.NewIDAction("v.io/v23/security/access.TooBig", verror.NoRetry)
	ErrNoPermissions             = verror.NewIDAction("v.io/v23/security/access.NoPermissions", verror.NoRetry)
	ErrAccessListMatch           = verror.NewIDAction("v.io/v23/security/access.AccessListMatch", verror.NoRetry)
	ErrUnenforceablePatterns     = verror.NewIDAction("v.io/v23/security/access.UnenforceablePatterns", verror.NoRetry)
	ErrInvalidOpenAccessList     = verror.NewIDAction("v.io/v23/security/access.InvalidOpenAccessList", verror.NoRetry)
	ErrAccessTagCaveatValidation = verror.NewIDAction("v.io/v23/security/access.AccessTagCaveatValidation", verror.NoRetry)
	ErrMultipleTags              = verror.NewIDAction("v.io/v23/security/access.MultipleTags", verror.NoRetry)
	ErrNoTags                    = verror.NewIDAction("v.io/v23/security/access.NoTags", verror.NoRetry)
)

// NewErrTooBig returns an error with the ErrTooBig ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfTooBig or MessageTooBig instead.
func NewErrTooBig(ctx *context.T) error {
	return verror.New(ErrTooBig, ctx)
}

// ErrorfTooBig calls ErrTooBig.Errorf with the supplied arguments.
func ErrorfTooBig(ctx *context.T, format string) error {
	return ErrTooBig.Errorf(ctx, format)
}

// MessageTooBig calls ErrTooBig.Message with the supplied arguments.
func MessageTooBig(ctx *context.T, message string) error {
	return ErrTooBig.Message(ctx, message)
}

// NewErrNoPermissions returns an error with the ErrNoPermissions ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfNoPermissions or MessageNoPermissions instead.
func NewErrNoPermissions(ctx *context.T, validBlessings []string, rejectedBlessings []security.RejectedBlessing, tag string) error {
	return verror.New(ErrNoPermissions, ctx, validBlessings, rejectedBlessings, tag)
}

// ErrorfNoPermissions calls ErrNoPermissions.Errorf with the supplied arguments.
func ErrorfNoPermissions(ctx *context.T, format string, validBlessings []string, rejectedBlessings []security.RejectedBlessing, tag string) error {
	return ErrNoPermissions.Errorf(ctx, format, validBlessings, rejectedBlessings, tag)
}

// MessageNoPermissions calls ErrNoPermissions.Message with the supplied arguments.
func MessageNoPermissions(ctx *context.T, message string, validBlessings []string, rejectedBlessings []security.RejectedBlessing, tag string) error {
	return ErrNoPermissions.Message(ctx, message, validBlessings, rejectedBlessings, tag)
}

// NewErrAccessListMatch returns an error with the ErrAccessListMatch ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfAccessListMatch or MessageAccessListMatch instead.
func NewErrAccessListMatch(ctx *context.T, validBlessings []string, rejectedBlessings []security.RejectedBlessing) error {
	return verror.New(ErrAccessListMatch, ctx, validBlessings, rejectedBlessings)
}

// ErrorfAccessListMatch calls ErrAccessListMatch.Errorf with the supplied arguments.
func ErrorfAccessListMatch(ctx *context.T, format string, validBlessings []string, rejectedBlessings []security.RejectedBlessing) error {
	return ErrAccessListMatch.Errorf(ctx, format, validBlessings, rejectedBlessings)
}

// MessageAccessListMatch calls ErrAccessListMatch.Message with the supplied arguments.
func MessageAccessListMatch(ctx *context.T, message string, validBlessings []string, rejectedBlessings []security.RejectedBlessing) error {
	return ErrAccessListMatch.Message(ctx, message, validBlessings, rejectedBlessings)
}

// NewErrUnenforceablePatterns returns an error with the ErrUnenforceablePatterns ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfUnenforceablePatterns or MessageUnenforceablePatterns instead.
func NewErrUnenforceablePatterns(ctx *context.T, rejectedPatterns []security.BlessingPattern) error {
	return verror.New(ErrUnenforceablePatterns, ctx, rejectedPatterns)
}

// ErrorfUnenforceablePatterns calls ErrUnenforceablePatterns.Errorf with the supplied arguments.
func ErrorfUnenforceablePatterns(ctx *context.T, format string, rejectedPatterns []security.BlessingPattern) error {
	return ErrUnenforceablePatterns.Errorf(ctx, format, rejectedPatterns)
}

// MessageUnenforceablePatterns calls ErrUnenforceablePatterns.Message with the supplied arguments.
func MessageUnenforceablePatterns(ctx *context.T, message string, rejectedPatterns []security.BlessingPattern) error {
	return ErrUnenforceablePatterns.Message(ctx, message, rejectedPatterns)
}

// NewErrInvalidOpenAccessList returns an error with the ErrInvalidOpenAccessList ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfInvalidOpenAccessList or MessageInvalidOpenAccessList instead.
func NewErrInvalidOpenAccessList(ctx *context.T) error {
	return verror.New(ErrInvalidOpenAccessList, ctx)
}

// ErrorfInvalidOpenAccessList calls ErrInvalidOpenAccessList.Errorf with the supplied arguments.
func ErrorfInvalidOpenAccessList(ctx *context.T, format string) error {
	return ErrInvalidOpenAccessList.Errorf(ctx, format)
}

// MessageInvalidOpenAccessList calls ErrInvalidOpenAccessList.Message with the supplied arguments.
func MessageInvalidOpenAccessList(ctx *context.T, message string) error {
	return ErrInvalidOpenAccessList.Message(ctx, message)
}

// NewErrAccessTagCaveatValidation returns an error with the ErrAccessTagCaveatValidation ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfAccessTagCaveatValidation or MessageAccessTagCaveatValidation instead.
func NewErrAccessTagCaveatValidation(ctx *context.T, methodTags []string, caveatTags []Tag) error {
	return verror.New(ErrAccessTagCaveatValidation, ctx, methodTags, caveatTags)
}

// ErrorfAccessTagCaveatValidation calls ErrAccessTagCaveatValidation.Errorf with the supplied arguments.
func ErrorfAccessTagCaveatValidation(ctx *context.T, format string, methodTags []string, caveatTags []Tag) error {
	return ErrAccessTagCaveatValidation.Errorf(ctx, format, methodTags, caveatTags)
}

// MessageAccessTagCaveatValidation calls ErrAccessTagCaveatValidation.Message with the supplied arguments.
func MessageAccessTagCaveatValidation(ctx *context.T, message string, methodTags []string, caveatTags []Tag) error {
	return ErrAccessTagCaveatValidation.Message(ctx, message, methodTags, caveatTags)
}

// NewErrMultipleTags returns an error with the ErrMultipleTags ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfMultipleTags or MessageMultipleTags instead.
func NewErrMultipleTags(ctx *context.T, suffix string, method string, tag string) error {
	return verror.New(ErrMultipleTags, ctx, suffix, method, tag)
}

// ErrorfMultipleTags calls ErrMultipleTags.Errorf with the supplied arguments.
func ErrorfMultipleTags(ctx *context.T, format string, suffix string, method string, tag string) error {
	return ErrMultipleTags.Errorf(ctx, format, suffix, method, tag)
}

// MessageMultipleTags calls ErrMultipleTags.Message with the supplied arguments.
func MessageMultipleTags(ctx *context.T, message string, suffix string, method string, tag string) error {
	return ErrMultipleTags.Message(ctx, message, suffix, method, tag)
}

// NewErrNoTags returns an error with the ErrNoTags ID.
// WARNING: this function is deprecated and will be removed in the future,
// use ErrorfNoTags or MessageNoTags instead.
func NewErrNoTags(ctx *context.T, suffix string, method string, tag string) error {
	return verror.New(ErrNoTags, ctx, suffix, method, tag)
}

// ErrorfNoTags calls ErrNoTags.Errorf with the supplied arguments.
func ErrorfNoTags(ctx *context.T, format string, suffix string, method string, tag string) error {
	return ErrNoTags.Errorf(ctx, format, suffix, method, tag)
}

// MessageNoTags calls ErrNoTags.Message with the supplied arguments.
func MessageNoTags(ctx *context.T, message string, suffix string, method string, tag string) error {
	return ErrNoTags.Message(ctx, message, suffix, method, tag)
}

// Hold type definitions in package-level variables, for better performance.
//nolint:unused
var (
	vdlTypeStruct1 *vdl.Type
	vdlTypeList2   *vdl.Type
	vdlTypeList3   *vdl.Type
	vdlTypeString4 *vdl.Type
	vdlTypeMap5    *vdl.Type
	vdlTypeString6 *vdl.Type
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
	vdl.Register((*AccessList)(nil))
	vdl.Register((*Permissions)(nil))
	vdl.Register((*Tag)(nil))

	// Initialize type definitions.
	vdlTypeStruct1 = vdl.TypeOf((*AccessList)(nil)).Elem()
	vdlTypeList2 = vdl.TypeOf((*[]security.BlessingPattern)(nil))
	vdlTypeList3 = vdl.TypeOf((*[]string)(nil))
	vdlTypeString4 = vdl.TypeOf((*security.BlessingPattern)(nil))
	vdlTypeMap5 = vdl.TypeOf((*Permissions)(nil))
	vdlTypeString6 = vdl.TypeOf((*Tag)(nil))

	// Set error format strings.
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrTooBig.ID), "{1:}{2:} AccessList is too big")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNoPermissions.ID), "{1:}{2:} {3} does not have {5} access (rejected blessings: {4})")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrAccessListMatch.ID), "{1:}{2:} {3} does not match the access list (rejected blessings: {4})")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrUnenforceablePatterns.ID), "{1:}{2:} AccessList contains the following invalid or unrecognized patterns in the In list: {3}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrInvalidOpenAccessList.ID), "{1:}{2:} AccessList with the pattern ... in its In list must have no other patterns in the In or NotIn lists")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrAccessTagCaveatValidation.ID), "{1:}{2:} access tags on method ({3}) do not include any of the ones in the caveat ({4}), or the method is using a different tag type")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrMultipleTags.ID), "{1:}{2:} authorizer on {3}.{4} cannot handle multiple tags of type {5}; this is likely unintentional")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNoTags.ID), "{1:}{2:} authorizer on {3}.{4} has no tags of type {5}; this is likely unintentional")

	return struct{}{}
}

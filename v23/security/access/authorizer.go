// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package access

import (
	"bytes"
	"fmt"
	"os"

	"v.io/v23/context"
	"v.io/v23/security"
	"v.io/v23/vdl"
	"v.io/v23/verror"
)

// PermissionsAuthorizer implements an authorization policy where access is
// granted if the remote end presents blessings included in the Access Control
// Lists (AccessLists) associated with the set of relevant tags.
//
// The set of relevant tags is the subset of tags associated with the
// method (security.Call.MethodTags) that have the same type as tagType.
// Currently, tagType.Kind must be reflect.String, i.e., only tags that are
// named string types are supported.
//
// PermissionsAuthorizer expects exactly one tag of tagType to be associated
// with the method. If there are multiple, it fails authorization and returns
// an error. However, if multiple tags become a common occurrence, then this
// behavior may change.
//
// If the Permissions provided is nil, then a nil authorizer is returned.
//
// Sample usage:
//
// (1) Attach tags to methods in the VDL (eg. myservice.vdl)
//
//	package myservice
//
//	type MyTag string
//	const (
//	  ReadAccess  = MyTag("R")
//	  WriteAccess = MyTag("W")
//	)
//
//	type MyService interface {
//	  Get() ([]string, error)       {ReadAccess}
//	  GetIndex(int) (string, error) {ReadAccess}
//
//	  Set([]string) error           {WriteAccess}
//	  SetIndex(int, string) error   {WriteAccess}
//	}
//
// (2) Configure the rpc.Dispatcher to use the PermissionsAuthorizer
//
//	import (
//	  "reflect"
//
//	  "v.io/v23/rpc"
//	  "v.io/v23/security"
//	  "v.io/v23/security/access"
//	)
//
//	type dispatcher struct{}
//	func (d dispatcher) Lookup(suffix, method) (rpc.Invoker, security.Authorizer, error) {
//	   perms := access.Permissions{
//	     "R": access.AccessList{In: []security.BlessingPattern{"alice:friends", "alice:family"} },
//	     "W": access.AccessList{In: []security.BlessingPattern{"alice:family", "alice:colleagues" } },
//	   }
//	   typ := reflect.TypeOf(ReadAccess)  // equivalently, reflect.TypeOf(WriteAccess)
//	   return newInvoker(), access.PermissionsAuthorizer(perms, typ), nil
//	}
//
// With the above dispatcher, the server will grant access to a peer with the
// blessing "alice:friend:bob" access only to the "Get" and "GetIndex" methods.
// A peer presenting the blessing "alice:colleague:carol" will get access only
// to the "Set" and "SetIndex" methods. A peer presenting "alice:family:mom"
// will get access to all methods.
func PermissionsAuthorizer(perms Permissions, tagType *vdl.Type) (security.Authorizer, error) {
	if err := validateTagType(tagType); err != nil {
		return nil, err
	}
	return &authorizer{perms, tagType}, nil
}

// TypicalTagTypePermissionsAuthorizer is like PermissionsAuthorizer, but
// assumes TypicalTagType and thus avoids returning an error.
func TypicalTagTypePermissionsAuthorizer(perms Permissions) security.Authorizer {
	return &authorizer{perms, TypicalTagType()}
}

// PermissionsAuthorizerFromFile applies the same authorization policy as
// PermissionsAuthorizer, with the Permissions to be used sourced from a file
// named filename.
//
// Changes to the file are monitored and affect subsequent calls to Authorize.
// Currently, this is achieved by re-reading the file on every call to
// Authorize.
// TODO(ashankar,ataly): Use inotify or a similar mechanism to watch for
// changes.
func PermissionsAuthorizerFromFile(filename string, tagType *vdl.Type) (security.Authorizer, error) {
	if err := validateTagType(tagType); err != nil {
		return nil, err
	}
	return &fileAuthorizer{filename, tagType}, nil
}

func validateTagType(tt *vdl.Type) error {
	if tt.Kind() != vdl.String {
		return fmt.Errorf("tag type(%v) must be backed by a string not %v", tt, tt.Kind())
	}
	return nil
}

// PermissionsSpec represents a specification for permissions derived
// from command line flags or some other means.
type PermissionsSpec struct {
	// ExplicitlySpecified is true if any part of the specification was obtained
	// from an explicitly specified command line flag.
	ExplicitlySpecified bool
	// Files represents a set of named files that contain permissions.
	// The name 'runtime' is reserved for use by the runtime.
	Files map[string]string
	// Literal represents a literal, ie. json, permissions specification.
	Literal string
}

// Copy returns a copy of the PermissionSpec.
func (ps *PermissionsSpec) Copy() PermissionsSpec {
	files := make(map[string]string, len(ps.Files))
	for k, v := range ps.Files {
		files[k] = v
	}
	return PermissionsSpec{ps.ExplicitlySpecified, files, ps.Literal}
}

// AuthorizerFromSpec creates an authorizer as specified by the
// supplied specification. If no permissions are specified then the
// default, (ie. nil) Authorizer is returned.
func AuthorizerFromSpec(ps PermissionsSpec, name string, tagType *vdl.Type) (security.Authorizer, error) {
	if len(ps.Literal) == 0 {
		filename := ps.Files[name]
		if len(filename) == 0 {
			return nil, nil
		}
		return PermissionsAuthorizerFromFile(ps.Files[name], tagType)
	}
	perms, err := ReadPermissions(bytes.NewBufferString(ps.Literal))
	if err != nil {
		return nil, fmt.Errorf("failed to parse permissions literal: %v: %v", ps.Literal, err)
	}
	return PermissionsAuthorizer(perms, tagType)
}

type authorizer struct {
	perms   Permissions
	tagType *vdl.Type
}

func (a *authorizer) Authorize(ctx *context.T, call security.Call) error {
	blessings, invalid := security.RemoteBlessingNames(ctx, call)
	hastag := false
	for _, tag := range call.MethodTags() {
		if tag.Type() == a.tagType {
			if hastag {
				return ErrorfMultipleTags(ctx, "authorizer on %v.%v cannot handle multiple tags of type %v; this is likely unintentional", call.Suffix(), call.Method(), a.tagType.String())
			}
			hastag = true
			if acl, exists := a.perms[tag.RawString()]; !exists || !acl.Includes(blessings...) {
				return ErrorfNoPermissions(ctx, "%v does not have %[3]v access (rejected blessings: %[2]v)", blessings, invalid, tag.RawString())
			}
		}
	}
	if !hastag {
		return ErrorfNoTags(ctx, "authorizer on %v.%v has no tags of type %v; this is likely unintentional", call.Suffix(), call.Method(), a.tagType.String())
	}
	return nil
}

type fileAuthorizer struct {
	filename string
	tagType  *vdl.Type
}

func (a *fileAuthorizer) Authorize(ctx *context.T, call security.Call) error {
	perms, err := loadPermissionsFromFile(a.filename)
	if err != nil {
		ctx.Infof("failed to read Permissions file: %v: %v", a.filename, err)
		return verror.ErrInternal.Errorf(ctx, "internal error: failed to read Permissions from file")
	}
	return (&authorizer{perms, a.tagType}).Authorize(ctx, call)
}

func loadPermissionsFromFile(filename string) (Permissions, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadPermissions(file)
}

// Copyright 2021 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vxray integrates amazon's xray distributed tracing system
// with vanadium's vtrace. When used this implementation of vtrace
// will publish to xray in addition to the normal vtrace functionality.
package vxray

import (
	"fmt"

	"github.com/aws/aws-xray-sdk-go/awsplugins/beanstalk"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ecs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"v.io/v23/context"
	"v.io/v23/logging"
	"v.io/v23/vtrace"
	"v.io/x/ref/lib/flags"
	libvtrace "v.io/x/ref/lib/vtrace"
)

type options struct {
	mergeLogging bool
}

// Option represents an option to InitXRay.
type Option func(o *options)

// EC2Plugin initializes the EC2 plugin.
func EC2Plugin() Option {
	return func(o *options) {
		ec2.Init()
	}
}

// ECSPlugin initializes the ECS plugin.
func ECSPlugin() Option {
	return func(o *options) {
		ecs.Init()
	}
}

// BeanstalkPlugin initializes the BeanstalkPlugin plugin.
func BeanstalkPlugin() Option {
	return func(o *options) {
		beanstalk.Init()
	}
}

// MergeLogging configures captureing xray logging and mergeing it with the
// logs produced from the context used to initialize Xray.
func MergeLogging(v bool) Option {
	return func(o *options) {
		o.mergeLogging = v
	}
}

type xraylogger struct {
	logging.Logger
}

func (xl *xraylogger) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	switch level {
	case xraylog.LogLevelInfo, xraylog.LogLevelWarn:
		xl.Info(msg.String())
	case xraylog.LogLevelError:
		xl.Error(msg.String())
	case xraylog.LogLevelDebug:
		xl.VI(1).Info(msg.String())
	}
}

func initXRay(ctx *context.T, config xray.Config, opts []Option) (*context.T, error) {
	o := options{}
	for _, fn := range opts {
		fn(&o)
	}
	if err := xray.Configure(config); err != nil {
		ctx.Errorf("failed to configure xray context: %v", err)
		return ctx, err
	}
	if o.mergeLogging {
		xray.SetLogger(&xraylogger{context.LoggerFromContext(ctx)})
	}
	return WithConfig(ctx, config)
}

func initVTraceForXRay(ctx *context.T, opts flags.VtraceFlags) (*context.T, error) {
	store, err := libvtrace.NewStore(opts)
	if err != nil {
		return ctx, err
	}
	mgr := &manager{}
	ctx = vtrace.WithManager(ctx, mgr)
	ctx = vtrace.WithStore(ctx, store)
	return ctx, nil
}

// InitXRay configures the AWS xray service and returns a context containing
// the xray configuration. This should only be called once.
func InitXRay(ctx *context.T, vopts flags.VtraceFlags, config xray.Config, opts ...Option) (*context.T, error) {
	if !vopts.EnableAWSXRay {
		return ctx, nil
	}
	ctx, err := initXRay(ctx, config, opts)
	if err != nil {
		return ctx, err
	}
	return initVTraceForXRay(ctx, vopts)
}
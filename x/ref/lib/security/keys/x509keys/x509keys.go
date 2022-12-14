// Copyright 2021 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package x509 provides support for using x509/ssl keys with the security/keys
// package.
package x509keys

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"v.io/v23/security"
	"v.io/x/ref/lib/security/keys"
	"v.io/x/ref/lib/security/keys/internal"
)

// Make the key functions local to this package available to this package.
var keyRegistrar = keys.NewRegistrar()

func init() {
	keys.MustRegister(keyRegistrar)
}

// MustRegister is like Register but panics on error.
func MustRegister(r *keys.Registrar) {
	if err := Register(r); err != nil {
		panic(err)
	}
}

// Register registers the required functions for handling ssl public and
// private key files via the x/ref/security/keys package.
func Register(r *keys.Registrar) error {
	r.RegisterPublicKeyParser(parseCertificateBlock, "CERTIFICATE", nil)
	r.RegisterPublicKeyMarshaler(marshalCertificate, (*x509.Certificate)(nil))
	return r.RegisterAPI((*x509CertAPI)(nil), (*x509.Certificate)(nil))
}

type x509CertAPI struct{}

func (*x509CertAPI) Signer(ctx context.Context, key crypto.PrivateKey) (security.Signer, error) {
	// Note that this method will never get called since it is only
	// registered for x509.Certificates.
	return nil, fmt.Errorf("x509keys.Signer: should never be called, has been called for unsupported key type %T", key)
}

func (*x509CertAPI) PublicKey(key interface{}) (security.PublicKey, error) {
	if c, ok := key.(*x509.Certificate); ok {
		return security.NewPublicKey(c)
	}
	return nil, fmt.Errorf("x509keys.Signer: unsupported key type %T", key)
}

func (*x509CertAPI) CryptoPublicKey(key interface{}) (crypto.PublicKey, error) {
	if c, ok := key.(*x509.Certificate); ok {
		return c.PublicKey, nil
	}
	return nil, fmt.Errorf("x509keys.Signer: unsupported key type %T", key)
}

func parseCertificateBlock(block *pem.Block) (crypto.PublicKey, error) {
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, err
}

func marshalCertificate(key crypto.PublicKey) ([]byte, error) {
	cert, ok := key.(*x509.Certificate)
	if !ok {
		return nil, fmt.Errorf("x509keys.marshalCertificate: unsupported key type %T", key)
	}
	return internal.EncodePEM("CERTIFICATE", cert.Raw, nil)
}

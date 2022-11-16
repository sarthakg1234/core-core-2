// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conn

import (
	"bytes"
	"crypto/rand"
	"time"

	"golang.org/x/crypto/nacl/box"
	v23 "v.io/v23"
	"v.io/v23/context"
	"v.io/v23/flow"
	"v.io/v23/flow/message"
	"v.io/v23/naming"
	"v.io/v23/rpc/version"
	"v.io/v23/security"
	"v.io/v23/verror"
	slib "v.io/x/ref/lib/security"
	iflow "v.io/x/ref/runtime/internal/flow"
)

var (
	authDialerTag   = []byte("AuthDial\x00")
	authAcceptorTag = []byte("AuthAcpt\x00")
)

type dialHandshakeResult struct {
	names    []string
	rejected []security.RejectedBlessing
	rtt      time.Duration
	err      error
}

func (c *Conn) dialHandshake(
	ctx *context.T,
	versions version.RPCVersionRange,
<<<<<<< HEAD
	auth flow.PeerAuthorizer,
	handshakeCh chan<- dialHandshakeResult) {

	defer c.loopWG.Done()

	c.mu.Lock()
	// We only send our real blessings if we are a server in addition to being a client.
	// Otherwise, we only send our public key through a nameless blessings object.
	// TODO(suharshs): Should we reveal server blessings if we are connecting to proxy here.
	if c.handler != nil {
		c.localBlessings, c.localValid = v23.GetPrincipal(ctx).BlessingStore().Default()
	} else {
		c.localBlessings, _ = security.NamelessBlessing(v23.GetPrincipal(ctx).PublicKey())
	}
	c.mu.Unlock()
=======
	auth flow.PeerAuthorizer) (names []string, rejected []security.RejectedBlessing, rtt time.Duration, err error) {
>>>>>>> cos-cleanup-accept-handshake

	binding, remoteEndpoint, rttstart, err := c.setup(ctx, versions, true, c.mtu)
	if err != nil {
		handshakeCh <- dialHandshakeResult{err: err}
		return
	}

	c.mu.Lock()
	dialedEP := c.remote
	c.remote.RoutingID = remoteEndpoint.RoutingID
	c.blessingsFlow = newBlessingsFlow(c)
	c.mu.Unlock()

	rttend, err := c.readRemoteAuth(ctx, binding, true)
	if err != nil {
		handshakeCh <- dialHandshakeResult{err: err}
		return
	}
	rtt := rttend.Sub(rttstart)

	c.mu.Lock()
	// Note that the remoteBlessings and discharges are stored in data
	// structures in the blessingsFlow implementation.
	rBlessings := c.remoteBlessings
	rDischarges := c.remoteDischarges
	c.mu.Unlock()

	if rBlessings.IsZero() {
		handshakeCh <- dialHandshakeResult{err: ErrAcceptorBlessingsMissing.Errorf(ctx, "dial: acceptor did not send blessings")}
		return
	}
	var names []string
	var rejected []security.RejectedBlessing
	if c.MatchesRID(dialedEP) {
		// If we hadn't reached the endpoint we expected we would have treated this connection as
		// a proxy, and proxies aren't authorized.  In this case we didn't find a proxy, so go ahead
		// and authorize the connection.
		names, rejected, err = auth.AuthorizePeer(ctx, c.local, c.remote, rBlessings, rDischarges)
		if err != nil {
			handshakeCh <- dialHandshakeResult{
				names:    names,
				rejected: rejected,
				rtt:      rtt,
				err:      iflow.MaybeWrapError(verror.ErrNotTrusted, ctx, err),
			}
			return
		}
	}
	signedBinding, err := v23.GetPrincipal(ctx).Sign(append(authDialerTag, binding...))
	if err != nil {
		handshakeCh <- dialHandshakeResult{
			names:    names,
			rejected: rejected,
			rtt:      rtt,
			err:      err,
		}
		return
	}
	lAuth := message.Auth{ChannelBinding: signedBinding}
	// The client sends its blessings without any blessing-pattern encryption to the
	// server as it has already authorized the server. Thus the 'peers' argument to
	// blessingsFlow.send is nil.
	if lAuth.BlessingsKey, _, err = c.blessingsFlow.send(ctx, c.localBlessings, nil, nil); err != nil {
		handshakeCh <- dialHandshakeResult{
			names:    names,
			rejected: rejected,
			rtt:      rtt,
			err:      err,
		}
		return
	}
	err = c.sendAuthMessage(ctx, lAuth)
	handshakeCh <- dialHandshakeResult{
		names:    names,
		rejected: rejected,
		rtt:      rtt,
		err:      err,
	}
}

// MatchesRID returns true if the given endpoint matches the routing
// ID of the remote server.  Also returns true if the given ep has
// the null routing id (in which case it is assumed that any connected
// server must be the target since nothing has been specified).
func (c *Conn) MatchesRID(ep naming.Endpoint) bool {
	return ep.RoutingID == naming.NullRoutingID ||
		c.remote.RoutingID == ep.RoutingID
}

type acceptHandshakeResult struct {
	rtt         time.Duration
	refreshTime time.Time
	err         error
}

func (c *Conn) acceptHandshake(
	ctx *context.T,
	versions version.RPCVersionRange,
	authorizedPeers []security.BlessingPattern,
	handshakeCh chan<- acceptHandshakeResult) {
	defer close(handshakeCh)
	defer c.loopWG.Done()

	principal := v23.GetPrincipal(ctx)
	localBlessings, localValid := principal.BlessingStore().Default()
	if localBlessings.IsZero() {
		localBlessings, _ = security.NamelessBlessing(principal.PublicKey())
	}
	// PrepareDischarges may issue RPCs to validate 3rd party caveats.
	localDischarges, refreshTime := slib.PrepareDischarges(ctx, localBlessings, nil, "", nil)

	binding, remoteEndpoint, _, err := c.setup(ctx, versions, false, c.mtu)
	if err != nil {
		handshakeCh <- acceptHandshakeResult{err: err}
		return
	}

	c.mu.Lock()
	c.localBlessings = localBlessings
	c.localValid = localValid
	c.localDischarges = localDischarges
	c.remote = remoteEndpoint
	c.blessingsFlow = newBlessingsFlow(c)
	c.mu.Unlock()

	signedBinding, err := v23.GetPrincipal(ctx).Sign(append(authAcceptorTag, binding...))
	if err != nil {
		handshakeCh <- acceptHandshakeResult{err: err}
		return
	}

	lAuth := message.Auth{
		ChannelBinding: signedBinding,
	}
	lAuth.BlessingsKey, lAuth.DischargeKey, err = c.blessingsFlow.send(
		ctx, c.localBlessings, c.localDischarges, authorizedPeers)
	if err != nil {
		handshakeCh <- acceptHandshakeResult{err: err}
		return
	}

	rttstart := time.Now()
	err = c.sendAuthMessage(ctx, lAuth)
	if err != nil {
		handshakeCh <- acceptHandshakeResult{err: err}
		return
	}
	rttend, err := c.readRemoteAuth(ctx, binding, false)
	handshakeCh <- acceptHandshakeResult{rttend.Sub(rttstart), refreshTime, err}
}

var emptyNaClPublicKey [32]byte

func (c *Conn) setup(ctx *context.T, versions version.RPCVersionRange, dialer bool, mtu uint64) ([]byte, naming.Endpoint, time.Time, error) { //nolint:gocyclo
	var rttstart time.Time
	pk, sk, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, naming.Endpoint{}, rttstart, err
	}
	lSetup := message.Setup{
		Versions:          versions,
		PeerLocalEndpoint: c.local,
		Mtu:               c.mtu,
		SharedTokens:      c.flowControl.bytesBufferedPerFlow,
	}
	copy(lSetup.PeerNaClPublicKey[:], (*pk)[:])
	if !c.remote.IsZero() {
		lSetup.PeerRemoteEndpoint = c.remote
	}
	ch := make(chan error, 1)
	go func() {
		rttstart = time.Now()
		ch <- c.sendSetupMessage(ctx, lSetup)
	}()
	rSetup, nBuf, err := c.mp.readSetup(ctx)
	defer putNetBuf(nBuf)
	if err != nil {
		<-ch
		return nil, naming.Endpoint{}, rttstart, ErrRecv.Errorf(ctx, "conn.setup: recv: %v", err)
	}
	if err := <-ch; err != nil {
		return nil, naming.Endpoint{}, rttstart, ErrSend.Errorf(ctx, "conn.setup: remote %v: %v", c.remoteEndpointForError(), err)
	}
	if c.version, err = version.CommonVersion(ctx, lSetup.Versions, rSetup.Versions); err != nil {
		return nil, naming.Endpoint{}, rttstart, err
	}
	if c.local.IsZero() {
		c.local = rSetup.PeerRemoteEndpoint
	}

	if rSetup.Mtu == 0 {
		rSetup.Mtu = mtu
	}
	if lSetup.Mtu == 0 {
		lSetup.Mtu = mtu
	}

	// Pick the smaller of the two MTUs.
	if rSetup.Mtu > lSetup.Mtu {
		c.mtu = lSetup.Mtu
	} else {
		c.mtu = rSetup.Mtu
	}

	lshared := lSetup.SharedTokens
	if rSetup.SharedTokens != 0 && rSetup.SharedTokens < lshared {
		lshared = rSetup.SharedTokens
	}

	c.flowControl.configure(c.mtu, lshared)

	if bytes.Equal(rSetup.PeerNaClPublicKey[:], emptyNaClPublicKey[:]) {
		return nil, naming.Endpoint{}, rttstart, ErrMissingSetupOption.Errorf(ctx, "conn.setup: missing required setup option: peerNaClPublicKey")
	}
	binding, err := c.mp.enableEncryption(ctx, pk, sk, &rSetup.PeerNaClPublicKey, c.version)
	if err != nil {
		return nil, naming.Endpoint{}, rttstart, err
	}
	c.mp.setMTU(c.mtu, c.flowControl.bytesBufferedPerFlow)

	if c.version >= version.RPCVersion14 {
		// We include the setup messages in the channel binding to prevent attacks
		// where a man in the middle changes fields in the Setup message (e.g. a
		// downgrade attack wherein a MITM attacker changes the Version field of
		// the Setup message to a lower-security version.)
		// We always put the dialer first in the binding.
		if dialer {
			if binding, err = lSetup.Append(ctx, nil); err != nil {
				return nil, naming.Endpoint{}, rttstart, err
			}
			if binding, err = rSetup.Append(ctx, binding); err != nil {
				return nil, naming.Endpoint{}, rttstart, err
			}
		} else {
			if binding, err = rSetup.Append(ctx, nil); err != nil {
				return nil, naming.Endpoint{}, rttstart, err
			}
			if binding, err = lSetup.Append(ctx, binding); err != nil {
				return nil, naming.Endpoint{}, rttstart, err
			}
		}
	}
	// if we're encapsulated in another flow, tell that flow to stop
	// encrypting now that we've started.
	c.mp.disableEncryptionOnEncapsulatedFlow()
	return binding, rSetup.PeerLocalEndpoint, rttstart, nil
}

// readRemoteAuth is used to read the auth handshake messages from the remote
// endpoint. This is a sequence of Data messages followed by an Auth message.
// readRemoteAuth runs asynchronously on the both the dialer and acceptor.
// On successful completion, the connection has accepted the remote's
// blessings and discharges and verified the channel binding. The remote's
// public key is non-nil and recorded in the connection and will never be
// changed from hereonin.
func (c *Conn) readRemoteAuth(ctx *context.T, binding []byte, dialer bool) (time.Time, error) {
	rauth, err := c.readRemoteAuthLoop(ctx)
	if err != nil {
		return time.Time{}, err
	}
	rttend := time.Now()

	tag := authDialerTag
	if dialer {
		tag = authAcceptorTag
	}

	if rauth.BlessingsKey != 0 {
		rBlessings, rDischarges, err := c.blessingsFlow.getRemote(
			ctx, rauth.BlessingsKey, rauth.DischargeKey)
		if err != nil {
			return rttend, err
		}
		// The first blessing that's received is 'bound' to this conn. All
		// subsequet blessings received must have the same public key.
		rpk := rBlessings.PublicKey()

		c.mu.Lock()
		c.rPublicKey = rpk
		c.remoteBlessings = rBlessings
		c.remoteDischarges = rDischarges
		c.remoteValid = make(chan struct{})
		c.mu.Unlock()
		c.blessingsFlow.setPublicKeyBinding(rpk)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rPublicKey == nil {
		return rttend, ErrNoPublicKey.Errorf(ctx, "conn.readRemoteAuth: no public key received")
	}

	if !rauth.ChannelBinding.Verify(c.rPublicKey, append(tag, binding...)) {
		return rttend, ErrInvalidChannelBinding.Errorf(ctx, "conn.readRemoteAuth: the channel binding was invalid")
	}
	return rttend, nil
}

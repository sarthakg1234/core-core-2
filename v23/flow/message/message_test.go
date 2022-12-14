// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message_test

import (
	"errors"
	"reflect"
	"runtime"
	"testing"

	v23 "v.io/v23"
	"v.io/v23/context"
	"v.io/v23/flow/message"
	"v.io/v23/naming"
	"v.io/v23/rpc/version"
	"v.io/v23/security"
	"v.io/x/ref/lib/security/keys"
	_ "v.io/x/ref/runtime/factories/fake"
	"v.io/x/ref/test"
	"v.io/x/ref/test/sectestdata"
)

func testMessages(t *testing.T, ctx *context.T, cases []message.Message) {
	testMessagesWithResults(t, ctx, cases, cases)
}

func testMessagesWithResults(t *testing.T, ctx *context.T, cases []message.Message, results []message.Message) {
	_, _, line, _ := runtime.Caller(2)
	for i, orig := range cases {
		want := results[i]
		encoded, err := orig.Append(ctx, nil)
		if err != nil {
			t.Errorf("%v: line: %v, unexpected error for %#v: %v", i, line, orig, err)
		}
		got, err := message.Read(ctx, encoded)
		if err != nil {
			t.Errorf("%v: line: %v, unexpected error reading %#v: %v", i, line, want, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%v: line: %v, got: %#v, want %#v", i, line, got, want)
		}
		got = got.Copy()
		for i := range encoded {
			encoded[i] = 0xff
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%v: line: %v, got: %#v, want %#v", i, line, got, want)
		}
	}
}

func TestSetup(t *testing.T) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	ep1, err := naming.ParseEndpoint(
		"@6@tcp@foo.com:1234@a,b@00112233445566778899aabbccddeeff@m@v.io/foo")
	if err != nil {
		t.Fatal(err)
	}
	ep2, err := naming.ParseEndpoint(
		"@6@tcp@bar.com:1234@a,b@00112233445566778899aabbccddeeff@m@v.io/bar")
	if err != nil {
		t.Fatal(err)
	}
	testMessages(t, ctx, []message.Message{
		message.Setup{Versions: version.RPCVersionRange{Min: 1, Max: 5}},
		message.Setup{
			Versions: version.RPCVersionRange{Min: 1, Max: 5},
			PeerNaClPublicKey: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
				14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
			PeerRemoteEndpoint: ep1,
			PeerLocalEndpoint:  ep2,
		},
		message.Setup{
			Versions:     version.RPCVersionRange{Min: 1, Max: 5},
			Mtu:          1 << 16,
			SharedTokens: 1 << 20,
		},
		message.Setup{},
	})
}

func TestTearDown(t *testing.T) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	testMessages(t, ctx, []message.Message{
		message.TearDown{Message: "foobar"},
		message.TearDown{},
	})
}

func TestAuth(t *testing.T) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	for _, kt := range sectestdata.SupportedKeyAlgos {
		signer := sectestdata.V23Signer(kt, sectestdata.V23KeySetA)
		p, err := security.CreatePrincipal(signer, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		sig, err := p.Sign([]byte("message"))
		if err != nil {
			t.Fatal(err)
		}
		msg := message.Auth{BlessingsKey: 1, DischargeKey: 5, ChannelBinding: sig}
		switch kt {
		case keys.ECDSA256, keys.ECDSA384, keys.ECDSA521:
			msg = message.ExposeSetAuthMessageType(msg, true, false, false)
		case keys.ED25519:
			msg = message.ExposeSetAuthMessageType(msg, false, true, false)
		case keys.RSA2048, keys.RSA4096:
			msg = message.ExposeSetAuthMessageType(msg, false, false, true)
		}
		testMessages(t, ctx, []message.Message{msg})

		encoded, err := msg.Append(ctx, nil)
		if err != nil {
			t.Errorf("unexpected error for %#v: %v", msg, err)
		}
		decoded, err := message.Read(ctx, encoded)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		authMsg, ok := decoded.(message.Auth)
		if !ok {
			t.Errorf("unexpected message type: %T", authMsg)
			continue
		}
		if !authMsg.ChannelBinding.Verify(p.PublicKey(), []byte("message")) {
			t.Errorf("failed to verify signature")
		}
	}
}

func TestOpenFlow(t *testing.T) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	testMessages(t, ctx, []message.Message{
		message.OpenFlow{
			ID:              23,
			InitialCounters: 1 << 20,
			BlessingsKey:    42,
			DischargeKey:    55,
			Flags:           message.CloseFlag,
			Payload:         []byte("fake payload"),
		},
		message.OpenFlow{ID: 23, InitialCounters: 1 << 20, BlessingsKey: 42, DischargeKey: 55},
	})
}

func TestMissingBlessings(t *testing.T) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	cases := []message.Message{
		message.OpenFlow{},
		message.Auth{},
	}
	for _, m := range cases {
		encoded, err := m.Append(ctx, nil)
		if err != nil {
			t.Errorf("unexpected error for %#v: %v", m, err)
		}
		_, err = message.Read(ctx, encoded)
		if !errors.Is(err, message.ErrMissingBlessings) {
			t.Errorf("unexpected error for %#v: got %v want MissingBlessings", m, err)
		}
	}
}

func TestAddReceiveBuffers(t *testing.T) {
	ctx, shutdown := test.V23Init()
	defer shutdown()
	testMessages(t, ctx, []message.Message{
		message.Release{},
		message.Release{Counters: []message.Counter{
			{4, 233},
			{9, 423242},
		}},
	})
}

func TestData(t *testing.T) {
	ctx, shutdown := test.V23Init()
	defer shutdown()
	testMessages(t, ctx, []message.Message{
		message.Data{ID: 1123, Flags: message.CloseFlag, Payload: []byte("fake payload")},
		message.Data{},
	})
	testMessagesWithResults(t, ctx,
		[]message.Message{
			message.Data{ID: 1123, Flags: message.DisableEncryptionFlag, Payload: []byte("fake payload")},
		},
		[]message.Message{
			message.Data{ID: 1123, Flags: message.DisableEncryptionFlag},
		})
}

func TestDataReuse(t *testing.T) {
	ctx, shutdown := test.V23Init()
	defer shutdown()
	m := message.Data{ID: 1123, Flags: message.CloseFlag, Payload: []byte("fake payload")}
	buf := make([]byte, 0, 1024)
	buf, err := m.Append(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	am, err := message.Read(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	rm := am.(message.Data)
	if got, want := rm.ID, uint64(1123); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rm.Flags, uint64(message.CloseFlag); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rm.Payload), "fake payload"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	m = message.Data{ID: 1124}
	buf = make([]byte, 0, 1024)
	buf, err = m.Append(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}

	am, err = message.Read(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}
	rm = am.(message.Data)
	if got, want := rm.ID, uint64(1124); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rm.Flags, uint64(0); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := len(rm.Payload), 0; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}

func plaintextPayload(m message.Message) ([]byte, bool) {
	switch msg := m.(type) {
	case message.Data:
		if msg.Flags&message.DisableEncryptionFlag != 0 {
			return msg.Payload, true
		}
	case message.OpenFlow:
		if msg.Flags&message.DisableEncryptionFlag != 0 {
			return msg.Payload, true
		}
	}
	return nil, false
}

func setPlaintextPayload(m message.Message, payload []byte) message.Message {
	switch msg := m.(type) {
	case message.Data:
		msg.Payload = payload
		return msg
	case message.OpenFlow:
		msg.Payload = payload
		return msg
	}
	return m
}

func expectsPlaintextPayload(m message.Message) bool {
	switch msg := m.(type) {
	case message.Data:
		return msg.Flags&message.DisableEncryptionFlag != 0
	case message.OpenFlow:
		return msg.Flags&message.DisableEncryptionFlag != 0
	}
	return false
}

func clearDisableEncryptionFlag(m message.Message) message.Message {
	switch msg := m.(type) {
	case message.Data:
		msg.Flags &^= message.DisableEncryptionFlag
		return msg
	case message.OpenFlow:
		msg.Flags &^= message.DisableEncryptionFlag
		return msg
	}
	return m
}

func TestPlaintextPayloads(t *testing.T) {

	encrypted := []message.Message{
		message.Data{Flags: message.CloseFlag, Payload: []byte("fake payload")},
		message.OpenFlow{Flags: message.CloseFlag, Payload: []byte("fake payload")},
	}
	for _, m := range encrypted {
		payload, ok := plaintextPayload(m)
		if got, want := ok, false; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if payload != nil {
			t.Errorf("unexpected payload")
		}

	}

	disabled := []message.Message{
		message.Data{Flags: message.DisableEncryptionFlag, Payload: []byte("fake payload")},
		message.OpenFlow{
			Flags:   message.DisableEncryptionFlag,
			Payload: []byte("fake payload"),
		},
	}
	for _, m := range disabled {
		payload, ok := plaintextPayload(m)
		if got, want := ok, true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if payload == nil {
			t.Errorf("expected a payload")
		}

		newPayload := []byte("hello")
		m = setPlaintextPayload(m, newPayload)
		m = m.Copy()
		copy(newPayload, []byte("world"))
		p, _ := plaintextPayload(m)
		if got, want := string(p), "hello"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	empty := []message.Message{
		message.Data{Flags: message.DisableEncryptionFlag},
		message.OpenFlow{Flags: message.DisableEncryptionFlag},
	}

	for _, m := range empty {
		payload, ok := plaintextPayload(m)
		if got, want := ok, true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if payload != nil {
			t.Errorf("unexpected payload")
		}
		if got, want := expectsPlaintextPayload(m), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		m = clearDisableEncryptionFlag(m)
		_, ok = plaintextPayload(m)
		if got, want := ok, false; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := expectsPlaintextPayload(m), false; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestProxy(t *testing.T) {
	ctx, shutdown := test.V23Init()
	defer shutdown()
	ep1, err := naming.ParseEndpoint(
		"@6@tcp@foo.com:1234@a,b@00112233445566778899aabbccddeeff@m@v.io/foo")
	if err != nil {
		t.Fatal(err)
	}
	ep2, err := naming.ParseEndpoint(
		"@6@tcp@bar.com:1234@a,b@00112233445566778899aabbccddeeff@m@v.io/bar")
	if err != nil {
		t.Fatal(err)
	}
	testMessages(t, ctx, []message.Message{
		message.MultiProxyRequest{},
		message.ProxyServerRequest{},
		message.ProxyResponse{},
		message.ProxyResponse{Endpoints: []naming.Endpoint{ep1, ep2}},
		message.ProxyErrorResponse{},
		message.ProxyErrorResponse{Error: "error"},
	})
}

func benchmarkMessageAppend(b *testing.B, ctx *context.T, m message.Message) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.Append(ctx, make([]byte, 0, 4096))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkMessageRead(b *testing.B, ctx *context.T, buf []byte) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := message.Read(ctx, buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func setupMessage(b *testing.B) message.Message {
	ep1, err := naming.ParseEndpoint(
		"@6@tcp@foo.com:1234@a,b@00112233445566778899aabbccddeeff@m@v.io/foo")
	if err != nil {
		b.Fatal(err)
	}
	ep2, err := naming.ParseEndpoint(
		"@6@tcp@bar.com:1234@a,b@00112233445566778899aabbccddeeff@m@v.io/bar")
	if err != nil {
		b.Fatal(err)
	}
	return message.Setup{
		Versions: version.RPCVersionRange{Min: 1, Max: 5},
		PeerNaClPublicKey: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
			14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		PeerRemoteEndpoint: ep1,
		PeerLocalEndpoint:  ep2,
	}
}

func BenchmarkSetupAppend(b *testing.B) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	m := setupMessage(b)
	benchmarkMessageAppend(b, ctx, m)
}

func BenchmarkSetupRead(b *testing.B) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	m := setupMessage(b)
	buf := make([]byte, 0, 2048)
	buf, _ = m.Append(ctx, buf)
	benchmarkMessageRead(b, ctx, buf)
}

func BenchmarkDataAppend(b *testing.B) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	m := message.Data{ID: 1123, Flags: message.CloseFlag, Payload: []byte("fake payload")}
	benchmarkMessageAppend(b, ctx, m)
}

func BenchmarkDataRead(b *testing.B) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	m := message.Data{ID: 1123, Flags: message.CloseFlag, Payload: []byte("fake payload")}
	buf := make([]byte, 0, 2048)
	buf, _ = m.Append(ctx, buf)
	benchmarkMessageRead(b, ctx, buf)
}

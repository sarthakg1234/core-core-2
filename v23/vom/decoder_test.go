// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vom_test

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"v.io/v23/vdl"
	"v.io/v23/verror"
	"v.io/v23/vom"
	"v.io/v23/vom/testdata/types"
	"v.io/v23/vom/vomtest"
	// Import verror to ensure that interface tests result in *verror.E
)

var (
	rtIface = reflect.TypeOf((*interface{})(nil)).Elem()
	//nolint:unused
	rtValue = reflect.TypeOf(vdl.Value{})
)

func encode(t *testing.T, in interface{}) *vom.Decoder {
	buf := &bytes.Buffer{}
	enc := vom.NewEncoder(buf)
	if err := enc.Encode(in); err != nil {
		t.Errorf("encode: %v", err)
		return vom.NewDecoder(buf)
	}
	return vom.NewDecoder(buf)
}
func decodeAny(t *testing.T, dec *vom.Decoder) interface{} {
	var out interface{}
	if err := dec.Decode(&out); err != nil {
		t.Errorf("decode: %v", err)
	}
	return out
}

func decodeVerror(t *testing.T, dec *vom.Decoder) error {
	var out verror.E
	if err := dec.Decode(&out); err != nil {
		t.Errorf("decode: %v", err)
	}
	return out
}

func decodeError(t *testing.T, dec *vom.Decoder) error {
	var out error
	if err := dec.Decode(&out); err != nil {
		t.Errorf("decode: %v", err)
	}
	return out
}

func TestDecoderNativeAnyTypes(t *testing.T) {
	assertType := func(any interface{}, typ string) {
		if got, want := fmt.Sprintf("%T", any), typ; got != want {
			_, _, line, _ := runtime.Caller(1)
			t.Errorf("line: %v: got %v, want %v", line, got, want)
		}
	}
	assertValue := func(any, val interface{}) {
		if got, want := any, val; !reflect.DeepEqual(got, want) {
			_, _, line, _ := runtime.Caller(1)
			t.Errorf("line: %v: got %#v, want %#v", line, got, want)
		}
	}

	// Create a verror.E directly to avoid having the stack and other
	// fields filled in which would mess up comparing the sent and
	// received values since the stack is never sent over the wire.
	verrorVal := verror.E{
		ID:     "foobar",
		Action: verror.RetryBackoff,
		Msg:    "verror slice",
	}

	{
		any := decodeAny(t, encode(t, verror.E{}))
		assertType(any, "verror.E")
		assertValue(any, verror.E{})

		verr := decodeVerror(t, encode(t, verror.E{}))
		assertType(verr, "verror.E")
		assertValue(verr, verror.E{})

		err := decodeError(t, encode(t, os.ErrClosed))
		assertType(err, "*verror.E")
		nerr := &verror.E{
			ID:        verror.ErrUnknown.ID,
			Action:    verror.NoRetry,
			Msg:       os.ErrClosed.Error(),
			ParamList: []interface{}{"", "", os.ErrClosed.Error()},
		}
		assertValue(err, nerr)
	}

	{
		any := decodeAny(t, encode(t, (*verror.E)(nil)))
		assertType(any, "*verror.E")
		assertValue(any, (*verror.E)(nil))

		any = decodeAny(t, encode(t, &verrorVal))
		assertType(any, "*verror.E")
		assertValue(any, &verrorVal)

		verr := decodeAny(t, encode(t, (*verror.E)(nil)))
		assertType(verr, "*verror.E")
		assertValue(verr, (*verror.E)(nil))

		verr = decodeAny(t, encode(t, &verrorVal))
		assertType(verr, "*verror.E")
		assertValue(verr, &verrorVal)
	}

	{
		any := decodeAny(t, encode(t, []*verror.E{}))
		assertType(any, "[]error")
		assertValue(any, []error(nil))

		any = decodeAny(t, encode(t, []*verror.E(nil)))
		assertType(any, "[]error")
		assertValue(any, []error(nil))

		any = decodeAny(t, encode(t, []*verror.E{&verrorVal}))
		assertType(any, "[]error")
		assertValue(any, []error{&verrorVal})
	}

	{
		// These currently translate to error rather than verror.E.
		// See TODO/comment in reflect_reader.go:readIntoAny:238.
		any := decodeAny(t, encode(t, []verror.E{}))
		assertType(any, "[]error")
		assertValue(any, []error(nil))

		any = decodeAny(t, encode(t, []verror.E(nil)))
		assertType(any, "[]error")
		assertValue(any, []error(nil))

		any = decodeAny(t, encode(t, []verror.E{verrorVal}))
		assertType(any, "[]error")
		assertValue(any, []error{&verrorVal})
	}
}

func TestDecoderIfcAndValue(t *testing.T) {
	// The decoder tests take a long time, so we run them concurrently and
	// split them across two tests.
	var pending sync.WaitGroup
	allPass := vomtest.AllPass()
	errCh := make(chan error, len(allPass))
	pending.Add(len(allPass))
	for i, test := range allPass {
		go func(i int, test vomtest.Entry) {
			defer pending.Done()
			if err := testDecoder("[go value]", test, rvPtrValue(test.Value)); err != nil {
				errCh <- err
				return
			}
			if err := testDecoder("[go iface]", test, rvPtrIface(test.Value)); err != nil {
				errCh <- err
				return
			}
		}(i, test)
	}
	pending.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDecoderNewValues(t *testing.T) {
	// The decoder tests take a long time, so we run them concurrently and
	// split them across two tests.
	var pending sync.WaitGroup
	allPass := vomtest.AllPass()
	errCh := make(chan error, len(allPass))
	pending.Add(len(allPass))
	for i, test := range allPass {
		go func(i int, test vomtest.Entry) {
			defer pending.Done()
			vv, err := vdl.ValueFromReflect(test.Value)
			if err != nil {
				errCh <- fmt.Errorf("%s: ValueFromReflect failed: %v", test.Name(), err)
				return
			}
			vvWant := reflect.ValueOf(vv)
			if err := testDecoder("[new *vdl.Value]", test, vvWant); err != nil {
				errCh <- err
				return
			}
			err = testDecoderFunc("[zero vdl.Value]", test, vvWant, func() reflect.Value {
				return reflect.ValueOf(vdl.ZeroValue(vv.Type()))
			})
			if err != nil {
				errCh <- err
				return
			}
		}(i, test)
	}
	pending.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Error(err)
		}
	}
}

func rvPtrValue(rv reflect.Value) reflect.Value {
	result := reflect.New(rv.Type())
	result.Elem().Set(rv)
	return result
}

func rvPtrIface(rv reflect.Value) reflect.Value {
	result := reflect.New(rtIface)
	result.Elem().Set(rv)
	return result
}

func testDecoder(pre string, test vomtest.Entry, rvWant reflect.Value) error {
	return testDecoderFunc(pre, test, rvWant, func() reflect.Value {
		return reflect.New(rvWant.Type().Elem())
	})
	// TODO(toddw): Add tests that start with a randomly-set value.
}

var printDetailsMu sync.Mutex

func printDetails(name string, rvGot, rvWant reflect.Value) {
	printDetailsMu.Lock()
	defer printDetailsMu.Unlock()
	fmt.Printf("FAILED test: %v\n", name)
	spew.Printf("as reflect.Value:\nGOT  %v\nWANT %v\n", rvGot, rvWant)
	gotVal, wantVal := rvGot.Interface(), rvWant.Interface()
	spew.Printf("as interface{}:\nGOT  %v\nWANT %v\n", gotVal, wantVal)
	fmt.Printf("IsValid: got %v, want %v\n", rvGot.IsValid(), rvWant.IsValid())
	fmt.Printf("IsZero: got %v, want %v\n", rvGot.IsZero(), rvWant.IsZero())
	fmt.Printf("IsNil: got %v, want %v\n", rvGot.IsNil(), rvWant.IsNil())
	fmt.Printf("Type: got %v, want %v\n", rvGot.Type(), rvWant.Type())
	fmt.Printf("Kind: got %v, want %v\n", rvGot.Kind(), rvWant.Kind())
}

func testDecoderFunc(pre string, test vomtest.Entry, rvWant reflect.Value, rvNew func() reflect.Value) error {
	readEOF := make([]byte, 1)
	for _, mode := range vom.AllReadModes {
		// Test vom.NewDecoder.
		{
			name := fmt.Sprintf("%s (%s) %s", pre, mode, test.Name())
			rvGot := rvNew()
			reader := mode.TestReader(bytes.NewReader(test.Bytes()))
			dec := vom.NewDecoder(reader)
			if err := dec.Decode(rvGot.Interface()); err != nil {
				return fmt.Errorf("%s: Decode failed: %v", name, err)
			}
			if !vdl.DeepEqualReflect(rvGot, rvWant) {
				printDetails(name, rvGot, rvWant)
				return fmt.Errorf("%s\nGOT  %v\nWANT %v", name, rvGot, rvWant)
			}
			if n, err := reader.Read(readEOF); n != 0 || err != io.EOF {
				return fmt.Errorf("%s: reader got (%d,%v), want (0,EOF)", name, n, err)
			}
		}
		// Test vom.NewDecoderWithTypeDecoder
		{
			name := fmt.Sprintf("%s (%s with TypeDecoder) %s", pre, mode, test.Name())
			rvGot := rvNew()
			readerT := mode.TestReader(bytes.NewReader(test.TypeBytes()))
			decT := vom.NewTypeDecoder(readerT)
			reader := mode.TestReader(bytes.NewReader(test.ValueBytes()))
			dec := vom.NewDecoderWithTypeDecoder(reader, decT)
			err := dec.Decode(rvGot.Interface())
			// TODO(cnicolaou): this test is racy, since decT.Stop() will
			// not cause the internal goroutine used by the type decoder
			// decT to actually finish, in fact, it's still blocked
			// reading from readerT, which is subsequently accessed below
			// which the race detector will correctly identify as a race.
			// Ideally, we should remove the use of a goroutine within the
			// type decoder to both avoid the race and hopefully simply
			// that code. It will also make it possible to more easily
			// reuse type decoders which will be performance improvement.
			if err != nil {
				return fmt.Errorf("%s: Decode failed: %v", name, err)
			}
			if !vdl.DeepEqualReflect(rvGot, rvWant) {
				printDetails(name, rvGot, rvWant)
				return fmt.Errorf("%s\nGOT  %v\nWANT %v", name, rvGot, rvWant)
			}
			if n, err := reader.Read(readEOF); n != 0 || err != io.EOF {
				return fmt.Errorf("%s: reader got (%d,%v), want (0,EOF)", name, n, err)
			}
			if n, err := readerT.Read(readEOF); n != 0 || err != io.EOF {
				return fmt.Errorf("%s: readerT got (%d,%v), want (0,EOF)", name, n, err)
			}
		}
	}
	// Test single-shot vom.Decode twice, to ensure we test the cache hit case.
	for i := 0; i < 2; i++ {
		name := fmt.Sprintf("%s (single-shot %d) %s", pre, i, test.Name())
		rvGot := rvNew()
		if err := vom.Decode(test.Bytes(), rvGot.Interface()); err != nil {
			return fmt.Errorf("%s: Decode failed: %v", name, err)
		}
		if !vdl.DeepEqualReflect(rvGot, rvWant) {
			return fmt.Errorf("%s\nGOT  %v\nWANT %v", name, rvGot, rvWant)
		}
	}
	return nil
}

// TestRoundtrip* tests test encoding and then decoding results in various modes.
func TestRoundtrip(t *testing.T) { testRoundtrip(t, false, 1) }

func TestRoundtrip_5(t *testing.T) { testRoundtrip(t, false, 5) }

func TestRoundtripWithTypeDecoder_1(t *testing.T) { testRoundtrip(t, true, 1) }
func TestRoundtripWithTypeDecoder_2(t *testing.T) { testRoundtrip(t, true, 5) }

func TestRoundtripWithTypeDecoder_5(t *testing.T)  { testRoundtrip(t, true, 5) }
func TestRoundtripWithTypeDecoder_10(t *testing.T) { testRoundtrip(t, true, 10) }
func TestRoundtripWithTypeDecoder_20(t *testing.T) { testRoundtrip(t, true, 20) }

func testRoundtrip(t *testing.T, withTypeEncoderDecoder bool, concurrency int) {
	mp := runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(mp)

	tests := []struct {
		In, Want interface{}
	}{
		// Test that encoding nil/empty composites leads to nil.
		{[]byte(nil), []byte(nil)},
		{[]byte{}, []byte(nil)},
		{[]int64(nil), []int64(nil)},
		{[]int64{}, []int64(nil)},
		{map[string]int64(nil), map[string]int64(nil)},
		{map[string]int64{}, map[string]int64(nil)},
		{struct{}{}, struct{}{}},
		{struct{ A []byte }{nil}, struct{ A []byte }{}},
		{struct{ A []byte }{[]byte{}}, struct{ A []byte }{}},
		{struct{ A []int64 }{nil}, struct{ A []int64 }{}},
		{struct{ A []int64 }{[]int64{}}, struct{ A []int64 }{}},
		// Test that encoding nil typeobject leads to AnyType.
		{(*vdl.Type)(nil), vdl.AnyType},
		// Test that both encoding and decoding ignore unexported fields.
		{struct{ a, X, b string }{"a", "XYZ", "b"}, struct{ d, X, e string }{X: "XYZ"}},
		{
			struct {
				a bool
				X string
				b int64
			}{true, "XYZ", 123},
			struct {
				a complex64
				X string
				b []byte
			}{X: "XYZ"},
		},
		// Test for array encoding/decoding.
		{[3]byte{1, 2, 3}, [3]byte{1, 2, 3}},
		{[3]int64{1, 2, 3}, [3]int64{1, 2, 3}},
		// Test for zero value struct/union field encoding/decoding.
		{struct{ A int64 }{0}, struct{ A int64 }{}},
		{struct{ T *vdl.Type }{nil}, struct{ T *vdl.Type }{vdl.AnyType}},
		{struct{ M map[uint64]struct{} }{make(map[uint64]struct{})}, struct{ M map[uint64]struct{} }{}},
		{struct{ M map[uint64]string }{make(map[uint64]string)}, struct{ M map[uint64]string }{}},
		{struct{ N struct{ A int64 } }{struct{ A int64 }{0}}, struct{ N struct{ A int64 } }{}},
		{struct{ N *types.NStruct }{&types.NStruct{A: false, B: "", C: 0}}, struct{ N *types.NStruct }{&types.NStruct{}}},
		{struct{ N *types.NStruct }{nil}, struct{ N *types.NStruct }{}},
		{types.NUnion(types.NUnionA{Value: false}), types.NUnion(types.NUnionA{})},
		{types.RecA{}, types.RecA(nil)},
		{types.RecStruct{}, types.RecStruct{}},
		// Test for verifying correctness when encoding/decoding shared types concurrently.
		{types.MStruct{}, types.MStruct{E: vdl.AnyType, F: vdl.AnyValue(nil)}},
		{types.NStruct{}, types.NStruct{}},
		{types.XyzStruct{}, types.XyzStruct{}},
		{types.YzStruct{}, types.YzStruct{}},
		{types.MBool(false), types.MBool(false)},
		{types.NString(""), types.NString("")},
		{vdl.ValueOf(uint16(5)), vdl.ValueOf(uint16(5))},
		{vdl.ValueOf([]interface{}{uint16(5)}).Index(0), vdl.ValueOf(uint16(5))},
	}

	var (
		typeenc *vom.TypeEncoder
		typedec *vom.TypeDecoder
	)
	if withTypeEncoderDecoder {
		r, w := newPipe()
		typeenc = vom.NewTypeEncoder(w)
		typedec = vom.NewTypeDecoder(r)
	}

	var wg sync.WaitGroup
	for n := 0; n < concurrency; n++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for _, n := range rand.Perm(len(tests) * 10) {
				test := tests[n%len(tests)]
				name := fmt.Sprintf("[%d]:%+v,%+v", n, test.In, test.Want)

				var (
					encoder *vom.Encoder
					decoder *vom.Decoder
					buf     bytes.Buffer
				)
				if withTypeEncoderDecoder {
					encoder = vom.NewEncoderWithTypeEncoder(&buf, typeenc)
					decoder = vom.NewDecoderWithTypeDecoder(&buf, typedec)
				} else {
					encoder = vom.NewEncoder(&buf)
					decoder = vom.NewDecoder(&buf)
				}

				if err := encoder.Encode(test.In); err != nil {
					t.Errorf("%s: Encode(%+v) failed: %v", name, test.In, err)
					return
				}
				rv := reflect.New(reflect.TypeOf(test.Want))
				if err := decoder.Decode(rv.Interface()); err != nil {
					t.Errorf("%s: Decode failed: %v", name, err)
					return
				}
				if got := rv.Elem().Interface(); !vdl.DeepEqual(got, test.Want) {
					t.Errorf("%s: Decode mismatch\nGOT  %T %+v\nWANT %T %+v", name, got, got, test.Want, test.Want)
					return
				}
			}
		}(n)
	}
	wg.Wait()
}

func hex2Bin(t *testing.T, hex string) []byte {
	var bin string
	if _, err := fmt.Sscanf(hex, "%x", &bin); err != nil {
		t.Fatalf("error converting %q to binary: %v", hex, err)
	}
	return []byte(bin)
}

// Test that no EOF is returned from Decode() if the type stream finished before the
// value stream. This can only happen if the value stream contains a type that
// does not exist in the type stream, in which case a specific error is returned.
func TestTypeStreamEndsFirst(t *testing.T) {
	hexversion := "81"
	hextype := "5133060025762e696f2f7632332f766f6d2f74657374646174612f74797065732e537472756374416e7901010003416e79010fe1e1533b060023762e696f2f7632332f766f6d2f74657374646174612f74797065732e4e53747275637401030001410101e10001420103e10001430109e1e1"
	hexvalue := "52012a0103070000000001e1e1"
	// Add a second value with a different type id that is not encoded in
	// the type stream.
	hexvalue += "62012a0103070000000001e1e1"
	binversion := string(hex2Bin(t, hexversion))
	bintype := string(hex2Bin(t, hextype))
	binvalue := string(hex2Bin(t, hexvalue))
	tr := strings.NewReader(binversion + bintype)
	typedec := vom.NewTypeDecoder(tr)
	wr := strings.NewReader(binversion + binvalue)
	decoder := vom.NewDecoderWithTypeDecoder(wr, typedec)
	var v interface{}
	if err := decoder.Decode(&v); err != nil {
		t.Errorf("expected no error in decode, but got: %v", err)
	}
	err := decoder.Decode(&v)
	if err == nil || !strings.Contains(err.Error(), "vom: failed to find type") {
		t.Errorf("expected a specific error, not %v", err)
	}
}

func TestTypeStreamTruncationAndErrors(t *testing.T) {
	hexversion := "81"
	hextype := "5133060025762e696f2f7632332f766f6d2f74657374646174612f74797065732e537472756374416e7901010003416e79010fe1e1533b060023762e696f2f7632332f766f6d2f74657374646174612f74797065732e4e53747275637401030001410101e10001420103e10001430109e1e1"
	hexvalue := "52012a0103070000000001e1e1"
	binversion := string(hex2Bin(t, hexversion))
	binvalue := string(hex2Bin(t, hexvalue))

	for _, typeStream := range []string{
		hextype[:2],
		hextype[:4],
		hextype[:32],
		"00",
		"04",
		"7285",
		"5177",
	} {
		bintype := string(hex2Bin(t, typeStream))
		tr := strings.NewReader(binversion + bintype)
		typedec := vom.NewTypeDecoder(tr)
		wr := strings.NewReader(binversion + binvalue)
		decoder := vom.NewDecoderWithTypeDecoder(wr, typedec)
		var v interface{}
		if err := decoder.Decode(&v); err == nil {
			t.Errorf("expected an error")
		}
	}
}

// Return an error on all reads.
type errorReader struct {
}

func (er *errorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("errorReader error")
}

// Test that non-EOF errors on the value stream are returned from Decode() calls.
func TestReceiveTypeStreamError(t *testing.T) {
	hexversion := "80"
	hexvalue := "5206002a000001e1"
	binversion := string(hex2Bin(t, hexversion))
	binvalue := string(hex2Bin(t, hexvalue))
	// Ensure EOF isn't returned if the type decode stream ends first
	typedec := vom.NewTypeDecoder(&errorReader{})
	decoder := vom.NewDecoderWithTypeDecoder(strings.NewReader(binversion+binvalue), typedec)
	var v interface{}
	if err := decoder.Decode(&v); err == nil {
		t.Errorf("expected error in decode, but got none")
	}
}

// Test that using the type decoder incorrectly does not result
// in a deadlock.
func TestFuzzTypeDecodeDeadlock(t *testing.T) {
	var v interface{}
	d := vom.NewDecoder(strings.NewReader("\x81\x30"))
	// Before the fix, this line caused a deadlock and panic.
	d.Decode(&v) //nolint:errcheck
}

// Tests that an input go-fuzz found will no longer cause a
// panic over in package vdl.
func TestFuzzVdlPanic(t *testing.T) {
	var v interface{}
	d := vom.NewDecoder(strings.NewReader("\x81S*\x00\x00$000000000000000000000000000000000000\x01*\xe1U(\x05\x00 00000000000000000000000000000000\x01*\x02+\xe1"))
	// Before this fix this line caused a panic.
	d.Decode(&v) //nolint:errcheck
}

// In concurrent modes, one goroutine may try to read vom types before they are
// actually sent by other goroutine. We use a simple buffered pipe to provide
// blocking read since bytes.Buffer will return EOF in this case.
type pipe struct {
	b         bytes.Buffer
	m         sync.Mutex
	c         sync.Cond
	cancelled bool
}

func newPipe() (io.ReadCloser, io.WriteCloser) {
	p := &pipe{}
	p.c.L = &p.m
	return p, p
}

func (p *pipe) Read(buf []byte) (n int, err error) {
	p.m.Lock()
	defer p.m.Unlock()
	for p.b.Len() == 0 || p.cancelled {
		p.c.Wait()
	}
	return p.b.Read(buf)
}

func (p *pipe) Close() error {
	p.m.Lock()
	p.cancelled = true
	p.c.Broadcast()
	p.m.Unlock()
	return nil
}

func (p *pipe) Write(buf []byte) (n int, err error) {
	p.m.Lock()
	defer p.m.Unlock()
	defer p.c.Signal()
	return p.b.Write(buf)
}

// Test that input found by go-fuzz cannot cause a stack overflow.
func TestFuzzDecodeOverflow(t *testing.T) {
	var v interface{}
	d := vom.NewDecoder(strings.NewReader("\x81\x51\x04\x03\x01\x29\xe1"))

	// Before the fix, this line caused a stack overflow.  After the fix, we
	// expect an error.
	if err := d.Decode(&v); err == nil {
		t.Fatal("unexpected success")
	}
}

// Test that input discovered by go-fuzz does not result in a hang anymore.
func TestFuzzTypeDecodeHang(t *testing.T) {
	var v interface{}
	d := vom.NewDecoder(strings.NewReader(
		"\x81W&\x03\x00 v.io/v23/vom/t" +
			"estdata/types.Rec4\x01)" +
			"\xe1U&\x03\x00 v.io/v23/vom/t" +
			"estdata/types.Rec3\x01," +
			"\xe1S&\x00\x00 v.io/v23/v\x04\x00/t" +
			"estdata/types.Rec2\x01*" +
			"\xe1Q&\x00\x00 v.io/v23/vom/t" +
			"estdataPtypesIRec1\x01*" +
			"\xe1"))

	// Before the fix, this line caused a hang. With the fix, it should
	// give an error.
	err := d.Decode(&v)
	if err == nil {
		t.Fatal("unexpected success")
	}
}

// Test that truncated input is handled gracefully.
func TestDecodeTruncatedInput(t *testing.T) {
	type bar struct {
		I []int64
	}
	type foo struct {
		S string
		I int32
		F float64
		X []byte
		B bar
	}
	encoded, err := vom.Encode(foo{"hello", 42, 3.14159265359, []byte{1, 2, 3, 4, 5, 6}, bar{[]int64{0}}})
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	for x := 0; x < len(encoded)-1; x++ {
		var f interface{}
		if err := vom.Decode(encoded[:x], &f); err == nil {
			t.Errorf("Decode did not fail with x=%d", x)
		}
	}
}

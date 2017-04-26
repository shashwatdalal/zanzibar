// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// TODO: (lu) to be generated

package baz

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/zanzibar/test/lib/bench_gateway"
	"github.com/uber/zanzibar/test/lib/test_gateway"

	"github.com/uber/zanzibar/examples/example-gateway/build/clients"
	bazServer "github.com/uber/zanzibar/examples/example-gateway/build/clients/baz"
	"github.com/uber/zanzibar/examples/example-gateway/build/endpoints"
	"github.com/uber/zanzibar/examples/example-gateway/build/gen-code/clients/baz/baz"
)

var testCallCounter int

func call(
	ctx context.Context, reqHeaders map[string]string, args *baz.SimpleService_Call_Args,
) (map[string]string, error) {
	testCallCounter++
	r := args.Arg
	if r.B1 && r.S2 == "hello" && r.I3 == 42 {
		return nil, nil
	}
	return nil, errors.New("Wrong Args")
}

func TestCallSuccessfulRequestOKResponse(t *testing.T) {
	gateway, err := testGateway.CreateGateway(t, map[string]interface{}{
		"clients.baz.serviceName": "Qux",
	}, &testGateway.Options{
		KnownTChannelBackends: []string{"baz"},
		TestBinary: filepath.Join(
			getDirName(), "..", "..", "..",
			"examples", "example-gateway", "build", "main.go",
		),
	})
	if !assert.NoError(t, err, "got bootstrap err") {
		return
	}
	defer gateway.Close()

	gateway.TChannelBackends()["baz"].Register(
		"SimpleService",
		"Call",
		bazServer.NewSimpleServiceCallHandler(call),
	)

	headers := map[string]string{}

	res, err := gateway.MakeRequest(
		"POST",
		"/baz/call",
		headers,
		bytes.NewReader([]byte(`{"arg":{"b1":true,"s2":"hello","i3":42}}`)),
	)

	if !assert.NoError(t, err, "got http error") {
		return
	}

	assert.Equal(t, 1, testCallCounter)
	assert.Equal(t, "204 No Content", res.Status)
}

func BenchmarkCall(b *testing.B) {
	gateway, err := benchGateway.CreateGateway(
		map[string]interface{}{
			"clients.baz.serviceName": "Qux",
		},
		&testGateway.Options{
			KnownTChannelBackends: []string{"baz"},
		},
		clients.CreateClients,
		endpoints.Register,
	)
	if err != nil {
		b.Error("got bootstrap err: " + err.Error())
		return
	}

	gateway.TChannelBackends()["baz"].Register(
		"SimpleService",
		"Call",
		bazServer.NewSimpleServiceCallHandler(call),
	)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res, err := gateway.MakeRequest(
				"POST", "/baz/call", nil,
				bytes.NewReader([]byte(`{"arg":{"b1":true,"s2":"hello","i3":42}}`)),
			)
			if err != nil {
				b.Error("got http error: " + err.Error())
				break
			}
			if res.Status != "204 No Content" {
				b.Error("got bad status error: " + res.Status)
				break
			}
		}
	})

	b.StopTimer()
	gateway.Close()
}
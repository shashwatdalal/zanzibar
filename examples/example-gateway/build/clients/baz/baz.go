// Code generated by zanzibar
// @generated

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

// Package bazClient is generated code used to make or handle TChannel calls using Thrift.
package bazClient

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/uber/zanzibar/runtime"

	clientsBazBaz "github.com/uber/zanzibar/examples/example-gateway/build/gen-code/clients/baz/baz"
)

// NewClient returns a new TChannel client for service baz.
func NewClient(gateway *zanzibar.Gateway) *BazClient {
	serviceName := gateway.Config.MustGetString("clients.baz.serviceName")
	sc := gateway.Channel.GetSubChannel(serviceName)

	ip := gateway.Config.MustGetString("clients.baz.ip")
	port := gateway.Config.MustGetInt("clients.baz.port")
	sc.Peers().Add(ip + ":" + strconv.Itoa(int(port)))

	timeout := time.Millisecond * time.Duration(
		gateway.Config.MustGetInt("clients.baz.timeout"),
	)
	timeoutPerAttempt := time.Millisecond * time.Duration(
		gateway.Config.MustGetInt("clients.baz.timeoutPerAttempt"),
	)

	client := zanzibar.NewTChannelClient(gateway.Channel,
		&zanzibar.TChannelClientOption{
			ServiceName:       serviceName,
			Timeout:           timeout,
			TimeoutPerAttempt: timeoutPerAttempt,
		},
	)

	return &BazClient{
		thriftService: "SimpleService",
		client:        client,
	}
}

// BazClient is the TChannel client for downstream service.
type BazClient struct {
	thriftService string
	client        zanzibar.TChannelClient
}

// Call ...
func (c *BazClient) Call(
	ctx context.Context,
	reqHeaders map[string]string,
	args *clientsBazBaz.SimpleService_Call_Args,
) (map[string]string, error) {
	var result clientsBazBaz.SimpleService_Call_Result

	success, respHeaders, err := c.client.Call(
		ctx, c.thriftService, "Call", reqHeaders, args, &result,
	)

	if err == nil && !success {
		switch {
		case result.AuthErr != nil:
			err = result.AuthErr
		default:
			err = errors.New("BazClient received no result or unknown exception for Call")
		}
	}
	if err != nil {
		return nil, err
	}

	return respHeaders, err
}

// Compare ...
func (c *BazClient) Compare(
	ctx context.Context,
	reqHeaders map[string]string,
	args *clientsBazBaz.SimpleService_Compare_Args,
) (*clientsBazBaz.BazResponse, map[string]string, error) {
	var result clientsBazBaz.SimpleService_Compare_Result

	success, respHeaders, err := c.client.Call(
		ctx, c.thriftService, "Compare", reqHeaders, args, &result,
	)

	if err == nil && !success {
		switch {
		case result.AuthErr != nil:
			err = result.AuthErr
		default:
			err = errors.New("BazClient received no result or unknown exception for Compare")
		}
	}
	if err != nil {
		return nil, nil, err
	}

	resp, err := clientsBazBaz.SimpleService_Compare_Helper.UnwrapResponse(&result)
	return resp, respHeaders, err
}

// Ping ...
func (c *BazClient) Ping(
	ctx context.Context,
	reqHeaders map[string]string,
) (*clientsBazBaz.BazResponse, map[string]string, error) {
	var result clientsBazBaz.SimpleService_Ping_Result

	args := &clientsBazBaz.SimpleService_Ping_Args{}
	success, respHeaders, err := c.client.Call(
		ctx, c.thriftService, "Ping", reqHeaders, args, &result,
	)

	if err == nil && !success {
		switch {
		default:
			err = errors.New("BazClient received no result or unknown exception for Ping")
		}
	}
	if err != nil {
		return nil, nil, err
	}

	resp, err := clientsBazBaz.SimpleService_Ping_Helper.UnwrapResponse(&result)
	return resp, respHeaders, err
}

// SillyNoop ...
func (c *BazClient) SillyNoop(
	ctx context.Context,
	reqHeaders map[string]string,
) (map[string]string, error) {
	var result clientsBazBaz.SimpleService_SillyNoop_Result

	args := &clientsBazBaz.SimpleService_SillyNoop_Args{}
	success, respHeaders, err := c.client.Call(
		ctx, c.thriftService, "SillyNoop", reqHeaders, args, &result,
	)

	if err == nil && !success {
		switch {
		case result.AuthErr != nil:
			err = result.AuthErr
		case result.ServerErr != nil:
			err = result.ServerErr
		default:
			err = errors.New("BazClient received no result or unknown exception for SillyNoop")
		}
	}
	if err != nil {
		return nil, err
	}

	return respHeaders, err
}

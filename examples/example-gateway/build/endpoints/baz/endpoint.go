// Code generated by zanzibar
// @generated

// Copyright (c) 2018 Uber Technologies, Inc.
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

package bazendpoint

import (
	module "github.com/uber/zanzibar/examples/example-gateway/build/endpoints/baz/module"
	zanzibar "github.com/uber/zanzibar/runtime"
)

// Endpoint registers a request handler on a gateway
type Endpoint interface {
	Register(*zanzibar.Gateway) error
}

// NewEndpoint returns a collection of endpoints that can be registered on
// a gateway
func NewEndpoint(deps *module.Dependencies) Endpoint {
	return &EndpointHandlers{
		SimpleServiceCallHandler:             NewSimpleServiceCallHandler(deps),
		SimpleServiceCompareHandler:          NewSimpleServiceCompareHandler(deps),
		SimpleServicePingHandler:             NewSimpleServicePingHandler(deps),
		SimpleServiceSillyNoopHandler:        NewSimpleServiceSillyNoopHandler(deps),
		SimpleServiceTransHandler:            NewSimpleServiceTransHandler(deps),
		SimpleServiceTransHeadersHandler:     NewSimpleServiceTransHeadersHandler(deps),
		SimpleServiceTransHeadersTypeHandler: NewSimpleServiceTransHeadersTypeHandler(deps),
		SimpleServiceHeaderSchemaHandler:     NewSimpleServiceHeaderSchemaHandler(deps),
	}
}

// EndpointHandlers is a collection of individual endpoint handlers
type EndpointHandlers struct {
	SimpleServiceCallHandler             *SimpleServiceCallHandler
	SimpleServiceCompareHandler          *SimpleServiceCompareHandler
	SimpleServicePingHandler             *SimpleServicePingHandler
	SimpleServiceSillyNoopHandler        *SimpleServiceSillyNoopHandler
	SimpleServiceTransHandler            *SimpleServiceTransHandler
	SimpleServiceTransHeadersHandler     *SimpleServiceTransHeadersHandler
	SimpleServiceTransHeadersTypeHandler *SimpleServiceTransHeadersTypeHandler
	SimpleServiceHeaderSchemaHandler     *SimpleServiceHeaderSchemaHandler
}

// Register registers the endpoint handlers with the gateway
func (handlers *EndpointHandlers) Register(gateway *zanzibar.Gateway) error {
	err0 := handlers.SimpleServiceCallHandler.Register(gateway)
	if err0 != nil {
		return err0
	}
	err1 := handlers.SimpleServiceCompareHandler.Register(gateway)
	if err1 != nil {
		return err1
	}
	err2 := handlers.SimpleServicePingHandler.Register(gateway)
	if err2 != nil {
		return err2
	}
	err3 := handlers.SimpleServiceSillyNoopHandler.Register(gateway)
	if err3 != nil {
		return err3
	}
	err4 := handlers.SimpleServiceTransHandler.Register(gateway)
	if err4 != nil {
		return err4
	}
	err5 := handlers.SimpleServiceTransHeadersHandler.Register(gateway)
	if err5 != nil {
		return err5
	}
	err6 := handlers.SimpleServiceTransHeadersTypeHandler.Register(gateway)
	if err6 != nil {
		return err6
	}
	err7 := handlers.SimpleServiceHeaderSchemaHandler.Register(gateway)
	if err7 != nil {
		return err7
	}
	return nil
}

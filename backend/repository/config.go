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

package repository

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/uber/zanzibar/codegen"
	zanzibar "github.com/uber/zanzibar/runtime"
)

// Config stores configuration for a gateway.
type Config struct {
	ID                  string
	Repository          string
	Team                string
	Tier                int
	ThriftRootDir       string
	PackageRoot         string
	GenCodePackage      string
	TargetGenDir        string
	ClientConfigDir     string
	EndpointConfigDir   string
	MiddlewareConfigDir string
	// Maps endpointID to configuration.
	Endpoints map[string]*EndpointConfig
	// Maps clientID to configuration.
	Clients        map[string]*ClientConfig
	ThriftServices ThriftServiceMap
	Middlewares    map[string]*MiddlewareConfig
}

// ThriftServiceMap maps thrift file -> service name -> *ThriftService
type ThriftServiceMap map[string]map[string]*ThriftService

// EndpointConfig stores configuration for an endpoint.
type EndpointConfig struct {
	ID                string
	Type              ProtocolType
	HandleID          string
	ConfigFile        string
	ThriftFile        string
	ThriftServiceName string
	MethodName        string
	WorkflowType      string
	ClientID          string
	ClientMethod      string
	TestFixture       string
	Middlewares       []*EndptMidConfig
}

// EndptMidConfig represents configuration for a middleware.
type EndptMidConfig struct {
	Name    string
	Options map[string]interface{}
}

// ClientConfig stores configuration for an client.
type ClientConfig struct {
	ID             string
	Type           ProtocolType
	ThriftFile     string
	ExposedMethods map[string]string
}

// MiddlewareConfig represents configuration for a middleware.
type MiddlewareConfig struct {
	ID     string
	Schema map[string]string
}

// ThriftService is a service defined in Thrift file.
type ThriftService struct {
	Name    string
	Path    string
	Methods []ThriftMethod
}

// ThriftMethod is a method defined in a Thrift Service.
type ThriftMethod struct {
	Name string
	Type ProtocolType
}

// ThriftMeta is the meta about a thrift file.
type ThriftMeta struct {
	// relative path under thrift root directory.
	Path string `json:"path"`
	// commited version
	Version string `json:"version"`
	// content of the thrift file
	Content string `json:"content,omitempty"`
}

// ProtocolType represents tranportation protocal type.
type ProtocolType string

const (
	gatewayConfigFile = "gateway.json"

	// HTTP type
	HTTP ProtocolType = "http"
	// TCHANNEL type
	TCHANNEL ProtocolType = "tchannel"
	// Custom type
	CUSTOM ProtocolType = "custom"
	// UNKNOWN type
	UNKNOWN ProtocolType = "unknown"
)

// GatewayConfig returns the cached gateway configuration of the repository if
// the repository has not been updated for a certain time.
func (r *Repository) GatewayConfig() (*Config, error) {
	r.RLock()
	curCfg, curCfgError := r.gatewayConfig, r.gatewayConfigError
	r.RUnlock()
	if (curCfg != nil || curCfgError != nil) && !r.Update() {
		return curCfg, curCfgError
	}
	return r.LatestGatewayConfig()
}

// LatestGatewayConfig returns the configuration of current repository.
func (r *Repository) LatestGatewayConfig() (*Config, error) {
	// The newCfg won't be changed once created.
	newCfg, newCfgError := r.newGatewayConfig()
	r.Lock()
	r.gatewayConfig, r.gatewayConfigError = newCfg, newCfgError
	r.Unlock()
	return newCfg, newCfgError
}

// newGatewayConfig regenerates the configuration for a repository.
func (r *Repository) newGatewayConfig() (configuration *Config, cfgErr error) {
	defer func() {
		if p := recover(); p != nil {
			cfgErr = errors.Errorf(
				"panic when getting configuration for gateway %s: %+v",
				r.LocalDir(), p,
			)
		}
	}()
	// Get the read lock for reading the content of repository.
	r.RLock()
	defer r.RUnlock()
	configDir := r.absPath(r.LocalDir())
	cfg := zanzibar.NewStaticConfigOrDie([]string{
		filepath.Join(configDir, gatewayConfigFile),
	}, nil)
	config := &Config{
		ID:                  cfg.MustGetString("gatewayName"),
		Repository:          r.remote,
		ThriftRootDir:       cfg.MustGetString("thriftRootDir"),
		PackageRoot:         cfg.MustGetString("packageRoot"),
		GenCodePackage:      cfg.MustGetString("genCodePackage"),
		TargetGenDir:        cfg.MustGetString("targetGenDir"),
		ClientConfigDir:     cfg.MustGetString("clientConfig"),
		EndpointConfigDir:   cfg.MustGetString("endpointConfig"),
		MiddlewareConfigDir: cfg.MustGetString("middlewareConfig"),
	}
	pkgHelper, err := codegen.NewPackageHelper(
		config.PackageRoot,
		configDir,
		config.MiddlewareConfigDir,
		config.ThriftRootDir,
		config.GenCodePackage,
		config.TargetGenDir,
		"",
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create package helper")
	}

	moduleSystem, err := codegen.NewDefaultModuleSystem(pkgHelper)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create module system")
	}

	moduleInstances, err := moduleSystem.ResolveModules(
		pkgHelper.PackageRoot(),
		configDir,
		pkgHelper.CodeGenTargetPath(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve module instances")
	}

	gatewaySpec, err := codegen.NewGatewaySpec(
		moduleInstances,
		pkgHelper,
		configDir,
		config.EndpointConfigDir,
		config.MiddlewareConfigDir,
		config.ID,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read gateway spec")
	}

	if config.ThriftServices, err = r.thriftservices(config.ThriftRootDir, pkgHelper); err != nil {
		return nil, errors.Wrapf(err, "failed to read thrift services")
	}

	config.Clients = r.clientConfigs(config.ThriftRootDir, gatewaySpec)
	config.Endpoints = r.endpointConfigs(config.ThriftRootDir, gatewaySpec)
	return config, nil
}

func (r *Repository) clientConfigs(thriftRootDir string, gatewaySpec *codegen.GatewaySpec) map[string]*ClientConfig {
	cfgs := make(map[string]*ClientConfig, len(gatewaySpec.ClientModules))
	for _, spec := range gatewaySpec.ClientModules {
		clientType := ProtocolTypeFromString(spec.ClientType)
		clientConfig := &ClientConfig{
			ID:             spec.ClientID,
			Type:           clientType,
			ExposedMethods: spec.ExposedMethods,
		}
		if clientType != CUSTOM {
			clientConfig.ThriftFile = r.relativePath(thriftRootDir, spec.ThriftFile)
		}
		cfgs[clientConfig.ID] = clientConfig
	}
	return cfgs
}

func (r *Repository) endpointConfigs(thriftRootDir string, gatewaySpec *codegen.GatewaySpec) map[string]*EndpointConfig {
	cfgs := make(map[string]*EndpointConfig, len(gatewaySpec.EndpointModules))
	for file, spec := range gatewaySpec.EndpointModules {
		endpointID := spec.EndpointID + "." + spec.HandleID
		cfgs[endpointID] = &EndpointConfig{
			ID:                endpointID,
			ConfigFile:        strings.TrimPrefix(file, r.localDir),
			Type:              ProtocolTypeFromString(spec.EndpointType),
			HandleID:          spec.HandleID,
			ThriftFile:        r.relativePath(thriftRootDir, spec.ThriftFile),
			ThriftServiceName: spec.ThriftServiceName,
			MethodName:        spec.ThriftMethodName,
			WorkflowType:      spec.WorkflowType,
			ClientID:          spec.ClientID,
			ClientMethod:      spec.ClientMethod,
			// TODO(zw): add test fixtures and middleware config.
		}
	}
	return cfgs
}

func (r *Repository) thriftservices(thriftRootDir string, packageHelper *codegen.PackageHelper) (ThriftServiceMap, error) {
	idlMap := make(ThriftServiceMap)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		// Parse service module as tchannel service.
		mspec, err := codegen.NewModuleSpec(path, false, false, packageHelper)
		if err != nil {
			return errors.Wrapf(err, "failed to genenerate module spec for thrift %s", path)
		}
		serviceType := TCHANNEL
		// Parse HTTP annotations.
		if _, err := codegen.NewModuleSpec(path, true, false, packageHelper); err == nil {
			serviceType = HTTP
		}
		relativePath := r.relativePath(thriftRootDir, path)
		idlMap[relativePath] = map[string]*ThriftService{}
		for _, service := range mspec.Services {
			tservice := &ThriftService{
				Name: service.Name,
				Path: relativePath,
			}
			tservice.Methods = make([]ThriftMethod, len(service.Methods))
			for i, method := range service.Methods {
				tservice.Methods[i].Name = method.Name
				tservice.Methods[i].Type = serviceType
			}
			idlMap[relativePath][service.Name] = tservice
		}
		return nil
	}
	if err := filepath.Walk(r.absPath(thriftRootDir), walkFn); err != nil {
		return nil, errors.Wrapf(err, "failed to traverse IDL dir")
	}
	return idlMap, nil
}

func (r *Repository) absPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	path, err := filepath.Abs(filepath.Join(r.localDir, path))
	if err != nil {
		return err.Error()
	}
	return path
}

func (r *Repository) relativePath(rootDir string, filePath string) string {
	rootAbsDir := r.absPath(rootDir)
	fileAbsPath := r.absPath(filePath)
	relative := strings.TrimPrefix(fileAbsPath, rootAbsDir)
	return strings.TrimPrefix(relative, "/")
}

// ProtocolTypeFromString converts a string to ProtocolType.
func ProtocolTypeFromString(str string) ProtocolType {
	switch str {
	case "http":
		return HTTP
	case "tchannel":
		return TCHANNEL
	case "custom":
		return CUSTOM
	default:
		return UNKNOWN
	}
}

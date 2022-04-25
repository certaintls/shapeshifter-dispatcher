/*
 * Copyright (c) 2014, Yawning Angel <yawning at torproject dot org>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

// Package transports provides a interface to query supported pluggable
// transports.
package transports

import (
	"encoding/json"
	"errors"
	Optimizer "github.com/OperatorFoundation/shapeshifter-transports/transports/Optimizer/v3"
	replicant "github.com/OperatorFoundation/shapeshifter-transports/transports/Replicant/v3"
	"golang.org/x/net/proxy"
)

// Transports returns the list of registered transport protocols.
func Transports() []string {
	return []string{"obfs2", "shadow", "Dust", "meeklite", "Replicant", "obfs4", "Optimizer"}
}

func ParseArgsObfs2(args string) (*obfs2.OptimizerTransport, error) {
	var config obfs2.Config
	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("obfs2 options json decoding error")
	}
	transport := obfs2.New(config.Address, proxy.Direct)

	return transport, nil
}

func ParseArgsObfs4(args string, dialer proxy.Dialer) (*obfs4.TransportClient, error) {
	var config obfs4.Config

	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("obfs4 options json decoding error")
	}

	iatMode := 0
	if config.IatMode == "1" {
		iatMode = 1
	}

	transport := obfs4.TransportClient{
		CertString: config.CertString,
		IatMode:    iatMode,
		Address:    config.Address,
		Dialer:     dialer,
	}

	return &transport, nil
}

func ParseArgsShadow(args string) (*shadow.Transport, error) {
	var config shadow.ClientConfig
	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("shadow options json decoding error")
	}
	transport := shadow.NewTransport(config.Password, config.CipherName, config.Address)

	return &transport, nil
}

func ParseArgsShadowServer(args string) (*shadow.ServerConfig, error) {
	var config shadow.ServerConfig

	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("shadow server options json decoding error")
	}

	return &config, nil
}

func ParseArgsDust(args string, dialer proxy.Dialer) (*Dust.Transport, error) {
	var config Dust.Config

	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("dust options json decoding error")
	}

	transport := Dust.Transport{
		ServerPublic: config.ServerPublic,
		Address:      config.Address,
		Dialer:       dialer,
	}

	return &transport, nil
}

func CreateDefaultReplicantServer() replicant.ServerConfig {
	config := replicant.ServerConfig{
		Toneburst: nil,
		Polish:    nil,
	}

	return config
}

func ParseArgsReplicantClient(args string, dialer proxy.Dialer) (*replicant.TransportClient, error) {
	var config *replicant.ClientConfig

	var ReplicantConfig replicant.ClientJSONConfig
	if args == "" {
		return nil, errors.New("must specify transport options when using replicant")
	}
	argsBytes := []byte(args)
	unmarshalError := json.Unmarshal(argsBytes, &ReplicantConfig)
	if unmarshalError != nil {
		return nil, errors.New("could not unmarshal Replicant args")
	}
	var parseErr error
	config, parseErr = replicant.DecodeClientConfig(ReplicantConfig.Config)
	if parseErr != nil {
		return nil, errors.New("could not parse config")
	}

	transport := replicant.TransportClient{
		Config:  *config,
		Address: (*config).Address,
		Dialer:  dialer,
	}

	return &transport, nil
}

//  target string, dialer proxy.Dialer
func ParseArgsReplicantServer(args string) (*replicant.ServerConfig, error) {
	var config *replicant.ServerConfig

	type replicantJsonConfig struct {
		Config string
	}
	var ReplicantConfig replicantJsonConfig
	if args == "" {
		transport := CreateDefaultReplicantServer()
		return &transport, nil
	}
	argsBytes := []byte(args)
	unmarshalError := json.Unmarshal(argsBytes, &ReplicantConfig)
	if unmarshalError != nil {
		return nil, errors.New("could not unmarshal Replicant args")
	}
	var parseErr error
	config, parseErr = replicant.DecodeServerConfig(ReplicantConfig.Config)
	if parseErr != nil {
		return nil, parseErr
	}

	return config, nil
}

func ParseArgsStarBridgeClient(args string, target string, dialer proxy.Dialer) (*StarBridge.Transport, error) {
	var config StarBridge.ClientConfig
	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("starbridge options json decoding error")
	}
	transport := StarBridge.Transport{
		Config:  config,
		Address: target,
		Dialer:  dialer,
	}

	return &transport, nil
}

func ParseArgsStarBridgeServer(args string) (*StarBridge.ServerConfig, error) {
	var config StarBridge.ServerConfig

	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("starbridge server options json decoding error")
	}

	return &config, nil
}

func ParseArgsMeeklite(args string, dialer proxy.Dialer) (*meeklite.Transport, error) {
	var config meeklite.Config

	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("meeklite options json decoding error")
	}

	transport := meeklite.Transport{
		URL:    config.URL,
		Front:  config.Front,
		Dialer: dialer,
	}

	return &transport, nil
}

func ParseArgsMeekliteServer(args string) (*meekserver.Config, error) {
	var config meekserver.Config

	bytes := []byte(args)
	jsonError := json.Unmarshal(bytes, &config)
	if jsonError != nil {
		return nil, errors.New("meeklite options json decoding error")
	}

	return &config, nil
}

type OptimizerConfig struct {
	Transports []interface{} `json:"transports"`
	Strategy   string        `json:"strategy"`
}

type OptimizerArgs struct {
	Address string                 `json:"address"`
	Name    string                 `json:"name"`
	Config  map[string]interface{} `json:"config"`
}

func ParseArgsOptimizer(jsonConfig string, dialer proxy.Dialer) (*Optimizer.Client, error) {
	var config OptimizerConfig
	var transports []Optimizer.TransportDialer
	var strategy Optimizer.Strategy
	jsonByte := []byte(jsonConfig)
	parseErr := json.Unmarshal(jsonByte, &config)
	if parseErr != nil {
		return nil, errors.New("could not marshal optimizer config")
	}
	transports, parseErr = parseTransports(config.Transports, dialer)
	if parseErr != nil {
		println("this is the returned error from parseTransports:", parseErr)
		return nil, errors.New("could not parse transports")
	}

	strategy, parseErr = parseStrategy(config.Strategy, transports)
	if parseErr != nil {
		return nil, errors.New("could not parse strategy")
	}

	transport := Optimizer.NewOptimizerClient(transports, strategy)

	return transport, nil
}

func parseStrategy(strategyString string, transports []Optimizer.TransportDialer) (Optimizer.Strategy, error) {
	switch strategyString {
	case "first":
		strategy := Optimizer.NewFirstStrategy(transports)
		return strategy, nil
	case "random":
		strategy := Optimizer.NewRandomStrategy(transports)
		return strategy, nil
	case "rotate":
		strategy := Optimizer.NewRotateStrategy(transports)
		return strategy, nil
	case "track":
		return Optimizer.NewTrackStrategy(transports), nil
	case "minimizeDialDuration":
		return Optimizer.NewMinimizeDialDuration(transports), nil

	default:
		return nil, errors.New("invalid strategy")
	}
}

func parseTransports(otcs []interface{}, dialer proxy.Dialer) ([]Optimizer.TransportDialer, error) {
	transports := make([]Optimizer.TransportDialer, len(otcs))
	for index, untypedOtc := range otcs {
		switch untypedOtc.(type) {
		case map[string]interface{}:
			otc := untypedOtc.(map[string]interface{})
			transport, err := parsedTransport(otc, dialer)
			if err != nil {
				return nil, errors.New("transport could not parse config")
				//this error sucks and is uninformative
			}
			transports[index] = transport
		default:
			return nil, errors.New("unsupported type for transport")
		}

	}
	return transports, nil
}

func parsedTransport(otc map[string]interface{}, dialer proxy.Dialer) (Optimizer.TransportDialer, error) {
	var config map[string]interface{}

	type PartialOptimizerConfig struct {
		Name string `json:"name"`
	}
	jsonString, MarshalErr := json.Marshal(otc)
	if MarshalErr != nil {
		return nil, errors.New("error marshalling optimizer otc")
	}
	var PartialConfig PartialOptimizerConfig
	unmarshalError := json.Unmarshal(jsonString, &PartialConfig)
	if unmarshalError != nil {
		return nil, errors.New("error unmarshalling optimizer otc")
	}
	//on to parsing the config
	untypedConfig, ok3 := otc["config"]
	if !ok3 {
		return nil, errors.New("missing config in transport parser")
	}

	switch untypedConfig.(type) {

	case map[string]interface{}:
		config = untypedConfig.(map[string]interface{})

	default:
		return nil, errors.New("unsupported type for optimizer config option")
	}

	jsonConfigBytes, configMarshalError := json.Marshal(config)
	if configMarshalError != nil {
		return nil, errors.New("could not marshal Optimizer config")
	}
	jsonConfigString := string(jsonConfigBytes)
	switch PartialConfig.Name {
	case "shadow":
		shadowTransport, parseErr := ParseArgsShadow(jsonConfigString)
		if parseErr != nil {
			return nil, errors.New("could not parse shadow Args")
		}
		return shadowTransport, nil
	case "obfs2":
		obfs2Transport, parseErr := ParseArgsObfs2(jsonConfigString)
		if parseErr != nil {
			return nil, errors.New("could not parse obfs2 Args")
		}
		return obfs2Transport, nil
	case "obfs4":
		obfs4Transport, parseErr := ParseArgsObfs4(jsonConfigString, dialer)
		if parseErr != nil {
			return nil, errors.New("could not parse obfs4 Args")
		}
		return obfs4Transport, nil
	case "meeklite":
		meekliteTransport, parseErr := ParseArgsMeeklite(jsonConfigString, dialer)
		if parseErr != nil {
			return nil, errors.New("could not parse meeklite Args")
		}
		return meekliteTransport, nil
	case "Dust":
		DustTransport, parseErr := ParseArgsDust(jsonConfigString, dialer)
		if parseErr != nil {
			return nil, errors.New("could not parse dust Args")
		}
		return DustTransport, nil
	case "Replicant":
		replicantTransport, parseErr := ParseArgsReplicantClient(jsonConfigString, dialer)
		if parseErr != nil {
			return nil, errors.New("could not parse replicant Args")
		}
		return replicantTransport, nil
	case "StarBridge":
		starbridgeTransport, parseErr := ParseArgsStarBridgeClient(jsonConfigString, PartialConfig.Address, dialer)
		if parseErr != nil {
			return nil, errors.New("could not parse starbridge Args")
		}
		return starbridgeTransport, nil
	case "Optimizer":
		optimizerTransport, parseErr := ParseArgsOptimizer(jsonConfigString, dialer)
		if parseErr != nil {
			return nil, errors.New("could not parse Optimizer Args")
		}
		return optimizerTransport, nil
	default:
		return nil, errors.New("unsupported transport name")
	}
}

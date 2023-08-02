# venomoid
[![Main](https://github.com/ranglust/venomoid/actions/workflows/main.yaml/badge.svg?branch=main)](https://github.com/ranglust/venomoid/actions/workflows/main.yaml) ![test coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ranglust/720f73912f5b895dcc9b63d2a872cc00/raw/venomoid__.json)


## Overview
[Venomoid](https://en.wikipedia.org/wiki/Venomoid) utlizes the builder pattern in order to simplify the use of [viper](https://github.com/spf13/viper) and reduce
the builderplate footprint in every project

## Download
To download the package run:
```go
go get github.com/ranglust/venomoid
```

## Usage
### With config file lookup

```go
package main

import (
	"fmt"
	"github.com/ranglust/venomoid"
	"os"
)

type myConfig struct {
	key1 string
	key2 bool
	key3 int
}

func main() {

	var paths []string = []string{"/etc", "."}
	defaults := map[string]interface{}{
		"key1": "value1",
		"key2": false,
		"key3": 512,
	}
	config := myConfig{}

	// lookup config file in given paths
	// configfile must example name + "." + type
	// i.e. shineyconfig.yaml
	err := venomoid.Config().WithName("shineyconfig").
		WithPath(paths).
		WithType("yaml").
		WithErrorOnMissing(false).
		WithConfigLookup(true).
		WithDefaults(defaults).
		Build(config)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// lookup multiple config file in given paths and merge them
	// configfile must example name + "." + type
	// i.e. shineyconfig.yaml and shineyconfig-secrets.yaml
	err = venomoid.Config().WithName("shineyconfig").
		WithName("shineyconfig-secrets").
		WithPath(paths).
		WithType("yaml").
		WithErrorOnMissing(false).
		WithConfigLookup(true).
		WithDefaults(defaults).
		Build(config)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Load config file directly, will not lookup 
	// configuration file so the WithPaths() method can we skipped
	configFile := "/tmp/my_shiney_config.yaml"
	err = venomoid.Config().WithName("shineyconfig").
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(false).
		WithDefaults(defaults).
		WithFile(configFile).
		Build(config)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Load multiple config files directly and merge them.
	// Will not look up configuration file so the WithPaths() method can we skipped
	configFile1 := "/tmp/my_shiney_config.yaml"
	configFile2 := "/tmp/my_shiney_secrets.yaml"
	err = venomoid.Config().WithName("shineyconfig").
		WithType("yaml").
		WithErrorOnMissing(true).
		WithConfigLookup(false).
		WithDefaults(defaults).
		WithFile(configFile1).
		WithFile(configFile2).
		Build(config)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Load environment variables
	// Using only WithAutomaticEnv, only allows you to access using `viper.Get(envKey)`
	// Using WithBindEnv, allows you to access using `config.envKey`
	envKey := "my_key"
	os.Setenv("MY_KEY", "my_value") // required key in the env
	err = venomoid.Config().WithName("shineyconfig").
		WithErrorOnMissing(true).
		WithDefaults(defaults).
		WithAutomaticEnv(true).
		WithBindEnv(envKey).
		Build(config)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Load environment variables
	// EnvPrefix adds given prefix to all env variables when looking them up
	envKey = "my_key"
	keyPrefix := "pre"
	os.Setenv("PRE_MY_KEY", "my_value") // required key in the env
	err = venomoid.Config().WithName("shineyconfig").
		WithErrorOnMissing(true).
		WithDefaults(defaults).
		WithAutomaticEnv(true).
		WithBindEnv(envKey).
		WithEnvPrefix(keyPrefix).
		Build(config)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// can be also be called using the exported struct direcly
	configBuilder := &venomoid.ConfigBuilder{
		Name:               "shineyconfig",
		Defaults:           defaults,
		AutomaticEnv:       true,
		BindEnv:            []string{envKey},
		EnvPrefix:          keyPrefix,
		ErrorOnMissingFile: true,
	}
	if err := configBuilder.Build(config); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}

```

## Limitation
This module currently only support a subset of functions provides by [viper](https://github.com/spf13/viper).
As this is whatever i require to boilerplate at the time being.
the viper module can still be accessed directly for further (more complex?) configuration. 
However, if one would like to add more functions, you're never further than a PR away

## License
This work is published under the MIT license.

Please see the `LICENSE` file for details.


# venomoid

## Overview
[Venomoid](https://en.wikipedia.org/wiki/Venomoid) is a builder pattern that simplifies the use of [viper](https://github.com/spf13/viper) and reduce
the need for builderplating viper in every project

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


 <pre>
                                          _
 _ __   _____  _____ _ __ __ _ _ __   ___| |_
| '_ \ / _ \ \/ / _ \ '__/ _` | '_ \ / _ \ __|
| | | |  __/>  <  __/ | | (_| | | | |  __/ |_
|_| |_|\___/_/\_\___|_|  \__,_|_| |_|\___|\__|
</pre>

# url_checker

```go
package main

import (
	"fmt"
	"time"
	"github.com/nexeranet/url_checker/pkg/urlchecker"
)

func main() {
    checker := urlchecker.NewURLChecker()
	log.Println("Lister and server :3000")
	if err := checker.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
```

## Pre-requisites

###### Install dependencies

Golang must be installed before working with the repository.

Follow [these instructions](https://go.dev/doc/install) for setting up Golang.

## Build

[Golang](#pre-requisites) must be installed to build a service. Once the dependencies are installed:

- use `go` command
```shell
go build -C ${BUILD_PATH} -o ${BUILD_NAME}
```
Replace `BUILD_PATH` with the path to service folder.
Replace `BUILD_NAME` with the desired name of the compiled binary.

## Run

[Golang](#pre-requisites) must be installed to build a service. Once the dependencies are installed:

- Use terminal:
```shell
./${BUILD_PATH}/${BUILD_NAME}
```
Replace `BUILD_PATH` with the path to service folder.
Replace `BUILD_NAME` with the name of the compiled binary.

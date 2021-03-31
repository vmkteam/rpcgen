# rpcgen: JSON-RPC 2.0 Client Generator Implementation for zenrpc

[![Go Report Card](https://goreportcard.com/badge/github.com/vmkteam/rpcgen)](https://goreportcard.com/report/github.com/vmkteam/rpcgen) [![Go Reference](https://pkg.go.dev/badge/github.com/vmkteam/rpcgen.svg)](https://pkg.go.dev/github.com/vmkteam/rpcgen)

`rpcgen` is a JSON-RPC 2.0 client library generator for [zenrpc](https://github.com/vmkteam/zenrpc). It supports client generation for following languages:
- Golang
- TypeScript

## Examples

### Basic usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/vmkteam/rpcgen"
	"github.com/vmkteam/zenrpc/v2"
)

func main() {
	rpc := zenrpc.NewServer(zenrpc.Options{})

	generated, err := rpcgen.FromSMD(rpc.SMD()).GoClient().Generate()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%s", generated)
}
```

### Generate in HTTP handler

```go
package main 

import (
	"net/http"
	
	"github.com/vmkteam/rpcgen"
	"github.com/vmkteam/zenrpc/v2"
)

func main () {
	rpc := zenrpc.NewServer(zenrpc.Options{})

	gen := rpcgen.FromSMD(rpc.SMD())

	http.HandleFunc("/client.go", rpcgen.Handler(gen.GoClient()))
	http.HandleFunc("/client.ts", rpcgen.Handler(gen.TSClient(nil)))
}
```

### Add custom TypeScript type mapper

```go
package main

import (
	"net/http"

	"github.com/vmkteam/rpcgen"
	"github.com/vmkteam/rpcgen/typescript"
	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/smd"
)

func main() {
	rpc := zenrpc.NewServer(zenrpc.Options{})

	gen := rpcgen.FromSMD(rpc.SMD())

	typeMapper := func(in smd.JSONSchema, tsType typescript.Type) typescript.Type {
		if in.Type == "object" {
			if in.Description == "Group" && in.Name == "groups" {
				tsType.Type = fmt.Sprintf("Record<number, I%s>", in.Description)
			}
		}
		
		return tsType
	}

	http.HandleFunc("/client.ts", rpcgen.Handler(gen.TSClient(typeMapper)))
}
```

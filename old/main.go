// package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/rpc"

// 	"qyLsp/lsp"
// )

// type EmailArgs struct {
// 	To, Subject, Content string
// }

// type Response struct {
// 	Result string
// }

// type SmsService struct{}
// type EmailService struct{}

// func (t *SmsService) SendSMS(r *http.Request, args *SmsArgs, result *Response) error {
// 	*result = Response{Result: fmt.Sprintf("Sms sent to %s", args.Number)}
// 	return nil
// }

// func (t *EmailService) SendEmail(r *http.Request, args *EmailArgs, result *Response) error {
// 	*result = Response{Result: fmt.Sprintf("Email sent to %s", args.To)}
// 	return nil
// }

// type Initialize struct{}

// func (i *Initialize) Get(payload int, reply *lsp.InitializeResult) error {
// 	reply = &lsp.InitializeResult{
// 		Capabilities: lsp.ServerCapabilities{
// 			HoverProvider: true,
// 		},
// 	}
// 	return nil
// }

// type TextDocument struct{}

// func main() {
// 	rpcServer := rpc.NewServer()

// 	initialize := new(Initialize)
// 	rpcServer.Register(initialize)

// 	rpc.HandleHTTP()

// 	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
// 		io.WriteString(res, "RPC SERVER LIVE!")
// 	})

// 	http.ListenAndServe(":3000", nil)
// }

// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2015 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"qyLsp/lsp"

	"github.com/AdamSLevy/jsonrpc2/v14"
)

type Logger struct{}

func (l *Logger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (l *Logger) Println(a ...interface{}) {
	fmt.Println(a...)
}

func version() jsonrpc2.MethodFunc {
	return func(ctx context.Context, params json.RawMessage) interface{} {
		return "1"
	}
}

type SmsArgs struct {
	Number, Content string
}

func sms() jsonrpc2.MethodFunc {
	return func(ctx context.Context, params json.RawMessage) interface{} {
		var sms SmsArgs
		if err := json.Unmarshal(params, &sms); err != nil {
			return jsonrpc2.ErrorInvalidParams(err)
		}
		if sms.Number == "0" {
			return jsonrpc2.NewError(
				jsonrpc2.ErrorCodeInvalidParams,
				"phone number not valid",
				sms.Number,
			)
		}
		return fmt.Sprintf("Sms sent to %s", sms.Number)
	}
}

func onInitialize() jsonrpc2.MethodFunc {
	return func(ctx context.Context, params json.RawMessage) interface{} {
		var args lsp.InitializeParams

		if params != nil {
			if err := json.Unmarshal(params, &args); err != nil {
				return jsonrpc2.ErrorInvalidParams(err)
			}
		}

		return &lsp.InitializeResult{
			Capabilities: lsp.ServerCapabilities{
				HoverProvider: true,
			},
		}
	}
}

func onHover() jsonrpc2.MethodFunc {
	return func(ctx context.Context, params json.RawMessage) interface{} {
		var args lsp.Hover

		if params != nil {
			if err := json.Unmarshal(params, &args); err != nil {
				return jsonrpc2.ErrorInvalidParams(err)
			}
		}

		return &lsp.Hover{
			Contents: []lsp.MarkedString{
				{
					Language: "Quanty",
					Value:    "Hello World",
				},
			},
		}
	}
}

// func initialize() jsonrpc2.MethodFunc {
// 	return func(ctx context.Context, params json.RawMessage) interface{} {
// 		var u User
// 		if err := json.Unmarshal(params, &u); err != nil {
// 			return jsonrpc2.InvalidParams(err)
// 		}
// 		conn, err := mydbpkg.GetDBConn()
// 		if err != nil {
// 			// The handler will recover, print debug info if enabled, and
// 			// return an Internal Error to the client.
// 			panic(err)
// 		}
// 		if err := u.Select(conn); err != nil {
// 			return jsonrpc2.NewError(-30000, "user not found", u.ID)
// 		}
// 		return u
// 	}
// }

func main() {
	logger := log.New(os.Stdout, "[LSP] ", log.LstdFlags)
	logger.Println("Server is starting...")

	methods := jsonrpc2.MethodMap{
		"version":       version(),
		"sms":           sms(),
		"initialize":    onInitialize(),
		"hoverProvider": onHover(),
	}

	http.ListenAndServe(":6012", jsonrpc2.HTTPRequestHandler(methods, logger))
}

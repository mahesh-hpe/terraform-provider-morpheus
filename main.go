// Copyright (c) 2019 Morpheus Data https://www.morpheusdata.com, All rights reserved.
// terraform-provider-morpheus source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/gomorpheus/terraform-provider-morpheus/morpheusv3"

	"github.com/gomorpheus/terraform-provider-morpheus/morpheus"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func main() {
	// plugin.Serve(&plugin.ServeOpts{
	// 	ProviderFunc: func() *schema.Provider {
	// 		return morpheus.Provider()
	// 	},
	// })
	// providerserver.Serve(context.Background())
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(morpheusv3.New("dev")()), // Example terraform-plugin-framework provider
		morpheus.Provider().GRPCProvider,                     // Example terraform-plugin-sdk provider
	}

	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/gomorpheus/morpheus",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}

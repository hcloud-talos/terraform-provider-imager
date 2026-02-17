package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hcloud-talos/terraform-provider-imager/internal/provider"
)

var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/hcloud-talos/imager",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err)
	}
}

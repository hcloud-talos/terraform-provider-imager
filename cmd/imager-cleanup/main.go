package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func main() {
	token := os.Getenv("HCLOUD_TOKEN")
	if token == "" {
		log.Fatal("HCLOUD_TOKEN must be set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	hcloudClient := hcloud.NewClient(hcloud.WithToken(token))
	client := hcloudimages.NewClient(hcloudClient)

	if err := client.CleanupTempResources(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("cleanup completed")
}

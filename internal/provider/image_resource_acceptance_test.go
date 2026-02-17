package provider_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	imagerprovider "github.com/hcloud-talos/terraform-provider-imager/internal/provider"
)

func TestAccImagerImage_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}
	if os.Getenv("HCLOUD_TOKEN") == "" {
		t.Skip("HCLOUD_TOKEN is not set")
	}
	imageURL := os.Getenv("IMAGER_TEST_IMAGE_URL")
	if imageURL == "" {
		t.Skip("IMAGER_TEST_IMAGE_URL is not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"imager": providerserver.NewProtocol6WithError(imagerprovider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(imageURL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("imager_image.test", "id", regexp.MustCompile(`^[0-9]+$`)),
				),
			},
		},
	})
}

func testAccConfig(imageURL string) string {
	return fmt.Sprintf(`
provider "imager" {}

resource "imager_image" "test" {
  image_url     = %q
  architecture  = "x86"
}
`, imageURL)
}

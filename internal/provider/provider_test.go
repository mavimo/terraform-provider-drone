package drone

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testDroneUser string = os.Getenv("DRONE_USER")
	testProviders map[string]*schema.Provider
	testProvider  *schema.Provider
)

func init() {
	testProvider = Provider()
	testProviders = map[string]*schema.Provider{
		"drone": testProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DRONE_SERVER"); v == "" {
		t.Fatal("DRONE_SERVER must be set for acceptance tests")
	}
	if v := os.Getenv("DRONE_TOKEN"); v == "" {
		t.Fatal("DRONE_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("DRONE_USER"); v == "" {
		t.Fatal("DRONE_USER must be set for acceptance tests")
	}
}

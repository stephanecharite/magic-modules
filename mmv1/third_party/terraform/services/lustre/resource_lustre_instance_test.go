package lustre_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccLustreInstance_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "default-vpc"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLustreInstance_full(context),
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccLustreInstance_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_lustre_instance.description",
							plancheck.ResourceActionUpdate,
						),
						plancheck.ExpectResourceAction(
							"google_lustre_instance.labels",
							plancheck.ResourceActionUpdate,
						),
						plancheck.ExpectResourceAction(
							"google_lustre_instance.gke_support_enabled",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccLustreInstance_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  instance_id  = "tf-test-my-instance%{random_suffix}"
  location     = "us-central1-a"
  filesystem   = "testfs"
  capacity_gib = 18000
  network      = data.google_compute_network.lustre-network.id
	timeouts {
		create = "240m"
	}
}

// This example assumes this network already exists.
// The API creates a tenant network per network authorized for a
// Lustre instance and that network is not deleted when the user-created
// network (authorized_network) is deleted, so this prevents issues
// with tenant network quota.
// If this network hasn't been created and you are using this example in your
// config, add an additional network resource or change
// this from "data"to "resource"
data "google_compute_network" "lustre-network" {
  name = "%{network_name}"
}
`, context)
}

func testAccLustreInstance_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  instance_id         = "tf-test-my-instance%{random_suffix}"
  location            = "us-central1-a"
  description         = "description updated"
  filesystem          = "testfs"
  capacity_gib        = 18000
  network             = data.google_compute_network.lustre-network.id
  gke_support_enabled = true
  labels              = {
    test = "newLabel"
  }
	timeouts {
		create = "240m"
  }
}

// This example assumes this network already exists.
// The API creates a tenant network per network authorized for a
// Lustre instance and that network is not deleted when the user-created
// network (authorized_network) is deleted, so this prevents issues
// with tenant network quota.
// If this network hasn't been created and you are using this example in your
// config, add an additional network resource or change
// this from "data"to "resource"
data "google_compute_network" "lustre-network" {
  name = "%{network_name}"
}
`, context)
}

package codefresh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCodefreshGitopsEnvironmentsResource(t *testing.T) {
	resourceName := "codefresh_gitops_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCodefreshGitopsEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCodefreshGitopsEnvironmentConfig(
					"test-for-tf",
					"NON_PROD",
					[]cfclient.GitopsEnvironmentCluster{{
						Name:        "in-cluster2",
						RuntimeName: "test-runtime",
						Namespaces:  []string{"test-ns-1", "test-ns2"},
					}},
					[]string{"codefresh.io/environment=test-for-tf", "codefresh.io/environment-1=test-for-tf1"},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshGitopsEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-for-tf"),
					resource.TestCheckResourceAttr(resourceName, "kind", "NON_PROD"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.name", "in-cluster2"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.server", "https://kubernetes.default.svc"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.runtime_name", "test-runtime"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.namespaces.0", "test-ns-1"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.namespaces.1", "test-ns2"),
					resource.TestCheckResourceAttr(resourceName, "label_pairs.0", "codefresh.io/environment=test-for-tf"),
					resource.TestCheckResourceAttr(resourceName, "label_pairs.1", "codefresh.io/environment-1=test-for-tf1"),
				),
			},
			{
				Config: testAccCodefreshGitopsEnvironmentConfig(
					"test-for-tf",
					"NON_PROD",
					[]cfclient.GitopsEnvironmentCluster{
						{
							Name:        "in-cluster2",
							RuntimeName: "test-runtime",
							Namespaces:  []string{"test-ns-1", "test-ns2"},
						},
						{
							Name:        "in-cluster3",
							RuntimeName: "test-runtime-2",
							Namespaces:  []string{"test-ns-3"},
						},
					},
					[]string{"codefresh.io/environment=test-for-tf", "codefresh.io/environment-1=test-for-tf1"},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshGitopsEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-for-tf"),
					resource.TestCheckResourceAttr(resourceName, "kind", "NON_PROD"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.name", "in-cluster2"),
					resource.TestCheckResourceAttr(resourceName, "cluster.1.name", "in-cluster3"),
					resource.TestCheckResourceAttr(resourceName, "cluster.1.server", "https://kubernetes2.default.svc"),
					resource.TestCheckResourceAttr(resourceName, "cluster.1.runtime_name", "test-runtime-2"),
					resource.TestCheckResourceAttr(resourceName, "cluster.1.namespaces.0", "test-ns-3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCodefreshGitopsEnvironmentExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		envID := rs.Primary.ID
		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetGitopsEnvironmentById(envID)
		if err != nil {
			return fmt.Errorf("error fetching gitops environment with ID %s. %s", envID, err)
		}
		return nil
	}
}

func testAccCheckCodefreshGitopsEnvironmentDestroy(state *terraform.State) error {
	// Implement destroy check if needed
	return nil
}

// CONFIG
func testAccCodefreshGitopsEnvironmentConfig(name, kind string, clusters []cfclient.GitopsEnvironmentCluster, labelPairs []string) string {
	var clusterBlocks []string
	for _, c := range clusters {
		ns := fmt.Sprintf("[\"%s\"]", strings.Join(c.Namespaces, "\", \""))
		block := fmt.Sprintf(`  cluster {
    name        = "%s"
    runtime_name = "%s"
    namespaces  = %s
  }`, c.Name, c.RuntimeName, ns)
		clusterBlocks = append(clusterBlocks, block)
	}
	labelsStr := fmt.Sprintf("[\"%s\"]", strings.Join(labelPairs, "\", \""))
	return fmt.Sprintf(`
resource "codefresh_gitops_environment" "test" {
  name        = "%s"
  kind        = "%s"
%s
  label_pairs = %s
}
`, name, kind, strings.Join(clusterBlocks, "\n"), labelsStr)
}

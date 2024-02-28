package codefresh

import (
	"fmt"
	"testing"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/cfclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Check create, update and delete of all supported IDP types
func TestAccountIDPCodefreshProject_AllSupportedTypes(t *testing.T) {
	uniqueId := acctest.RandString(10)
	resourceName := "codefresh_account_idp.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckCodefreshProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccountIDPCodefreshConfig("onelogin", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-onelogin-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "onelogin"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// For existing resources codefresh returns the secrets in encrypted format, to avoid constant diff we store those in _encrypted, hence on import the secrets will be empty
				ImportStateVerifyIgnore: []string{"onelogin.0.client_secret"},
			},
			{
				Config: testAccountIDPCodefreshConfig("auth0", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-auth0-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "auth0"),
				),
			},
			{
				Config: testAccountIDPCodefreshConfig("azure", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-azure-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "azure"),
				),
			},
			{
				Config: testAccountIDPCodefreshConfig("google", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-google-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "google"),
				),
			},
			{
				Config: testAccountIDPCodefreshConfig("keycloak", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-keycloak-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "keycloak"),
				),
			},
			{
				Config: testAccountIDPCodefreshConfig("ldap", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-ldap-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "ldap"),
				),
			},
			{
				Config: testAccountIDPCodefreshConfig("okta", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-okta-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "okta"),
				),
			},
			{
				Config: testAccountIDPCodefreshConfig("saml", uniqueId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodefreshAccountIDPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("tf-test-saml-%s", uniqueId)),
					resource.TestCheckResourceAttr(resourceName, "client_type", "saml"),
				),
			},
		},
	})
}

func testAccCheckCodefreshAccountIDPExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		idpID := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*cfclient.Client)
		_, err := apiClient.GetAccountIdpByID(idpID)

		if err != nil {
			return fmt.Errorf("error fetching project with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccountIDPCodefreshConfig(idpType string, uniqueId string) string {

	idpResource := ""

	if idpType == "onelogin" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" { 
			display_name = "tf-test-onelogin-%s"

			onelogin {
				client_id = "onelogin-%s"
				client_secret = "myoneloginsecret1"
				domain = "myonelogindomain"
				app_id = "myappid"
				api_client_id = "myonelogindomain"
				api_client_secret = "myapiclientsecret1"
			}
		}`, uniqueId, uniqueId)
	}

	if idpType == "auth0" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" {
			display_name = "tf-test-auth0-%s"
			name = "tf-auth0-test34"
			auth0 {
			  client_id = "blah-auth0-%s"
			  client_secret = "asdddd"
			  domain = "codefresh.auth0.com"
			}
		  }`, uniqueId, uniqueId)
	}

	if idpType == "azure" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" {
			display_name = "tf-test-azure-%s"
			name = "tf-azure-test"
			azure {
			  app_id = "azure-codefresh-test-%s"
			  client_secret = "mysecret99"
			  object_id = "myobjectidtest"
			  autosync_teams_and_users = true
			  sync_interval = 7
			}
		  }`, uniqueId, uniqueId)
	}

	if idpType == "google" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" {
			display_name = "tf-test-google-%s"
			name = "tf-google-test"
			google {
			  client_id = "tf-test-google-%s"
			  client_secret = "mysecret99"
			  admin_email = "admin@codefresh.io"
			  sync_field = "myfield"
			  json_keyfile = <<EOT
			  {  
				"installed":{  
					"client_id":"clientid",
					"project_id":"projectname",
					"auth_uri":"https://accounts.google.com/o/oauth2/auth",
					"token_uri":"https://accounts.google.com/o/oauth2/token",
					"auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
				}
			  }
			  EOT
			}
		  }`, uniqueId, uniqueId)
	}

	if idpType == "keycloak" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" {
			display_name = "tf-test-keycloak-%s"
			name = "tf-keycloak-test1"
		  
			keycloak {
			  client_id = "tf-test-keycloak-%s"
			  client_secret = "mykeycloaksecret1"
			  realm = "myrealm"
			  host = "https://myhost.com"
			}
		  }`, uniqueId, uniqueId)
	}

	if idpType == "ldap" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" {
			display_name = "tf-test-ldap-%s"
			
			ldap {
			  url = "ldaps://myldap.server.com:389"
			  password = "mypassword"
			  distinguished_name = "cn=admin,dc=mydomain,dc=com"
			  search_base = "ou=people,dc=mydomain,dc=com"
			  search_filter = "sAMAccountName"
			  allowed_groups_for_sync = "mygroup"
			  search_base_for_sync = "ou=codefresh-users,ou=people,dc=mydomain,dc=com"
			  certificate = <<EOT
		  -----BEGIN CERTIFICATE-----
		  MIIGKjCCBRKgAwIBAgIQAWjXYVVrFZwRKaeAGTxhSTANBgkqhkiG9w0BAQsFADBg
		  MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
		  d3cuZGlnaWNlcnQuY29tMR8wHQYDVQQDExZHZW9UcnVzdCBUTFMgUlNBIENBIEcx
		  UhEEs2oMtYNbp2jzJKpcv13MpJPdKqTPU6Jb7hsOTKNjLSpPL4QhpezZ5sYhzg==
		  -----END CERTIFICATE-----
			  EOT
			}
		  }`, uniqueId)
	}

	if idpType == "okta" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" {
			display_name = "tf-test-okta-%s"
			name = "tf-okta-test3"

			okta {
			  client_id = "tf-test-okta-%s"
			  client_secret = "asdddd"
			  client_host = "http://asdd.okta.com"
			  app_id = "test1"
			  access_token = "myaccesstoken1"
			}
		  }`, uniqueId, uniqueId)
	}

	if idpType == "saml" {
		idpResource = fmt.Sprintf(` 
		resource "codefresh_account_idp" "test" {
			display_name = "tf-test-saml-%s"
			name = "tf-samlt-okta-test1"
		  
			saml {
			  provider = "okta"
			  endpoint = "https://example.com/endpoint"
			  access_token = "myaccesstoken"
			  autosync_teams_and_users = true
			  activate_users_after_sync = true
			  sync_interval = 7
			  client_host = "https://codefresh-example.okta.com"
			  app_id = "tf-test-saml-%s"
			  application_certificate = <<EOT
		  -----BEGIN CERTIFICATE-----
		  MIIGKjCCBRKgAwIBAgIQAWjXYVVrFZwRKaeAGTxhSTANBgkqhkiG9w0BAQsFADBg
		  MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
		  d3cuZGlnaWNlcnQuY29tMR8wHQYDVQQDExZHZW9UcnVzdCBUTFMgUlNBIENBIEcx
		  MB4XDTIzMTIyNTAwMDAwMFoXDTI1MDExMTIzNTk1OVowGTEXMBUGA1UEAwwOKi5j
		  0PXGSsbWQTI/ItYe1yK1VebDODp3Cx7IeMJqpsUI0IjxdPjQObsHKqmVNAd8kOXi
		  UhEEs2oMtYNbp2jzJKpcv13MpJPdKqTPU6Jb7hsOTKNjLSpPL4QhpezZ5sYhzg==
		  -----END CERTIFICATE-----
			  EOT
			}
		  }`, uniqueId, uniqueId)
	}

	return idpResource
}

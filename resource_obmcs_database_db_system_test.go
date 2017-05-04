// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"bytes"
	"testing"
	"text/template"
	"time"

	"github.com/MustWin/baremetal-sdk-go"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/oracle/terraform-provider-baremetal/client/mocks"
	"github.com/stretchr/testify/suite"
)

type DatabaseDBSystemTestSuite struct {
	suite.Suite
	Client       *mocks.BareMetalClient
	Provider     terraform.ResourceProvider
	Providers    map[string]terraform.ResourceProvider
	TimeCreated  baremetal.Time
	Config       string
	ResourceName string
	Res          *baremetal.DBSystem
	DeletedRes   *baremetal.DBSystem
}

func (s *DatabaseDBSystemTestSuite) SetupTest() {
	s.Client = &mocks.BareMetalClient{}
	s.TimeCreated = baremetal.Time{Time: time.Now()}

	s.Provider = Provider(
		func(d *schema.ResourceData) (interface{}, error) { return s.Client, nil },
	)
	s.Providers = map[string]terraform.ResourceProvider{"baremetal": s.Provider}

	dbHomeOpts := &baremetal.DisplayNameOptions{}
	dbHomeOpts.DisplayName = "db_home_display_name"
	dbHome := baremetal.NewCreateDBHomeDetails(
		"admin_password", "db_name", "db_version", dbHomeOpts,
	)
	opts := &baremetal.LaunchDBSystemOptions{}
	opts.DisplayName = "display_name"
	opts.DatabaseEdition = baremetal.DatabaseEditionStandard
	opts.DBHome = dbHome
	opts.DiskRedundancy = baremetal.DiskRedundancyNormal
	opts.Domain = "domain.com"
	opts.Hostname = "hostname"

	s.Res = &baremetal.DBSystem{
		AvailabilityDomain: "availability_domain",
		CompartmentID:      "compartment_id",
		Shape:              "shape",
		SubnetID:           "subnet_id",
		SSHPublicKeys:      []string{"ansshkey"},
		CPUCoreCount:       2,

		DisplayName:     opts.DisplayName,
		DatabaseEdition: opts.DatabaseEdition,
		DBHome:          opts.DBHome,
		DiskRedundancy:  opts.DiskRedundancy,
		Domain:          opts.Domain,
		Hostname:        opts.Hostname,

		ID:               "id",
		LifecycleDetails: "lifecycle_details",
		ListenerPort:     1,
		State:            baremetal.ResourceAvailable,
		TimeCreated:      s.TimeCreated,
	}
	s.Res.ETag = "etag"
	s.Res.RequestID = "opcrequestid"
	s.Client.On("LaunchDBSystem",
		s.Res.AvailabilityDomain, s.Res.CompartmentID, s.Res.Shape, s.Res.SubnetID,
		s.Res.SSHPublicKeys, s.Res.CPUCoreCount, opts,
	).Return(s.Res, nil)

	deletedRes := *s.Res
	s.DeletedRes = &deletedRes
	s.DeletedRes.State = baremetal.ResourceTerminated
	s.Client.On("TerminateDBSystem", "id", (*baremetal.IfMatchOptions)(nil)).Return(nil)

	tmpl := `
		resource "baremetal_database_db_system" "t" {
			availability_domain = "{{.AvailabilityDomain}}"
			compartment_id = "{{.CompartmentID}}"
			shape = "{{.Shape}}"
			subnet_id = "{{.SubnetID}}"
			ssh_public_keys = ["{{index .SSHPublicKeys 0}}"]
			cpu_core_count = {{.CPUCoreCount}}
			display_name = "{{.DisplayName}}"
			database_edition = "{{.DatabaseEdition}}"
			db_home {
				database {
					"admin_password" = "{{.DBHome.Database.AdminPassword}}"
					"db_name" = "{{.DBHome.Database.DBName}}"
				}
				db_version = "{{.DBHome.DBVersion}}"
				display_name = "{{.DBHome.DisplayName}}"
			}
			disk_redundancy = "{{.DiskRedundancy}}"
			domain = "{{.Domain}}"
			hostname = "{{.Hostname}}"
		}
	`
	var buf bytes.Buffer
	parsed := template.Must(template.New("config").Parse(tmpl))
	parsed.Execute(&buf, s.Res)
	s.Config = buf.String()
	s.Config += testProviderConfig

	s.ResourceName = "baremetal_database_db_system.t"
}

func (s *DatabaseDBSystemTestSuite) TestCreateDBSystem() {
	s.Client.On("GetDBSystem", "id").Return(s.Res, nil).Times(2)
	s.Client.On("GetDBSystem", "id").Return(s.DeletedRes, nil)

	resource.UnitTest(s.T(), resource.TestCase{
		Providers: s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            s.Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(s.ResourceName, "compartment_id", s.Res.CompartmentID),
					resource.TestCheckResourceAttr(s.ResourceName, "db_home.0.db_version", s.Res.DBHome.DBVersion),
					resource.TestCheckResourceAttr(s.ResourceName, "db_home.0.database.0.db_name", s.Res.DBHome.Database.DBName),
				),
			},
		},
	})
}

func (s *DatabaseDBSystemTestSuite) TestTerminateDBSystem() {
	s.Client.On("GetDBSystem", "id").Return(s.Res, nil).Times(2)
	s.Client.On("GetDBSystem", "id").Return(s.DeletedRes, nil)

	resource.UnitTest(s.T(), resource.TestCase{
		Providers: s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            s.Config,
			},
			{
				Config:  s.Config,
				Destroy: true,
			},
		},
	})

	s.Client.AssertCalled(s.T(), "TerminateDBSystem", "id", (*baremetal.IfMatchOptions)(nil))
}

func TestDatabaseDBSystemTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseDBSystemTestSuite))
}

func TestAccOBMCSDBSystem_thing(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccOBMASDBSystem_basic(ri)

	resourceName := "baremetal_database_db_system.t"

	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", "ocid1.compartment.FIXME"),
				),
			},
		},
	})
}

func testAccOBMASDBSystem_basic(ri int) string {
	config := `
provider "baremetal" {
  tenancy_ocid     = "ocid1.tenancy.FIXME"
  user_ocid        = "ocid1.user.FIXME"
  fingerprint      = "FIXME_fingerprint"
  private_key_path = "auth_key"
}

resource "baremetal_database_db_system" "TFDBNode" {
  domain              = "enkitec"
  availability_domain = "foTc:PHX-AD-1"
  display_name        = "RandomDBSystem"
  compartment_id      = "ocid1.compartment.FIXME"
  subnet_id           = "ocid1.subnet.FIXME"
  hostname            = "node1"
  database_edition    = "ENTERPRISE_EDITION_EXTREME_PERFORMANCE"

  ssh_public_keys = [
    FIXME_ssh_public_key",
  ]

  shape           = "BM.DenseIO1.36"
  disk_redundancy = "HIGH"
  cpu_core_count  = "2"

  db_home = {
    db_version   = "12.1.0.2"
    display_name = "DBfromAPI"

    database {
      db_name        = "TESTBMC"
      admin_password = "stub_admin_password"
    }
  }

  lifecycle {
    ignore_changes = [
      "db_home.0.database.0.admin_password",
      "db_home.0.database.0.db_name",
      "db_home.0.db_version",
      "db_home.0.display_name",
    ]
  }
}

data "baremetal_identity_availability_domains" "ADs" {
  compartment_id = "${var.tenancy_ocid}"
}

data "baremetal_database_db_node" "DBNodeDetails" {
  db_node_id = "${baremetal_database_db_system.TFDBNode.id}"
}

data "baremetal_core_vnic" "DBNodeVnic" {
  vnic_id = "${lookup(data.baremetal_database_db_node.DBNodeDetails.vnic_id,"vnic_id")}"
}

output "instance_id" {
  value = "${baremetal_database_db_system.TFDBNode.id}"
}

output "DBNodePublicIP" {
  value = ["${data.baremetal_core_vnic.DBNodeVnic.public_ip_address}"]
}

output "DBNodePrivateIP" {
  value = ["${data.baremetal_core_vnic.DBNodeVnic.private_ip_address}"]
}
`
	return config
}

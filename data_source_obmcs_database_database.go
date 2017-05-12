// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"github.com/MustWin/baremetal-sdk-go"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/oracle/terraform-provider-baremetal/client"
	"github.com/oracle/terraform-provider-baremetal/crud"
)

func DatabaseDatasource() *schema.Resource {
	return &schema.Resource{
		Read: readDatabase,
		Schema: map[string]*schema.Schema{
			"compartment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_home_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_unique_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readDatabase(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(client.BareMetalClient)
	sync := &DatabaseDatasourceCrud{}
	sync.D = d
	sync.Client = client
	return crud.ReadResource(sync)
}

type DatabaseDatasourceCrud struct {
	crud.BaseCrud
	Res *baremetal.Database
}

func (s *DatabaseDatasourceCrud) Get() (e error) {
	id := s.D.Get("database_id").(string)
	s.Res, e = s.Client.GetDatabase(id)
	return
}

func (s *DatabaseDatasourceCrud) SetData() {
	if s.Res != nil {
		s.D.SetId(s.Res.ID)
		s.D.Set("compartment_id", s.Res.CompartmentID)
		s.D.Set("db_home_id", s.Res.DBHomeID)
		s.D.Set("db_name", s.Res.DBName)
		s.D.Set("db_unique_name", s.Res.DBUniqueName)
		s.D.Set("id", s.Res.ID)
		s.D.Set("state", s.Res.State)
		s.D.Set("time_created", s.Res.TimeCreated.String())
	}
	return
}

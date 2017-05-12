// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"time"

	"github.com/MustWin/baremetal-sdk-go"

	"github.com/oracle/terraform-provider-baremetal/crud"
)

type LoadBalancerShapeDatasourceCrud struct {
	crud.BaseCrud
	Res *baremetal.ListLoadBalancerShapes
}

func (s *LoadBalancerShapeDatasourceCrud) Get() (e error) {
	cID := s.D.Get("compartment_id").(string)
	s.Res, e = s.Client.ListLoadBalancerShapes(cID, nil)
	return
}

func (s *LoadBalancerShapeDatasourceCrud) SetData() {
	if s.Res != nil {
		s.D.SetId(time.Now().UTC().String())
		resources := []map[string]interface{}{}

		for _, v := range s.Res.LoadBalancerShapes {
			res := map[string]interface{}{
				"name": v.Name,
			}
			resources = append(resources, res)

		}
		s.D.Set("shapes", resources)
	}
	return
}

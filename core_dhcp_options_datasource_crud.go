// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"time"

	"github.com/MustWin/baremetal-sdk-go"

	"github.com/oracle/terraform-provider-baremetal/crud"
	"github.com/oracle/terraform-provider-baremetal/options"
)

type DHCPOptionsDatasourceCrud struct {
	crud.BaseCrud
	Res *baremetal.ListDHCPOptions
}

func (s *DHCPOptionsDatasourceCrud) Get() (e error) {
	compartmentID := s.D.Get("compartment_id").(string)
	vcnID := s.D.Get("vcn_id").(string)

	opts := &baremetal.ListOptions{}
	options.SetListOptions(s.D, opts)

	s.Res = &baremetal.ListDHCPOptions{DHCPOptions: []baremetal.DHCPOptions{}}

	for {
		var list *baremetal.ListDHCPOptions
		if list, e = s.Client.ListDHCPOptions(compartmentID, vcnID, opts); e != nil {
			break
		}

		s.Res.DHCPOptions = append(s.Res.DHCPOptions, list.DHCPOptions...)

		if hasNextPage := options.SetNextPageOption(list.NextPage, &opts.PageListOptions); !hasNextPage {
			break
		}
	}

	return
}

func (s *DHCPOptionsDatasourceCrud) SetData() {
	if s.Res != nil {
		s.D.SetId(time.Now().UTC().String())

		stateObjs := []map[string]interface{}{}
		for _, res := range s.Res.DHCPOptions {

			nestedStateObjs := []map[string]interface{}{}

			for _, nestedRes := range res.Options {
				nestedStateObj := map[string]interface{}{
					"type":               nestedRes.Type,
					"custom_dns_servers": nestedRes.CustomDNSServers,
					"server_type":        nestedRes.ServerType,
				}
				nestedStateObjs = append(nestedStateObjs, nestedStateObj)
			}

			stateObj := map[string]interface{}{
				"compartment_id": res.CompartmentID,
				"display_name":   res.DisplayName,
				"id":             res.ID,
				"options":        nestedStateObjs,
				"state":          res.State,
				"time_created":   res.TimeCreated.String(),
			}
			stateObjs = append(stateObjs, stateObj)
		}
		s.D.Set("options", stateObjs)
	}
	return
}

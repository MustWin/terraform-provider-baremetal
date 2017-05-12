// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"time"

	"github.com/MustWin/baremetal-sdk-go"

	"github.com/oracle/terraform-provider-baremetal/crud"
)

type NamespaceDatasourceCrud struct {
	crud.BaseCrud
	Res *baremetal.Namespace
}

func (s *NamespaceDatasourceCrud) Get() (e error) {
	s.Res, e = s.Client.GetNamespace()
	return
}

func (s *NamespaceDatasourceCrud) SetData() {
	if s.Res != nil {
		// Important, if you don't have an ID, make one up for your datasource
		// or things will end in tears
		s.D.SetId(time.Now().UTC().String())
		s.D.Set("namespace", string(*s.Res))
	}
	return
}

// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"time"

	"github.com/MustWin/baremetal-sdk-go"

	"github.com/oracle/terraform-provider-baremetal/crud"
)

type ObjectHeadDatasourceCrud struct {
	crud.BaseCrud
	Res *baremetal.HeadObject
}

func (s *ObjectHeadDatasourceCrud) Get() (e error) {
	namespace := s.D.Get("namespace").(string)
	bucket := s.D.Get("bucket").(string)
	object := s.D.Get("object").(string)

	s.Res, e = s.Client.HeadObject(baremetal.Namespace(namespace), bucket, object, &baremetal.HeadObjectOptions{})
	return
}

func (s *ObjectHeadDatasourceCrud) SetData() {
	// Important, if you don't have an ID, make one up for your datasource
	// or things will end in tears
	s.D.SetId(time.Now().UTC().String())
	s.D.Set("metadata", s.Res.Metadata)
	s.D.Set("content-length", string(s.Res.ContentLength))
	s.D.Set("content-type", s.Res.ContentType)
	return
}

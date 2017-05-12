// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"github.com/MustWin/baremetal-sdk-go"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/oracle/terraform-provider-baremetal/client"
	"github.com/oracle/terraform-provider-baremetal/crud"
)

func ObjectResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: crud.DefaultTimeout,
		Create:   createObject,
		Read:     readObject,
		Update:   updateObject,
		Delete:   deleteObject,
		Schema:   objectSchema,
	}
}

func createObject(d *schema.ResourceData, m interface{}) (e error) {
	sync := &ObjectResourceCrud{}
	sync.D = d
	sync.Client = m.(client.BareMetalClient)
	return crud.CreateResource(d, sync)
}

func readObject(d *schema.ResourceData, m interface{}) (e error) {
	sync := &ObjectResourceCrud{}
	sync.D = d
	sync.Client = m.(client.BareMetalClient)
	return crud.ReadResource(sync)
}

func updateObject(d *schema.ResourceData, m interface{}) (e error) {
	sync := &ObjectResourceCrud{}
	sync.D = d
	sync.Client = m.(client.BareMetalClient)
	return crud.UpdateResource(d, sync)
}

func deleteObject(d *schema.ResourceData, m interface{}) (e error) {
	sync := &ObjectResourceCrud{}
	sync.D = d
	sync.Client = m.(client.BareMetalClient)
	return crud.DeleteResource(d, sync)
}

type ObjectResourceCrud struct {
	crud.BaseCrud
	Res *baremetal.Object
}

func (s *ObjectResourceCrud) ID() string {
	return "tfobm-object-" + string(s.Res.Namespace) + "/" + s.Res.Bucket + "/" + s.Res.ID
}

func (s *ObjectResourceCrud) SetData() {
	s.D.Set("namespace", s.Res.Namespace)
	s.D.Set("bucket", s.Res.Bucket)
	s.D.Set("object", s.Res.ID)
	s.D.Set("content", s.Res.Body)
	s.D.Set("metadata", s.Res.Metadata)
}

func (s *ObjectResourceCrud) Create() (e error) {
	e = s.Update()
	return
}

func (s *ObjectResourceCrud) Get() (e error) {
	namespace := s.D.Get("namespace").(string)
	bucket := s.D.Get("bucket").(string)
	object := s.D.Get("object").(string)
	s.Res, e = s.Client.GetObject(baremetal.Namespace(namespace), bucket, object, &baremetal.GetObjectOptions{})
	return
}

func (s *ObjectResourceCrud) Update() (e error) {
	namespace := s.D.Get("namespace").(string)
	bucket := s.D.Get("bucket").(string)
	object := s.D.Get("object").(string)
	content := s.D.Get("content").(string)
	opts := &baremetal.PutObjectOptions{}

	if rawMetadata, ok := s.D.GetOk("metadata"); ok {
		metadata := resourceObjectStorageMapToMetadata(rawMetadata.(map[string]interface{}))
		opts.Metadata = metadata
	}
	s.Res, e = s.Client.PutObject(baremetal.Namespace(namespace), bucket, object, []byte(content), opts)
	return
}

func (s *ObjectResourceCrud) Delete() (e error) {
	namespace := s.D.Get("namespace").(string)
	bucket := s.D.Get("bucket").(string)
	object := s.D.Get("object").(string)
	opts := &baremetal.DeleteObjectOptions{}

	_, e = s.Client.DeleteObject(baremetal.Namespace(namespace), bucket, object, opts)
	return
}

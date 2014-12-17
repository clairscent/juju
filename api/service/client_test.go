// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api/service"
	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	jujutesting "github.com/juju/juju/juju/testing"
)

type serviceSuite struct {
	jujutesting.JujuConnSuite

	client *service.Client
}

var _ = gc.Suite(&serviceSuite{})

func (s *serviceSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
	s.client = service.NewClient(s.APIState)
	c.Assert(s.client, gc.NotNil)
}

func (s *serviceSuite) TestSetServiceMetricCredentials(c *gc.C) {
	var called bool
	service.PatchFacadeCall(s, s.client, func(request string, a, response interface{}) error {
		called = true
		c.Assert(request, gc.Equals, "SetMetricCredentials")
		args, ok := a.(params.ServiceMetricCredentials)
		c.Assert(ok, jc.IsTrue)
		c.Assert(args.Creds, gc.HasLen, 1)
		c.Assert(args.Creds[0].ServiceName, gc.Equals, "serviceA")
		c.Assert(args.Creds[0].MetricCredentials, gc.DeepEquals, []byte("creds 1"))

		result := response.(*params.ErrorResults)
		result.Results = make([]params.ErrorResult, 1)
		return nil
	})
	err := s.client.SetMetricCredentials("serviceA", []byte("creds 1"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(called, jc.IsTrue)
}

func (s *serviceSuite) TestSetServiceMetricCredentialsFails(c *gc.C) {
	var called bool
	service.PatchFacadeCall(s, s.client, func(request string, args, response interface{}) error {
		called = true
		c.Assert(request, gc.Equals, "SetMetricCredentials")
		result := response.(*params.ErrorResults)
		result.Results = make([]params.ErrorResult, 1)
		result.Results[0].Error = common.ServerError(common.ErrPerm)
		return result.OneError()
	})
	err := s.client.SetMetricCredentials("service", []byte("creds"))
	c.Assert(err, gc.ErrorMatches, "permission denied")
	c.Assert(called, jc.IsTrue)
}

func (s *serviceSuite) TestSetServiceMetricCredentialsNoMocks(c *gc.C) {
	service := s.Factory.MakeService(c, nil)
	err := s.client.SetMetricCredentials(service.Name(), []byte("creds"))
	c.Assert(err, jc.ErrorIsNil)
	err = service.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(service.MetricCredentials(), gc.DeepEquals, []byte("creds"))
}

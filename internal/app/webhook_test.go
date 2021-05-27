package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceImageRegistryHost(t *testing.T) {
	type testCase struct {
		name            string
		inputImage      string
		expectedImage   string
		replacementHost string
	}
	awsRegistryHost := "xxx.dkr.ecr.eu-west-1.amazonaws.com"
	pathBasedRegistryHost := "foo.com/images"
	testCases := []testCase{
		{name: "replacing image that has no registry host with standard registry host", replacementHost: awsRegistryHost,
			inputImage: "centos:8", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/centos:8"},
		{name: "replacing image that has no registry host with path based registry host", replacementHost: pathBasedRegistryHost,
			inputImage: "centos:8", expectedImage: "foo.com/images/centos:8"},
		{name: "replacing image that has the same standard registry host", replacementHost: awsRegistryHost,
			inputImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/centos:8", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/centos:8"},
		{name: "replacing image that has the same path based registry host", replacementHost: pathBasedRegistryHost,
			inputImage: "foo.com/images/centos:8", expectedImage: "foo.com/images/centos:8"},
		{name: "replacing image that has different registry host with standard registry host", replacementHost: awsRegistryHost,
			inputImage: "cr.l5d.io/linkerd/controller:stable-2.9.5", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/linkerd/controller:stable-2.9.5"},
		{name: "replacing image that has different registry host with path based registry host", replacementHost: pathBasedRegistryHost,
			inputImage: "cr.l5d.io/linkerd/controller:stable-2.9.5", expectedImage: "foo.com/images/linkerd/controller:stable-2.9.5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualImage := replaceImageRegistryHost(tc.replacementHost, tc.inputImage)
			assert.Equal(t, tc.expectedImage, actualImage)
		})
	}
}

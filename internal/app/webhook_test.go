package app

import (
	"fmt"
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
		{name: "replacing image that has different registry host with standard registry with double digit number in version", replacementHost: awsRegistryHost,
			inputImage: "cr.l5d.io/linkerd/controller:stable-2.10.2", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/linkerd/controller:stable-2.10.2"},
		{name: "replacing image that has no registry host and a sub path with standard registry with double digit number in version", replacementHost: awsRegistryHost,
			inputImage: "kiwigrid/k8s-sidecar:0.1.151", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/kiwigrid/k8s-sidecar:0.1.151"},
		{name: "replacing image that has no registry host and a sub path with standard registry host", replacementHost: awsRegistryHost,
			inputImage: "kiwigrid/k8s-sidecar:0.1.151", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/kiwigrid/k8s-sidecar:0.1.151"},
		{name: "replacing image that has no registry host and a sub path with path based registry host", replacementHost: pathBasedRegistryHost,
			inputImage: "kiwigrid/k8s-sidecar:0.1.151", expectedImage: "foo.com/images/kiwigrid/k8s-sidecar:0.1.151"},
		{name: "replacing image that has a registry host in a different AWS region", replacementHost: awsRegistryHost,
			inputImage: "yyy.dkr.ecr.eu-west-2.amazonaws.com/openzipkin/zipkin-aws:0.21.2", expectedImage: fmt.Sprint(awsRegistryHost, "/openzipkin/zipkin-aws:0.21.2")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualImage := replaceImageRegistryHost(tc.replacementHost, tc.inputImage)
			assert.Equal(t, tc.expectedImage, actualImage)
		})
	}
}

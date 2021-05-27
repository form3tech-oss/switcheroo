package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceImageRegistryHost(t *testing.T) {

	type testCase struct {
		inputImage    string
		expectedImage string
	}

	newRegistryHost := "xxx.dkr.ecr.eu-west-1.amazonaws.com"

	testCases := []testCase{
		{inputImage: "centos:8", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/centos:8"},
		{inputImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/centos:8", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/centos:8"},
		{inputImage: "cr.l5d.io/linkerd/controller:stable-2.9.5", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/linkerd/controller:stable-2.9.5"},
		{inputImage: "docker.io/amazon/aws-alb-ingress-controller:v1.1.5", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/amazon/aws-alb-ingress-controller:v1.1.5"},
		{inputImage: "quay.io/jetstack/cert-manager-controller:v1.1.1", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/jetstack/cert-manager-controller:v1.1.1"},
		{inputImage: "k8s.gcr.io/metrics-server-amd64:v0.3.6", expectedImage: "xxx.dkr.ecr.eu-west-1.amazonaws.com/metrics-server-amd64:v0.3.6"},
	}

	for _, tc := range testCases {
		actualImage := replaceImageRegistryHost(newRegistryHost, tc.inputImage)
		assert.Equal(t, tc.expectedImage, actualImage)
	}

}

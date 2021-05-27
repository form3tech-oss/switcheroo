package test

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"testing"
)

var ServerPort int

type mutationStage struct {
	t         *testing.T
	k8sClient client.Client
	pod       *v1.Pod
}

func NewMutationScenario(t *testing.T, client *client.Client) (*mutationStage, *mutationStage, *mutationStage) {
	stage := &mutationStage{
		t:         t,
		k8sClient: *client,
	}
	return stage, stage, stage
}

func (s *mutationStage) preparedPod(podProvider func() *v1.Pod) *mutationStage {
	s.pod = podProvider()
	return s
}

func (s *mutationStage) prepareUpdatedPod(podUpdater func(pod *v1.Pod) *v1.Pod) *mutationStage {
	s.pod = podUpdater(s.pod)
	return s
}

func (s *mutationStage) podCreated() *mutationStage {
	err := s.k8sClient.Create(context.TODO(), s.pod)
	if err != nil {
		s.t.Fatalf("Error creating pod: %v", err)
	}
	return s
}

func (s *mutationStage) podUpdated() *mutationStage {
	err := s.k8sClient.Update(context.TODO(), s.pod)
	if err != nil {
		s.t.Fatalf("Error creating pod: %v", err)
	}
	return s
}

func (s *mutationStage) and() *mutationStage {
	return s
}

func (s *mutationStage) podImagesHasRequiredRegistry() *mutationStage {
	pod := v1.Pod{}
	err := s.k8sClient.Get(context.TODO(), client.ObjectKey{Name: s.pod.Name, Namespace: s.pod.Namespace}, &pod)
	if err != nil {
		s.t.Fatalf("Error getting pod: %v", err)
	}
	for i, container := range pod.Spec.Containers {
		if !strings.HasPrefix(container.Image, NewRegistryHost) {
			s.t.Fatalf("Container number %d does not start with required prefix: %s", i, container.Image)
		}
	}
	for i, container := range pod.Spec.InitContainers {
		if !strings.HasPrefix(container.Image, NewRegistryHost) {
			s.t.Fatalf("Container number %d does not start with required prefix: %s", i, container.Image)
		}
	}
	return s
}

func updateImageOfPod(pod *v1.Pod) *v1.Pod {
	pod.Spec.Containers[0].Image = "quay.io/jetstack/cert-manager-controller:v2.0.0"
	return pod
}

func testPodWithName(name string) func() *v1.Pod {
	var f func() *v1.Pod
	f = func() *v1.Pod {
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      name,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Image: "quay.io/jetstack/cert-manager-controller:v1.1.1",
						Name:  "cert-manager",
					},
					{
						Image: "centos:8",
						Name:  "centos",
					},
				},
				InitContainers: []v1.Container{
					{
						Image: "docker.io/amazon/aws-alb-ingress-controller:v1.1.5",
						Name:  "aws-alb-ingress",
					},
					{
						Image: "cr.l5d.io/linkerd/controller:stable-2.9.5",
						Name:  "linkerd",
					},
				},
			},
		}
		return &pod
	}
	return f
}

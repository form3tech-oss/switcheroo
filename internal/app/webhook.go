package app

import (
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

type podMutator struct {
	newRegistryHost string
	Client          client.Client
	decoder         *admission.Decoder
}

// podAnnotator implements admission.Handler.
var _ admission.Handler = (*podMutator)(nil)

// podAnnotator implements inject.Decoder.
var _ admission.DecoderInjector = &podMutator{}

func (p *podMutator) Handle(ctx context.Context, request admission.Request) admission.Response {
	pod := &corev1.Pod{}
	err := p.decoder.Decode(request, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	for i := range pod.Spec.InitContainers {
		replacedImage := replaceImageRegistryHost(p.newRegistryHost, pod.Spec.InitContainers[i].Image)
		pod.Spec.InitContainers[i].Image = replacedImage
	}

	for i := range pod.Spec.Containers {
		replacedImage := replaceImageRegistryHost(p.newRegistryHost, pod.Spec.Containers[i].Image)
		pod.Spec.Containers[i].Image = replacedImage
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(request.Object.Raw, marshaledPod)
}

func (p *podMutator) InjectDecoder(d *admission.Decoder) error {
	p.decoder = d
	return nil
}

func replaceImageRegistryHost(newRegistryHost string, image string) string {
	if strings.HasPrefix(image, newRegistryHost) {
		return image
	}
	slashIndex := strings.Index(image, "/")
	if slashIndex > -1 {
		return spew.Sprintf("%s/%s", newRegistryHost, image[slashIndex+1:])
	} else {
		return spew.Sprintf("%s/%s", newRegistryHost, image)
	}
}

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
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

	if slashIndex == -1 {
		return fmt.Sprintf("%s/%s", newRegistryHost, image)
	}

	domain := image[0:slashIndex]

	if isDomain(domain) {
		return fmt.Sprintf("%s/%s", newRegistryHost, image[slashIndex+1:])
	} else {
		return fmt.Sprintf("%s/%s", newRegistryHost, image)
	}
}

func isDomain(value string) bool {
	// https://www.oreilly.com/library/view/regular-expressions-cookbook/9781449327453/ch08s15.html
	domainMatchingPattern := `^((?=[a-z0-9-]{1,63}\.)(xn--)?[a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,63}$`
	// regexp2 allows for a positive lookahead so we can make sure each part of the domain is 1 to 63 characters as per spec
	expression := regexp2.MustCompile(domainMatchingPattern, regexp2.None)
	isMatch, _ := expression.MatchString(value)
	return isMatch
}

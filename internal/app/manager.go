package app

import (
	"fmt"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const MutatingWebhookPath string = "/mutate-v1-pod"

func NewRegistryHostMutatorManager(config *rest.Config, options manager.Options, newRegistryHost string) (manager.Manager, error) {
	entryLog := options.Logger
	entryLog.Info("setting up manager")

	mgr, err := manager.New(config, options)
	if err != nil {
		return nil, fmt.Errorf("unable to set up overall controller manager:%s", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("unable to set up health check: %s", err)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("unable to set up ready check: %s", err)
	}

	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()
	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register(MutatingWebhookPath, &webhook.Admission{Handler: &podMutator{Client: mgr.GetClient(), newRegistryHost: newRegistryHost}})

	return mgr, nil
}

package main

import (
	"errors"
	"github.com/form3tech/switcheroo/internal/app"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"strconv"
)

func init() {
	log.SetLogger(zap.New())
}

const DefaultPort = 9543

func main() {
	entryLog := log.Log.WithName("entrypoint")
	newRegistryHost, ok := os.LookupEnv("NEW_REGISTRY_HOST")
	if !ok {
		entryLog.Error(errors.New("NEW_REGISTRY_HOST is not set"), "")
		os.Exit(1)
	}
	certDirectory, ok := os.LookupEnv("CERT_DIRECTORY")
	if !ok {
		entryLog.Error(errors.New("CERT_DIRECTORY is not set"), "")
		os.Exit(1)
	}
	webhookPortString, ok := os.LookupEnv("WEBHOOK_PORT")
	var appPort int
	if ok {
		parsedAppPort, err := strconv.Atoi(webhookPortString)
		if err != nil {
			entryLog.Error(errors.New("WEBHOOK_PORT has invalid value: "+webhookPortString), "")
			os.Exit(1)
		}
		appPort = parsedAppPort
	} else {
		appPort = DefaultPort
	}

	entryLog.Info("setting up manager")
	manager, err := app.NewRegistryHostMutatorManager(config.GetConfigOrDie(), manager.Options{
		Port:    appPort,
		Logger:  entryLog,
		CertDir: certDirectory,
	}, newRegistryHost)
	if err != nil {
		entryLog.Error(err, "")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := manager.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

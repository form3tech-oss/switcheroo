package test

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/form3tech/switcheroo/internal/app"
	"github.com/go-logr/logr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"log"
	"net"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"time"
)

type testContext struct {
	options     *manager.Options
	client      *client.Client
	config      *rest.Config
	environment *envtest.Environment
	manager     *manager.Manager
	freePort    int
	logger      logr.Logger
}

func newTestContext() *testContext {
	logf.SetLogger(zap.New())
	logger := logf.Log.WithName("testing")
	return &testContext{
		logger: logger,
		options: &manager.Options{
			Logger: logger,
		},
	}
}

const NewRegistryHost = "xxx.dkr.ecr.eu-west-1.amazonaws.com"

func (c *testContext) a_local_instance_of_kubernetes() *testContext {
	os.Setenv("KUBEBUILDER_ASSETS", "./../bin")
	os.Setenv("KUBEBUILDER_ATTACH_CONTROL_PLANE_OUTPUT", "true")
	apiServerFlags := envtest.DefaultKubeAPIServerFlags[0 : len(envtest.DefaultKubeAPIServerFlags)-1]
	apiServerFlags = append(apiServerFlags, "--enable-admission-plugins=MutatingAdmissionWebhook")
	environment := &envtest.Environment{
		//CRDDirectoryPaths:  []string{filepath.Join("..", "config", "crd", "bases")},
		KubeAPIServerFlags:    apiServerFlags,
		WebhookInstallOptions: webhookOptions(),
	}

	config, err := environment.Start()
	failIfError("not able to start local test kubernetes environment", err)

	scheme := runtime.NewScheme()
	err = clientgoscheme.AddToScheme(scheme)
	failIfError("Error adding scheme", err)
	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	failIfError("not able to create k8sClient for local test kubernetes environment", err)

	c.client = &k8sClient
	c.config = config
	c.environment = environment
	return c
}

func (c *testContext) the_webhook_api_ready_to_receive_mutations() *testContext {
	webhookOptions := c.environment.WebhookInstallOptions
	manager, err := app.NewRegistryHostMutatorManager(c.config, manager.Options{
		Logger:         logf.Log.WithName("testlogger"),
		LeaderElection: false,
		Port:           webhookOptions.LocalServingPort,
		Host:           webhookOptions.LocalServingHost,
		CertDir:        webhookOptions.LocalServingCertDir,
	}, NewRegistryHost)

	failIfError("Unable to create the webhook manager", err)
	c.manager = &manager
	go func() {
		err := manager.Start(context.Background())
		failIfError("not able to start web hook API", err)
	}()
	waitUntilWebhookServerIsRunning(c)

	return c
}

func waitUntilWebhookServerIsRunning(c *testContext) {
	d := &net.Dialer{Timeout: time.Second}
	err := retry.Do(
		func() error {
			serverURL := fmt.Sprintf("%s:%d", c.environment.WebhookInstallOptions.LocalServingHost,
				c.environment.WebhookInstallOptions.LocalServingPort)
			conn, err := tls.DialWithDialer(d, "tcp", serverURL, &tls.Config{
				InsecureSkipVerify: true,
			})
			if err != nil {
				return err
			}
			conn.Close()
			return nil
		}, retry.Attempts(10), retry.Delay(time.Second))
	failIfError("not able to start web hook API", err)
}

func webhookOptions() envtest.WebhookInstallOptions {
	servicePath := app.MutatingWebhookPath
	failPolicy := admissionregistrationv1.Fail
	sideEffects := admissionregistrationv1.SideEffectClassNone
	webhookInstallOptions := envtest.WebhookInstallOptions{
		MutatingWebhooks: []client.Object{
			&admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: "switcheroo",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "MutatingWebhookConfiguration",
					APIVersion: "admissionregistration.k8s.io/v1",
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name:                    "image-registry-mutating-hook.form3.tech",
						AdmissionReviewVersions: []string{"v1"},
						FailurePolicy:           &failPolicy,
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Path: &servicePath,
							},
						},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{
									admissionregistrationv1.Create,
									admissionregistrationv1.Update,
								},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{""},
									APIVersions: []string{"v1"},
									Resources:   []string{"pods"},
								},
							},
						},
						SideEffects: &sideEffects,
					},
				},
			},
		},
	}
	return webhookInstallOptions
}

func (c *testContext) and() *testContext {
	return c
}

func (c *testContext) tearDown() {
	c.environment.Stop()
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func failIfError(format string, err error) {
	if err != nil {
		log.Fatalf(format + fmt.Sprintf(": %s", err))
	}
}

package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/agent-sandbox/agent-sandbox/pkg/activator"
	"github.com/agent-sandbox/agent-sandbox/pkg/config"
	"github.com/agent-sandbox/agent-sandbox/pkg/handler"
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"
	"github.com/agent-sandbox/agent-sandbox/pkg/scaler"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	configmapinformer "knative.dev/pkg/configmap/informer"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/version"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

func main() {
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//init config for global
	cfg := config.InitConfig()

	var fs flag.FlagSet
	klog.InitFlags(&fs)
	fs.Set("v", "2")

	klog.Infof("Loaded config %+v", config.Cfg)

	klog.Info("Setup k8s cluster connection and start informers")

	kubecfg := injection.ParseAndGetRESTConfigOrDie()
	klog.Info("Cluster info ", "host=", kubecfg.Host)

	log.Printf("Registering %d clients", len(injection.Default.GetClients()))
	log.Printf("Registering %d informer factories", len(injection.Default.GetInformerFactories()))
	log.Printf("Registering %d informers", len(injection.Default.GetInformers()))
	log.Printf("Registering %d filtered informers", len(injection.Default.GetFilteredInformers()))

	kubecfg.QPS = 2 * rest.DefaultQPS
	kubecfg.Burst = 2 * rest.DefaultBurst
	rootCtx = injection.WithNamespaceScope(rootCtx, config.Cfg.SandboxNamespace)
	rootCtx, informers := injection.Default.SetupInformers(rootCtx, kubecfg)

	kubeClient := kubeclient.Get(rootCtx)

	//load template for sandbox deployment and pool replicaSet
	cfg.KubeClient = kubeClient
	cfg.LoadSandboxRSTemplate()
	cfg.CheckConfigmap()
	// watch configmap for dynamic update
	configMapWatcher := configmapinformer.NewInformedWatcher(kubeClient, cfg.SandboxNamespace)
	configMapWatcher.Watch(config.TemplatesConfigMapName, config.WatchConfigMap())
	if err := configMapWatcher.Start(rootCtx.Done()); err != nil {
		klog.Fatal("Failed to start configuration manager", zap.Error(err))
	}

	// check k8s version is matched
	// We sometimes start up faster than we can reach kube-api. Poll on failure to prevent us terminating
	var err error
	if perr := wait.PollUntilContextTimeout(rootCtx, time.Second, 60*time.Second, true, func(ctx context.Context) (bool, error) {
		if err = version.CheckMinimumVersion(kubeClient.Discovery()); err != nil {
			ctx.Done()
			log.Print("Failed to get k8s version ", err)
		}
		return err == nil, nil
	}); perr != nil {
		log.Fatal("Timed out attempting to get k8s version: ", err)
	}

	if err := controller.StartInformers(rootCtx.Done(), informers...); err != nil {
		log.Fatalln("Failed to start informers", zap.Error(err))
	}
	log.Printf("Starting informers %v", len(informers))

	pl := sandbox.NewPoolManager(rootCtx)
	a := activator.NewActivator(rootCtx)
	c := sandbox.NewController(rootCtx, kubecfg, pl)

	// Start the autoscaler
	go func() {
		s := scaler.NewScaler(rootCtx, a, c)
		klog.Info("Starting timeout and idle timeout  scaler")
		s.RunScaling()
	}()

	go func() {
		// Start the pool syncer
		klog.Info("Starting pool syncer")
		pl.StartPoolSyncing()
	}()

	klog.Info("Starting the api server")
	apiServer := handler.New(rootCtx, a, c)
	if err := apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Print("Failed to run HTTP server", zap.Error(err))
	}

}

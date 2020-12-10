package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/cirocosta/pizza-controller/pkg/reconciler"
)

func init() {
	log.SetLogger(zap.New(zap.UseDevMode(true)))
}

func run() error {
	scheme := runtime.NewScheme()

	if err := reconciler.AddToScheme(scheme); err != nil {
		return fmt.Errorf("add to scheme: %w", err)
	}

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		MetricsBindAddress: "0",
		Scheme:             scheme,
	})
	if err != nil {
		return fmt.Errorf("new manager: %w", err)
	}

	if err := reconciler.RegisterReconcilers(mgr); err != nil {
		return fmt.Errorf("register reconcilers: %w", err)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		return fmt.Errorf("mgr start: %w", err)
	}

	return nil
}

func main() {
	entryLog := log.Log.WithName("entrypoint")
	entryLog.Info("initializing")

	if err := run(); err != nil {
		entryLog.Error(err, "failed to initialize controller")
		return
	}

	entryLog.Info("finished")
}

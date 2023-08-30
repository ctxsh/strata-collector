package main

import (
	"os"

	"ctx.sh/strata-collector/pkg/controller"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme = runtime.NewScheme()
)

func main() {
	ctx := ctrl.SetupSignalHandler()

	// TODO: Actually do a better job of configuring the logger.
	ctrl.SetLogger(zap.New())

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Port:   9443,
		// TODO: set up leader election once we have the dev cluster up
		// and the cert generators.
	})

	reconciler := &controller.Reconciler{
		Client: mgr.GetClient(),
		Log:    mgr.GetLogger().WithValues("controller", "strata"),
	}

	err := reconciler.SetupWithManager(mgr)
	if err != nil {
		// log and exit
		os.Exit(1)
	}

	if err != mgr.Start(ctx); err != nil {
		// log and exit
		os.Exit(1)
	}
}

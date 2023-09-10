package main

import (
	"flag"
	"os"

	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	// Temporary logger for initial setup
	setupLog = ctrl.Log.WithName("setup")
	scheme   = runtime.NewScheme()
)

func init() {
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
}

func main() {
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	log := zap.New(zap.UseFlagOptions(&opts))

	ctx := ctrl.SetupSignalHandler()

	// TODO: Actually do a better job of configuring the logger.
	ctrl.SetLogger(log)

	setupLog.Info("initializing manager")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		// TODO: set up leader election once we have the dev cluster up
		// and the cert generators.
	})
	if err != nil {
		log.Error(err, "unable to initialize manager")
		os.Exit(1)
	}

	controller := controller.New(mgr, &controller.ControllerOpts{
		Logger: mgr.GetLogger(),
	})

	err = controller.Setup()
	if err != nil {
		log.Error(err, "unable to setup controller")
		os.Exit(1)
	}

	// Start the manager process
	log.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		log.Error(err, "unable to start manager")
		os.Exit(1)
	}
}

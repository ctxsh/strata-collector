// Copyright 2023 Rob Lyon <rob@ctxswitch.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/controller/collector"
	"ctx.sh/strata-collector/pkg/controller/discovery"
	"ctx.sh/strata-collector/pkg/service"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ControllerOpts struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Controller struct {
	mgr      ctrl.Manager
	logger   logr.Logger
	metrics  *strata.Metrics
	services *service.Manager
}

func New(mgr ctrl.Manager, opts *ControllerOpts) *Controller {
	return &Controller{
		mgr:     mgr,
		logger:  opts.Logger,
		metrics: opts.Metrics,
		services: service.NewManager(mgr, &service.ManagerOpts{
			Logger:  opts.Logger,
			Metrics: opts.Metrics,
		}),
	}
}

func (c *Controller) Setup() error {
	// Set up collector controller.
	collectorController := &collector.Controller{
		Client:   c.mgr.GetClient(),
		Cache:    c.mgr.GetCache(),
		Log:      c.mgr.GetLogger().WithValues("controller", "collector"),
		Services: c.services,
	}

	err := collectorController.SetupWithManager(c.mgr)
	if err != nil {
		return err
	}

	// Set up discovery controller.
	discoveryController := &discovery.Controller{
		Client:   c.mgr.GetClient(),
		Log:      c.mgr.GetLogger().WithValues("controller", "discovery"),
		Services: c.services,
	}

	err = discoveryController.SetupWithManager(c.mgr)
	if err != nil {
		return err
	}

	return nil
}

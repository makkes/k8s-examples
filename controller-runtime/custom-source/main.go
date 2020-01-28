package main

import (
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type Reconciler struct {
	client client.Client
	logger logr.Logger
}

// Reconcile makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/reconcile.Reconciler interface
// that is used by the controller.
func (r Reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	r.logger.Info("reconciliation request", "req", req)
	return reconcile.Result{}, nil
}

// InjectClient makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/runtime/inject.Client interface
// so that a client is injected into the reconciler by controller-runtimee.
func (r *Reconciler) InjectClient(c client.Client) error {
	r.client = c
	return nil
}

func (r *Reconciler) timerEvent(ev event.GenericEvent) bool {
	r.logger.Info("timer event", "ev", ev)
	// returning false here will cause the event to not be passed to our Reconcile method.
	return false
}

func main() {
	controllerruntime.SetLogger(zap.Logger(true))
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		panic(err)
	}

	ch := make(chan event.GenericEvent)
	go func(ch chan event.GenericEvent) {
		ticker := time.NewTicker(3 * time.Second)
		for range ticker.C {
			ch <- event.GenericEvent{}
		}
	}(ch)
	src := source.Channel{
		Source: ch,
	}

	reconciler := &Reconciler{logger: controllerruntime.Log}

	controller, err := builder.ControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Build(reconciler)
	if err != nil {
		panic(err)
	}
	if err := controller.Watch(&src, &handler.EnqueueRequestForObject{}, predicate.Funcs{GenericFunc: reconciler.timerEvent}); err != nil {
		panic(err)
	}

	if err := mgr.Start(nil); err != nil {
		panic(err)
	}
}

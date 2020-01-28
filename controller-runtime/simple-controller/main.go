package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	client client.Client
}

// Reconcile makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/reconcile.Reconciler interface
// that is used by the controller.
func (r Reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	var pod corev1.Pod
	if err := r.client.Get(ctx, req.NamespacedName, &pod); err != nil {
		return reconcile.Result{}, err
	}
	fmt.Printf("reconciled pod %s/%s (%s)\n", pod.GetNamespace(), pod.GetName(), pod.GetCreationTimestamp())
	return reconcile.Result{}, nil
}

// InjectClient makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/runtime/inject.Client interface
// so that a client is injected into the reconciler by controller-runtimee.
func (r *Reconciler) InjectClient(c client.Client) error {
	r.client = c
	return nil
}

func main() {
	controllerruntime.SetLogger(zap.Logger(true))
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		panic(err)
	}
	err = builder.ControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(&Reconciler{})
	if err != nil {
		panic(err)
	}

	if err := mgr.Start(nil); err != nil {
		panic(err)
	}
}

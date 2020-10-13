package main

import (
	"context"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	client client.Client
	logger logr.Logger
}

// Reconcile makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/reconcile.Reconciler interface
// that is used by the controller.
func (r Reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	var cm corev1.ConfigMap
	if err := r.client.Get(ctx, req.NamespacedName, &cm); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if cm.DeletionTimestamp != nil && !cm.DeletionTimestamp.IsZero() {
		r.logger.Info("ConfigMap is being deleted", "name", cm.GetName())
		for idx, finalizer := range cm.Finalizers {
			if finalizer == "makk.es/ctr" {
				cm.Finalizers[len(cm.Finalizers)-1], cm.Finalizers[idx] = cm.Finalizers[idx], cm.Finalizers[len(cm.Finalizers)-1]
				cm.Finalizers = cm.Finalizers[:len(cm.Finalizers)-1]
				if err := r.client.Update(ctx, &cm); err != nil {
					return reconcile.Result{}, err
				}
			}
		}
		return reconcile.Result{}, nil
	}

	for _, finalizer := range cm.Finalizers {
		if finalizer == "makk.es/ctr" {
			return reconcile.Result{}, nil
		}
	}

	cm.Finalizers = append(cm.Finalizers, "makk.es/ctr")

	if err := r.client.Update(ctx, &cm); err != nil {
		return reconcile.Result{}, err
	}

	r.logger.Info("reconciled ConfigMap", "name", cm.GetName(), "creation", cm.GetCreationTimestamp())
	return reconcile.Result{}, nil
}

// InjectClient makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/runtime/inject.Client interface
// so that a client is injected into the reconciler by controller-runtimee.
func (r *Reconciler) InjectClient(c client.Client) error {
	r.client = c
	return nil
}

func (r *Reconciler) InjectLogger(l logr.Logger) error {
	r.logger = l
	return nil
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true), zap.Level(zapcore.InfoLevel)))

	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		panic(err)
	}

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Namespace:               string(namespace),
		MetricsBindAddress:      "0",
		LeaderElection:          false,
		LeaderElectionID:        "pod-manager",
		LeaderElectionNamespace: "default",
	})
	if err != nil {
		panic(err)
	}

	err = builder.ControllerManagedBy(mgr).
		WithOptions(controller.Options{}).
		For(&corev1.ConfigMap{}).
		Complete(&Reconciler{})
	if err != nil {
		panic(err)
	}

	if err := mgr.Start(nil); err != nil {
		panic(err)
	}
}

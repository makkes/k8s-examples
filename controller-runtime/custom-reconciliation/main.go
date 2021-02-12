package main

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
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
	log    logr.Logger
}

// Reconcile makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/reconcile.Reconciler interface
// that is used by the controller.
func (r Reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	r.log.Info("reconciling", "namespaced name", req.NamespacedName)
	ctx := context.Background()
	var pod corev1.Pod
	if err := r.client.Get(ctx, req.NamespacedName, &pod); err != nil {
		return reconcile.Result{}, err
	}
	r.log.Info("reconciled pod", "namespace", pod.GetNamespace(), "name", pod.GetName(), "created", pod.GetCreationTimestamp(), "version", pod.ResourceVersion)
	return reconcile.Result{}, nil
}

// InjectClient makes our reconciler implement the sigs.k8s.io/controller-runtime/pkg/runtime/inject.Client interface
// so that a client is injected into the reconciler by controller-runtime.
func (r *Reconciler) InjectClient(c client.Client) error {
	r.client = c
	return nil
}

func main() {
	logger := zap.Logger(true)
	controllerruntime.SetLogger(logger)
	// syncPeriod := 1 * time.Second
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		// SyncPeriod: &syncPeriod,
	})
	if err != nil {
		panic(err)
	}

	labelSelector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app.kubernetes.io/name":      "kubeaddons",
			"app.kubernetes.io/component": "controller",
		},
	})
	if err != nil {
		panic(err)
	}

	ctrl, err := builder.ControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Build(&Reconciler{
			log: logger,
		})
	if err != nil {
		panic(err)
	}

	ctrl.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(o handler.MapObject) []reconcile.Request {
			c := mgr.GetClient()
			pods := &corev1.PodList{}
			err := c.List(context.Background(), pods)
			if err != nil {
				logger.Error(err, "could not fetch pods list")
				return nil
			}
			requests := make([]reconcile.Request, 0)
			for _, pod := range pods.Items {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: pod.GetNamespace(),
						Name:      pod.GetName(),
					},
				})
			}
			return requests
		}),
	}, predicate.Funcs{
		CreateFunc: func(ev event.CreateEvent) bool {
			return labelSelector.Matches(labels.Set(ev.Meta.GetLabels()))
		},
		UpdateFunc: func(ev event.UpdateEvent) bool {
			return labelSelector.Matches(labels.Set(ev.MetaOld.GetLabels())) || labelSelector.Matches(labels.Set(ev.MetaNew.GetLabels()))
		},
		DeleteFunc: func(ev event.DeleteEvent) bool {
			return labelSelector.Matches(labels.Set(ev.Meta.GetLabels()))
		},
	})

	if err := mgr.Start(nil); err != nil {
		panic(err)
	}
}

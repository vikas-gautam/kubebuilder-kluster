/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demov1alpha1 "github.com/vikas-gautam/kubebuilder-kluster/api/v1alpha1"
	doutils "github.com/vikas-gautam/kubebuilder-kluster/doUtils"
)

// KlusterReconciler reconciles a Kluster object
type KlusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demo.golearning.dev,resources=klusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.golearning.dev,resources=klusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.golearning.dev,resources=klusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KlusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Enter Reconcile", "req", req)

	// TODO(user): your logic here
	//demo.golearning.dev - it picks only first octet

	kluster := &demov1alpha1.Kluster{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, kluster)
	l.Info("Enter Reconcile", "klusterSpec", kluster.Spec, "klusterStatus", kluster.Status)

	//check if object/resource has finalizer (deletetimestamp), after removal only u could delete object

	if err := doutils.HandleKlusterFinalizer(kluster, r.Client); err != nil {
		return ctrl.Result{}, err
	}

	//add finalizer to new CR, it will update resource and that will be a new trigger to reconciler
	if err := doutils.AddKlusterFinalizer(kluster, r.Client); err != nil {
		return ctrl.Result{}, err
	}

	//Beofore DO cluster creation, we will cross check the cluster.spec with DO APIs

	//create DO cluster
	doClusterID, err := doutils.Createk8sCluster(r.Client, kluster)
	if err != nil {
		l.Error(err, "Error while creating DO cluster")
		return ctrl.Result{}, err
	}
	l.Info("DO cluster id", "id", doClusterID)

	// //deleting k8s cluster from DO
	// err := doutils.Deletek8sCluster(c.k8sclient, "default/dosecret", name)
	// if err != nil {
	// 	l.Error(err, "Error while deleting DO cluster")
	// 	return ctrl.Result{}, err
	// }
	// l.Info("deleted cluster from DO: ", name)

	l.Info("Will reconcile redis operator in again 10 seconds")
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	//return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KlusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1alpha1.Kluster{}).
		Complete(r)
}

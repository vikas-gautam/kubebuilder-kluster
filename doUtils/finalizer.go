package doutils

import (
	"context"
	"log"

	demov1alpha1 "github.com/vikas-gautam/kubebuilder-kluster/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	klusterFinalizer string = "klusterFinalizer"
)

// // finalizeLogger will generate logging interface
// func finalizerLogger(namespace string, name string) logr.Logger {
// 	reqLogger := log.WithValues("Request.Service.Namespace", namespace, "Request.Finalizer.Name", name)
// 	return reqLogger
// }

// HandleRedisFinalizer finalize resource if instance is marked to be deleted
func HandleKlusterFinalizer(cr *demov1alpha1.Kluster, cl client.Client) error {
	if cr.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(cr, klusterFinalizer) {
			controllerutil.RemoveFinalizer(cr, klusterFinalizer)
			if err := cl.Update(context.TODO(), cr); err != nil {
				log.Println("Could not remove finalizer")
				return err
			}
			// if err := Deletek8sCluster(kubernetes.Interface, cr.Spec, cr.Spec.Name); err != nil {
			// 	logger.Error(err, "Could not delete cluster: "+cr.Spec.Name)
			// 	return err
			// }
		}
	}
	return nil
}

// AddRedisFinalizer add finalizer for graceful deletion
func AddKlusterFinalizer(cr *demov1alpha1.Kluster, cl client.Client) error {
	if !controllerutil.ContainsFinalizer(cr, klusterFinalizer) {
		controllerutil.AddFinalizer(cr, klusterFinalizer)
		return cl.Update(context.TODO(), cr)
	}
	return nil
}

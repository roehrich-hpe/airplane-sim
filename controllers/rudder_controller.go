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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	playv1alpha1 "github.com/roehrich-hpe/airplane-sim/api/v1alpha1"
)

// RudderReconciler reconciles a Rudder object
type RudderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=play.github.com,resources=rudders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=play.github.com,resources=rudders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=play.github.com,resources=rudders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Rudder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *RudderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Try this with and without the WithName().  See how this affects
	// readability of the logs when we add more controllers.
	log := log.FromContext(ctx).WithName("rudder")

	rudder := &playv1alpha1.Rudder{}
	if err := r.Get(ctx, req.NamespacedName, rudder); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if rudder.Status.Position != rudder.Spec.Position {
		log.Info("Resetting position")
		rudder.Status.Position = rudder.Spec.Position
		if err := r.Status().Update(ctx, rudder); err != nil {
			if apierrors.IsConflict(err) {
				// You may decide to not log these.  They can
				// be very common.
				log.Info("Conflict while setting position")
				return ctrl.Result{Requeue: true}, nil
			}
			log.Error(err, "Error while setting position")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RudderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&playv1alpha1.Rudder{}).
		Complete(r)
}

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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	playv1alpha1 "github.com/roehrich-hpe/airplane-sim/api/v1alpha1"
)

// PedalLinkageReconciler reconciles a PedalLinkage object
type PedalLinkageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=play.github.com,resources=pedals,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=play.github.com,resources=pedals/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=play.github.com,resources=pedals/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PedalLinkage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *PedalLinkageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	pedals := &playv1alpha1.Pedals{}
	if err := r.Get(ctx, req.NamespacedName, pedals); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var positionWanted string
	switch pedals.Spec.Pressed {
	case "none":
		if pedals.Status.LinkagePosition != "neutral" {
			positionWanted = "neutral"
		}
	default:
		if pedals.Status.LinkagePosition != pedals.Spec.Pressed {
			positionWanted = pedals.Spec.Pressed
		}
	}

	if len(positionWanted) > 0 {
		log.Info("Resetting position")
		pedals.Status.LinkagePosition = positionWanted
		if err := r.Status().Update(ctx, pedals); err != nil {
			if apierrors.IsConflict(err) {
				log.Info("Conflict while setting position")
				return ctrl.Result{Requeue: true}, nil
			}
			log.Error(err, "Error while setting position")
			return ctrl.Result{}, err
		}
	}

	// Get the rudder, move it if necessary.
	rudder := &playv1alpha1.Rudder{}
	// Rudder and Pedals have the same name.
	rudderKey := req.NamespacedName
	if err := r.Get(ctx, rudderKey, rudder); err != nil {
		log.Error(err, "Did not find rudder")
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	if rudder.Spec.Position != pedals.Status.LinkagePosition {
		rudder.Spec.Position = pedals.Status.LinkagePosition
		if err := r.Update(ctx, rudder); err != nil {
			if apierrors.IsConflict(err) {
				log.Info("Conflict on rudder")
				return ctrl.Result{Requeue: true}, nil
			}
			log.Error(err, "Error on rudder")
			return ctrl.Result{}, err
		}
		log.Info("rudder has been set")
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PedalLinkageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&playv1alpha1.Pedals{}).
		Complete(r)
}

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
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	playv1alpha1 "github.com/roehrich-hpe/airplane-sim/api/v1alpha1"
)

// AirplaneReconciler reconciles a Airplane object
type AirplaneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=play.github.com,resources=airplanes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=play.github.com,resources=airplanes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=play.github.com,resources=airplanes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Airplane object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *AirplaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithName("airplane")
	r.Log = log

	airplane := &playv1alpha1.Airplane{}
	if err := r.Get(ctx, req.NamespacedName, airplane); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Check parts")
	// Check the pedals.
	if requeue, err := r.verifyPedals(ctx, airplane); err != nil {
		return ctrl.Result{}, err
	} else if requeue {
		return ctrl.Result{Requeue: true}, nil
	}

	// Check the rudder.
	if requeue, err := r.verifyRudder(ctx, airplane); err != nil {
		return ctrl.Result{}, err
	} else if requeue {
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// Create the pedals resource if it doesn't aleady exist.  Hook up the pedals
// to the airplane.
func (r *AirplaneReconciler) verifyPedals(ctx context.Context, airplane *playv1alpha1.Airplane) (bool, error) {
	log := r.Log.WithName("pedals")

	pedals := &playv1alpha1.Pedals{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(airplane.Spec.TailNumber),
			Namespace: airplane.GetNamespace(),
		},
	}

	// First check whether it exists. Maybe it was orphaned
	// on an earlier pass.
	if err := r.Get(ctx, client.ObjectKeyFromObject(pedals), pedals); err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "Unable to verify existence of pedals")
			return false, err
		}

		// It doesn't exist, so create it.
		ctrl.SetControllerReference(airplane, pedals, r.Scheme)
		pedals.Spec.Pressed = "none"
		if err := r.Create(ctx, pedals); err != nil {
			log.Error(err, "Unable to create pedal`s")
			return false, err
		}
		log.Info("Created pedals", "pedals", pedals)
	}

	// Hook up the pedals to the airplane, if it isn't already.
	pedalsRef := corev1.ObjectReference{
		Kind:      reflect.TypeOf(playv1alpha1.Pedals{}).Name(),
		Name:      pedals.GetName(),
		Namespace: pedals.GetNamespace(),
	}
	if airplane.Status.Pedals == pedalsRef {
		// All good.
		return false, nil
	}
	airplane.Status.Pedals = pedalsRef
	if err := r.Status().Update(ctx, airplane); err != nil {
		// We created the pedals above, but we weren't able to
		// hook it up to the airplane this time.  We'll go
		// around again, find the pedals, and try to hook
		// it up.
		log.Error(err, "Unable to set pedals reference in airplane")
		return false, err
	}
	log.Info("Hooked up pedals to airplane")

	return true, nil
}

// Create the rudder resource if it doesn't aleady exist.  Hook up the rudder
// to the airplane.
func (r *AirplaneReconciler) verifyRudder(ctx context.Context, airplane *playv1alpha1.Airplane) (bool, error) {
	log := r.Log.WithName("rudder")

	rudder := &playv1alpha1.Rudder{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(airplane.Spec.TailNumber),
			Namespace: airplane.GetNamespace(),
		},
	}

	// First check whether it exists. Maybe it was orphaned
	// on an earlier pass.
	if err := r.Get(ctx, client.ObjectKeyFromObject(rudder), rudder); err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "Unable to verify existence of rudder")
			return false, err
		}

		// It doesn't exist, so create it.
		ctrl.SetControllerReference(airplane, rudder, r.Scheme)
		rudder.Spec.Position = "neutral"
		if err := r.Create(ctx, rudder); err != nil {
			log.Error(err, "Unable to create rudder")
			return false, err
		}
		log.Info("Created rudder", "rudder", rudder)
	}

	// Hook up the rudder to the airplane, if it isn't already.
	rudderRef := corev1.ObjectReference{
		Kind:      reflect.TypeOf(playv1alpha1.Rudder{}).Name(),
		Name:      rudder.GetName(),
		Namespace: rudder.GetNamespace(),
	}
	if airplane.Status.Rudder == rudderRef {
		// All good.
		return false, nil
	}
	airplane.Status.Rudder = rudderRef
	if err := r.Status().Update(ctx, airplane); err != nil {
		// We created the rudder above, but we weren't able to
		// hook it up to the airplane this time.  We'll go
		// around again, find the rudder, and try to hook
		// it up.
		log.Error(err, "Unable to set rudder reference in airplane")
		return false, err
	}
	log.Info("Hooked up rudder to airplane")

	return true, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AirplaneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&playv1alpha1.Airplane{}).
		Owns(&playv1alpha1.Pedals{}).
		Owns(&playv1alpha1.Rudder{}).
		Complete(r)
}

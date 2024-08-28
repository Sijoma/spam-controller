/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	sijomadevv1alpha1 "github.com/sijoma/spam-controller/api/v1alpha1"
)

// SpamReconciler reconciles a Spam object
type SpamReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sijoma.dev,resources=spams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sijoma.dev,resources=spams/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=sijoma.dev,resources=spams/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *SpamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var spam sijomadevv1alpha1.Spam
	err := r.Client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &spam)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for i := range 1000 {
		cm := core.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "spam" + strconv.Itoa(i),
				Namespace: req.Namespace,
			},
			Data: map[string]string{
				"spammer": "yolo",
			},
		}

		// Trigger:  Waited for 9.998942417s due to client-side throttling, not priority and fairness, request:
		op, err := ctrl.CreateOrUpdate(ctx, r.Client, &cm, func() error {
			cm.Data["spammer"] = "Modify-all-the-time" + time.Now().Format("2006-01-02 15:04:05")
			return nil
		})
		if err != nil {
			return ctrl.Result{}, err
		}
		fmt.Println(time.Now().Format(time.TimeOnly), req.Name, op, cm.Name)
	}

	return ctrl.Result{
		RequeueAfter: time.Second * 1,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sijomadevv1alpha1.Spam{}).
		Complete(r)
}

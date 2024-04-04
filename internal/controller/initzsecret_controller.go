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
	"log"
	"time"

	"github.com/go-logr/logr"

	initializv1alpha1 "github.com/initializ/K8s-Operator/api/v1alpha1"
	"github.com/initializ/K8s-Operator/package/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// InitzSecretReconciler reconciles a InitzSecret object
type InitzSecretReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=alpha.initializ.com,resources=initzsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=alpha.initializ.com,resources=initzsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=alpha.initializ.com,resources=initzsecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the InitzSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *InitzSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("initzsecret", req.NamespacedName)

	var initzSecret initializv1alpha1.InitzSecret
	if err := r.Get(ctx, req.NamespacedName, &initzSecret); err != nil {
		if errors.IsNotFound(err) {
			log.Info("InitzSecret resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get InitzSecret")
		return ctrl.Result{}, err
	}
	log.Info("Reconciling InitzSecret", "Name", initzSecret.Name)

	if err := r.ReconcileInitzSecret(ctx, initzSecret); err != nil {
		log.Error(err, "Failed to reconcile InitzSecret")
		return ctrl.Result{}, err
	}

	// Update status if necessary
	// if err := r.updateStatus(ctx, &initzSecret); err != nil {
	// 	log.Error(err, "Failed to update status")
	// 	return ctrl.Result{}, err
	// }
	requeueTime := time.Duration(initzSecret.Spec.ResyncInterval) * time.Second
	log.Info("Reconciliation completed successfully")
	fmt.Println("Reconciliation completed successfully")
	return ctrl.Result{
		RequeueAfter: requeueTime,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InitzSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&initializv1alpha1.InitzSecret{}).
		Complete(r)
}

// ReconcileInitzSecret reconciles the InitzSecret resource
func (r *InitzSecretReconciler) ReconcileInitzSecret(ctx context.Context, initzSecret initializv1alpha1.InitzSecret) error {
	// Get the service token secret reference details from the InitzSecret

	// Fetch the service token from the Kubernetes secret
	serviceToken := initzSecret.Spec.Authentication.ServiceToken.ServiceTokenSecretReference.Servicetoken

	// Get other details from the InitzSecret
	objectIds := initzSecret.Spec.Authentication.ServiceToken.SecretsScope.SecretVars
	orgID := initzSecret.Spec.Authentication.ServiceToken.SecretsScope.OrganisationID
	envSlug := initzSecret.Spec.Authentication.ServiceToken.SecretsScope.EnvSlug
	// Fetch plaintext secrets via service token
	plainTextSecrets, err := util.GetPlainTextSecretsViaServiceToken(serviceToken, orgID, envSlug, objectIds)
	if err != nil {
		return err
	}

	// Create or update the managed Kubernetes secret
	managedSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      initzSecret.Spec.ManagedSecretReference.SecretName,
			Namespace: initzSecret.Spec.ManagedSecretReference.SecretNamespace,
		},
		Data: map[string][]byte{},
	}

	for _, secret := range plainTextSecrets {
		managedSecret.Data[secret.Key] = []byte(secret.Value)
	}

	err = r.createOrUpdateManagedSecret(ctx, managedSecret)
	if err != nil {
		return err
	}

	return nil
}

// createOrUpdateManagedSecret creates or updates the managed Kubernetes secret
func (r *InitzSecretReconciler) createOrUpdateManagedSecret(ctx context.Context, secret *corev1.Secret) error {
	found := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the secret if it doesn't exist
		log.Printf("Creating a new managed secret, Name: %s, Namespace: %s", secret.Name, secret.Namespace)
		err = r.Create(ctx, secret)
		if err != nil {
			return err
		}
	} else if err == nil {
		// Update the secret if it exists
		log.Printf("Updating an existing managed secret, Name: %s, Namespace: %s", secret.Name, secret.Namespace)
		found.Data = secret.Data
		err = r.Update(ctx, found)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// updateStatus updates the status of the InitzSecret resource
// func (r *InitzSecretReconciler) updateStatus(ctx context.Context, initzSecret *initializv1alpha1.InitzSecret) error {
// 	// Perform any status update logic here
// 	// For example, updating the last reconcile time
// 	initzSecret.Status.LastReconcileTime = metav1.Now()

// 	// Update the status
// 	err := r.Status().Update(ctx, initzSecret)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

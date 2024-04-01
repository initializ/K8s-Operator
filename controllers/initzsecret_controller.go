package controllers

import (
	initializv1alpha1 "Initializ-Operator/api/v1alpha1"
	"Initializ-Operator/package/util"
	"context"
	"fmt"
	"log"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// InitzSecretReconciler reconciles the InitzSecret resource
type InitzSecretReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile handles reconciliation logic for InitzSecret
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
	if err := r.updateStatus(ctx, &initzSecret); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	log.Info("Reconciliation completed successfully")
	return ctrl.Result{}, nil
}

// Implement the SetupWithManager method to setup the controller with a manager
func (r *InitzSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&initializv1alpha1.InitzSecret{}).
		Complete(r)
}

// ReconcileInitzSecret reconciles the InitzSecret resource
func (r *InitzSecretReconciler) ReconcileInitzSecret(ctx context.Context, initzSecret initializv1alpha1.InitzSecret) error {
	// Get the service token secret reference details from the InitzSecret
	serviceTokenSecretName := initzSecret.Spec.Authentication.ServiceToken.ServiceTokenSecretReference.SecretName
	serviceTokenSecretNamespace := initzSecret.Spec.Authentication.ServiceToken.ServiceTokenSecretReference.SecretNamespace

	// Fetch the service token from the Kubernetes secret
	serviceToken, err := r.getServiceTokenFromSecret(ctx, serviceTokenSecretName, serviceTokenSecretNamespace)
	if err != nil {
		return err
	}

	// Get other details from the InitzSecret
	objectIds := initzSecret.Spec.Authentication.ServiceToken.SecretsScope.SecretVars
	orgID := initzSecret.Spec.Authentication.ServiceToken.SecretsScope.Workspace
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

// getServiceTokenFromSecret retrieves the service token from the specified Kubernetes secret
func (r *InitzSecretReconciler) getServiceTokenFromSecret(ctx context.Context, secretName, secretNamespace string) (string, error) {
	secret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: secretNamespace}, secret)
	if err != nil {
		return "", err
	}

	serviceToken, ok := secret.Data["serviceToken"]
	if !ok {
		return "", fmt.Errorf("service token not found in secret")
	}

	return string(serviceToken), nil
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
func (r *InitzSecretReconciler) updateStatus(ctx context.Context, initzSecret *initializv1alpha1.InitzSecret) error {
	// Perform any status update logic here
	// For example, updating the last reconcile time
	initzSecret.Status.LastReconcileTime = metav1.Now()

	// Update the status
	err := r.Status().Update(ctx, initzSecret)
	if err != nil {
		return err
	}

	return nil
}

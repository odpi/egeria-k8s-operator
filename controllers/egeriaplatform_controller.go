/*
Copyright 2021 Contributors to the Egeria project.

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
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"

	//"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// was egeriav1
	egeriav1alpha1 "github.com/odpi/egeria-k8s-operator/api/v1alpha1"
)

// EgeriaPlatformReconciler reconciles a EgeriaPlatform object
type EgeriaPlatformReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=egeria.egeria-project.org,resources=egeriaplatforms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=egeria.egeria-project.org,resources=egeriaplatforms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=egeria.egeria-project.org,resources=egeriaplatforms/finalizers,verbs=update
// -- Additional roles required to manage statefulsets & services
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// We compare the state specified by
// the EgeriaPlatform object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
//
func (r *EgeriaPlatformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Fetch the Egeria instance
	egeria := &egeriav1alpha1.EgeriaPlatform{}
	err := r.Get(ctx, req.NamespacedName, egeria)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.FromContext(ctx).Info("Egeria resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.FromContext(ctx).Error(err, "Failed to get Egeria")
		return ctrl.Result{}, err
	}

	// TODO: check what we do if the crd goes away
	// Check if the StatefulSet already exists, if not create a new one
	found := &appsv1.StatefulSet{}

	err = r.Get(ctx, types.NamespacedName{Name: egeria.Name, Namespace: egeria.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new StatefulSet
		dep := r.statefulsetForEgeriaPlatform(egeria)
		log.FromContext(ctx).Info("Creating a new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// StatefulSet created successfully - return and requeue
		// TODO Tag Statefulset with info about the config we used. See also https://cloud.redhat.com/blog/kubernetes-operators-best-practices
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.FromContext(ctx).Error(err, "Failed to get StatefulSet")
		return ctrl.Result{}, err
	}

	// TODO: Check we are using the same image as before(we only go by name)
	// TODO: Check this statefulset is using the same config document that we specced (including name of secret & date)
	// TODO: CHeck our storage spec is the same as before
	// TODO: Check our security certs (just the name - it's ok if the contents change) are same as before
	// TODO: Ensure the statefulset size is the same as the spec
	size := egeria.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// TODO: Check there is a service definition
	foundsvc := &corev1.Service{}
	// TODO: Compute service name from crd & modify to be unique
	err = r.Get(ctx, types.NamespacedName{Name: "egeria-service", Namespace: egeria.Namespace}, foundsvc)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Service
		svc := r.serviceForEgeriaPlatform(egeria)
		log.FromContext(ctx).Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to create new StatefulSet", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}
		// StatefulSet created successfully - return and requeue
		// TODO Tag Statefulset with info about the config we used. See also https://cloud.redhat.com/blog/kubernetes-operators-best-practices
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.FromContext(ctx).Error(err, "Failed to get StatefulSet")
		return ctrl.Result{}, err
	}
	// Update the Egeria status with the pod names
	// List the pods for this egeria's statefulset

	return ctrl.Result{}, nil
}

// statefulsetForEgeria returns an egeria StatefulSet object
func (r *EgeriaPlatformReconciler) statefulsetForEgeriaPlatform(m *egeriav1alpha1.EgeriaPlatform) *appsv1.StatefulSet {
	ls := labelsForEgeria(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		//TODO: Need to spec out storage
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						//TODO: Use latest image for now, and hardcode
						Image: "odpi/egeria:latest",
						Name:  "odpi",
						//Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						//TODO: Allow port to be overridden
						Ports: []corev1.ContainerPort{{
							ContainerPort: 9443,
							Name:          "egeria",
						}},
					}},
				},
			},
		},
	}
	// Set Egeria instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// servicetForEgeria returns an egeria Service  object
func (r *EgeriaPlatformReconciler) serviceForEgeriaPlatform(m *egeriav1alpha1.EgeriaPlatform) *corev1.Service {
	ls := labelsForEgeria(m.Name)
	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			// TODO: More flexible service types needed in future
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Port:     9443,
					Protocol: corev1.ProtocolTCP,
					Name:     "egeria-port",
				},
			},
			Selector: getServiceSelectorLabels("egeriaplatform", "egeria-service"),
		},
	}
	// Set Egeria instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForEgeria returns the labels for selecting the resources
// belonging to the given egeria CR name.
func labelsForEgeria(name string) map[string]string {
	return map[string]string{"app": "egeriaplatform", "egeria_cr": name}
}

// Service selector labels
func getServiceSelectorLabels(component, name string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      name,
		"app.kubernetes.io/component": component,
		"app.kubernetes.io/part-of":   "egeria",
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *EgeriaPlatformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&egeriav1alpha1.EgeriaPlatform{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}

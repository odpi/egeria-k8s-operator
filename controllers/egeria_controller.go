/*


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

	"github.com/go-logr/logr"
	egeriav1 "github.com/odpi/egeria-k8s-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EgeriaReconciler reconciles a Egeria object
type EgeriaReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=egeria.odpi.org,resources=egeria,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=egeria.odpi.org,resources=egeria/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=egeria.odpi.org,resources=egeria/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

func (r *EgeriaReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("egeria", req.NamespacedName)

	// Fetch the Egeria instance
	egeria := &egeriav1.Egeria{}
	err := r.Get(ctx, req.NamespacedName, egeria)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Egeria resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Egeria")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: egeria.Name, Namespace: egeria.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForEgeria(egeria)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment is as expected - compare spec with observed
	for i := range egeria.Spec.PlatformNames
	{
		// Is platform found - if so, great!

		// Is platform in the list, but not found - we need to make it happen
	}

	// Are there any platforms that we're not expecting - we need to remove
	for j := range found.Spec.
	size := egeria.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the Egeria status with the pod names
	// List the pods for this egeria's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(egeria.Namespace),
		client.MatchingLabels(labelsForEgeria(egeria.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Egeria.Namespace", egeria.Namespace, "Egeria.Name", egeria.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, egeria.Status.Nodes) {
		egeria.Status.Nodes = podNames
		err := r.Status().Update(ctx, egeria)
		if err != nil {
			log.Error(err, "Failed to update Egeria status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// deploymentForEgeria returns an egeria Deployment object
func (r *EgeriaReconciler) deploymentForEgeria(m *egeriav1.Egeria) *appsv1.Deployment {
	ls := labelsForEgeria(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
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
						Image: "odpi/egeria:2.3-SNAPSHOT",
						Name:  "odpi",
						//Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
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

// labelsForEgeria returns the labels for selecting the resources
// belonging to the given egeria CR name.
func labelsForEgeria(name string) map[string]string {
	return map[string]string{"app": "egeria", "egeria_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func (r *EgeriaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&egeriav1.Egeria{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

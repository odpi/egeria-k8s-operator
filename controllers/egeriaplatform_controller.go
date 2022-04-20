/*
Copyright 2022 Contributors to the Egeria project.

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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	egeriav1alpha1 "github.com/odpi/egeria-k8s-operator/api/v1alpha1"
)

// EgeriaPlatformReconciler reconciles a EgeriaPlatform object
type EgeriaPlatformReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=egeria.egeria-project.org,resources=egeriaplatforms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=egeria.egeria-project.org,resources=egeriaplatforms/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=egeria.egeria-project.org,resources=egeriaplatforms/finalizers,verbs=update
// -- Additional roles required to manage deployments & services
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// We compare the state specified by
// the EgeriaPlatform object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// ctrl.Result - see https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
//
func (reconciler *EgeriaPlatformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Common returns
	var err error
	var requeue bool
	var egeria *egeriav1alpha1.EgeriaPlatform

	log.FromContext(ctx).Info("Reconciler called.")
	// Fetch the Egeria instance - ie the custom resource definition object.
	egeria, err = reconciler.getEgeriaPlatform(ctx, req)
	if egeria == nil {
		// There is no CRD, so as long as we've set the correct references, Garbage collection
		// will delete everything else we've created
		log.FromContext(ctx).Info("Egeria custom resource not found.", "err", err)
		return ctrl.Result{}, err
	}

	// See if we have the autostart config figured out
	log.FromContext(ctx).Info("Checking autostart configmap", "err", err)
	cm, err, requeue := reconciler.ensureAutostartConfigMap(ctx, egeria)
	if (err != nil) || (cm == nil) {
		log.FromContext(ctx).Info("Problems setting up autostart configmap.", "err", err)
		return ctrl.Result{}, err
	}
	if requeue {
		// We need to wait a little longer as this is creating a deployment in the background
		log.FromContext(ctx).Info("Requeueing after creating autostart configmap")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// TODO: Add finalizer for additional cleanup when CR deleted

	// Check if the Deployment already exists, if not create a new one
	log.FromContext(ctx).Info("Checking deployment", "err", err)
	deployment, err, requeue := reconciler.ensureDeployment(ctx, egeria)
	if (err != nil) || (deployment == nil) {
		log.FromContext(ctx).Info("Problems checking Deployment.", "err", err)
		return ctrl.Result{}, err
	}
	if requeue {
		// We need to wait a little longer as this is creating a deployment in the background
		log.FromContext(ctx).Info("Requeueing after creating deployment")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Update deployment name if needed
	log.FromContext(ctx).Info("Updating deployment name in cr", "err", err)
	err = reconciler.updateDeploymentName(ctx, egeria)
	if err != nil {
		log.FromContext(ctx).Info("Problems updating deployment name", "err", err)
		return ctrl.Result{}, err
	}

	// Update the status with the pod names
	log.FromContext(ctx).Info("Updating pod names in cr", "err", err)
	err = reconciler.updatePodNames(ctx, egeria)
	if err != nil {
		log.FromContext(ctx).Info("Problems updating pod names", "err", err)
		return ctrl.Result{}, err
	}

	// TODO: Check we are using the same image as before(we only go by name)

	// TODO: Extract version from container image metadata?

	// TODO: Check this deployment is using the same config document that we specced (including name of secret & date)
	// TODO: Check for added/removed/changed servers. Update configuration/startup server and restart as needed
	// TODO: Rolling restart

	// TODO: Check our security certs (just the name - it's ok if the contents change) are same as before

	// Ensure the deployment size is the same as the spec
	log.FromContext(ctx).Info("Checking Replica count")
	err, requeue = reconciler.checkReplicas(ctx, egeria, deployment)
	if err != nil {
		log.FromContext(ctx).Info("Problems checking replica count", "err", err)
		return ctrl.Result{}, err
	}
	// TODO: Simplify signature/handling error & requeue
	if requeue {
		log.FromContext(ctx).Info("Requeueing after updating replica count")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Check there is a service definition
	log.FromContext(ctx).Info("Checking service")
	err, requeue = reconciler.ensureService(ctx, egeria)
	if err != nil {
		log.FromContext(ctx).Info("Problems checking service", "err", err)
		return ctrl.Result{}, err
	}
	if requeue {
		// Allow time for the new service definition to be effective before rechecking
		log.FromContext(ctx).Info("Requeueing after creating service")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Update service name if needed
	log.FromContext(ctx).Info("Checking service name")
	err, requeue = reconciler.updateServiceName(ctx, egeria)
	if err != nil {
		log.FromContext(ctx).Info("Problems checking service name", "err", err)
		return ctrl.Result{}, err
	}

	//
	// TODO: Add status conditions https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/

	// TODO Check the services are running on the platforms
	// TODO Admission webhook may be needed to improve validation
	// TODO: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#finalizers
	// We are all done. Fully reconciled!

	log.FromContext(ctx).Info("Fully Reconciled!")

	return ctrl.Result{}, nil
}

//
// Retrieve the Custom Resource for Egeria which we are doing the reconciliation on
//
func (reconciler *EgeriaPlatformReconciler) getEgeriaPlatform(ctx context.Context, req ctrl.Request) (*egeriav1alpha1.EgeriaPlatform, error) {

	// TODO: Handle case where CR is in the process of being deleted
	egeria := &egeriav1alpha1.EgeriaPlatform{}
	err := reconciler.Get(ctx, req.NamespacedName, egeria)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			// TODO: Add framework for logging
			log.FromContext(ctx).Info("EgeriaPlatform custom resource + req.NamespacesName not found. Ignoring since object must be deleted.")
			return nil, nil
		}
		// Error reading the object - requeue the request.
		log.FromContext(ctx).Error(err, "Failed to retrieve EgeriaPlatform custom resource.")
		return nil, err
	}
	return egeria, nil
}

//
// Ensure that a deployment exists for this Egeria instance (ie a platform)
//
func (reconciler *EgeriaPlatformReconciler) ensureDeployment(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform) (*appsv1.Deployment, error, bool) {
	deployment := &appsv1.Deployment{}
	// TODO: Make object name generation configurable
	err := reconciler.Get(ctx, types.NamespacedName{Name: egeria.Name + "-deployment", Namespace: egeria.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Deployment
		dep := reconciler.deploymentForEgeriaPlatform(egeria)
		log.FromContext(ctx).Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = reconciler.Create(ctx, dep)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return deployment, err, false
		}
		// Deployment created successfully - return and requeue
		// TODO Tag deployment with info about the config we used. See also https://cloud.redhat.com/blog/kubernetes-operators-best-practices
		return deployment, nil, true
	} else if err != nil {
		log.FromContext(ctx).Error(err, "Failed to get Deployment")
		return deployment, nil, false
	}

	// TODO: Ensure health check is appropriately setup so that service is only routed when working
	return deployment, err, false
}

func (reconciler *EgeriaPlatformReconciler) ensureAutostartConfigMap(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform) (*corev1.ConfigMap, error, bool) {
	configmap := &corev1.ConfigMap{}
	// TODO: Make object name generation configurable
	err := reconciler.Get(ctx, types.NamespacedName{Name: egeria.Name + "-autostart", Namespace: egeria.Namespace}, configmap)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Deployment
		cm := reconciler.configmapForEgeriaPlatform(ctx, egeria)
		log.FromContext(ctx).Info("Creating a new autostart Configmap", "ConfigMap.Namespace", cm.Namespace, "Deployment.Name", cm.Name)
		err = reconciler.Create(ctx, cm)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "Deployment.Name", cm.Name)
			return cm, err, false
		}
		// Deployment created successfully - return and requeue
		// TODO Tag deployment with info about the config we used. See also https://cloud.redhat.com/blog/kubernetes-operators-best-practices
		return cm, nil, true
	} else if err != nil {
		log.FromContext(ctx).Error(err, "Failed to get Deployment")
		return nil, err, false
	}
	return configmap, nil, false
}

//
// Updates the CR with a summary of the deployment for visibility in the status report
//
func (reconciler *EgeriaPlatformReconciler) updateDeploymentName(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform) error {
	if !reflect.DeepEqual(egeria.Name+"-deployment", egeria.Status.ManagedDeployment) {
		egeria.Status.ManagedDeployment = egeria.Name + "-deployment"
		err := reconciler.Status().Update(ctx, egeria)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to update status")
			return err
		}
	}
	return nil
}

func (reconciler *EgeriaPlatformReconciler) updatePodNames(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform) error {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(egeria.Namespace),
		client.MatchingLabels(egeriaLabels(egeria.Name, "deployment")),
	}
	if err := reconciler.List(ctx, podList, listOpts...); err != nil {
		//TODO: Logs should be consistent
		log.FromContext(ctx).Error(err, "Failed to list pods", "Namespace", egeria.Namespace, "Name", egeria.Name)
		return err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, egeria.Status.Pods) {
		egeria.Status.Pods = podNames
		err2 := reconciler.Status().Update(ctx, egeria)
		if err2 != nil {
			log.FromContext(ctx).Error(err2, "Failed to update Egeria status")
			return err2
		}
	}
	return nil
}

func (reconciler *EgeriaPlatformReconciler) updateServiceName(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform) (error, bool) {
	if !reflect.DeepEqual(egeria.Name+"-service", egeria.Status.ManagedService) {
		egeria.Status.ManagedService = egeria.Name + "-service"
		err := reconciler.Status().Update(ctx, egeria)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to update status")
			return err, false
		}
	}
	return nil, false
}

func (reconciler *EgeriaPlatformReconciler) ensureService(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform) (error, bool) {
	service := &corev1.Service{}
	// TODO: Compute service name from crd & modify to be unique
	err := reconciler.Get(ctx, types.NamespacedName{Name: egeria.Name + "-service", Namespace: egeria.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		// We can't find the service, so will Define a new Service
		svc := reconciler.serviceForEgeriaPlatform(egeria)
		log.FromContext(ctx).Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = reconciler.Create(ctx, svc)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			// Theres a problem. do not requeue & return error
			return err, false
		}
		// Deployment created successfully - return and requeue
		// TODO Tag Deployment with info about the config we used. See also https://cloud.redhat.com/blog/kubernetes-operators-best-practices
		return err, true
	} else if err != nil {
		// Some other error occurred - return & do not requeue for further processing
		log.FromContext(ctx).Error(err, "Failed to get Deployment")
		return err, false
	}
	// We found the deployment ok - continue (no requeue)
	return err, false
}

func (reconciler *EgeriaPlatformReconciler) checkReplicas(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform, deployment *appsv1.Deployment) (error, bool) {
	size := egeria.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		err := reconciler.Update(ctx, deployment)
		if err != nil {
			log.FromContext(ctx).Error(err, "Failed to update Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return err, false
		}
		// Spec updated - return and requeue
		return nil, true
	}
	return nil, false
}

// TODO: Migrate to stateful set if identity needed
// deploymentForEgeria returns an egeria Deployment object
func (reconciler *EgeriaPlatformReconciler) deploymentForEgeriaPlatform(egeriaInstance *egeriav1alpha1.EgeriaPlatform) *appsv1.Deployment {
	labels := egeriaLabels(egeriaInstance.Name, "deployment")
	replicas := egeriaInstance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      egeriaInstance.Name + "-deployment",
			Namespace: egeriaInstance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					// The server configuration is stored in a configmap, which is mapped to a volume
					Volumes: getVolumes(egeriaInstance.Spec.ServerConfig),
					Containers: []corev1.Container{{
						Name:  "platform",
						Image: egeriaInstance.Spec.Image,
						//TODO: Allow port to be overridden
						Ports: []corev1.ContainerPort{{
							ContainerPort: 9443,
							Name:          "platformport",
						}},
						// Additional environment, including the autostart configuration for servers found in configmap
						EnvFrom: []corev1.EnvFromSource{{
							ConfigMapRef: &corev1.ConfigMapEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: egeriaInstance.Name + "-autostart",
								},
							},
						}},
						// Mountpoints are needed for egeria configuration
						//TODO: VolumeMounts: BmeMounts(),
						// This probe defines when to RESTART the container
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Port:   intstr.FromInt(9443),
									Scheme: "HTTPS",
									// TODO Replace hardcoded reference to garygeeke (relevant with platform security)
									Path: "/open-metadata/platform-services/users/garygeeke/server-platform/origin",
								},
							},
							InitialDelaySeconds: 30,
							PeriodSeconds:       30,
							FailureThreshold:    10,
						},
						// This probe defines if the pod can accept work via a service
						// Can help with long-running queries, ensuring routing to another replica
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Port:   intstr.FromInt(9443),
									Scheme: "HTTPS",
									// TODO Replace hardcoded reference to garygeeke (relevant with platform security)
									Path: "/open-metadata/platform-services/users/garygeeke/server-platform/origin",
								},
							},
							InitialDelaySeconds: 30,
							PeriodSeconds:       5,
							FailureThreshold:    3,
						},
						// This probe allows for some settling time at startup & overrides the other probes
						// TODO - Currently each probe is the same, ultimately should be unique
						StartupProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Port:   intstr.FromInt(9443),
									Scheme: "HTTPS",
									// TODO Replace hardcoded reference to garygeeke (relevant with platform security)
									Path: "/open-metadata/platform-services/users/garygeeke/server-platform/origin",
								},
							},
							InitialDelaySeconds: 30,
							PeriodSeconds:       10,
							FailureThreshold:    15,
						},
					},
					},
				},
			},
		},
	}
	// Set Egeria instance as the owner and controller
	// TODO: resolve management of references
	ctrl.SetControllerReference(egeriaInstance, deployment, reconciler.Scheme)
	return deployment
}

// TODO: Go equiv of javadoc?
// serviceForEgeria returns an egeria Service  object
func (reconciler *EgeriaPlatformReconciler) serviceForEgeriaPlatform(egeriaInstance *egeriav1alpha1.EgeriaPlatform) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      egeriaInstance.Name + "-service",
			Namespace: egeriaInstance.Namespace,
			Labels:    egeriaLabels(egeriaInstance.Name, "service"),
		},
		Spec: corev1.ServiceSpec{
			// TODO: More flexible service types needed in future (for now can be exposed manually after deployment-cloud dependent)
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Port:     9443,
					Protocol: corev1.ProtocolTCP,
					Name:     "egeria-port",
				},
			},
			Selector: getServiceSelectorLabels(egeriaInstance.Name),
		},
	}
	// Set Egeria instance as the owner and controller
	// TODO: resolve management of references
	ctrl.SetControllerReference(egeriaInstance, service, reconciler.Scheme)
	return service
}

// TODO: Go equiv of javadoc?
// serviceForEgeria returns an egeria Service  object
// If configmap changes we need to update this
func (reconciler *EgeriaPlatformReconciler) configmapForEgeriaPlatform(ctx context.Context, egeriaInstance *egeriav1alpha1.EgeriaPlatform) *corev1.ConfigMap {

	// Figure out autostart list
	var autostart = make(map[string]string)
	//TODO avoid hardcoding serverame
	//autostart["STARTUP_SERVER_LIST"] = "cocoMDS4, cocoMDS1, cocoView1"
	autostart["STARTUP_SERVER_LIST"], _ = reconciler.getServersFromConfigMap(ctx, egeriaInstance)

	configmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      egeriaInstance.Name + "-autostart",
			Namespace: egeriaInstance.Namespace,
			Labels:    egeriaLabels(egeriaInstance.Name, "autostart"),
		},
		Data: autostart,
	}
	// Set Egeria instance as the owner and controller
	// TODO: resolve management of references
	// TODO Ensure configmap gets deleted when no longer required
	_ = ctrl.SetControllerReference(egeriaInstance, configmap, reconciler.Scheme)
	return configmap
}

// egeriaLabels returns the labels we set on created resources
// belonging to the given egeria CR name.
// see also https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
func egeriaLabels(egeriaInstanceName string, componentName string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "egeriaplatform",
		"app.kubernetes.io/instance":   egeriaInstanceName, // name of Custom Resource
		"app.kubernetes.io/version":    "0.9",              // TODO Figure out version to use
		"app.kubernetes.io/component":  componentName,
		"app.kubernetes.io/part-of":    "Egeria",
		"app.kubernetes.io/managed-by": "Operator",
		"app.kubernetes.io/created-by": "egeriaplatform_controller",
	}
}

// Service selector labels - this is the criteria used for directing requests to pods
func getServiceSelectorLabels(crName string) map[string]string {
	return egeriaLabels(crName+"-deployment", "deployment")
}

// SetupWithManager sets up the controller with the Manager.
//
// watches our CRs primarily, but additionally Deployments & Services - as we create them
func (reconciler *EgeriaPlatformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: Can we make the watches more efficient with conditions?
	return ctrl.NewControllerManagedBy(mgr).
		For(&egeriav1alpha1.EgeriaPlatform{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(reconciler)
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

//func getVolumeMounts() []corev1.VolumeMount {
//	return []corev1.VolumeMount{
//		// This contains all of the server configuration documents for the platform
//		{
//			Name: "serverconfig",
//			// Do not permit updating the configuration
//			ReadOnly: true,
//			// TODO: Need to aggregate data from multiple files? Or combine configs into one map
//			MountPath: "/deployments/data/servers",
//		},
//	}
//}

func getVolumes(configname string) []corev1.Volume {
	// In future we may have multiple volumes, so extracted out to fn
	return []corev1.Volume{{
		Name: "serverconfig",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configname,
				},
			},
		},
	}}
}

// TODO : Automatically create autostart server list from querying config

func (reconciler *EgeriaPlatformReconciler) getServersFromConfigMap(ctx context.Context, egeria *egeriav1alpha1.EgeriaPlatform) (servers string, err error) {
	// Retrieve the list of servers from the configmap
	configmap := &corev1.ConfigMap{}
	err = reconciler.Get(ctx, types.NamespacedName{Name: egeria.Spec.ServerConfig, Namespace: egeria.Namespace}, configmap)
	if err != nil {
		if errors.IsNotFound(err) {
			log.FromContext(ctx).Error(err, "Configured server configmap not found.")
			return "", err
		}
		// Error reading the object - requeue the request.
		log.FromContext(ctx).Error(err, "Error reading server configmap")
		return "", err
	}
	// We now have the configmap object. need to extract server names
	for k := range configmap.Data {
		log.FromContext(ctx).Info("Found server config & adding to startup list: ", "servername", k)
		servers += k + ","
	}
	return servers, nil
}

// TODO : Create ingress As ervice

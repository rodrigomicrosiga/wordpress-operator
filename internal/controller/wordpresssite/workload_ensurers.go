package wordpresssite

import (
	"context"

	"github.com/cloud104/reconciler/v2"
	"github.com/rodrigomicrosiga/wordpress-operator/api/v1alpha1"
	"github.com/rodrigomicrosiga/wordpress-operator/internal/factory"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// --- 3. DatabaseStatefulSetEnsurer ---
type DatabaseStatefulSetEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *DatabaseStatefulSetEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	desired := factory.BuildMySQLStatefulSet(site)
	obj := &appsv1.StatefulSet{ObjectMeta: ctrl.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace}}

	_, err := controllerutil.CreateOrUpdate(ctx, e.Client, obj, func() error {
		obj.Labels = desired.Labels
		obj.Spec.Replicas = desired.Spec.Replicas
		obj.Spec.Selector = desired.Spec.Selector
		obj.Spec.Template = desired.Spec.Template
		if obj.CreationTimestamp.IsZero() {
			obj.Spec.VolumeClaimTemplates = desired.Spec.VolumeClaimTemplates // Imutável após criação
		}
		return controllerutil.SetControllerReference(site, obj, e.Scheme)
	})
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}
	return e.Next(ctx, site)
}

// --- 4. DatabaseServiceEnsurer ---
type DatabaseServiceEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *DatabaseServiceEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	desired := factory.BuildMySQLService(site)
	obj := &corev1.Service{ObjectMeta: ctrl.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace}}

	_, err := controllerutil.CreateOrUpdate(ctx, e.Client, obj, func() error {
		obj.Labels = desired.Labels
		obj.Spec.Selector = desired.Spec.Selector
		obj.Spec.Ports = desired.Spec.Ports
		obj.Spec.ClusterIP = desired.Spec.ClusterIP
		return controllerutil.SetControllerReference(site, obj, e.Scheme)
	})
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}
	return e.Next(ctx, site)
}

// --- 5. WordpressConfigMapEnsurer ---
type WordpressConfigMapEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *WordpressConfigMapEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	desired := factory.BuildWordpressConfigMap(site)
	obj := &corev1.ConfigMap{ObjectMeta: ctrl.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace}}

	_, err := controllerutil.CreateOrUpdate(ctx, e.Client, obj, func() error {
		obj.Labels = desired.Labels
		obj.Data = desired.Data
		return controllerutil.SetControllerReference(site, obj, e.Scheme)
	})
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}
	return e.Next(ctx, site)
}

// --- 6. WordpressPVCEnsurer ---
type WordpressPVCEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *WordpressPVCEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	desired := factory.BuildWordpressPVC(site)
	obj := &corev1.PersistentVolumeClaim{ObjectMeta: ctrl.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace}}

	_, err := controllerutil.CreateOrUpdate(ctx, e.Client, obj, func() error {
		obj.Labels = desired.Labels
		// AccessModes e Resources geralmente são imutáveis após a criação do PVC
		if obj.CreationTimestamp.IsZero() {
			obj.Spec.AccessModes = desired.Spec.AccessModes
			obj.Spec.Resources = desired.Spec.Resources
		}
		return controllerutil.SetControllerReference(site, obj, e.Scheme)
	})
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}
	return e.Next(ctx, site)
}

// --- 7. WordpressDeploymentEnsurer ---
type WordpressDeploymentEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *WordpressDeploymentEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	desired := factory.BuildWordpressDeployment(site)
	obj := &appsv1.Deployment{ObjectMeta: ctrl.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace}}

	_, err := controllerutil.CreateOrUpdate(ctx, e.Client, obj, func() error {
		obj.Labels = desired.Labels
		obj.Spec.Replicas = desired.Spec.Replicas
		obj.Spec.Selector = desired.Spec.Selector
		obj.Spec.Template = desired.Spec.Template
		return controllerutil.SetControllerReference(site, obj, e.Scheme)
	})
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}
	return e.Next(ctx, site)
}

// --- 8. WordpressServiceEnsurer ---
type WordpressServiceEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *WordpressServiceEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	desired := factory.BuildWordpressService(site)
	obj := &corev1.Service{ObjectMeta: ctrl.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace}}

	_, err := controllerutil.CreateOrUpdate(ctx, e.Client, obj, func() error {
		obj.Labels = desired.Labels
		obj.Spec.Selector = desired.Spec.Selector
		obj.Spec.Ports = desired.Spec.Ports
		obj.Spec.Type = desired.Spec.Type
		return controllerutil.SetControllerReference(site, obj, e.Scheme)
	})
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}
	return e.Next(ctx, site)
}

// --- 9. IngressEnsurer ---
type IngressEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *IngressEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	desired := factory.BuildWordpressIngress(site)
	obj := &networkingv1.Ingress{ObjectMeta: ctrl.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace}}

	_, err := controllerutil.CreateOrUpdate(ctx, e.Client, obj, func() error {
		obj.Labels = desired.Labels
		obj.Spec.Rules = desired.Spec.Rules
		obj.Spec.IngressClassName = desired.Spec.IngressClassName
		return controllerutil.SetControllerReference(site, obj, e.Scheme)
	})
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}
	return e.Next(ctx, site)
}

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloud104/reconciler/v2"
	wordpressv1alpha1 "github.com/rodrigomicrosiga/wordpress-operator/api/v1alpha1"
	"github.com/rodrigomicrosiga/wordpress-operator/internal/controller/wordpresssite"
)

type WordpressSiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// RBAC: Permissões para o Operator interagir com o cluster
// +kubebuilder:rbac:groups=wordpress.cloud104.io,resources=wordpresssites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wordpress.cloud104.io,resources=wordpresssites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets;services;configmaps;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

func (r *WordpressSiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	site := &wordpressv1alpha1.WordpressSite{}
	if err := r.Get(ctx, req.NamespacedName, site); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err) // Se foi deletado, ignora
	}

	// Dispara a Corrente de Responsabilidade
	return r.buildChain().Reconcile(ctx, site)
}

// buildChain monta os elos na ordem de dependência correta
func (r *WordpressSiteReconciler) buildChain() reconciler.Handler[*wordpressv1alpha1.WordpressSite] {
	return reconciler.Chain(
		&wordpresssite.DatabaseSecretEnsurer{Client: r.Client, Scheme: r.Scheme},
		&wordpresssite.DatabaseStatefulSetEnsurer{Client: r.Client, Scheme: r.Scheme},

		// NOVO ENSURER ADICIONADO AQUI:
		&wordpresssite.DatabasePVCOwnerEnsurer{Client: r.Client, Scheme: r.Scheme},

		&wordpresssite.DatabaseServiceEnsurer{Client: r.Client, Scheme: r.Scheme},
		&wordpresssite.WordpressConfigMapEnsurer{Client: r.Client, Scheme: r.Scheme},
		&wordpresssite.WordpressPVCEnsurer{Client: r.Client, Scheme: r.Scheme},
		&wordpresssite.WordpressDeploymentEnsurer{Client: r.Client, Scheme: r.Scheme},
		&wordpresssite.WordpressServiceEnsurer{Client: r.Client, Scheme: r.Scheme},
		&wordpresssite.IngressEnsurer{Client: r.Client, Scheme: r.Scheme},
		&wordpresssite.StatusEnsurer{Client: r.Client},
	)
}

// SetupWithManager registra o controller no cluster e define os gatilhos de observabilidade
func (r *WordpressSiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wordpressv1alpha1.WordpressSite{}).
		Owns(&corev1.Secret{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&appsv1.Deployment{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}

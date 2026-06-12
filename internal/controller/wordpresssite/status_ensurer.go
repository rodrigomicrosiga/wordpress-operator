package wordpresssite

import (
	"context"

	"github.com/cloud104/reconciler/v2"
	"github.com/rodrigomicrosiga/wordpress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatusEnsurer verifica a saúde dos recursos e atualiza o campo Status da CRD
type StatusEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
}

func (e *StatusEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	// 1. Inicializa os dados básicos com inteligência de protocolo
	protocol := "http://"
	if site.Spec.TLS != nil && site.Spec.TLS.Enabled {
		protocol = "https://"
	}
	site.Status.URL = protocol + site.Spec.Domain

	site.Status.ObservedGeneration = site.Generation

	// 2. Checa a saúde do Banco de Dados (StatefulSet)
	sts := &appsv1.StatefulSet{}
	errDB := e.Client.Get(ctx, types.NamespacedName{Name: site.Name + "-mysql", Namespace: site.Namespace}, sts)
	dbReady := errDB == nil && sts.Status.ReadyReplicas > 0

	// 3. Checa a saúde da Aplicação (Deployment)
	dep := &appsv1.Deployment{}
	errApp := e.Client.Get(ctx, types.NamespacedName{Name: site.Name, Namespace: site.Namespace}, dep)
	appReady := errApp == nil && dep.Status.ReadyReplicas > 0

	// 4. Agrega a fase final
	if dbReady && appReady {
		site.Status.Phase = "Ready"
	} else {
		site.Status.Phase = "Provisioning"
	}

	// 5. Salva no sub-recurso de Status do Kubernetes
	err := e.Client.Status().Update(ctx, site)
	if err != nil {
		return e.RequeueOnErr(ctx, err)
	}

	return e.Next(ctx, site)
}

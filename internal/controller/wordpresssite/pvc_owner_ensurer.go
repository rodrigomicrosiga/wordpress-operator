package wordpresssite

import (
	"context"
	"fmt"

	"github.com/cloud104/reconciler/v2"
	"github.com/rodrigomicrosiga/wordpress-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// DatabasePVCOwnerEnsurer garante que o PVC do banco de dados seja deletado junto com a aplicação
type DatabasePVCOwnerEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *DatabasePVCOwnerEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	// O Kubernetes nomeia o PVC dinâmico juntando: <nome-do-volume-claim>-<nome-do-statefulset>-0
	pvcName := fmt.Sprintf("mysql-data-%s-mysql-0", site.Name)
	pvc := &corev1.PersistentVolumeClaim{}

	err := e.Client.Get(ctx, types.NamespacedName{Name: pvcName, Namespace: site.Namespace}, pvc)
	if err != nil {
		if errors.IsNotFound(err) {
			// O StatefulSet acabou de ser criado e ainda não gerou o disco.
			// Passamos o bastão; o loop rodará novamente em breve quando o status do banco atualizar.
			return e.Next(ctx, site)
		}
		return e.RequeueOnErr(ctx, err)
	}

	original := pvc.DeepCopy()

	// Utilizamos SetOwnerReference em vez de SetControllerReference.
	// Isso injeta a paternidade para o Garbage Collection funcionar sem conflitar
	// com a gestão nativa que o próprio StatefulSet já faz do recurso.
	if err := controllerutil.SetOwnerReference(site, pvc, e.Scheme); err != nil {
		return e.RequeueOnErr(ctx, err)
	}

	// Aplicamos a mutação no cluster apenas se houver diferença (Patch)
	if err := e.Client.Patch(ctx, pvc, client.MergeFrom(original)); err != nil {
		return e.RequeueOnErr(ctx, err)
	}

	return e.Next(ctx, site)
}

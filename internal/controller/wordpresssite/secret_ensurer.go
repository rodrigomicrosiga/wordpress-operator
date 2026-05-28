package wordpresssite

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/cloud104/reconciler/v2"
	"github.com/rodrigomicrosiga/wordpress-operator/api/v1alpha1"
	"github.com/rodrigomicrosiga/wordpress-operator/internal/factory"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// DatabaseSecretEnsurer garante que o Secret com a senha do MySQL exista.
type DatabaseSecretEnsurer struct {
	reconciler.Funcs[*v1alpha1.WordpressSite]
	Client client.Client
	Scheme *runtime.Scheme
}

func (e *DatabaseSecretEnsurer) Reconcile(ctx context.Context, site *v1alpha1.WordpressSite) (ctrl.Result, error) {
	secretName := fmt.Sprintf("%s-mysql-secret", site.Name)
	existingSecret := &corev1.Secret{}

	// 1. Tenta buscar o Secret existente no cluster
	err := e.Client.Get(ctx, types.NamespacedName{Name: secretName, Namespace: site.Namespace}, existingSecret)

	var password string
	if err != nil {
		if apierrors.IsNotFound(err) {
			// 2A. Não existe: Gera uma nova senha forte
			password = generateRandomPassword(16)
		} else {
			return e.RequeueOnErr(ctx, err)
		}
	} else {
		// 2B. Já existe: Pega a senha atual para não sobrescrever (IDEMPOTÊNCIA)
		passwordBytes := existingSecret.Data["mysql-password"]
		password = string(passwordBytes)
	}

	// 3. Usa a Factory para montar como o Secret deve ser
	desired := factory.BuildMySQLSecret(site, password)

	// Objeto vazio que será preenchido ou atualizado
	sec := &corev1.Secret{}
	sec.Name = desired.Name
	sec.Namespace = desired.Namespace

	// 4. CreateOrUpdate faz a mágica de aplicar a diferença
	_, err = controllerutil.CreateOrUpdate(ctx, e.Client, sec, func() error {
		sec.Labels = desired.Labels
		if sec.Data == nil {
			sec.Data = make(map[string][]byte)
		}
		// A stringData não é lida de volta no Get, injetamos no Data
		sec.Data["mysql-password"] = []byte(password)
		sec.Data["mysql-root-password"] = []byte(password)

		// Owner Reference garante que se o WordpressSite for deletado, o Secret também será
		return controllerutil.SetControllerReference(site, sec, e.Scheme)
	})

	if err != nil {
		return e.RequeueOnErr(ctx, fmt.Errorf("falha ao reconciliar Secret: %w", err))
	}

	// 5. Passa o bastão para o próximo Ensurer na Chain
	return e.Next(ctx, site)
}

// generateRandomPassword cria uma senha aleatória base64
func generateRandomPassword(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

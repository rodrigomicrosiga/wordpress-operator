# WordPress Operator - Desafio de Onboarding

Desafio proposto por: Fabrizio Malta Di Napoli (https://github.com/fmnapoli) e Chris-Corinthians (https://github.com/christian-devops)

Este repositório contém a implementação de um Kubernetes Operator desenvolvido do zero para gerenciar um site WordPress completo (Aplicação + Banco de Dados MySQL + Rede + Storage) de forma declarativa.

## 🚀 Diário de Desenvolvimento

### [22/05/2026] - Desafio Recebido
**Objetivo:** Entendimento de todo escopo

### [23/05/2026] - Estudos
### [24/05/2026] - Estudos
### [25/05/2026] - Estudos

### [26/05/2026] - Passo 1: Setup Inicial e Fundação do Projeto
**Objetivo:** Estruturar o repositório, inicializar o projeto com Kubebuilder e definir o escopo do desafio.

**Ações Realizadas:**
* Criação do repositório no GitHub para rastreabilidade via commits descritivos.
* Definição da arquitetura base utilizando o padrão *Chain of Responsibility* com a biblioteca `cloud104/reconciler/v2`.
* Preparação do ambiente local com as ferramentas exigidas (Go 1.24+, Kubebuilder v4, Kustomize, Kind, Kubectl e Docker).

**Próximos Passos:** * Executar o *scaffold* do Kubebuilder.
* Criar a API (CRD) `WordpressSite` mapeando o *Spec* e o *Status*.

**Pontos Intermediários:**
* Montagem do esqueleto do operator, utilizando os comandos exigidos pelo desafio
1.Inicializar o módulo Go e o domínio do Operator:
kubebuilder init --domain cloud104.io --repo github.com/SEU_USUARIO/wordpress-operator
2.Criar a API (O Custom Resource Definition - CRD):
kubebuilder create api --group wordpress --version v1alpha1 --kind WordpressSite --resource --controller
3.Baixar a biblioteca de Reconciler exigida pelo time:
go get github.com/cloud104/reconciler/v2@latest




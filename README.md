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

**Pontos Importantes do Passo 1:**
* Montagem do esqueleto do operator, utilizando os comandos exigidos pelo desafio

1.Inicializar o módulo Go e o domínio do Operator:

kubebuilder init --domain cloud104.io --repo github.com/rodrigomicrosiga/wordpress-operator

2.Criar a API (O Custom Resource Definition - CRD):

kubebuilder create api --group wordpress --version v1alpha1 --kind WordpressSite --resource --controller

3.Baixar a biblioteca de Reconciler exigida pelo time:

go get github.com/cloud104/reconciler/v2@latest

### [26/05/2026] - Passo 2: O Contrato (API / CRD)
**Objetivo:** Modelar a API declarativa do `WordpressSite`, definindo o estado desejado (`Spec`) e o estado observado (`Status`).

**Ações Realizadas:**
* Edição do arquivo `api/v1alpha1/wordpresssite_types.go` para mapear os requisitos da Seção 4 do desafio.
* Utilização de *markers* do Kubebuilder (`//+kubebuilder:...`) para definir validações, valores default e colunas personalizadas no terminal (`kubectl get`).
* Geração dos manifestos YAML das CRDs através do comando `make manifests`.

**Pontos Importantes do Passo 2:**
* Modelagem da API

1.WordpressSite?

O Kubernetes precisa entender o que é um `WordpressSite`.
O arquivo que define isso é o `api/v1alpha1/wordpresssite_types.go`.
Serão estruturadas todas as structs(`Spec` e `Status`).
Serão adicionados os "Markers" que acabam sendo comentários especiais que ensinam o Kubernetes a validar os dados e incluir valores default (exemplo: `replicas: 1`)

**IMPORTANTE:**
* Problemas ao alterar o arquivo `api/v1alpha1/wordpresssite_types.go`.

`api/v1alpha1/zz_generate.deepcopy.go` passou a apresentar erros como "has no field or method Foo"

**O arquivo em destaque ficou desatualizado e provavelmente será atualizado após a execução do `make generate`.**

`go.mod` passou a apresentar erros como "github.com/cloud104/reconcilier/v2 is not used in this module"

**Inconsistência entre o que tenho de local e o que foi baixado, e provavelmente será atualizado após a execução do `go mod tidy` que irá varrer o projeto, remover o que não é necessário e realizar download do que está faltando.**

2.Manifestos

Nessa etapa todo código Go irá se transformar em YAML de CRD que será interpretado pelo cluster.
Se tudo for bem sucedido, novos arquivos deverão ser criados no diretório `config/crd/bases/`.

make generate

make manifests

go mod tidy

### [26/05/2026] - Passo 3: Factory Pattern (Geração Pura de Manifestos)
**Objetivo:** Isolar a definição da infraestrutura (O *Shape* dos recursos) da lógica de reconciliação (O *Action*).

**Ações Realizadas:**
* Criação do pacote `internal/factory/factory.go`.
* Implementação de funções puras que recebem a CRD `WordpressSite` e retornam objetos nativos do Kubernetes (`corev1`, `appsv1`, `networkingv1`).
* Mapeamento de relacionamentos como injeção de `Secrets` como variáveis de ambiente e definição de `VolumeClaimTemplates` para persistência de dados.



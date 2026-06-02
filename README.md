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

**Problemas:**
* Problemas ao alterar o arquivo `api/v1alpha1/wordpresssite_types.go`.
* `api/v1alpha1/zz_generate.deepcopy.go` passou a apresentar erros como "has no field or method Foo"

**Correção**
* Execução do `make generate`.

**Problemas**
* `go.mod` passou a apresentar erros como "github.com/cloud104/reconcilier/v2 is not used in this module"

**Correção**
* Inconsistência entre o que tenho de local e o que foi baixado, executado `go mod tidy` removendo o que não é necessário e realizando download do que está faltando

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

### [26/05/2026] - Passo 4: Implementação da Chain of Responsibility (Ensurers)
**Objetivo:** Criar os elos de reconciliação utilizando a biblioteca `cloud104/reconciler`.

**Ações Realizadas:**
* Criação do diretório `internal/controller/wordpresssite/` para abrigar os ensurers isolados.
* Implementação do `DatabaseSecretEnsurer`, garantindo a geração segura e idempotente da senha do MySQL (gerada apenas na primeira execução e preservada nas seguintes).
* Uso de `controllerutil.CreateOrUpdate` para mutação segura de estado e `SetControllerReference` para Garbage Collection automático.

### [27/05/2026] - Estudos e Entendimento das Correções
**Objetivo:** Entender a quebra ocorrida após a alteração no arquivo `internal/controller/wordpresssite/secret_ensurer.go`

**IMPORTANTE**
* O erro acima provavelmente se refere a alguma incompatibilidade entre a função `Reconcile` e a biblioteca `cloud104/reconciler`. Nesse tópico de resolução utilizei IA para entender mais a fundo e com isso foi possível verificar que o retorno estava definido como `reconciler.Result, error` e a biblioteca do time reaproveita o tipo nativo do Kubernetes que é `Result` e faz parte do pacote `controller-runtime`.
* No item 2 da seçaõ 3 é exigido o uso do pacote `ctrl`. E as funções `e.Next()` e `e.RequeueOnErr()` retornam `ctrl.Result` houve uma quebra ao tentar devolver isso como um `reconciler.Result`.

**CORREÇÃO**
* `internal/controller/wordpresssite/secret_ensurer.go`
* Inclusão do alias `ctrl` no import.
* Alteração da função substituindo `reconciler.Result` por `ctrl.Result`. 
* Realizado `go mod tidy` para atualização de referências.

### [28/05/2026] - Passo 5: Implementação dos Workload Ensurers (Idempotência Estrutural)
**Objetivo:** Consolidar os elos da corrente de reconciliação responsáveis pela infraestrutura base da aplicação (Banco de Dados, Aplicação, Storage e Rede) garantindo idempotência e ciclo de vida atrelado ao CRD.

**Ações Realizadas:**
* Criação do arquivo `internal/controller/wordpresssite/workload_ensurers.go` agrupando os ensurers estruturais (StatefulSet, Deployment, Services, ConfigMap, PVC e Ingress).
* Utilização da função `controllerutil.CreateOrUpdate` em todos os recursos para garantir que o estado real do cluster sempre convirja para o estado desejado definido pela *Factory*.
* Implementação de proteções lógicas para campos imutáveis da API do Kubernetes (como `VolumeClaimTemplates` no StatefulSet e especificações de PVC) através da validação `obj.CreationTimestamp.IsZero()`.
* Vinculação de todos os recursos criados ao objeto pai (`WordpressSite`) via `SetControllerReference`, assegurando o *Garbage Collection* nativo na deleção do operator.

### [28/05/2026] - Passo 6: Observabilidade e Fechamento do Controlador
**Objetivo:** Implementar o reporte de status em tempo real e registrar o Operator no Kubernetes com seus gatilhos de escuta.

**Ações Realizadas:**
* Criação do `StatusEnsurer` para avaliar se os *Pods* do banco e da aplicação atingiram o estado `Ready` antes de atualizar a CRD.
* Refatoração do `wordpresssite_controller.go` para inicializar a *Chain of Responsibility* via `r.buildChain()`.
* Mapeamento de RBAC (`//+kubebuilder:rbac`) para conceder privilégios ao Operator para manipular os recursos *core*, *apps* e *networking*.
* Configuração do `SetupWithManager` utilizando `Owns()` para garantir a re-reconciliação imediata caso um usuário ou processo modifique os recursos gerenciados indevidamente.

**Problemas**
* `unnecessary type arguments` passou a ser apresentado após a alteração do `wordpresssite_controller.go`

**Correção**
* Na linha 45 informamos qual objeto a função `reconciler.Chain` iria manipular. Porém todos os `ensurers` já tem o tipo `*WordpressSite` embutido e com isso praticamente o Go avisa que não precisa reescrever o tipo.
* `internal/controller/wordpresssite_controller.go` foi removida a declaração explícita `*wordpressv1alpha1.WordpressSite]`

make manifests (que fará a leitura dos novos RBACS)

make generate

**Execução via Cluster**

* Preciso garantir que o Kind esteja rodando

kind get clusters

* Como nenhum cluster havia sido criado

kind create cluster (criação do cluster)

kubectl get nodes (para verificação)

* Apresentar ao Kubernetes o que será compilado

make install

* Iniciar o Operator

make run

Nessa etapa o terminal fica apresentando logs de inicialização.

**Manifesto**

* Em uma outra sessão de terminal foi realizada a criação de um diretorio `mkdir teste` e a criação do arquivo `touch teste.yaml`.

kubectl apply -f teste.yaml

Com essa ação foi possível verificar que o cluster respondeu corretamente e atualizou as informações:

* Nesse ponto solicitei ajuda para a IA com o propósito de entender tudo que de fato estava sendo realizado pelo cluster nesse momento. Então foi possível entender que a `Chain of Responsibility` foi ativada em milisegundos, loop de reconciliação passando por cada `Ensurer`, gerando a senha forte, subindo o `StatefulSet` do banco, montando os PVCs e por fim o Deployment do WordPress.

* Visualizar a tabela customizada criada via api

kubectl get wordpresssite

* Visualizar toda infraestrutura materializada via código

kubectl get all,pvc,ingress,secret

* Testar a deleção (Garbage Collection)

kubectl delete wordpresssite meu-blog

**Religando o Operator e Recriando a Infraestrutura**

* Religando o operator

make run (terminal principal)

* Recriando a infra

kubectl apply -f teste.yaml (terminal secundário)

* Aguardando a inicialização

kubectl get pods (aqui o status esperado é `Running` e `Ready`)

**Definition Of Done (Seção 9)**

* Tentativa de acesso via navegador:

kubectl port-forward svc/meu-blog 8080:80

**Problema**

* Ao tentar realizar o acesso via `http://127.0.0.1:8080` não foi apresentada a tela de instalação do WordPress (apresentado `Error establishing a database connection`)

**Entendimento do Problema**

O entendimento inicial foi superficial e na linha "teoricamente o database ainda não estava pronto no momento em que tentei acessar".

Indo por essa linha realizei a parada do direcionamento de porta (`ctrl+c`) e forcei a exclusão do pod `kubectl delete pod -l app.kubernetes.io/component=application`.

Aguardei alguns segundos para a subida de um novo pod, realizei a consulta via `kubectl get pods` e obtive o status `Running`em ambos, e na sequencia realizei `kubectl port-forward svc/meu-blog 8080:80`.

Porém ao tentar acessar o endereço `http://localhost:8080` ainda estava sendo apresentado o erro.

Com a ajuda da IA, foi possível obter informações importantes sobre `PVC` e o ciclo de vida do `StatefulSet`.

Como isso entendi que ao utilizar o `delete` praticamente o operator:
* eliminou o Deployment
* eliminou o Service
* eliminou o Secret (senha)
* eliminou o StatefulSet do MySQL

Ao deletar o `StatefulSet` o Kubernetes não elimina os discos (PVC) criados pelo `VolumeClaimTemplates`.

Ao realizar novamente a aplicação do `teste.yaml` foi gerado um novo `Secret`, o `StatefulSet` novo subir e se reconectou ao antigo disco que ainda existia. O MySQL percebeu que o disco ainda existia, ignorou a nova senha, mantendo o que já estava gravado no database.

O WordPress leu o novo `Secret`, tentou conectar com a senha nova, porém o banco exigia a senha antiga, resultando no erro `Error establishing a database connection`.

**Correção**

* Deletar a aplicação

kubectl delete wordpresssite meu-blog

* Verificar os PVCs que ainda podem existir

kubectl get pvc

Nesse momento foi possível verificar o PVC que ainda existia.

* Deletar o PVC

kubectl delete pvc mysql-data-meu-blog-mysql-0

* Aplicar novamente

kubectl apply -f teste.yaml

* Aguardar a subida e consultar os pods (estado esperado `Running`)

kubectl get pods

* Direcionamento de portas

kubectl port-forward svc/meu-blog 8080:80

Sucesso ao acessar o endereço `http://localhost:8080` onde foi apresentada a tela de instalação do WordPress.

### [01/06/2026] - Estudos

### [02/06/2026] - Novas Implementações

* Implementar regra para que o OwnerReference não permita PVC orfão ao realizar o delete da aplicação

Nessa situação, foi possível entender que o `StatefulSet`cria os PVCs de forma dinâmica através do `VolumeClaimTemplates` porém não repassa de forma automática a `OwnerReference` do `Wordpresssite` para os discos.

E para atender ao requisito, foi necessário:

* Criar o arquivo `internal/controller/wordpresssite/pvc_owner_ensurer.go` para que seja criado um elo de paternidade.

Feito isso, se tornou necessário informar ao operator em que momento isso deverá ocorrer, e o mais ideal é que seja após a reconciliação do `Statefulset`.

E para isso foi necessário: 

* Realizar a alteração da função `buildChain` no arquivo `internal/controller/wordpresssite_controller.go` inserindo o novo `DatabasePVCOwnerEnsurer`.

**Testes**

Para facilitar a utilização de "n" terminais foi adoto a utilização do `Tilix`.

* Terminal 1.

Executado `make run` para compilação do código atualizado.

* Terminal 2.

Executado `kubectl apply -f teste.yaml` para recriar a aplicação aplicando o manifesto para criação de toda infraestrutura.

* Terminal 3.

Aguardar alguns segundos para criação via `StatefulSet` e o operator injetar a referência. Com isso poderemos inspecionar o manifesto do PVC gerado no cluster.

Executado `kubectl get pvc mysql-data-meu-blog-mysql-0 -o yaml`

Agora devemos procurar no retorno a seção `metadata` e verificar se existe o bloco `ownerReferences` que relaciona o disco ao CRD.

* Terminal 4.

Vamos simular a remoção da aplicação.

Executado `kubectl delete wordpresssite meu-blog`

* Terminal 5.

Vamos verificar e garantir que o PVC orfão não existe mais no cluster.

Executado `kubectl get pvc`

O retorno `No resources found in default namespace` indica que não existe nenhum PVC orfão.

Com isso foi possível entender que o `Kubernetes` ao deletar o `WordpressSite` leu o `OwnerReference` que foi injetado e eliminou automaticamente junto com os demais itens da infraestrutura.








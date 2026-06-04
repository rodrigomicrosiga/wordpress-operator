# 🚀 WordPress Operator for Kubernetes

Este repositório contém a implementação de um **Kubernetes Operator** desenvolvido do zero para gerenciar um site WordPress completo (Aplicação + Banco de Dados MySQL + Rede + Storage + Certificados TLS) de forma declarativa.

> **Desafio de Onboarding** proposto por [Fabrizio Malta Di Napoli](https://github.com/fmnapoli) e [Chris-Corinthians](https://github.com/christian-devops).

---

## 🏗️ Arquitetura e Fluxo de Reconciliação (Chain of Responsibility)

O Operator utiliza o padrão de design **Chain of Responsibility** (através da biblioteca `cloud104/reconciler/v2`). O fluxograma abaixo demonstra como o Operator reage à criação de um Custom Resource (`WordpressSite`) e encadeia a criação de cada componente no cluster, garantindo idempotência e segurança:

```mermaid
graph TD
    CRD((Criar CRD<br>WordpressSite)) --> Operator{Operator Watcher}
    Operator --> S[Database Secret]
    S --> SS[Database StatefulSet]
    SS --> PVC[Injeta OwnerReference<br>no PVC Dinâmico]
    PVC --> DBSVC[Database Service]
    DBSVC --> WPCM[Wordpress ConfigMap]
    WPCM --> WPPVC[Wordpress PVC]
    WPPVC --> WPD[Wordpress Deployment]
    WPD --> WPSVC[Wordpress Service]
    WPSVC --> ING[Ingress com TLS<br>Cert-Manager]
    ING --> F((Fim do Ciclo<br>Aplicação Pronta))

    style CRD fill:#2ecc71,stroke:#27ae60,stroke-width:2px,color:#fff
    style Operator fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#fff
    style F fill:#2ecc71,stroke:#27ae60,stroke-width:2px,color:#fff

📋 Pré-requisitos
Antes de iniciar, certifique-se de ter instalado em sua máquina:

Go (v1.24+)

Docker

kubectl

Kind (para criação do cluster local)

🚀 1. Utilizando em Cluster Local (Desenvolvimento)
Ideal para testes rápidos e desenvolvimento na sua própria máquina.

1. Suba o cluster local e baixe as dependências:

Bash
kind create cluster
go mod tidy
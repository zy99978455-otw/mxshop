<div align="center">

# MxShop - Cloud-Native Microservices Architecture

![Go](https://img.shields.io/badge/Go-1.20%2B-blue?logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-24.0%2B-2496ED?logo=docker&logoColor=white)
![Gin](https://img.shields.io/badge/Framework-Gin-00ADD8?logo=go&logoColor=white)
![gRPC](https://img.shields.io/badge/RPC-gRPC-244c5a?logo=google&logoColor=white)
![Elasticsearch](https://img.shields.io/badge/Search-Elasticsearch%208.x-005571?logo=elasticsearch&logoColor=white)

A high-concurrency e-commerce microservice system built with **Golang (Gin, gRPC)**, orchestrated via **Docker**.

[Architecture] • [Getting Started] • [Documentation]

</div>

---

## 🏗 Architecture Overview

MxShop adopts a modern cloud-native approach, separating the API gateway from core service logic, integrating full-text search capabilities.

### 🛠 Technology Stack

| Component | Technology | Modules / Description |
| :--- | :--- | :--- |
| **API Gateway** | Gin Framework | `User-Web`, `Goods-Web` (HTTP Entry) |
| **Service Layer** | gRPC + Protobuf | `User-Srv`, `Goods-Srv` (Core Logic) |
| **Config Center** | Nacos 2.x | Dynamic Configuration Management |
| **Service Discovery** | Consul | Service Registration & Health Checks |
| **Database** | MySQL 8.0 | Database-per-service pattern |
| **Search Engine** | Elasticsearch 8.x | High-performance Goods Search |
| **Containerization** | Docker | Docker & Docker Compose Orchestration |

### 🔒 Port Allocation Strategy

To ensure security and avoid port conflicts, we implement a strict port isolation strategy.
> **Note:** RPC services and Databases are NOT exposed to the host network directly; they communicate within the Docker network.

| Service | Type | Container Port | Host Port | Description |
| :--- | :---: | :---: | :---: | :--- |
| **Nacos** | Infra | `8848` | `8848` | Configuration Console |
| **Consul** | Infra | `8500` | `8500` | Service Discovery Console |
| **Elasticsearch** | Infra | `9200` | `59200` | Search Engine API (**Mapped**) |
| **Kibana** | Infra | `5601` | `5601` | Search Visualization |
| **User Web** | API | `8080` | `8080` | Public HTTP Entry for **Users** |
| **Goods Web** | API | `8081` | `8081` | Public HTTP Entry for **Goods** |
| **User Srv** | RPC | `50051` | `-` | **Internal Network Only** |
| **Goods Srv** | RPC | `50052` | `-` | **Internal Network Only** |

---

## 🚀 Quick Start

### 1. Infrastructure Setup
Start the base services (MySQL, Nacos, Consul, ES, Network).

> **⚠️ Pre-requisites for Elasticsearch:**
> * **Memory Limit:** Ensure your `docker-compose.es.yml` sets `ES_JAVA_OPTS: "-Xms512m -Xmx512m"` to prevent OOM kills.
> * **Linux Users:** You must run `sudo sysctl -w vm.max_map_count=262144` on your host machine before starting, or ES will fail to bootstrap.

```bash
docker-compose -f docker-compose.base.yml \
               -f docker-compose.nacos.yml \
               -f docker-compose.consul.yml \
               -f docker-compose.es.yml up -d
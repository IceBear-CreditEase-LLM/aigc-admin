# AIGC 管理平台

AIGC平台是一个综合了模型管理、模型部署、模型微调、渠道管理等功能的平台，通过该平台可以快速的部署模型、微调模型、管理模型、管理渠道等功能。

## 项目简介

前端UI是一个独立的项目，点击[aigc-admin-web](https://github.com/IceBear-CreditEase-LLM/aigc-admin-web)查看

### 系统架构设计

系统是前后端分离的架构。

#### 模型推理框架

我们使用的是[FastChat](https://github.com/lm-sys/FastChat)作为模型推理框架，FastChat是一个非常优秀的开源项目。

> [FastChat](https://github.com/lm-sys/FastChat) 是一个开放平台，用于训练、服务和评估基于大型语言模型的聊天机器人。

**FastChat我们主要用其三个服务**

`controller` 用于模型的注册中心及健康检查

`worker` 服务启动模型并将当前模型注册到controller

`api` 从controller获取模型的地址代理到worker并提供标准API

我们主要通过它来实现大模型的高可用，高可扩展性。

![img.png](https://github.com/lm-sys/FastChat/raw/main/assets/server_arch.png)

模型部署的操作可以参考[模型部署](docs/model/list.md)

### 模型微调

为了实现模型的微调，您可以参考我们的详细指南：[模型微调](docs/model/finetune.md)。

### 模型部署与微调

您可以将模型部署到任意配备GPU的节点上，无论是私有的K8s集群、Docker集群，还是云服务商提供的K8s集群，均能轻松对接。

### 本系统组成

本系统主要由以下几个部分组成：

- **HTTP服务**：提供Web服务接口，方便用户进行交互。
- **定时任务**：执行预定任务，如模型训练、数据预处理等。
- **训练镜像**：包含所有必要的环境和依赖，用于模型的训练和微调。

- 通过这些组件的协同工作，我们能够提供一个灵活、高效的模型微调和部署解决方案。

#### 部署流程

```mermaid
graph LR
    A[aigc] --> B[点击部署]
    B --> C[创建部署模版]
    C --> D[使用Docker或k8s进行调度]
    D --> E[挂载相应配置有模型]
    E --> F[启动模型]
    F --> G[注册到fschat-controller]
```

#### 微调训练流程

```mermaid
graph LR
    A[aigc] --> B[上传微调文件]
    B --> C[生成微调模版]
    C --> D[使用Docker或k8s进行调度]
    D --> E[挂载相应配置有模型]
    E --> F[启动训练脚本]
    F --> G[输出日志]
```

## 使用手册

[AIGC平台使用手册](docs/SUMMARY.md)

### 安装使用步骤

- 克隆项目: `git clone https://github.com/IceBear-CreditEase-LLM/aigc-admin.git`
- 进入项目: `cd aigc-admin`

该系统依赖**Mysql**、**Redis**和**Docker**需要安装此服务

推理或训练节点只需要安装**Docker**和**Nvidia-Docker**
即可。[NVIDIA Container Toolkit](https://github.com/NVIDIA/nvidia-container-toolkit)

#### 本地开发

[golang](https://github.com/golang/go)版本请安装go1.21以上版本

- 本地启动: `make run`
- build成x86 Linux可执行文件: `make build-linux`
- build成当前电脑可执行文件: `make build`

build完通常会保存在 `$(GOPATH)/bin/` 目录下

#### Docker部署

安装docker和docker-compose可以参考官网教程：[Install Docker Engine](https://docs.docker.com/engine/install/)

执行命令启动全部服务

```
$ docker-compose up
```

### 项目配置

项目配置可以通过两种方式进行配置

#### 通过命令先传参

**需要注意的是，如果即设置了环境变量也设置了命令行参数，那么命令行参数的值会覆盖环境变量的值**

执行: `./aigc-admin start --help` 查看命令行参数

```bash
# Aigc Admin服务
有关本系统的相关概述，请参阅 http://github.com/IceBear-CreditEase-LLM/aigc-admin

Usage:
  aigc-admin [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  cronjob     定时任务
  generate    生成命令
  help        Help about any command
  job         任务命令
  start       启动服务

Flags:
  -c, --config.path string                配置文件路径，如果没有传入配置文件路径则默认使用环境变量
      --db.drive string                   数据库驱动 (default "mysql")
      --db.mysql.database string          mysql数据库 (default "aigc")
      --db.mysql.host string              mysql数据库地址: mysql (default "localhost")
      --db.mysql.metrics                  是否启GORM的Metrics
      --db.mysql.password string          mysql数据库密码
      --db.mysql.port int                 mysql数据库端口 (default 3306)
      --db.mysql.user string              mysql数据库用户 (default "aigc")
  -h, --help                              help for aigc-admin
      --ldap.base.dn string               LDAP Base DN (default "OU=HABROOT,DC=ORG,DC=corp")
      --ldap.bind.pass string             LDAP Bind Password
      --ldap.bind.user string             LDAP Bind User (default "aigc_ldap")
      --ldap.group.filter string          LDAP Group Filter
      --ldap.host string                  LDAP地址 (default "ldap://ldap")
      --ldap.port int                     LDAP端口 (default 389)
      --ldap.use.ssl                      LDAP Base DN
      --ldap.user.attr strings            LDAP Attributes (default [name,mail,userPrincipalName,displayName,sAMAccountName])
      --ldap.user.filter string           LDAP User Filter (default "(userPrincipalName=%s)")
  -n, --namespace string                  命名空间 (default "aigc")
      --redis.auth string                 连接Redis密码
      --redis.db int                      连接Redis DB
      --redis.hosts string                连接Redis地址 (default "redis:6379")
      --redis.prefix string               Redis写入Cache的前缀 (default "aigc")
      --server.admin.pass string          系统管理员密码 (default "admin")
      --server.admin.user string          系统管理员账号 (default "admin")
      --server.debug                      是否开启Debug模式
      --server.key string                 本系统服务密钥 (default "Aigcfj@202401")
      --server.log.drive string           本系统日志驱动, 支持syslog,term (default "term")
      --server.log.level string           本系统日志级别 (default "all")
      --server.log.name string            本系统日志名称 (default "aigc-admin.log")
      --server.log.path string            本系统日志路径
  -a, --server.name string                本系统服务名称 (default "aigc-admin")
      --service.alarm.token string        告警中心服务地址 (default "http://alarm:8080")
      --service.chat.host string          ChatApi 地址 (default "http://chat-api:8080")
      --service.chat.token string         ChatApi Token
      --service.gpt.host string           Chat-Api 地址 (default "http://chat-api:8080/v1")
      --service.openai.enable             是否启用OpenAI服务
      --service.openai.host string        OpenAI服务地址 (default "https://api.openai.com/v1")
      --service.openai.model string       OpenAI模型名称 (default "gpt-3.5-turbo")
      --service.openai.org.id string      OpenAI OrgId
      --service.s3.access.key string      S3 AccessKey
      --service.s3.bucket string          S3 Bucket (default "aigc")
      --service.s3.bucket.public string   S3 Bucket Public (default "aigc")
      --service.s3.cluster string         S3 集群 (default "ceph-c2")
      --service.s3.host string            S3服务地址 (default "http://s3")
      --service.s3.project.name string    S3 项目名称 (default "aigc")
      --service.s3.region string          S3 Bucket (default "default")
      --service.s3.secret.key string      S3 SecretKey
```

#### 系统公共环境变量配置

可以修改`.env`调整相关配置

##### 数据库配置

| 环境变量                  | 值           | 描述               |
|-----------------------|-------------|------------------|
| `AIGC_DB_DRIVER`      | `mysql`     | 数据库驱动类型（可能是遗留错误） |
| `AIGC_MYSQL_DRIVE`    | `mysql`     | 数据库驱动类型          |
| `AIGC_MYSQL_HOST`     | `localhost` | 数据库主机地址          |
| `AIGC_MYSQL_PORT`     | `3306`      | 数据库端口号           |
| `AIGC_MYSQL_USER`     | `aigc`      | 数据库用户名           |
| `AIGC_MYSQL_PASSWORD` | `admin`     | 数据库密码            |
| `AIGC_MYSQL_DATABASE` | `aigc`      | 数据库名             |

##### Redis 配置

| 环境变量                  | 值            | 描述                   |
|-----------------------|--------------|----------------------|
| `AIGC_REDIS_HOSTS`    | `redis:6379` | Redis 服务地址和端口        |
| `AIGC_REDIS_PREFIX`   | `aigc`       | Redis 前缀，用于区分不同的数据集合 |
| `AIGC_REDIS_PASSWORD` |              | Redis 访问密码           |

##### Tracer 链路追踪配置

| 环境变量                           | 值        | 描述          |
|--------------------------------|----------|-------------|
| `AIGC_TRACER_ENABLE`           | `false`  | 是否启用链路追踪    |
| `AIGC_TRACER_DRIVE`            | `jaeger` | 链路追踪驱动类型    |
| `AIGC_TRACER_JAEGER_HOST`      |          | Jaeger 服务地址 |
| `AIGC_TRACER_JAEGER_PARAM`     | `1`      | Jaeger 采样参数 |
| `AIGC_TRACER_JAEGER_TYPE`      | `const`  | Jaeger 采样类型 |
| `AIGC_TRACER_JAEGER_LOG_SPANS` | `false`  | 是否记录追踪日志    |

##### 跨域配置

| 环境变量               | 值       | 描述               |
|--------------------|---------|------------------|
| `AIGC_ENABLE_CORS` | `false` | 是否启用CORS（跨源资源共享） |

##### 外部服务调用配置

| 环境变量                         | 值        | 描述          |
|------------------------------|----------|-------------|
| `AIGC_SERVICE_ALARM_HOST`    |          | 报警服务地址      |
| `AIGC_SERVICE_CHAT_API_HOST` |          | 聊天API服务地址   |
| `AIGC_SERVICE_OPENAI_TOKEN`  | `sk-***` | API Key     |
| `AIGC_SERVICE_OPENAI_ORG_ID` |          | OpenAI 组织ID |

##### S3 存储配置

| 环境变量                            | 值 | 描述         |
|---------------------------------|---|------------|
| `AIGC_SERVICE_S3_HOST`          |   | S3 服务地址    |
| `AIGC_SERVICE_S3_ACCESS_KEY`    |   | S3 访问密钥    |
| `AIGC_SERVICE_S3_SECRET_KEY`    |   | S3 访问密钥密码  |
| `AIGC_SERVICE_S3_BUCKET`        |   | S3 存储桶名称   |
| `AIGC_SERVICE_S3_BUCKET_PUBLIC` |   | S3 公共存储桶名称 |
| `AIGC_SERVICE_S3_PROJECT_NAME`  |   | S3 项目名称    |

##### 聊天API配置

| 环境变量                      | 值                      | 描述       |
|---------------------------|------------------------|----------|
| `AIGC_SERVICE_CHAT_HOST`  | `http://chat-api:8080` | 聊天服务地址   |
| `AIGC_SERVICE_CHAT_TOKEN` | `sk-001`               | 聊天服务访问令牌 |

##### LDAP 配置

| 环境变量                  | 值                    | 描述          |
|-----------------------|----------------------|-------------|
| `AIGC_LDAP_HOST`      | `ldap`               | LDAP 服务器地址  |
| `AIGC_LDAP_BASE_DN`   | `OU=HABROOT,DC=corp` | LDAP 基础DN   |
| `AIGC_LDAP_BIND_USER` |                      | LDAP 绑定用户   |
| `AIGC_LDAP_BIND_PASS` |                      | LDAP 绑定用户密码 |
| `AIGC_LDAP_USER_ATTR` | `mail,displayName`   | LDAP 用户属性   |

##### aigc-admin 环境变量配置

| 环境变量                                    | 值                | 描述       |
|-----------------------------------------|------------------|----------|
| `AIGC_ADMIN_SERVER_HTTP_PORT`           | `:8080`          | 服务HTTP端口 |
| `AIGC_ADMIN_SERVER_LOG_DRIVE`           | `term`           | 日志驱动类型   |
| `AIGC_ADMIN_SERVER_NAME`                | `aigc-admin`     | 服务名称     |
| `AIGC_ADMIN_SERVER_DEBUG`               | `true`           | 是否开启调试模式 |
| `AIGC_ADMIN_SERVER_LOG_LEVEL`           | `all`            | 日志级别     |
| `AIGC_ADMIN_SERVER_LOG_PATH`            |                  | 日志路径     |
| `AIGC_ADMIN_SERVER_LOG_NAME`            | `aigc-admin.log` | 日志文件名称   |
| `AIGC_ADMIN_SERVER_DEFAULT_CHANNEL_KEY` | `sk-001`         | 默认渠道密钥   |

## Docker镜像

我们提供了Docker镜像，您可以直接使用我们提供的镜像，也可以自行构建。

- [LLMOps](docker/llmops-deepspeed/README.md)
- [百川2](docker/baichuan2/README.md)
- [FastChat](docker/fastchat/README.md)
- [Qwen](docker/qwen/README.md)
- [Vicuna](docker/vicuna/README.md)

### 文件资源目录

```
.
├── CHANGELOG                   # 变更日志
├── Dockerfile                  # Dockerfile 构建
├── Makefile                    # Makefile 构建
├── README.md                   # 项目说明
├── aigc-admin.service          # systemd 服务配置
├── cmd                         # 项目入口
│ ├── data                      # 数据目录
│ ├── main.go
│ ├── service                   # 服务启动从这里开始
│ └── web
├── docker                      # docker 镜像构建
├── docs                        # 文档
├── go.mod                      # go mod 依赖
├── go.sum
├── src                         # 项目源码
│ ├── api                       # 调用外部API模块
│ ├── encode                    # 输出编码模块
│ ├── logging                   # 日志处理模块
│ ├── middleware                # 中间件模块
│ ├── pkg                       # 项目模块目录
│ │ ├── assistants              # 助手模块
│ │ ├── auth                    # 认证模块
│ │ ├── channels                # 渠道模块
│ │ ├── datasets                # 数据集模块
│ │ ├── files                   # 文件模块
│ │ ├── finetuning              # 微调模块
│ │ ├── models                  # 模型模块
│ │ ├── sys                     # 系统模块
│ │ └── tools                   # 工具模块
│ ├── repository                # 数据库操作模块
│ └── util                      # 工具模块
└── tests                       # 测试模块
```
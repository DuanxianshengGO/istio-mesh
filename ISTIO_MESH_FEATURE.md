# Istio Service Mesh 功能说明

## 概述

KubePi 现已集成 Istio Service Mesh 功能，提供对 Istio 资源的完整管理能力，包括 VirtualService、DestinationRule、Gateway 等核心资源的 CRUD 操作，以及流量分析功能。

## 功能特性

### 1. VirtualService 管理
- **列表展示**: 查看所有 VirtualService 资源
- **创建**: 支持表单创建和 YAML 编辑两种方式
- **编辑**: 支持 YAML 编辑模式
- **删除**: 支持单个和批量删除
- **详情查看**: 查看 VirtualService 详细配置

**主要配置项**:
- Hosts: 目标主机列表
- Gateways: 关联的 Gateway
- HTTP Routes: HTTP 路由规则
- Match Conditions: 匹配条件
- Destinations: 目标服务配置

### 2. DestinationRule 管理
- **列表展示**: 查看所有 DestinationRule 资源
- **创建**: 支持表单创建，包含 subset 配置
- **编辑**: 支持 YAML 编辑模式
- **删除**: 支持单个和批量删除
- **详情查看**: 查看 DestinationRule 详细配置

**主要配置项**:
- Host: 目标服务主机
- Load Balancer: 负载均衡策略
- Subsets: 服务子集定义
- Labels: 子集标签选择器

### 3. Gateway 管理
- **列表展示**: 查看所有 Gateway 资源
- **创建**: 支持表单创建，包含 TLS 配置
- **编辑**: 支持 YAML 编辑模式
- **删除**: 支持单个和批量删除
- **详情查看**: 查看 Gateway 详细配置

**主要配置项**:
- Selector: Gateway 选择器
- Servers: 服务器配置
- Ports: 端口和协议配置
- Hosts: 主机列表
- TLS: TLS 配置（HTTPS/TLS 协议）

### 4. Traffic Analytics 流量分析
- **流量概览**: 显示 VirtualService、DestinationRule、Gateway 统计信息
- **服务过滤**: 按命名空间和服务过滤
- **流量路由**: 展示当前流量路由规则
- **Pod 流量**: 显示 Pod 级别的流量信息
- **可视化图表**: 流量流向图表（开发中）

## 使用说明

### 前置条件
1. Kubernetes 集群已安装 Istio
2. KubePi 用户具有相应的 RBAC 权限

### 权限要求
用户需要具有以下 API 组的权限：
- `networking.istio.io/v1beta1` - VirtualService, DestinationRule, Gateway

### 访问路径
在 KubePi 主界面中，选择对应集群后，在左侧菜单中找到 "Service Mesh" 选项。

### 典型使用场景

#### 场景1: 配置流量路由
1. 创建 DestinationRule 定义服务子集
2. 创建 VirtualService 配置路由规则
3. 在 Traffic Analytics 中查看流量分布

#### 场景2: 配置入口网关
1. 创建 Gateway 定义入口配置
2. 创建 VirtualService 绑定 Gateway
3. 配置 TLS 证书（如需要）

## API 接口

### 后端 API 路径
```
/api/v1/istio/{cluster}/virtualservices
/api/v1/istio/{cluster}/destinationrules  
/api/v1/istio/{cluster}/gateways
/api/v1/istio/{cluster}/traffic-analytics
```

### 支持的操作
- `GET` - 列表查询
- `POST` - 创建资源
- `PUT` - 更新资源
- `DELETE` - 删除资源

## 文件结构

### 后端文件
```
internal/api/v1/istio/
├── istio.go                 # 主要 API 处理逻辑
```

### 前端文件
```
web/dashboard/src/
├── api/istio.js            # API 客户端函数
├── router/modules/istio.js # 路由配置
├── business/istio/         # 业务组件
│   ├── virtualservice/     # VirtualService 管理
│   ├── destinationrule/    # DestinationRule 管理
│   ├── gateway/           # Gateway 管理
│   └── traffic-analytics/ # 流量分析
└── components/
    └── namespace-select/   # 命名空间选择组件
```

## 国际化支持

支持中文和英文界面，相关翻译文件：
- `web/dashboard/src/i18n/lang/zh-CN.js`
- `web/dashboard/src/i18n/lang/en-US.js`

## 测试

运行测试脚本验证 API 功能：
```bash
chmod +x test_istio_api.sh
./test_istio_api.sh
```

## 注意事项

1. **权限控制**: 确保用户具有操作 Istio 资源的权限
2. **命名空间**: 资源操作需要指定正确的命名空间
3. **YAML 格式**: 编辑时注意 YAML 格式的正确性
4. **依赖关系**: 删除资源时注意依赖关系，避免影响正在运行的服务

## 故障排除

### 常见问题
1. **权限不足**: 检查用户是否具有 `networking.istio.io` API 组权限
2. **Istio 未安装**: 确认集群中已正确安装 Istio
3. **网络连接**: 检查 KubePi 与 Kubernetes API Server 的连接

### 日志查看
查看 KubePi 后端日志获取详细错误信息：
```bash
kubectl logs -f deployment/kubepi-server
```

## 后续计划

1. **增强流量分析**: 集成 Prometheus 指标，提供实时流量监控
2. **可视化改进**: 完善流量拓扑图和服务依赖关系图
3. **模板支持**: 提供常用配置模板
4. **批量操作**: 支持批量导入/导出配置

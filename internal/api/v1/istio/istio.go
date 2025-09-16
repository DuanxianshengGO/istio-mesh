package istio

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"

	"github.com/KubeOperator/kubepi/internal/api/v1/session"
	v1Cluster "github.com/KubeOperator/kubepi/internal/model/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"k8s.io/client-go/rest"
)

type Handler struct {
	clusterService        cluster.Service
	clusterBindingService clusterbinding.Service
}

func NewHandler() *Handler {
	return &Handler{
		clusterService:        cluster.NewService(),
		clusterBindingService: clusterbinding.NewService(),
	}
}

// ListVirtualServices 获取 VirtualService 列表
func (h *Handler) ListVirtualServices() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		// 优先从路径参数获取namespace，如果没有则从查询参数获取
		namespace := ctx.Params().GetString("namespace")
		if namespace == "" {
			namespace = ctx.URLParam("namespace")
		}

		// 构建 API 路径
		apiPath := "/apis/networking.istio.io/v1beta1"
		if namespace != "" {
			apiPath = fmt.Sprintf("%s/namespaces/%s/virtualservices", apiPath, namespace)
		} else {
			apiPath = fmt.Sprintf("%s/virtualservices", apiPath)
		}
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// GetVirtualService 获取单个 VirtualService
func (h *Handler) GetVirtualService() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/virtualservices/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// CreateVirtualService 创建 VirtualService
func (h *Handler) CreateVirtualService() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/virtualservices", namespace)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// UpdateVirtualService 更新 VirtualService
func (h *Handler) UpdateVirtualService() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/virtualservices/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// DeleteVirtualService 删除 VirtualService
func (h *Handler) DeleteVirtualService() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/virtualservices/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// ListDestinationRules 获取 DestinationRule 列表
func (h *Handler) ListDestinationRules() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		// 优先从路径参数获取namespace，如果没有则从查询参数获取
		namespace := ctx.Params().GetString("namespace")
		if namespace == "" {
			namespace = ctx.URLParam("namespace")
		}

		apiPath := "/apis/networking.istio.io/v1beta1"
		if namespace != "" {
			apiPath = fmt.Sprintf("%s/namespaces/%s/destinationrules", apiPath, namespace)
		} else {
			apiPath = fmt.Sprintf("%s/destinationrules", apiPath)
		}

		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// GetDestinationRule 获取单个 DestinationRule
func (h *Handler) GetDestinationRule() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/destinationrules/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// CreateDestinationRule 创建 DestinationRule
func (h *Handler) CreateDestinationRule() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/destinationrules", namespace)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// UpdateDestinationRule 更新 DestinationRule
func (h *Handler) UpdateDestinationRule() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/destinationrules/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// DeleteDestinationRule 删除 DestinationRule
func (h *Handler) DeleteDestinationRule() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/destinationrules/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// ListGateways 获取 Gateway 列表
func (h *Handler) ListGateways() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		// 优先从路径参数获取namespace，如果没有则从查询参数获取
		namespace := ctx.Params().GetString("namespace")
		if namespace == "" {
			namespace = ctx.URLParam("namespace")
		}

		apiPath := "/apis/networking.istio.io/v1beta1"
		if namespace != "" {
			apiPath = fmt.Sprintf("%s/namespaces/%s/gateways", apiPath, namespace)
		} else {
			apiPath = fmt.Sprintf("%s/gateways", apiPath)
		}

		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// GetGateway 获取单个 Gateway
func (h *Handler) GetGateway() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/gateways/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// CreateGateway 创建 Gateway
func (h *Handler) CreateGateway() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/gateways", namespace)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// UpdateGateway 更新 Gateway
func (h *Handler) UpdateGateway() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/gateways/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// DeleteGateway 删除 Gateway
func (h *Handler) DeleteGateway() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.Params().GetString("namespace")
		name := ctx.Params().GetString("name")

		apiPath := fmt.Sprintf("/apis/networking.istio.io/v1beta1/namespaces/%s/gateways/%s", namespace, name)
		h.proxyToKubernetes(ctx, clusterName, apiPath)
	}
}

// GetTrafficAnalytics 获取流量分析数据
func (h *Handler) GetTrafficAnalytics() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("cluster")
		namespace := ctx.URLParam("namespace")

		// 这里需要实现流量分析逻辑
		// 1. 获取所有 VirtualService
		// 2. 获取所有 DestinationRule
		// 3. 获取相关的 Pod 信息
		// 4. 分析流量路由关系

		analytics := h.analyzeTraffic(clusterName, namespace)
		ctx.JSON(analytics)
	}
}

// proxyToKubernetes 代理请求到 Kubernetes API
func (h *Handler) proxyToKubernetes(ctx *context.Context, clusterName, apiPath string) {
	// 获取集群信息
	c, err := h.clusterService.Get(clusterName, common.DBOptions{})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
		return
	}

	// 获取用户会话信息
	u := ctx.Values().Get("profile")
	profile := u.(session.UserProfile)

	// 生成传输层
	ts, err := h.generateTLSTransport(c, profile)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("message", err.Error())
		return
	}

	// 创建 HTTP 客户端
	httpClient := http.Client{Transport: ts}

	// 构建完整的 API URL
	fullURL := fmt.Sprintf("%s%s", c.Spec.Connect.Forward.ApiServer, apiPath)

	// 创建 HTTP 请求
	req, err := http.NewRequest(ctx.Request().Method, fullURL, ctx.Request().Body)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("message", err.Error())
		return
	}

	// 设置 Content-Type
	if ctx.Method() == "POST" || ctx.Method() == "PUT" || ctx.Method() == "PATCH" {
		req.Header.Set("Content-Type", "application/json")
	}
	if ctx.Method() == "PATCH" {
		req.Header.Set("Content-Type", "application/merge-patch+json")
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("message", err.Error())
		return
	}
	defer resp.Body.Close()

	// 复制响应头（排除Content-Length，因为我们会修改响应体）
	for key, values := range resp.Header {
		if key != "Content-Length" {
			for _, value := range values {
				ctx.Header(key, value)
			}
		}
	}

	// 设置状态码
	ctx.StatusCode(resp.StatusCode)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Values().Set("message", err.Error())
		return
	}

	// 如果是成功响应，包装成 KubePi 标准格式
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// 解析 Kubernetes API 响应
		var k8sResponse interface{}
		if err := json.Unmarshal(body, &k8sResponse); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		// 包装成 KubePi 标准格式
		response := map[string]interface{}{
			"data":    k8sResponse,
			"success": true,
		}

		ctx.JSON(response)
	} else {
		// 错误响应直接返回
		ctx.Write(body)
	}
}

// generateTLSTransport 生成 TLS 传输层，与其他 API 保持一致
func (h *Handler) generateTLSTransport(c *v1Cluster.Cluster, profile session.UserProfile) (http.RoundTripper, error) {
	if profile.IsAdministrator {
		k := kubernetes.NewKubernetes(c)
		adminConfig, err := k.Config()
		if err != nil {
			return nil, err
		}
		return rest.TransportFor(adminConfig)
	}

	binding, err := h.clusterBindingService.GetBindingByClusterNameAndUserName(c.Name, profile.Name, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	kubeConf := &rest.Config{
		Host: c.Spec.Connect.Forward.ApiServer,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
			CertData: binding.Certificate,
			KeyData:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: c.PrivateKey}),
		},
	}
	return rest.TransportFor(kubeConf)
}

// analyzeTraffic 分析流量路由关系
func (h *Handler) analyzeTraffic(clusterName, namespace string) interface{} {
	// TODO: 实现流量分析逻辑
	// 这里返回一个示例结构
	return map[string]interface{}{
		"virtualServices":  []interface{}{},
		"destinationRules": []interface{}{},
		"gateways":         []interface{}{},
		"trafficRoutes":    []interface{}{},
		"podTraffic":       []interface{}{},
	}
}

// Install 安装路由
func Install(parent iris.Party) {
	handler := NewHandler()

	// Istio 资源管理路由
	istioParty := parent.Party("/istio/:cluster")

	// VirtualService 路由
	istioParty.Get("/virtualservices", handler.ListVirtualServices())
	istioParty.Get("/namespaces/:namespace/virtualservices", handler.ListVirtualServices())
	istioParty.Get("/namespaces/:namespace/virtualservices/:name", handler.GetVirtualService())
	istioParty.Post("/namespaces/:namespace/virtualservices", handler.CreateVirtualService())
	istioParty.Put("/namespaces/:namespace/virtualservices/:name", handler.UpdateVirtualService())
	istioParty.Delete("/namespaces/:namespace/virtualservices/:name", handler.DeleteVirtualService())

	// DestinationRule 路由
	istioParty.Get("/destinationrules", handler.ListDestinationRules())
	istioParty.Get("/namespaces/:namespace/destinationrules", handler.ListDestinationRules())
	istioParty.Get("/namespaces/:namespace/destinationrules/:name", handler.GetDestinationRule())
	istioParty.Post("/namespaces/:namespace/destinationrules", handler.CreateDestinationRule())
	istioParty.Put("/namespaces/:namespace/destinationrules/:name", handler.UpdateDestinationRule())
	istioParty.Delete("/namespaces/:namespace/destinationrules/:name", handler.DeleteDestinationRule())

	// Gateway 路由
	istioParty.Get("/gateways", handler.ListGateways())
	istioParty.Get("/namespaces/:namespace/gateways", handler.ListGateways())
	istioParty.Get("/namespaces/:namespace/gateways/:name", handler.GetGateway())
	istioParty.Post("/namespaces/:namespace/gateways", handler.CreateGateway())
	istioParty.Put("/namespaces/:namespace/gateways/:name", handler.UpdateGateway())
	istioParty.Delete("/namespaces/:namespace/gateways/:name", handler.DeleteGateway())

	// Traffic Analytics 路由
	istioParty.Get("/traffic-analytics", handler.GetTrafficAnalytics())
}

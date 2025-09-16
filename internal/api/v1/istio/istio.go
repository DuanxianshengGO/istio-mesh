package istio

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"

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
		// 错误响应，先记录原始错误用于调试
		fmt.Printf("Kubernetes API Error: %s\n", string(body))

		// 解析Kubernetes错误并包装成KubePi格式
		var k8sError map[string]interface{}
		if err := json.Unmarshal(body, &k8sError); err != nil {
			// 如果无法解析为JSON，直接返回原始错误
			response := map[string]interface{}{
				"code":    resp.StatusCode,
				"message": string(body),
				"success": false,
			}
			ctx.JSON(response)
			return
		}

		// 提取Kubernetes错误信息
		var message string
		if msg, ok := k8sError["message"].(string); ok {
			message = msg
		} else if reason, ok := k8sError["reason"].(string); ok {
			message = reason
		} else {
			message = "Unknown error"
		}

		// 包装成KubePi标准错误格式
		response := map[string]interface{}{
			"code":    resp.StatusCode,
			"message": message,
			"success": false,
		}
		ctx.JSON(response)
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
	// 获取集群信息
	c, err := h.clusterService.Get(clusterName, common.DBOptions{})
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("get cluster failed: %s", err.Error()),
		}
	}

	// 生成 TLS 传输层
	profile := session.UserProfile{Name: "admin", IsAdministrator: true}
	ts, err := h.generateTLSTransport(c, profile)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("generate TLS transport failed: %s", err.Error()),
		}
	}

	httpClient := http.Client{Transport: ts}
	baseURL := c.Spec.Connect.Forward.ApiServer

	// 获取所有 VirtualService
	virtualServices, err := h.fetchVirtualServices(httpClient, baseURL, namespace)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("fetch VirtualServices failed: %s", err.Error()),
		}
	}

	// 获取所有 DestinationRule
	destinationRules, err := h.fetchDestinationRules(httpClient, baseURL, namespace)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("fetch DestinationRules failed: %s", err.Error()),
		}
	}

	// 获取所有 Pod
	pods, err := h.fetchPods(httpClient, baseURL, namespace)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("fetch Pods failed: %s", err.Error()),
		}
	}

	// 分析流量关系
	trafficAnalysis := h.analyzeTrafficFlow(virtualServices, destinationRules, pods)

	return trafficAnalysis
}

// fetchVirtualServices 获取所有 VirtualService
func (h *Handler) fetchVirtualServices(httpClient http.Client, baseURL, namespace string) ([]map[string]interface{}, error) {
	apiPath := "/apis/networking.istio.io/v1beta1"
	if namespace != "" {
		apiPath = fmt.Sprintf("%s/namespaces/%s/virtualservices", apiPath, namespace)
	} else {
		apiPath = fmt.Sprintf("%s/virtualservices", apiPath)
	}

	fullURL := fmt.Sprintf("%s%s", baseURL, apiPath)
	resp, err := httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	var virtualServices []map[string]interface{}
	for _, item := range items {
		if vs, ok := item.(map[string]interface{}); ok {
			virtualServices = append(virtualServices, vs)
		}
	}

	return virtualServices, nil
}

// fetchDestinationRules 获取所有 DestinationRule
func (h *Handler) fetchDestinationRules(httpClient http.Client, baseURL, namespace string) ([]map[string]interface{}, error) {
	apiPath := "/apis/networking.istio.io/v1beta1"
	if namespace != "" {
		apiPath = fmt.Sprintf("%s/namespaces/%s/destinationrules", apiPath, namespace)
	} else {
		apiPath = fmt.Sprintf("%s/destinationrules", apiPath)
	}

	fullURL := fmt.Sprintf("%s%s", baseURL, apiPath)
	resp, err := httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	var destinationRules []map[string]interface{}
	for _, item := range items {
		if dr, ok := item.(map[string]interface{}); ok {
			destinationRules = append(destinationRules, dr)
		}
	}

	return destinationRules, nil
}

// fetchPods 获取所有 Pod
func (h *Handler) fetchPods(httpClient http.Client, baseURL, namespace string) ([]map[string]interface{}, error) {
	apiPath := "/api/v1"
	if namespace != "" {
		apiPath = fmt.Sprintf("%s/namespaces/%s/pods", apiPath, namespace)
	} else {
		apiPath = fmt.Sprintf("%s/pods", apiPath)
	}

	fullURL := fmt.Sprintf("%s%s", baseURL, apiPath)
	resp, err := httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	var pods []map[string]interface{}
	for _, item := range items {
		if pod, ok := item.(map[string]interface{}); ok {
			pods = append(pods, pod)
		}
	}

	return pods, nil
}

// analyzeTrafficFlow 分析流量流向
func (h *Handler) analyzeTrafficFlow(virtualServices, destinationRules, pods []map[string]interface{}) map[string]interface{} {
	// 构建 DestinationRule 映射，key 为 host
	drMap := make(map[string]map[string]interface{})
	for _, dr := range destinationRules {
		if spec, ok := dr["spec"].(map[string]interface{}); ok {
			if host, ok := spec["host"].(string); ok {
				drMap[host] = dr
			}
		}
	}

	var trafficAnalysis []map[string]interface{}

	// 为每个Pod分析流量类型
	for _, pod := range pods {
		podName := ""
		podLabels := make(map[string]string)
		serviceName := ""

		// 获取Pod基本信息
		if metadata, ok := pod["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				podName = name
			}
			if labels, ok := metadata["labels"].(map[string]interface{}); ok {
				for k, v := range labels {
					if vStr, ok := v.(string); ok {
						podLabels[k] = vStr
					}
				}
				// 通过app标签确定服务名
				if app, ok := podLabels["app"]; ok {
					serviceName = app
				}
			}
		}

		// 分析这个Pod的流量类型
		trafficResults := h.analyzePodTrafficWithSubsets(podName, serviceName, podLabels, virtualServices, drMap)

		// 为每个匹配的subset创建一条记录
		for _, result := range trafficResults {
			trafficAnalysis = append(trafficAnalysis, map[string]interface{}{
				"podName":      podName,
				"serviceName":  serviceName,
				"vsName":       result["vsName"],
				"trafficType":  result["type"],
				"subset":       result["subset"],
				"matchContent": result["matchContent"],
			})
		}
	}

	return map[string]interface{}{
		"trafficAnalysis": trafficAnalysis,
		"summary": map[string]interface{}{
			"totalPods":    len(pods),
			"totalVS":      len(virtualServices),
			"totalDR":      len(destinationRules),
			"basicTraffic": h.countTrafficType(trafficAnalysis, "基础流量"),
			"grayTraffic":  h.countTrafficType(trafficAnalysis, "灰度流量"),
			"noTraffic":    h.countTrafficType(trafficAnalysis, "未进行灰度"),
		},
	}
}

// countTrafficType 统计指定类型的流量数量
func (h *Handler) countTrafficType(trafficAnalysis []map[string]interface{}, trafficType string) int {
	count := 0
	for _, analysis := range trafficAnalysis {
		if t, ok := analysis["trafficType"].(string); ok && t == trafficType {
			count++
		}
	}
	return count
}

// analyzePodTraffic 分析单个Pod的流量类型
func (h *Handler) analyzePodTraffic(podName, serviceName string, podLabels map[string]string, virtualServices []map[string]interface{}, drMap map[string]map[string]interface{}) map[string]interface{} {
	// 默认返回未进行灰度
	result := map[string]interface{}{
		"type":   "未进行灰度",
		"vsName": "",
		"subset": nil,
	}

	// 遍历所有VirtualService，查找匹配的路由规则
	for _, vs := range virtualServices {
		vsName := ""
		if metadata, ok := vs["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				vsName = name
			}
		}

		if spec, ok := vs["spec"].(map[string]interface{}); ok {
			// 获取hosts
			hosts := []string{}
			if hostsInterface, ok := spec["hosts"].([]interface{}); ok {
				for _, host := range hostsInterface {
					if hostStr, ok := host.(string); ok {
						hosts = append(hosts, hostStr)
					}
				}
			}

			// 检查这个VS是否适用于当前服务
			serviceMatched := false
			for _, host := range hosts {
				// 简化匹配：host可能是服务名或FQDN
				if host == serviceName || strings.Contains(host, serviceName) {
					serviceMatched = true
					break
				}
			}

			if !serviceMatched {
				continue
			}

			// 分析HTTP路由规则
			if httpRoutes, ok := spec["http"].([]interface{}); ok {
				for _, httpRoute := range httpRoutes {
					if route, ok := httpRoute.(map[string]interface{}); ok {
						// 检查是否有match条件（灰度流量）
						hasMatch := false
						if _, exists := route["match"]; exists {
							hasMatch = true
						}

						// 分析route目标
						if routeDestinations, ok := route["route"].([]interface{}); ok {
							for _, dest := range routeDestinations {
								if destination, ok := dest.(map[string]interface{}); ok {
									if destInfo, ok := destination["destination"].(map[string]interface{}); ok {
										if host, ok := destInfo["host"].(string); ok {
											// 检查subset是否匹配当前Pod
											if subset, ok := destInfo["subset"].(string); ok {
												if h.podMatchesSubset(podLabels, host, subset, drMap) {
													trafficType := "基础流量"
													if hasMatch {
														trafficType = "灰度流量"
													}
													return map[string]interface{}{
														"type":   trafficType,
														"vsName": vsName,
														"subset": subset,
													}
												}
											} else {
												// 没有指定subset，检查是否直接匹配服务
												if host == serviceName || strings.Contains(host, serviceName) {
													trafficType := "基础流量"
													if hasMatch {
														trafficType = "灰度流量"
													}
													return map[string]interface{}{
														"type":   trafficType,
														"vsName": vsName,
														"subset": nil,
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

// podMatchesSubset 检查Pod是否匹配DestinationRule中定义的subset
func (h *Handler) podMatchesSubset(podLabels map[string]string, host, subset string, drMap map[string]map[string]interface{}) bool {
	// 查找对应的DestinationRule
	dr, exists := drMap[host]
	if !exists {
		return false
	}

	if spec, ok := dr["spec"].(map[string]interface{}); ok {
		if subsets, ok := spec["subsets"].([]interface{}); ok {
			for _, sub := range subsets {
				if subsetMap, ok := sub.(map[string]interface{}); ok {
					if name, ok := subsetMap["name"].(string); ok && name == subset {
						// 找到了对应的subset，检查labels是否匹配
						if labels, ok := subsetMap["labels"].(map[string]interface{}); ok {
							allMatch := true
							for labelKey, labelValue := range labels {
								if labelValueStr, ok := labelValue.(string); ok {
									if podLabels[labelKey] != labelValueStr {
										allMatch = false
										break
									}
								}
							}
							return allMatch
						}
					}
				}
			}
		}
	}

	return false
}

// analyzePodTrafficWithSubsets 分析Pod在所有相关subset中的流量类型
func (h *Handler) analyzePodTrafficWithSubsets(podName, serviceName string, podLabels map[string]string, virtualServices []map[string]interface{}, drMap map[string]map[string]interface{}) []map[string]interface{} {
	var results []map[string]interface{}

	// 查找与此服务相关的DestinationRule
	var relevantDR map[string]interface{}
	for host, dr := range drMap {
		if host == serviceName || strings.Contains(host, serviceName) {
			relevantDR = dr
			break
		}
	}

	if relevantDR == nil {
		// 没有找到相关的DestinationRule，使用原来的逻辑
		trafficType := h.analyzePodTraffic(podName, serviceName, podLabels, virtualServices, drMap)
		return []map[string]interface{}{trafficType}
	}

	// 获取所有定义的subset
	var allSubsets []map[string]interface{}
	if spec, ok := relevantDR["spec"].(map[string]interface{}); ok {
		if subsets, ok := spec["subsets"].([]interface{}); ok {
			for _, subset := range subsets {
				if subsetMap, ok := subset.(map[string]interface{}); ok {
					allSubsets = append(allSubsets, subsetMap)
				}
			}
		}
	}

	// 检查Pod是否匹配每个subset
	for _, subset := range allSubsets {
		subsetName := ""
		if name, ok := subset["name"].(string); ok {
			subsetName = name
		}

		// 检查Pod标签是否匹配subset标签
		if h.podMatchesSubsetLabels(podLabels, subset) {
			// 检查这个subset是否有流量路由
			trafficType := h.getSubsetTrafficType(serviceName, subsetName, virtualServices)
			matchContent := h.getSubsetMatchContent(serviceName, subsetName, virtualServices)
			results = append(results, map[string]interface{}{
				"type":         trafficType,
				"vsName":       h.getVSNameForSubset(serviceName, subsetName, virtualServices),
				"subset":       subsetName,
				"matchContent": matchContent,
			})
		}
	}

	// 如果没有匹配任何subset，使用原来的逻辑
	if len(results) == 0 {
		trafficType := h.analyzePodTraffic(podName, serviceName, podLabels, virtualServices, drMap)
		return []map[string]interface{}{trafficType}
	}

	return results
}

// podMatchesSubsetLabels 检查Pod标签是否匹配subset标签
func (h *Handler) podMatchesSubsetLabels(podLabels map[string]string, subset map[string]interface{}) bool {
	if labels, ok := subset["labels"].(map[string]interface{}); ok {
		for labelKey, labelValue := range labels {
			if labelValueStr, ok := labelValue.(string); ok {
				if podLabels[labelKey] != labelValueStr {
					return false
				}
			}
		}
		return true
	}
	return false
}

// getSubsetTrafficType 获取subset的流量类型
func (h *Handler) getSubsetTrafficType(serviceName, subsetName string, virtualServices []map[string]interface{}) string {
	for _, vs := range virtualServices {
		if spec, ok := vs["spec"].(map[string]interface{}); ok {
			// 检查hosts是否匹配
			hosts := []string{}
			if hostsInterface, ok := spec["hosts"].([]interface{}); ok {
				for _, host := range hostsInterface {
					if hostStr, ok := host.(string); ok {
						hosts = append(hosts, hostStr)
					}
				}
			}

			serviceMatched := false
			for _, host := range hosts {
				if host == serviceName || strings.Contains(host, serviceName) {
					serviceMatched = true
					break
				}
			}

			if !serviceMatched {
				continue
			}

			// 检查HTTP路由
			if httpRoutes, ok := spec["http"].([]interface{}); ok {
				for _, httpRoute := range httpRoutes {
					if route, ok := httpRoute.(map[string]interface{}); ok {
						// 检查route目标
						if routeDestinations, ok := route["route"].([]interface{}); ok {
							for _, dest := range routeDestinations {
								if destination, ok := dest.(map[string]interface{}); ok {
									if destInfo, ok := destination["destination"].(map[string]interface{}); ok {
										if subset, ok := destInfo["subset"].(string); ok && subset == subsetName {
											// 检查是否有match条件
											if _, hasMatch := route["match"]; hasMatch {
												return "灰度流量"
											} else {
												return "基础流量"
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return "未进行灰度"
}

// getVSNameForSubset 获取subset对应的VirtualService名称
func (h *Handler) getVSNameForSubset(serviceName, subsetName string, virtualServices []map[string]interface{}) string {
	for _, vs := range virtualServices {
		vsName := ""
		if metadata, ok := vs["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				vsName = name
			}
		}

		if spec, ok := vs["spec"].(map[string]interface{}); ok {
			// 检查hosts是否匹配
			hosts := []string{}
			if hostsInterface, ok := spec["hosts"].([]interface{}); ok {
				for _, host := range hostsInterface {
					if hostStr, ok := host.(string); ok {
						hosts = append(hosts, hostStr)
					}
				}
			}

			serviceMatched := false
			for _, host := range hosts {
				if host == serviceName || strings.Contains(host, serviceName) {
					serviceMatched = true
					break
				}
			}

			if !serviceMatched {
				continue
			}

			// 检查HTTP路由
			if httpRoutes, ok := spec["http"].([]interface{}); ok {
				for _, httpRoute := range httpRoutes {
					if route, ok := httpRoute.(map[string]interface{}); ok {
						if routeDestinations, ok := route["route"].([]interface{}); ok {
							for _, dest := range routeDestinations {
								if destination, ok := dest.(map[string]interface{}); ok {
									if destInfo, ok := destination["destination"].(map[string]interface{}); ok {
										if subset, ok := destInfo["subset"].(string); ok && subset == subsetName {
											return vsName
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return ""
}

// getSubsetMatchContent 获取subset对应的匹配内容
func (h *Handler) getSubsetMatchContent(serviceName, subsetName string, virtualServices []map[string]interface{}) string {
	var matchContents []string

	for _, vs := range virtualServices {
		if spec, ok := vs["spec"].(map[string]interface{}); ok {
			// 检查hosts是否匹配
			hosts := []string{}
			if hostsInterface, ok := spec["hosts"].([]interface{}); ok {
				for _, host := range hostsInterface {
					if hostStr, ok := host.(string); ok {
						hosts = append(hosts, hostStr)
					}
				}
			}

			serviceMatched := false
			for _, host := range hosts {
				if host == serviceName || strings.Contains(host, serviceName) {
					serviceMatched = true
					break
				}
			}

			if !serviceMatched {
				continue
			}

			// 检查HTTP路由
			if httpRoutes, ok := spec["http"].([]interface{}); ok {
				for _, httpRoute := range httpRoutes {
					if route, ok := httpRoute.(map[string]interface{}); ok {
						// 检查route目标是否匹配subset
						if routeDestinations, ok := route["route"].([]interface{}); ok {
							subsetMatched := false
							for _, dest := range routeDestinations {
								if destination, ok := dest.(map[string]interface{}); ok {
									if destInfo, ok := destination["destination"].(map[string]interface{}); ok {
										if subset, ok := destInfo["subset"].(string); ok && subset == subsetName {
											subsetMatched = true
											break
										}
									}
								}
							}

							// 如果subset匹配，且有match条件，提取match内容
							if subsetMatched {
								if matches, hasMatch := route["match"].([]interface{}); hasMatch {
									for _, match := range matches {
										if matchMap, ok := match.(map[string]interface{}); ok {
											var conditions []string

											// 处理URI匹配
											if uri, ok := matchMap["uri"].(map[string]interface{}); ok {
												for key, value := range uri {
													conditions = append(conditions, fmt.Sprintf("URI %s: %v", key, value))
												}
											}

											// 处理Header匹配
											if headers, ok := matchMap["headers"].(map[string]interface{}); ok {
												for headerName, headerValue := range headers {
													if headerMap, ok := headerValue.(map[string]interface{}); ok {
														for key, value := range headerMap {
															conditions = append(conditions, fmt.Sprintf("Header %s %s: %v", headerName, key, value))
														}
													} else {
														conditions = append(conditions, fmt.Sprintf("Header %s: %v", headerName, headerValue))
													}
												}
											}

											if len(conditions) > 0 {
												matchContents = append(matchContents, strings.Join(conditions, ", "))
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return strings.Join(matchContents, "; ")
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

import {del, get, post, put} from "@/plugins/request"

const baseUrl = "/api/v1/istio"

// VirtualService API
export function listVirtualServices(clusterName, namespace, search, keywords, pageNum, pageSize) {
  let url = `${baseUrl}/${clusterName}/virtualservices`
  if (namespace && namespace.trim() !== "") {
    url = `${baseUrl}/${clusterName}/namespaces/${namespace}/virtualservices`
  }
  
  const params = {}
  if (search) {
    params["search"] = true
    if (keywords) {
      params["keywords"] = keywords
    }
    if (pageNum && pageSize) {
      params["pageNum"] = pageNum
      params["pageSize"] = pageSize
    }
  }
  return get(url, params)
}

export function getVirtualService(clusterName, namespace, name) {
  return get(`${baseUrl}/${clusterName}/namespaces/${namespace}/virtualservices/${name}`)
}

export function createVirtualService(clusterName, namespace, data) {
  return post(`${baseUrl}/${clusterName}/namespaces/${namespace}/virtualservices`, data)
}

export function updateVirtualService(clusterName, namespace, name, data) {
  return put(`${baseUrl}/${clusterName}/namespaces/${namespace}/virtualservices/${name}`, data)
}

export function deleteVirtualService(clusterName, namespace, name) {
  return del(`${baseUrl}/${clusterName}/namespaces/${namespace}/virtualservices/${name}`)
}

// DestinationRule API
export function listDestinationRules(clusterName, namespace, search, keywords, pageNum, pageSize) {
  let url = `${baseUrl}/${clusterName}/destinationrules`
  if (namespace && namespace.trim() !== "") {
    url = `${baseUrl}/${clusterName}/namespaces/${namespace}/destinationrules`
  }
  
  const params = {}
  if (search) {
    params["search"] = true
    if (keywords) {
      params["keywords"] = keywords
    }
    if (pageNum && pageSize) {
      params["pageNum"] = pageNum
      params["pageSize"] = pageSize
    }
  }
  return get(url, params)
}

export function getDestinationRule(clusterName, namespace, name) {
  return get(`${baseUrl}/${clusterName}/namespaces/${namespace}/destinationrules/${name}`)
}

export function createDestinationRule(clusterName, namespace, data) {
  return post(`${baseUrl}/${clusterName}/namespaces/${namespace}/destinationrules`, data)
}

export function updateDestinationRule(clusterName, namespace, name, data) {
  return put(`${baseUrl}/${clusterName}/namespaces/${namespace}/destinationrules/${name}`, data)
}

export function deleteDestinationRule(clusterName, namespace, name) {
  return del(`${baseUrl}/${clusterName}/namespaces/${namespace}/destinationrules/${name}`)
}

// Gateway API
export function listGateways(clusterName, namespace, search, keywords, pageNum, pageSize) {
  let url = `${baseUrl}/${clusterName}/gateways`
  if (namespace && namespace.trim() !== "") {
    url = `${baseUrl}/${clusterName}/namespaces/${namespace}/gateways`
  }
  
  const params = {}
  if (search) {
    params["search"] = true
    if (keywords) {
      params["keywords"] = keywords
    }
    if (pageNum && pageSize) {
      params["pageNum"] = pageNum
      params["pageSize"] = pageSize
    }
  }
  return get(url, params)
}

export function getGateway(clusterName, namespace, name) {
  return get(`${baseUrl}/${clusterName}/namespaces/${namespace}/gateways/${name}`)
}

export function createGateway(clusterName, namespace, data) {
  return post(`${baseUrl}/${clusterName}/namespaces/${namespace}/gateways`, data)
}

export function updateGateway(clusterName, namespace, name, data) {
  return put(`${baseUrl}/${clusterName}/namespaces/${namespace}/gateways/${name}`, data)
}

export function deleteGateway(clusterName, namespace, name) {
  return del(`${baseUrl}/${clusterName}/namespaces/${namespace}/gateways/${name}`)
}

// Traffic Analytics API
export function getTrafficAnalytics(clusterName, namespace) {
  const params = {}
  if (namespace) {
    params["namespace"] = namespace
  }
  return get(`${baseUrl}/${clusterName}/traffic-analytics`, params)
}

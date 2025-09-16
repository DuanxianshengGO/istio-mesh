import Layout from "@/business/app-layout/horizontal-layout"

const Istio = {
  path: "/istio",
  sort: 4,
  component: Layout,
  name: "IstioMesh",
  requirePermission: {
    apiGroup: "networking.istio.io",
    resource: "virtualservices",
    verb: "list",
    scope: "namespace"
  },
  meta: {
    title: "business.istio.mesh",
    icon: "iconfont iconnetwork",
    scope: "namespace"
  },
  children: [
    {
      path: "/istio/virtualservices",
      component: () => import("@/business/istio/virtualservice"),
      name: "VirtualServices",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "virtualservices",
        verb: "list",
        scope: "namespace"
      },
      meta: {
        title: "business.istio.virtualservice",
        activeMenu: "/istio/virtualservices",
      }
    },
    {
      path: "/istio/virtualservices/create",
      component: () => import("@/business/istio/virtualservice/create"),
      name: "VirtualServiceCreate",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "virtualservices",
        verb: "create",
        scope: "namespace"
      },
      hidden: true,
      meta: {
        activeMenu: "/istio/virtualservices",
      }
    },
    {
      path: "/istio/virtualservices/:namespace/:name/edit",
      component: () => import("@/business/istio/virtualservice/edit"),
      name: "VirtualServiceEdit",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "virtualservices",
        verb: "update",
        scope: "namespace"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/istio/virtualservices",
      }
    },
    {
      path: "/istio/virtualservices/:namespace/:name",
      component: () => import("@/business/istio/virtualservice/detail"),
      name: "VirtualServiceDetail",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "virtualservices",
        verb: "get",
        scope: "namespace"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/istio/virtualservices",
      }
    },
    {
      path: "/istio/destinationrules",
      component: () => import("@/business/istio/destinationrule"),
      name: "DestinationRules",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "destinationrules",
        verb: "list",
        scope: "namespace"
      },
      meta: {
        title: "business.istio.destinationrule",
        activeMenu: "/istio/destinationrules",
      }
    },
    {
      path: "/istio/destinationrules/create",
      component: () => import("@/business/istio/destinationrule/create"),
      name: "DestinationRuleCreate",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "destinationrules",
        verb: "create",
        scope: "namespace"
      },
      hidden: true,
      meta: {
        activeMenu: "/istio/destinationrules",
      }
    },
    {
      path: "/istio/destinationrules/:namespace/:name/edit",
      component: () => import("@/business/istio/destinationrule/edit"),
      name: "DestinationRuleEdit",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "destinationrules",
        verb: "update",
        scope: "namespace"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/istio/destinationrules",
      }
    },
    {
      path: "/istio/destinationrules/:namespace/:name",
      component: () => import("@/business/istio/destinationrule/detail"),
      name: "DestinationRuleDetail",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "destinationrules",
        verb: "get",
        scope: "namespace"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/istio/destinationrules",
      }
    },
    {
      path: "/istio/gateways",
      component: () => import("@/business/istio/gateway"),
      name: "Gateways",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "gateways",
        verb: "list",
        scope: "namespace"
      },
      meta: {
        title: "business.istio.gateway",
        activeMenu: "/istio/gateways",
      }
    },
    {
      path: "/istio/gateways/create",
      component: () => import("@/business/istio/gateway/create"),
      name: "GatewayCreate",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "gateways",
        verb: "create",
        scope: "namespace"
      },
      hidden: true,
      meta: {
        activeMenu: "/istio/gateways",
      }
    },
    {
      path: "/istio/gateways/:namespace/:name/edit",
      component: () => import("@/business/istio/gateway/edit"),
      name: "GatewayEdit",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "gateways",
        verb: "update",
        scope: "namespace"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/istio/gateways",
      }
    },
    {
      path: "/istio/gateways/:namespace/:name",
      component: () => import("@/business/istio/gateway/detail"),
      name: "GatewayDetail",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "gateways",
        verb: "get",
        scope: "namespace"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/istio/gateways",
      }
    },
    {
      path: "/istio/traffic-analytics",
      component: () => import("@/business/istio/traffic-analytics"),
      name: "TrafficAnalytics",
      requirePermission: {
        apiGroup: "networking.istio.io",
        resource: "virtualservices",
        verb: "list",
        scope: "namespace"
      },
      meta: {
        title: "business.istio.traffic_analytics",
        activeMenu: "/istio/traffic-analytics",
      }
    }
  ]
}

export default Istio

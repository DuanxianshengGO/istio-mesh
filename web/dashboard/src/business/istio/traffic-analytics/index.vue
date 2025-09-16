<template>
  <layout-content :header="$t('business.istio.traffic_analytics')">
    <div v-loading="loading">
      <!-- 过滤器 -->
      <el-card style="margin-bottom: 20px;">
        <div slot="header">
          <span>{{ $t('business.istio.filter') }}</span>
        </div>
        <el-form :inline="true" :model="filterForm" size="small">
          <el-form-item :label="$t('commons.table.namespace')">
            <namespace-select v-model="filterForm.namespace" :cluster="cluster" @change="onNamespaceChange"></namespace-select>
          </el-form-item>
          <el-form-item :label="$t('business.istio.service')">
            <el-select v-model="filterForm.service" clearable placeholder="Select service">
              <el-option
                v-for="service in services"
                :key="service"
                :label="service"
                :value="service">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadTrafficAnalytics">{{ $t('commons.button.search') }}</el-button>
            <el-button @click="resetFilter">{{ $t('commons.button.reset') }}</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 流量概览 -->
      <el-row :gutter="20" style="margin-bottom: 20px;">
        <el-col :span="6">
          <el-card>
            <div slot="header">{{ $t('business.istio.total_pods') }}</div>
            <div class="metric-value">{{ analytics.summary ? analytics.summary.totalPods : 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div slot="header">基础流量</div>
            <div class="metric-value basic-traffic">{{ analytics.summary ? analytics.summary.basicTraffic : 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div slot="header">灰度流量</div>
            <div class="metric-value gray-traffic">{{ analytics.summary ? analytics.summary.grayTraffic : 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div slot="header">无流量</div>
            <div class="metric-value no-traffic">{{ analytics.summary ? analytics.summary.noTraffic : 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 流量分析表格 -->
      <el-card>
        <div slot="header">
          <span>流量分析详情</span>
        </div>
        <el-table :data="filteredTrafficData" style="width: 100%">
          <el-table-column prop="podName" label="Pod名称" width="200"></el-table-column>
          <el-table-column prop="serviceName" label="Service" width="150"></el-table-column>
          <el-table-column prop="vsName" label="VirtualService" width="200"></el-table-column>
          <el-table-column prop="trafficType" label="流量类型" width="120">
            <template slot-scope="scope">
              <el-tag
                :type="getTrafficTypeColor(scope.row.trafficType)"
                size="small">
                {{ scope.row.trafficType }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="subset" label="Subset" width="120"></el-table-column>
        </el-table>
      </el-card>
          <el-card>
            <div slot="header">{{ $t('business.istio.traffic_routes') }}</div>
            <div class="metric-value">{{ analytics.trafficRoutes ? analytics.trafficRoutes.length : 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 流量路由图 -->
      <el-card style="margin-bottom: 20px;">
        <div slot="header">
          <span>{{ $t('business.istio.traffic_flow') }}</span>
        </div>
        <div id="traffic-flow-chart" style="height: 400px;"></div>
      </el-card>

      <!-- 服务流量详情 -->
      <el-card>
        <div slot="header">
          <span>{{ $t('business.istio.service_traffic_details') }}</span>
        </div>
        
        <el-tabs v-model="activeTab">
          <el-tab-pane :label="$t('business.istio.virtual_services')" name="virtualservices">
            <el-table :data="analytics.virtualServices" style="width: 100%">
              <el-table-column prop="metadata.name" :label="$t('commons.table.name')" width="200"></el-table-column>
              <el-table-column prop="metadata.namespace" :label="$t('commons.table.namespace')" width="150"></el-table-column>
              <el-table-column :label="$t('business.istio.hosts')" width="200">
                <template v-slot:default="{row}">
                  <el-tag v-for="host in row.spec.hosts" :key="host" size="mini" style="margin: 2px;">
                    {{ host }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="$t('business.istio.destinations')" min-width="300">
                <template v-slot:default="{row}">
                  <div v-if="row.spec.http">
                    <div v-for="(route, index) in row.spec.http" :key="index">
                      <div v-for="(dest, destIndex) in route.route" :key="destIndex">
                        <el-tag type="success" size="mini" style="margin: 2px;">
                          {{ dest.destination.host }}
                          <span v-if="dest.destination.subset">:{{ dest.destination.subset }}</span>
                          <span v-if="dest.weight"> ({{ dest.weight }}%)</span>
                        </el-tag>
                      </div>
                    </div>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          
          <el-tab-pane :label="$t('business.istio.destination_rules')" name="destinationrules">
            <el-table :data="analytics.destinationRules" style="width: 100%">
              <el-table-column prop="metadata.name" :label="$t('commons.table.name')" width="200"></el-table-column>
              <el-table-column prop="metadata.namespace" :label="$t('commons.table.namespace')" width="150"></el-table-column>
              <el-table-column prop="spec.host" :label="$t('business.istio.host')" width="200"></el-table-column>
              <el-table-column :label="$t('business.istio.subsets')" min-width="300">
                <template v-slot:default="{row}">
                  <div v-if="row.spec.subsets">
                    <el-tag v-for="subset in row.spec.subsets" :key="subset.name" size="mini" style="margin: 2px;">
                      {{ subset.name }}: {{ JSON.stringify(subset.labels) }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          
          <el-tab-pane :label="$t('business.istio.pod_traffic')" name="podtraffic">
            <el-table :data="analytics.podTraffic" style="width: 100%">
              <el-table-column prop="podName" :label="$t('business.workload.pod_name')" width="200"></el-table-column>
              <el-table-column prop="namespace" :label="$t('commons.table.namespace')" width="150"></el-table-column>
              <el-table-column prop="service" :label="$t('business.istio.service')" width="150"></el-table-column>
              <el-table-column prop="version" :label="$t('business.istio.version')" width="100"></el-table-column>
              <el-table-column prop="trafficRule" :label="$t('business.istio.traffic_rule')" min-width="300"></el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import NamespaceSelect from "@/components/namespace-select"
import { getTrafficAnalytics } from "@/api/istio"
import { listServicesWithNs } from "@/api/services"

export default {
  name: "TrafficAnalytics",
  components: { 
    LayoutContent,
    NamespaceSelect
  },
  data() {
    return {
      cluster: "",
      loading: false,
      activeTab: "virtualservices",
      filterForm: {
        namespace: "",
        service: ""
      },
      services: [],
      analytics: {
        trafficAnalysis: [],
        summary: {
          totalPods: 0,
          totalVS: 0,
          totalDR: 0,
          basicTraffic: 0,
          grayTraffic: 0,
          noTraffic: 0
        }
      }
    }
  },
  computed: {
    filteredTrafficData() {
      if (!this.analytics.trafficAnalysis) return []

      let data = this.analytics.trafficAnalysis

      // 按服务过滤
      if (this.filterForm.service) {
        data = data.filter(item => item.serviceName === this.filterForm.service)
      }

      return data
    }
  },
  methods: {
    loadTrafficAnalytics() {
      this.loading = true
      getTrafficAnalytics(this.cluster, this.filterForm.namespace)
        .then(data => {
          // 处理 KubePi API 返回的数据结构
          if (data && data.data) {
            this.analytics = data.data
          } else {
            this.analytics = data
          }
          this.updateServicesList()
          this.loading = false
        })
        .catch(() => {
          this.loading = false
        })
    },
    updateServicesList() {
      // 从流量分析数据中提取服务列表
      if (this.analytics.trafficAnalysis) {
        const serviceSet = new Set()
        this.analytics.trafficAnalysis.forEach(item => {
          if (item.serviceName) {
            serviceSet.add(item.serviceName)
          }
        })
        this.services = Array.from(serviceSet)
      }
    },
    getTrafficTypeColor(trafficType) {
      switch (trafficType) {
        case '基础流量':
          return 'success'
        case '灰度流量':
          return 'warning'
        case '无流量':
          return 'info'
        default:
          return ''
      }
    },
    onNamespaceChange() {
      this.loadTrafficAnalytics()
    },
    resetFilter() {
      this.filterForm = {
        namespace: "",
        service: ""
      }
      this.loadTrafficAnalytics()
    },
    renderTrafficFlowChart() {
      // TODO: 使用 ECharts 或其他图表库渲染流量流向图
      // 这里先显示一个简单的文本提示
      const chartContainer = document.getElementById('traffic-flow-chart')
      if (chartContainer) {
        chartContainer.innerHTML = `
          <div style="display: flex; align-items: center; justify-content: center; height: 100%; color: #666;">
            <div>
              <i class="el-icon-loading" style="font-size: 24px; margin-bottom: 10px;"></i>
              <div>${this.$t('business.istio.traffic_flow_chart_placeholder')}</div>
            </div>
          </div>
        `
      }
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.filterForm.namespace = this.$route.query.namespace || ""
    this.loadServices()
    this.loadTrafficAnalytics()
  },
  mounted() {
    this.renderTrafficFlowChart()
  }
}
</script>

<style scoped>
.metric-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
}

.metric-value.basic-traffic {
  color: #67C23A;
}

.metric-value.gray-traffic {
  color: #E6A23C;
}

.metric-value.no-traffic {
  color: #909399;
}

.metric-value:not(.basic-traffic):not(.gray-traffic):not(.no-traffic) {
  color: #409EFF;
}
</style>

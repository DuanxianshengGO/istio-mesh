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
            <div slot="header">{{ $t('business.istio.basic_traffic') }}</div>
            <div class="metric-value basic-traffic">{{ analytics.summary ? analytics.summary.basicTraffic : 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div slot="header">{{ $t('business.istio.gray_traffic') }}</div>
            <div class="metric-value gray-traffic">{{ analytics.summary ? analytics.summary.grayTraffic : 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div slot="header">{{ $t('business.istio.no_traffic') }}</div>
            <div class="metric-value no-traffic">{{ analytics.summary ? analytics.summary.noTraffic : 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 流量分析表格 -->
      <el-card>
        <div slot="header">
          <span>流量分析详情</span>
        </div>
        <el-table :data="filteredTrafficData" style="width: 100%" size="small" :show-header="true">
          <el-table-column prop="podName" label="Pod名称" min-width="200" :filterable="false"></el-table-column>
          <el-table-column prop="serviceName" label="Service" min-width="150" :filterable="false"></el-table-column>
          <el-table-column prop="vsName" label="VirtualService" min-width="200" :filterable="false"></el-table-column>
          <el-table-column prop="trafficType" label="流量类型" width="120" :filterable="false">
            <template slot-scope="scope">
              <el-tag
                :type="getTrafficTypeColor(scope.row.trafficType)"
                size="small">
                {{ scope.row.trafficType }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="subset" label="Subset" width="120" :filterable="false"></el-table-column>
          <el-table-column prop="matchContent" :label="$t('business.istio.match_content')" min-width="200" :filterable="false">
            <template slot-scope="scope">
              <el-tooltip v-if="scope.row.matchContent" :content="scope.row.matchContent" placement="top">
                <span class="match-content-text">{{ scope.row.matchContent }}</span>
              </el-tooltip>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
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
        case '未进行灰度':
          return 'info'
        default:
          return ''
      }
    },
    onNamespaceChange() {
      this.loadTrafficAnalytics()
    },


  },
  created() {
    this.cluster = this.$route.query.cluster
    this.filterForm.namespace = this.$route.query.namespace || ""
    this.loadTrafficAnalytics()
  },

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

.match-content-text {
  display: inline-block;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

/* 隐藏所有可能的重置按钮 */
::v-deep .el-button {
  &:contains('commons.button.reset') {
    display: none !important;
  }
}

/* 隐藏表格过滤器相关按钮 */
::v-deep .el-table__column-filter-trigger {
  display: none !important;
}

::v-deep .el-table-filter__bottom {
  display: none !important;
}

/* 隐藏包含特定文本的按钮 */
::v-deep button[data-v-27185e19] {
  &:contains('commons.button.reset') {
    display: none !important;
  }
}
</style>

<template>
  <layout-content :header="$t('business.istio.destinationrule_detail')" :back-to="{name: 'DestinationRules', query: $route.query}">
    <div v-loading="loading">
      <el-descriptions :title="item.metadata ? item.metadata.name : ''" :column="2" border>
        <el-descriptions-item :label="$t('commons.table.name')">
          {{ item.metadata ? item.metadata.name : '' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('commons.table.namespace')">
          {{ item.metadata ? item.metadata.namespace : '' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('commons.table.created_time')">
          {{ item.metadata ? (item.metadata.creationTimestamp | dateFormat) : '' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('business.istio.host')">
          {{ item.spec ? item.spec.host : '' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('business.istio.load_balancer')" :span="2">
          {{ item.spec && item.spec.trafficPolicy && item.spec.trafficPolicy.loadBalancer ? 
               item.spec.trafficPolicy.loadBalancer.simple : 'ROUND_ROBIN' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider></el-divider>

      <h3>{{ $t('business.istio.subsets') }}</h3>
      <div v-if="item.spec && item.spec.subsets">
        <el-card v-for="(subset, index) in item.spec.subsets" :key="index" style="margin-bottom: 10px;">
          <div slot="header">
            <span>{{ subset.name }}</span>
          </div>
          
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item :label="$t('business.common.label')">
              <div v-if="subset.labels">
                <el-tag v-for="(value, key) in subset.labels" :key="key" size="small" style="margin: 2px;">
                  {{ key }}: {{ value }}
                </el-tag>
              </div>
            </el-descriptions-item>
            
            <el-descriptions-item v-if="subset.trafficPolicy" :label="$t('business.istio.traffic_policy')">
              {{ JSON.stringify(subset.trafficPolicy) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>

      <el-divider></el-divider>

      <h3>{{ $t('business.istio.yaml_content') }}</h3>
      <yaml-editor :value="item" :read-only="true"></yaml-editor>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import { getDestinationRule } from "@/api/istio"

export default {
  name: "DestinationRuleDetail",
  components: { 
    LayoutContent,
    YamlEditor
  },
  props: {
    namespace: String,
    name: String
  },
  data() {
    return {
      cluster: "",
      loading: false,
      item: {}
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getDestinationRule(this.cluster, this.namespace, this.name)
        .then(res => {
          this.item = res
          this.loading = false
        })
        .catch(() => {
          this.loading = false
        })
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.getDetail()
  }
}
</script>

<style scoped>
</style>

<template>
  <layout-content :header="$t('business.istio.virtualservice_detail')" :back-to="{name: 'VirtualServices', query: $route.query}">
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
        <el-descriptions-item :label="$t('business.istio.hosts')">
          <div v-if="item.spec && item.spec.hosts">
            <el-tag v-for="host in item.spec.hosts" :key="host" size="small" style="margin: 2px;">
              {{ host }}
            </el-tag>
          </div>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('business.istio.gateways')" :span="2">
          <div v-if="item.spec && item.spec.gateways">
            <el-tag v-for="gateway in item.spec.gateways" :key="gateway" size="small" style="margin: 2px;">
              {{ gateway }}
            </el-tag>
          </div>
          <span v-else>mesh</span>
        </el-descriptions-item>
      </el-descriptions>

      <el-divider></el-divider>

      <h3>{{ $t('business.istio.http_routes') }}</h3>
      <div v-if="item.spec && item.spec.http">
        <el-card v-for="(route, index) in item.spec.http" :key="index" style="margin-bottom: 10px;">
          <div slot="header">
            <span>{{ $t('business.istio.route') }} {{ index + 1 }}</span>
          </div>
          
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item :label="$t('business.istio.match_conditions')">
              <div v-if="route.match">
                <div v-for="(match, matchIndex) in route.match" :key="matchIndex">
                  <el-tag v-if="match.uri" size="mini">
                    URI: {{ Object.keys(match.uri)[0] }} = {{ Object.values(match.uri)[0] }}
                  </el-tag>
                  <el-tag v-if="match.headers" size="mini" style="margin-left: 5px;">
                    Headers: {{ JSON.stringify(match.headers) }}
                  </el-tag>
                </div>
              </div>
            </el-descriptions-item>
            
            <el-descriptions-item :label="$t('business.istio.destinations')">
              <div v-if="route.route">
                <div v-for="(destination, destIndex) in route.route" :key="destIndex" style="margin-bottom: 5px;">
                  <el-tag type="success" size="small">
                    {{ destination.destination.host }}
                    <span v-if="destination.destination.subset">:{{ destination.destination.subset }}</span>
                    <span v-if="destination.weight"> ({{ destination.weight }}%)</span>
                  </el-tag>
                </div>
              </div>
            </el-descriptions-item>
            
            <el-descriptions-item v-if="route.fault" :label="$t('business.istio.fault_injection')">
              {{ JSON.stringify(route.fault) }}
            </el-descriptions-item>
            
            <el-descriptions-item v-if="route.timeout" :label="$t('business.istio.timeout')">
              {{ route.timeout }}
            </el-descriptions-item>
            
            <el-descriptions-item v-if="route.retries" :label="$t('business.istio.retries')">
              {{ JSON.stringify(route.retries) }}
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
import { getVirtualService } from "@/api/istio"

export default {
  name: "VirtualServiceDetail",
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
      getVirtualService(this.cluster, this.namespace, this.name)
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

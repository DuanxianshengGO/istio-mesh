<template>
  <layout-content :header="$t('business.istio.gateway_detail')" :back-to="{name: 'Gateways', query: $route.query}">
    <div v-loading="loading">
      <el-descriptions :title="item.metadata ? item.metadata.name : ''" :column="2" border>
        <el-descriptions-item :label="$t('commons.table.name')">
          {{ item.metadata ? item.metadata.name : '' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('commons.table.namespace')">
          {{ item.metadata ? item.metadata.namespace : '' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('commons.table.created_time')">
          {{ item.metadata ? $dateFormat(item.metadata.creationTimestamp) : '' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('business.istio.selector')" :span="2">
          <div v-if="item.spec && item.spec.selector">
            <el-tag v-for="(value, key) in item.spec.selector" :key="key" size="small" style="margin: 2px;">
              {{ key }}: {{ value }}
            </el-tag>
          </div>
        </el-descriptions-item>
      </el-descriptions>

      <el-divider></el-divider>

      <h3>{{ $t('business.istio.servers') }}</h3>
      <div v-if="item.spec && item.spec.servers">
        <el-card v-for="(server, index) in item.spec.servers" :key="index" style="margin-bottom: 10px;">
          <div slot="header">
            <span>{{ $t('business.istio.server') }} {{ index + 1 }}</span>
          </div>
          
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item :label="$t('business.istio.port')">
              {{ server.port.number }}/{{ server.port.protocol }}
              <span v-if="server.port.name"> ({{ server.port.name }})</span>
            </el-descriptions-item>
            
            <el-descriptions-item :label="$t('business.istio.hosts')">
              <div v-if="server.hosts">
                <el-tag v-for="host in server.hosts" :key="host" size="small" style="margin: 2px;">
                  {{ host }}
                </el-tag>
              </div>
            </el-descriptions-item>
            
            <el-descriptions-item v-if="server.tls" :label="$t('business.istio.tls')" :span="2">
              <div>
                <el-tag type="info" size="small">{{ $t('business.istio.tls_mode') }}: {{ server.tls.mode }}</el-tag>
                <el-tag v-if="server.tls.credentialName" type="success" size="small" style="margin-left: 5px;">
                  {{ $t('business.istio.credential_name') }}: {{ server.tls.credentialName }}
                </el-tag>
              </div>
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
import { getGateway } from "@/api/istio"

export default {
  name: "GatewayDetail",
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
      getGateway(this.cluster, this.namespace, this.name)
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

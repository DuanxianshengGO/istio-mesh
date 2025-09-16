<template>
  <layout-content :header="$t('business.istio.create_gateway')" :back-to="{name: 'Gateways', query: $route.query}">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="16">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="form" :rules="rules" label-width="120px" label-position="left">
            <el-form-item :label="$t('commons.table.name')" prop="name" required>
              <el-input v-model="form.name" clearable></el-input>
            </el-form-item>
            
            <el-form-item :label="$t('commons.table.namespace')" prop="namespace" required>
              <namespace-select v-model="form.namespace" :cluster="cluster"></namespace-select>
            </el-form-item>
            
            <el-form-item :label="$t('business.istio.selector')">
              <div v-for="(selector, index) in form.selectors" :key="index" style="display: flex; margin-bottom: 5px;">
                <el-input v-model="selector.key" placeholder="Key" style="margin-right: 10px;"></el-input>
                <el-input v-model="selector.value" placeholder="Value" style="margin-right: 10px;"></el-input>
                <el-button size="mini" type="danger" @click="removeSelector(index)">
                  {{ $t('commons.button.delete') }}
                </el-button>
              </div>
              <el-button size="small" @click="addSelector">
                {{ $t('business.common.add_label') }}
              </el-button>
            </el-form-item>
            
            <el-form-item :label="$t('business.istio.servers')">
              <div v-for="(server, index) in form.servers" :key="index" style="border: 1px solid #ddd; padding: 10px; margin-bottom: 10px;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                  <span>{{ $t('business.istio.server') }} {{ index + 1 }}</span>
                  <el-button size="mini" type="danger" @click="removeServer(index)">
                    {{ $t('commons.button.delete') }}
                  </el-button>
                </div>
                
                <el-form-item :label="$t('business.istio.port')" required>
                  <el-row :gutter="10">
                    <el-col :span="8">
                      <el-input-number v-model="server.port.number" :min="1" :max="65535" placeholder="80"></el-input-number>
                    </el-col>
                    <el-col :span="8">
                      <el-input v-model="server.port.name" placeholder="http"></el-input>
                    </el-col>
                    <el-col :span="8">
                      <el-select v-model="server.port.protocol" placeholder="Protocol">
                        <el-option label="HTTP" value="HTTP"></el-option>
                        <el-option label="HTTPS" value="HTTPS"></el-option>
                        <el-option label="TCP" value="TCP"></el-option>
                        <el-option label="TLS" value="TLS"></el-option>
                      </el-select>
                    </el-col>
                  </el-row>
                </el-form-item>
                
                <el-form-item :label="$t('business.istio.hosts')" required>
                  <el-select v-model="server.hosts" multiple allow-create filterable placeholder="Enter hosts">
                    <el-option
                      v-for="host in commonHosts"
                      :key="host"
                      :label="host"
                      :value="host">
                    </el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item v-if="server.port.protocol === 'HTTPS' || server.port.protocol === 'TLS'" :label="$t('business.istio.tls')">
                  <el-form-item :label="$t('business.istio.tls_mode')">
                    <el-select v-model="server.tls.mode" placeholder="TLS Mode">
                      <el-option label="SIMPLE" value="SIMPLE"></el-option>
                      <el-option label="MUTUAL" value="MUTUAL"></el-option>
                      <el-option label="PASSTHROUGH" value="PASSTHROUGH"></el-option>
                      <el-option label="ISTIO_MUTUAL" value="ISTIO_MUTUAL"></el-option>
                    </el-select>
                  </el-form-item>
                  
                  <el-form-item v-if="server.tls.mode === 'SIMPLE' || server.tls.mode === 'MUTUAL'" :label="$t('business.istio.credential_name')">
                    <el-input v-model="server.tls.credentialName" placeholder="tls-secret"></el-input>
                  </el-form-item>
                </el-form-item>
              </div>
              
              <el-button type="primary" size="small" @click="addServer">
                {{ $t('business.istio.add_server') }}
              </el-button>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="onSubmit" :loading="loading">
                {{ $t("commons.button.submit") }}
              </el-button>
              <el-button @click="onCancel">{{ $t("commons.button.cancel") }}</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import NamespaceSelect from "@/components/namespace-select"
import { createGateway } from "@/api/istio"

export default {
  name: "GatewayCreate",
  components: { 
    LayoutContent,
    NamespaceSelect
  },
  data() {
    return {
      cluster: "",
      loading: false,
      form: {
        name: "",
        namespace: "",
        selectors: [
          { key: "istio", value: "ingressgateway" }
        ],
        servers: [
          {
            port: {
              number: 80,
              name: "http",
              protocol: "HTTP"
            },
            hosts: ["*"],
            tls: {
              mode: "SIMPLE",
              credentialName: ""
            }
          }
        ]
      },
      rules: {
        name: [
          { required: true, message: this.$t("commons.rule.name"), trigger: "blur" },
          { pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/, message: this.$t("commons.rule.name_rule"), trigger: "blur" }
        ],
        namespace: [
          { required: true, message: this.$t("commons.rule.namespace"), trigger: "blur" }
        ]
      },
      commonHosts: [
        "*",
        "example.com",
        "api.example.com",
        "www.example.com"
      ]
    }
  },
  methods: {
    addSelector() {
      this.form.selectors.push({ key: "", value: "" })
    },
    removeSelector(index) {
      if (this.form.selectors.length > 1) {
        this.form.selectors.splice(index, 1)
      }
    },
    addServer() {
      this.form.servers.push({
        port: {
          number: 443,
          name: "https",
          protocol: "HTTPS"
        },
        hosts: ["*"],
        tls: {
          mode: "SIMPLE",
          credentialName: ""
        }
      })
    },
    removeServer(index) {
      if (this.form.servers.length > 1) {
        this.form.servers.splice(index, 1)
      }
    },
    onSubmit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          
          const gateway = {
            apiVersion: "networking.istio.io/v1beta1",
            kind: "Gateway",
            metadata: {
              name: this.form.name,
              namespace: this.form.namespace
            },
            spec: {
              selector: this.form.selectors.reduce((acc, selector) => {
                if (selector.key && selector.value) {
                  acc[selector.key] = selector.value
                }
                return acc
              }, {}),
              servers: this.form.servers.map(server => {
                const serverSpec = {
                  port: server.port,
                  hosts: server.hosts
                }
                
                if (server.port.protocol === 'HTTPS' || server.port.protocol === 'TLS') {
                  serverSpec.tls = {}
                  if (server.tls.mode) {
                    serverSpec.tls.mode = server.tls.mode
                  }
                  if (server.tls.credentialName) {
                    serverSpec.tls.credentialName = server.tls.credentialName
                  }
                }
                
                return serverSpec
              })
            }
          }
          
          createGateway(this.cluster, this.form.namespace, gateway)
            .then(() => {
              this.$message({
                type: "success",
                message: this.$t("commons.msg.create_success")
              })
              this.$router.push({name: "Gateways", query: this.$route.query})
            })
            .finally(() => {
              this.loading = false
            })
        }
      })
    },
    onCancel() {
      this.$router.push({name: "Gateways", query: this.$route.query})
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.form.namespace = this.$route.query.namespace || "default"
  }
}
</script>

<style scoped>
</style>

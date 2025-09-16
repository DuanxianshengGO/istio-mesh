<template>
  <layout-content :header="$t('business.istio.create_virtualservice')" :back-to="{name: 'VirtualServices', query: $route.query}">
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
            
            <el-form-item :label="$t('business.istio.hosts')" prop="hosts" required>
              <el-select v-model="form.hosts" multiple allow-create filterable placeholder="Enter hosts">
                <el-option
                  v-for="host in commonHosts"
                  :key="host"
                  :label="host"
                  :value="host">
                </el-option>
              </el-select>
            </el-form-item>
            
            <el-form-item :label="$t('business.istio.gateways')">
              <el-select v-model="form.gateways" multiple allow-create filterable placeholder="Select gateways">
                <el-option
                  v-for="gateway in availableGateways"
                  :key="gateway"
                  :label="gateway"
                  :value="gateway">
                </el-option>
              </el-select>
            </el-form-item>
            
            <el-form-item :label="$t('business.istio.http_routes')">
              <div v-for="(route, index) in form.http" :key="index" style="border: 1px solid #ddd; padding: 10px; margin-bottom: 10px;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                  <span>{{ $t('business.istio.route') }} {{ index + 1 }}</span>
                  <el-button size="mini" type="danger" @click="removeRoute(index)">
                    {{ $t('commons.button.delete') }}
                  </el-button>
                </div>
                
                <el-form-item :label="$t('business.istio.match_uri')">
                  <el-input v-model="route.match[0].uri.prefix" placeholder="/path"></el-input>
                </el-form-item>
                
                <el-form-item :label="$t('business.istio.destination_host')" required>
                  <el-input v-model="route.route[0].destination.host" placeholder="service-name"></el-input>
                </el-form-item>
                
                <el-form-item :label="$t('business.istio.destination_subset')">
                  <el-input v-model="route.route[0].destination.subset" placeholder="v1"></el-input>
                </el-form-item>
                
                <el-form-item :label="$t('business.istio.weight')">
                  <el-input-number v-model="route.route[0].weight" :min="0" :max="100"></el-input-number>
                </el-form-item>
              </div>
              
              <el-button type="primary" size="small" @click="addRoute">
                {{ $t('business.istio.add_route') }}
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
import { createVirtualService } from "@/api/istio"

export default {
  name: "VirtualServiceCreate",
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
        hosts: [],
        gateways: [],
        http: [
          {
            match: [
              {
                uri: {
                  prefix: "/"
                }
              }
            ],
            route: [
              {
                destination: {
                  host: "",
                  subset: ""
                },
                weight: 100
              }
            ]
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
        ],
        hosts: [
          { required: true, message: this.$t("business.istio.rule.hosts"), trigger: "blur" }
        ]
      },
      commonHosts: [
        "*",
        "example.com",
        "api.example.com"
      ],
      availableGateways: [
        "mesh",
        "istio-system/gateway"
      ]
    }
  },
  methods: {
    addRoute() {
      this.form.http.push({
        match: [
          {
            uri: {
              prefix: "/"
            }
          }
        ],
        route: [
          {
            destination: {
              host: "",
              subset: ""
            },
            weight: 100
          }
        ]
      })
    },
    removeRoute(index) {
      if (this.form.http.length > 1) {
        this.form.http.splice(index, 1)
      }
    },
    onSubmit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          
          const virtualService = {
            apiVersion: "networking.istio.io/v1beta1",
            kind: "VirtualService",
            metadata: {
              name: this.form.name,
              namespace: this.form.namespace
            },
            spec: {
              hosts: this.form.hosts,
              http: this.form.http.map(route => ({
                match: route.match,
                route: route.route.filter(r => r.destination.host)
              }))
            }
          }
          
          if (this.form.gateways && this.form.gateways.length > 0) {
            virtualService.spec.gateways = this.form.gateways
          }
          
          createVirtualService(this.cluster, this.form.namespace, virtualService)
            .then(() => {
              this.$message({
                type: "success",
                message: this.$t("commons.msg.create_success")
              })
              this.$router.push({name: "VirtualServices", query: this.$route.query})
            })
            .finally(() => {
              this.loading = false
            })
        }
      })
    },
    onCancel() {
      this.$router.push({name: "VirtualServices", query: this.$route.query})
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

<template>
  <layout-content :header="$t('business.istio.create_virtualservice')" :back-to="{name: 'VirtualServices', query: $route.query}">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="16">
        <div class="grid-content bg-purple-light">
          <!-- 创建方式选择 -->
          <el-tabs v-model="createMode" @tab-click="onModeChange">
            <el-tab-pane label="表单创建" name="form">
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

                <el-form-item :label="$t('business.istio.match_headers')">
                  <div v-for="(header, headerIndex) in route.headers" :key="headerIndex" style="display: flex; margin-bottom: 5px;">
                    <el-input v-model="header.name" placeholder="Header名称" style="width: 150px; margin-right: 10px;"></el-input>
                    <el-select v-model="header.matchType" style="width: 100px; margin-right: 10px;">
                      <el-option label="exact" value="exact"></el-option>
                      <el-option label="prefix" value="prefix"></el-option>
                      <el-option label="regex" value="regex"></el-option>
                    </el-select>
                    <el-input v-model="header.value" placeholder="Header值" style="flex: 1; margin-right: 10px;"></el-input>
                    <el-button size="mini" type="danger" @click="removeHeader(index, headerIndex)">删除</el-button>
                  </div>
                  <el-button size="mini" type="primary" @click="addHeader(index)">添加Header匹配</el-button>
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
            </el-tab-pane>

            <el-tab-pane label="YAML创建" name="yaml">
              <div style="margin-bottom: 20px;">
                <el-button type="primary" size="small" @click="loadTemplate">加载模板</el-button>
                <el-button type="success" size="small" @click="validateYaml">验证YAML</el-button>
              </div>

              <el-form-item :label="$t('business.istio.yaml_content')">
                <div style="border: 1px solid #dcdfe6; border-radius: 4px;">
                  <textarea
                    ref="yamlEditor"
                    v-model="yamlContent"
                    style="width: 100%; height: 500px; border: none; padding: 10px; font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace; font-size: 14px; resize: vertical;"
                    placeholder="请输入VirtualService的YAML配置...">
                  </textarea>
                </div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" @click="onSubmitYaml" :loading="loading">
                  {{ $t("commons.button.submit") }}
                </el-button>
                <el-button @click="onCancel">{{ $t("commons.button.cancel") }}</el-button>
              </el-form-item>
            </el-tab-pane>
          </el-tabs>
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
      createMode: "form",
      yamlContent: "",
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
            headers: [],
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
        headers: [],
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
    addHeader(routeIndex) {
      this.form.http[routeIndex].headers.push({
        name: "",
        matchType: "exact",
        value: ""
      })
    },
    removeHeader(routeIndex, headerIndex) {
      this.form.http[routeIndex].headers.splice(headerIndex, 1)
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
              http: this.form.http.map(route => {
                const httpRoute = {
                  route: route.route.filter(r => r.destination.host)
                }

                // 处理match条件
                if (route.match[0].uri.prefix || route.headers.length > 0) {
                  httpRoute.match = [{
                    ...route.match[0]
                  }]

                  // 添加headers匹配
                  if (route.headers.length > 0) {
                    httpRoute.match[0].headers = {}
                    route.headers.forEach(header => {
                      if (header.name && header.value) {
                        httpRoute.match[0].headers[header.name] = {
                          [header.matchType]: header.value
                        }
                      }
                    })
                  }
                }

                return httpRoute
              })
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
    onModeChange() {
      // 切换创建模式时的处理
    },
    loadTemplate() {
      this.yamlContent = `apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: example-vs
  namespace: default
spec:
  hosts:
  - example-service
  http:
  - match:
    - uri:
        prefix: /api
    - headers:
        version:
          exact: v1
    route:
    - destination:
        host: example-service
        subset: v1
      weight: 100`
    },
    validateYaml() {
      try {
        // 简单的YAML格式验证
        if (!this.yamlContent.trim()) {
          this.$message.error('YAML内容不能为空')
          return false
        }

        if (!this.yamlContent.includes('apiVersion') ||
            !this.yamlContent.includes('kind') ||
            !this.yamlContent.includes('metadata')) {
          this.$message.error('YAML格式不正确，缺少必要字段')
          return false
        }

        this.$message.success('YAML格式验证通过')
        return true
      } catch (error) {
        this.$message.error('YAML格式验证失败: ' + error.message)
        return false
      }
    },
    onSubmitYaml() {
      if (!this.validateYaml()) {
        return
      }

      this.loading = true

      try {
        // 解析YAML内容为JSON对象
        const yamlData = this.parseYamlToJson(this.yamlContent)

        createVirtualService(this.cluster, yamlData.metadata.namespace, yamlData)
          .then(() => {
            this.$message({
              type: "success",
              message: this.$t("commons.msg.create_success")
            })
            this.$router.push({name: "VirtualServices", query: this.$route.query})
          })
          .catch(() => {
            this.loading = false
          })
      } catch (error) {
        this.loading = false
        this.$message.error('YAML解析失败: ' + error.message)
      }
    },
    parseYamlToJson(yamlStr) {
      // 简单的YAML到JSON转换（实际项目中建议使用js-yaml库）
      try {
        const lines = yamlStr.split('\n')
        const result = {}
        let currentObj = result
        const stack = [result]
        let currentIndent = 0

        for (let line of lines) {
          line = line.replace(/^\s+/, match => match)
          if (line.trim() === '' || line.trim().startsWith('#')) continue

          const indent = line.length - line.trimLeft().length
          const content = line.trim()

          if (content.includes(':')) {
            const [key, value] = content.split(':').map(s => s.trim())
            if (value) {
              // 处理简单的键值对
              if (value.startsWith('"') && value.endsWith('"')) {
                currentObj[key] = value.slice(1, -1)
              } else if (!isNaN(value)) {
                currentObj[key] = Number(value)
              } else {
                currentObj[key] = value
              }
            } else {
              // 处理对象
              currentObj[key] = {}
              stack.push(currentObj[key])
              currentObj = currentObj[key]
            }
          } else if (content.startsWith('-')) {
            // 处理数组
            const value = content.substring(1).trim()
            if (!Array.isArray(currentObj)) {
              // 需要将当前对象转换为数组
              const parent = stack[stack.length - 2]
              const keys = Object.keys(parent)
              const lastKey = keys[keys.length - 1]
              parent[lastKey] = []
              currentObj = parent[lastKey]
            }
            currentObj.push(value)
          }
        }

        return result
      } catch (error) {
        throw new Error('YAML解析失败，请检查格式')
      }
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

<template>
  <layout-content :header="$t('business.istio.create_gateway')" :back-to="{name: 'Gateways', query: $route.query}">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="16">
        <div class="grid-content bg-purple-light">
          <!-- YAML创建 -->
          <el-card>
            <div slot="header">
              <span>{{ $t('business.istio.yaml_create') }}</span>
            </div>

            <div style="margin-bottom: 20px;">
              <el-button type="primary" size="small" @click="loadTemplate">{{ $t('business.istio.load_template') }}</el-button>
              <el-button type="success" size="small" @click="validateYaml">{{ $t('business.istio.validate_yaml') }}</el-button>
            </div>

            <yaml-editor ref="yaml_editor" :value="yamlData" :is-edit="false"></yaml-editor>

            <div style="text-align: right; margin-top: 20px;">
              <el-button @click="onCancel">{{ $t("commons.button.cancel") }}</el-button>
              <el-button type="primary" @click="onSubmitYaml" :loading="loading">
                {{ $t("commons.button.submit") }}
              </el-button>
            </div>
          </el-card>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import { createGateway } from "@/api/istio"

export default {
  name: "GatewayCreate",
  components: {
    LayoutContent,
    YamlEditor
  },
  data() {
    return {
      cluster: "",
      loading: false,
      yamlData: {}
    }
  },
  methods: {
    loadTemplate() {
      this.yamlContent = `apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: example-gateway
  namespace: default
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - example.com
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      mode: SIMPLE
      credentialName: example-tls
    hosts:
    - example.com`
    },
    validateYaml() {
      try {
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
        const yamlData = this.parseYamlToJson(this.yamlContent)

        createGateway(this.cluster, yamlData.metadata.namespace, yamlData)
          .then(() => {
            this.$message({
              type: "success",
              message: this.$t("commons.msg.create_success")
            })
            this.$router.push({name: "Gateways", query: this.$route.query})
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
      // 改进的YAML到JSON转换
      try {
        const lines = yamlStr.split('\n')
        const result = {}
        const stack = [{obj: result, indent: -1}]

        for (let i = 0; i < lines.length; i++) {
          const line = lines[i]
          if (line.trim() === '' || line.trim().startsWith('#')) continue

          const indent = line.length - line.trimLeft().length
          const content = line.trim()

          // 调整栈深度
          while (stack.length > 1 && stack[stack.length - 1].indent >= indent) {
            stack.pop()
          }

          const currentContext = stack[stack.length - 1]

          if (content.startsWith('- ')) {
            // 数组项
            const arrayItem = content.substring(2).trim()

            if (!Array.isArray(currentContext.obj)) {
              // 转换为数组
              const keys = Object.keys(currentContext.obj)
              if (keys.length > 0) {
                const lastKey = keys[keys.length - 1]
                currentContext.obj[lastKey] = []
                currentContext.obj = currentContext.obj[lastKey]
              }
            }

            if (arrayItem.includes(':')) {
              // 数组中的对象
              const newObj = {}
              currentContext.obj.push(newObj)
              stack.push({obj: newObj, indent: indent})

              const [key, value] = arrayItem.split(':').map(s => s.trim())
              if (value) {
                newObj[key] = this.parseValue(value)
              } else {
                newObj[key] = {}
                stack.push({obj: newObj[key], indent: indent + 2})
              }
            } else {
              // 简单数组项
              currentContext.obj.push(this.parseValue(arrayItem))
            }
          } else if (content.includes(':')) {
            // 键值对
            const colonIndex = content.indexOf(':')
            const key = content.substring(0, colonIndex).trim()
            const value = content.substring(colonIndex + 1).trim()

            if (value) {
              currentContext.obj[key] = this.parseValue(value)
            } else {
              currentContext.obj[key] = {}
              stack.push({obj: currentContext.obj[key], indent: indent})
            }
          }
        }

        return result
      } catch (error) {
        console.error('YAML解析错误:', error)
        throw new Error('YAML解析失败，请检查格式: ' + error.message)
      }
    },
    parseValue(value) {
      if (value.startsWith('"') && value.endsWith('"')) {
        return value.slice(1, -1)
      } else if (value.startsWith("'") && value.endsWith("'")) {
        return value.slice(1, -1)
      } else if (value === 'true') {
        return true
      } else if (value === 'false') {
        return false
      } else if (value === 'null') {
        return null
      } else if (!isNaN(value) && !isNaN(parseFloat(value))) {
        return parseFloat(value)
      } else {
        return value
      }
    },
    onCancel() {
      this.$router.push({name: "Gateways", query: this.$route.query})
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>
.yaml-editor-container {
  border: 1px solid #4a5568;
  border-radius: 4px;
  margin-bottom: 20px;
  min-height: 500px;
  background: #1a202c;
}

.yaml-editor {
  width: 100%;
  height: 500px;
  min-height: 500px;
  border: none;
  padding: 15px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  resize: vertical;
  outline: none;
  background: #1a202c;
  color: #e2e8f0;
  box-sizing: border-box;
  line-height: 1.5;
}

.yaml-editor:focus {
  border: none;
  outline: none;
  background: #1a202c;
}

.yaml-editor::placeholder {
  color: #718096;
}
</style>

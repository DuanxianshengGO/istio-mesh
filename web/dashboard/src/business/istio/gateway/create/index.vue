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
      yamlData: {},
      yamlContent: ""
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

      // 设置YAML编辑器的值
      this.$refs.yaml_editor.setValue(this.yamlContent)
    },
    validateYaml() {
      try {
        const yamlData = this.$refs.yaml_editor.getValue()

        if (!yamlData) {
          this.$message.error('YAML内容不能为空')
          return false
        }

        if (!yamlData.apiVersion || !yamlData.kind || !yamlData.metadata) {
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
      this.loading = true

      try {
        // 直接从YAML编辑器获取解析后的对象
        const yamlData = this.$refs.yaml_editor.getValue()

        if (!yamlData || !yamlData.metadata || !yamlData.metadata.namespace) {
          this.$message.error('YAML格式不正确，缺少必要的metadata.namespace字段')
          this.loading = false
          return
        }

        createGateway(this.cluster, yamlData.metadata.namespace, yamlData)
          .then(() => {
            this.$message({
              type: "success",
              message: this.$t("commons.msg.create_success")
            })
            this.$router.push({name: "Gateways", query: this.$route.query})
          })
          .catch((error) => {
            this.loading = false
            console.error('Gateway创建失败:', error)
            this.$message.error('创建失败: ' + (error.response?.data?.message || error.message))
          })
      } catch (error) {
        this.loading = false
        console.error('YAML解析失败:', error)
        this.$message.error('YAML解析失败: ' + error.message)
      }
    },
    onCancel() {
      this.$router.push({name: "Gateways", query: this.$route.query})
    },

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

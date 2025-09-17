<template>
  <layout-content :header="$t('business.istio.create_destinationrule')" :back-to="{name: 'DestinationRules', query: $route.query}">
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
import { createDestinationRule } from "@/api/istio"

export default {
  name: "DestinationRuleCreate",
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
      this.yamlData = {
        apiVersion: "networking.istio.io/v1beta1",
        kind: "DestinationRule",
        metadata: {
          name: "example-dr",
          namespace: "default"
        },
        spec: {
          host: "example-service",
          trafficPolicy: {
            loadBalancer: {
              simple: "ROUND_ROBIN"
            }
          },
          subsets: [
            {
              name: "v1",
              labels: {
                version: "v1"
              }
            },
            {
              name: "v2",
              labels: {
                version: "v2"
              }
            }
          ]
        }
      }
    },
    validateYaml() {
      try {
        const data = this.$refs.yaml_editor.getValue()

        if (!data) {
          this.$message.error('YAML内容不能为空')
          return false
        }

        // 基本验证
        if (!data.apiVersion || !data.kind || !data.metadata || !data.spec) {
          this.$message.error('YAML格式不完整，缺少必要字段')
          return false
        }

        if (data.kind !== 'DestinationRule') {
          this.$message.error('资源类型必须是DestinationRule')
          return false
        }

        if (!data.metadata.name) {
          this.$message.error('metadata.name不能为空')
          return false
        }

        this.$message.success('YAML格式验证通过')
        return true
      } catch (error) {
        this.$message.error('YAML验证失败: ' + error.message)
        return false
      }
    },
    onSubmitYaml() {
      if (!this.validateYaml()) {
        return
      }

      this.loading = true

      try {
        const data = this.$refs.yaml_editor.getValue()
        const namespace = data.metadata.namespace || 'default'

        createDestinationRule(this.cluster, namespace, data)
          .then(() => {
            this.$message({
              type: "success",
              message: this.$t("commons.msg.create_success")
            })
            this.$router.push({name: "DestinationRules", query: this.$route.query})
          })
          .catch(error => {
            console.error('创建失败:', error)
            this.$message.error('创建失败: ' + (error.response?.data?.message || error.message))
          })
          .finally(() => {
            this.loading = false
          })
      } catch (error) {
        this.loading = false
        this.$message.error('YAML解析失败: ' + error.message)
      }
    },

    onCancel() {
      this.$router.push({name: "DestinationRules", query: this.$route.query})
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>
</style>

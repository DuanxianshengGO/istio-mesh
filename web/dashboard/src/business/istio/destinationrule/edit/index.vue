<template>
  <layout-content :header="$t('business.istio.edit_destinationrule')" :back-to="{name: 'DestinationRules', query: $route.query}">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="16">
        <div class="grid-content bg-purple-light">
          <yaml-editor ref="yaml_editor" :value="yaml" :is-edit="true"></yaml-editor>
          <div style="margin-top: 10px;">
            <el-button type="primary" @click="onSubmit" :loading="loading">
              {{ $t("commons.button.submit") }}
            </el-button>
            <el-button @click="onCancel">{{ $t("commons.button.cancel") }}</el-button>
          </div>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import { getDestinationRule, updateDestinationRule } from "@/api/istio"

export default {
  name: "DestinationRuleEdit",
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
      item: {},
      yaml: {}
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getDestinationRule(this.cluster, this.namespace, this.name)
        .then(res => {
          this.item = res
          this.yaml = JSON.parse(JSON.stringify(this.item))
          this.loading = false
        })
        .catch(() => {
          this.loading = false
        })
    },
    onCancel() {
      this.$router.push({name: "DestinationRules", query: this.$route.query})
    },
    onSubmit() {
      this.loading = true
      const data = this.$refs.yaml_editor.getValue()
      updateDestinationRule(this.cluster, this.namespace, this.name, data)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success")
          })
          this.$router.push({name: "DestinationRules", query: this.$route.query})
        })
        .finally(() => {
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

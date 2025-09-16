<template>
  <layout-content :header="$t('business.istio.create_destinationrule')" :back-to="{name: 'DestinationRules', query: $route.query}">
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
            
            <el-form-item :label="$t('business.istio.host')" prop="host" required>
              <el-input v-model="form.host" placeholder="service-name" clearable></el-input>
            </el-form-item>
            
            <el-form-item :label="$t('business.istio.load_balancer')">
              <el-select v-model="form.loadBalancer" placeholder="Select load balancer">
                <el-option label="ROUND_ROBIN" value="ROUND_ROBIN"></el-option>
                <el-option label="LEAST_CONN" value="LEAST_CONN"></el-option>
                <el-option label="RANDOM" value="RANDOM"></el-option>
                <el-option label="PASSTHROUGH" value="PASSTHROUGH"></el-option>
              </el-select>
            </el-form-item>
            
            <el-form-item :label="$t('business.istio.subsets')">
              <div v-for="(subset, index) in form.subsets" :key="index" style="border: 1px solid #ddd; padding: 10px; margin-bottom: 10px;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                  <span>{{ $t('business.istio.subset') }} {{ index + 1 }}</span>
                  <el-button size="mini" type="danger" @click="removeSubset(index)">
                    {{ $t('commons.button.delete') }}
                  </el-button>
                </div>
                
                <el-form-item :label="$t('commons.table.name')" required>
                  <el-input v-model="subset.name" placeholder="v1"></el-input>
                </el-form-item>
                
                <el-form-item :label="$t('business.common.label')">
                  <div v-for="(label, labelIndex) in subset.labels" :key="labelIndex" style="display: flex; margin-bottom: 5px;">
                    <el-input v-model="label.key" placeholder="Key" style="margin-right: 10px;"></el-input>
                    <el-input v-model="label.value" placeholder="Value" style="margin-right: 10px;"></el-input>
                    <el-button size="mini" type="danger" @click="removeLabel(index, labelIndex)">
                      {{ $t('commons.button.delete') }}
                    </el-button>
                  </div>
                  <el-button size="small" @click="addLabel(index)">
                    {{ $t('business.common.add_label') }}
                  </el-button>
                </el-form-item>
              </div>
              
              <el-button type="primary" size="small" @click="addSubset">
                {{ $t('business.istio.add_subset') }}
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
import { createDestinationRule } from "@/api/istio"

export default {
  name: "DestinationRuleCreate",
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
        host: "",
        loadBalancer: "ROUND_ROBIN",
        subsets: [
          {
            name: "v1",
            labels: [
              { key: "version", value: "v1" }
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
        host: [
          { required: true, message: this.$t("business.istio.rule.host"), trigger: "blur" }
        ]
      }
    }
  },
  methods: {
    addSubset() {
      this.form.subsets.push({
        name: "",
        labels: [
          { key: "", value: "" }
        ]
      })
    },
    removeSubset(index) {
      if (this.form.subsets.length > 1) {
        this.form.subsets.splice(index, 1)
      }
    },
    addLabel(subsetIndex) {
      this.form.subsets[subsetIndex].labels.push({ key: "", value: "" })
    },
    removeLabel(subsetIndex, labelIndex) {
      if (this.form.subsets[subsetIndex].labels.length > 1) {
        this.form.subsets[subsetIndex].labels.splice(labelIndex, 1)
      }
    },
    onSubmit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          
          const destinationRule = {
            apiVersion: "networking.istio.io/v1beta1",
            kind: "DestinationRule",
            metadata: {
              name: this.form.name,
              namespace: this.form.namespace
            },
            spec: {
              host: this.form.host,
              trafficPolicy: {
                loadBalancer: {
                  simple: this.form.loadBalancer
                }
              },
              subsets: this.form.subsets.map(subset => ({
                name: subset.name,
                labels: subset.labels.reduce((acc, label) => {
                  if (label.key && label.value) {
                    acc[label.key] = label.value
                  }
                  return acc
                }, {})
              })).filter(subset => subset.name && Object.keys(subset.labels).length > 0)
            }
          }
          
          createDestinationRule(this.cluster, this.form.namespace, destinationRule)
            .then(() => {
              this.$message({
                type: "success",
                message: this.$t("commons.msg.create_success")
              })
              this.$router.push({name: "DestinationRules", query: this.$route.query})
            })
            .finally(() => {
              this.loading = false
            })
        }
      })
    },
    onCancel() {
      this.$router.push({name: "DestinationRules", query: this.$route.query})
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

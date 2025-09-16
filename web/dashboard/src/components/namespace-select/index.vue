<template>
  <el-select v-model="selectedNamespace" v-bind="$attrs" v-on="$listeners" :disabled="disabled" @change="onChange">
    <el-option v-for="namespace in namespaces"
               :key="namespace"
               :label="namespace"
               :value="namespace">
    </el-option>
  </el-select>
</template>

<script>
import { getNamespaces } from "@/api/auth"

export default {
  name: "NamespaceSelect",
  props: {
    value: String,
    cluster: String
  },
  data() {
    return {
      selectedNamespace: "",
      namespaces: [],
      disabled: false
    }
  },
  methods: {
    initData(namespaces) {
      this.namespaces = namespaces
      if (this.value !== "") {
        this.selectedNamespace = this.value
        return
      }
      this.selectedNamespace = namespaces.find(item => item === "default")
      if (this.selectedNamespace === undefined || this.selectedNamespace === "") {
        this.selectedNamespace = namespaces[0]
      }
      this.$emit("input", this.selectedNamespace)
      this.$emit("change", this.selectedNamespace)
    },
    onChange(value) {
      this.$emit("input", value)
      this.$emit("change", value)
    },
    loadNamespaces() {
      if (!this.cluster) {
        return
      }
      
      if (!sessionStorage.getItem("namespace")) {
        getNamespaces(this.cluster).then(res => {
          this.initData(res.data)
        })
        this.disabled = false
      } else {
        this.disabled = true
        this.initData([sessionStorage.getItem("namespace")])
      }
    }
  },
  watch: {
    value(newValue) {
      this.selectedNamespace = newValue
    },
    cluster() {
      this.loadNamespaces()
    }
  },
  created() {
    this.selectedNamespace = this.value
    this.loadNamespaces()
  }
}
</script>

<style scoped>
</style>

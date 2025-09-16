<template>
  <layout-content :header="$t('business.istio.destinationrule_list')">
    <div style="float: left">
      <el-button 
        v-has-permissions="{apiGroup:'networking.istio.io',resource:'destinationrules',verb:'create'}" 
        type="primary" 
        size="small" 
        @click="onCreate">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button 
        v-has-permissions="{apiGroup:'networking.istio.io',resource:'destinationrules',verb:'delete'}" 
        :disabled="selects.length===0" 
        type="primary" 
        size="small" 
        @click="onDelete()">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>

    <complex-table 
      :search-config="searchConfig" 
      :selects.sync="selects" 
      :data="data"
      :pagination-config="paginationConfig" 
      @search="search">
      
      <el-table-column type="selection" fix></el-table-column>
      
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" min-width="100" fix>
        <template v-slot:default="{row}">
          <el-link @click="onDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      
      <el-table-column :label="$t('commons.table.namespace')" prop="metadata.namespace" min-width="100">
        <template v-slot:default="{row}">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      
      <el-table-column :label="$t('business.istio.host')" min-width="150">
        <template v-slot:default="{row}">
          {{ row.spec ? row.spec.host : '' }}
        </template>
      </el-table-column>
      
      <el-table-column :label="$t('business.istio.subsets')" min-width="200">
        <template v-slot:default="{row}">
          <div v-if="row.spec && row.spec.subsets">
            <el-tag v-for="subset in row.spec.subsets" :key="subset.name" size="mini" style="margin: 2px;">
              {{ subset.name }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column :label="$t('business.istio.load_balancer')" min-width="120">
        <template v-slot:default="{row}">
          {{ row.spec && row.spec.trafficPolicy && row.spec.trafficPolicy.loadBalancer ? 
               row.spec.trafficPolicy.loadBalancer.simple : 'ROUND_ROBIN' }}
        </template>
      </el-table-column>
      
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" min-width="100">
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | dateFormat }}
        </template>
      </el-table-column>
      
      <ko-table-operations 
        :buttons="buttons" 
        :label="$t('commons.table.action')">
      </ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import { listDestinationRules, deleteDestinationRule } from "@/api/istio"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "DestinationRuleList",
  components: { 
    LayoutContent, 
    ComplexTable, 
    KoTableOperations 
  },
  data() {
    return {
      data: [],
      selects: [],
      cluster: "",
      loading: false,
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.onEdit(row)
          },
          disabled: (row) => {
            return !checkPermissions({
              apiGroup: "networking.istio.io",
              resource: "destinationrules", 
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: (row) => {
            return !checkPermissions({
              apiGroup: "networking.istio.io",
              resource: "destinationrules", 
              verb: "delete"
            })
          }
        }
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      }
    }
  },
  methods: {
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      const {currentPage, pageSize} = this.paginationConfig
      listDestinationRules(this.cluster, this.$route.query.namespace, true, this.searchConfig.keywords, currentPage, pageSize)
        .then(response => {
          this.loading = false
          // 处理 KubePi API 返回的数据结构
          if (response && response.data && response.data.items) {
            this.data = response.data.items
            this.paginationConfig.total = response.data.items.length
          } else {
            this.data = []
            this.paginationConfig.total = 0
          }
        })
        .catch(() => {
          this.loading = false
          this.data = []
          this.paginationConfig.total = 0
        })
    },
    onCreate() {
      this.$router.push({name: "DestinationRuleCreate", query: this.$route.query})
    },
    onEdit(row) {
      this.$router.push({
        name: "DestinationRuleEdit", 
        params: {
          namespace: row.metadata.namespace,
          name: row.metadata.name
        },
        query: this.$route.query
      })
    },
    onDetail(row) {
      this.$router.push({
        name: "DestinationRuleDetail", 
        params: {
          namespace: row.metadata.namespace,
          name: row.metadata.name
        },
        query: this.$route.query
      })
    },
    onDelete(row) {
      let items = []
      if (row) {
        items = [row]
      } else {
        items = this.selects
      }
      
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning"
      }).then(() => {
        this.loading = true
        const promises = items.map(item => 
          deleteDestinationRule(this.cluster, item.metadata.namespace, item.metadata.name)
        )
        
        Promise.all(promises).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.delete_success")
          })
          this.search()
        }).catch(() => {
          this.loading = false
        })
      })
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>
</style>

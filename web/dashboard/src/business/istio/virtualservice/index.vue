<template>
  <layout-content :header="$t('business.istio.virtualservice_list')">
    <div style="float: left">
      <el-button 
        v-has-permissions="{apiGroup:'networking.istio.io',resource:'virtualservices',verb:'create'}" 
        type="primary" 
        size="small" 
        @click="onCreate">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button 
        v-has-permissions="{apiGroup:'networking.istio.io',resource:'virtualservices',verb:'delete'}" 
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
      
      <el-table-column :label="$t('business.istio.hosts')" min-width="150">
        <template v-slot:default="{row}">
          <div v-if="row.spec && row.spec.hosts">
            <el-tag v-for="host in row.spec.hosts" :key="host" size="mini" style="margin: 2px;">
              {{ host }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column :label="$t('business.istio.gateways')" min-width="150">
        <template v-slot:default="{row}">
          <div v-if="row.spec && row.spec.gateways">
            <el-tag v-for="gateway in row.spec.gateways" :key="gateway" size="mini" style="margin: 2px;">
              {{ gateway }}
            </el-tag>
          </div>
        </template>
      </el-table-column>

      <el-table-column :label="$t('business.istio.basic_subsets')" min-width="120">
        <template v-slot:default="{row}">
          <div v-if="getBasicSubsets(row).length > 0">
            <el-tag v-for="subset in getBasicSubsets(row)" :key="subset" size="mini" type="success" style="margin: 2px;">
              {{ subset }}
            </el-tag>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column :label="$t('business.istio.gray_subsets')" min-width="120">
        <template v-slot:default="{row}">
          <div v-if="getGraySubsets(row).length > 0">
            <el-tag v-for="subset in getGraySubsets(row)" :key="subset" size="mini" type="warning" style="margin: 2px;">
              {{ subset }}
            </el-tag>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column :label="$t('business.istio.match_content')" min-width="200">
        <template v-slot:default="{row}">
          <div v-if="getMatchContent(row)">
            <el-tooltip :content="getMatchContent(row)" placement="top">
              <span class="match-content-text">{{ getMatchContent(row) }}</span>
            </el-tooltip>
          </div>
          <span v-else>-</span>
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
import { listVirtualServices, deleteVirtualService } from "@/api/istio"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "VirtualServiceList",
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
              resource: "virtualservices", 
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
              resource: "virtualservices", 
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
      listVirtualServices(this.cluster, this.$route.query.namespace, true, this.searchConfig.keywords, currentPage, pageSize)
        .then(response => {
          this.loading = false
          console.log('VirtualService API Response:', response)
          // 处理 KubePi API 返回的数据结构
          if (response && response.data && response.data.items) {
            console.log('Found items:', response.data.items.length)
            this.data = response.data.items
            this.paginationConfig.total = response.data.items.length
          } else {
            console.log('No items found, response structure:', response)
            this.data = []
            this.paginationConfig.total = 0
          }
        })
        .catch(error => {
          console.error('VirtualService API Error:', error)
          this.loading = false
          this.data = []
          this.paginationConfig.total = 0
        })
    },
    onCreate() {
      this.$router.push({name: "VirtualServiceCreate", query: this.$route.query})
    },
    onEdit(row) {
      this.$router.push({
        name: "VirtualServiceEdit", 
        params: {
          namespace: row.metadata.namespace,
          name: row.metadata.name
        },
        query: this.$route.query
      })
    },
    onDetail(row) {
      this.$router.push({
        name: "VirtualServiceDetail", 
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
          deleteVirtualService(this.cluster, item.metadata.namespace, item.metadata.name)
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
    },
    // 获取基础路由的子集
    getBasicSubsets(row) {
      const subsets = new Set()
      if (row.spec && row.spec.http) {
        row.spec.http.forEach(httpRoute => {
          // 没有match条件的是基础路由
          if (!httpRoute.match && httpRoute.route) {
            httpRoute.route.forEach(route => {
              if (route.destination && route.destination.subset) {
                subsets.add(route.destination.subset)
              }
            })
          }
        })
      }
      return Array.from(subsets)
    },
    // 获取灰度路由的子集
    getGraySubsets(row) {
      const subsets = new Set()
      if (row.spec && row.spec.http) {
        row.spec.http.forEach(httpRoute => {
          // 有match条件的是灰度路由
          if (httpRoute.match && httpRoute.route) {
            httpRoute.route.forEach(route => {
              if (route.destination && route.destination.subset) {
                subsets.add(route.destination.subset)
              }
            })
          }
        })
      }
      return Array.from(subsets)
    },
    // 获取匹配内容
    getMatchContent(row) {
      const matchContents = []
      if (row.spec && row.spec.http) {
        row.spec.http.forEach(httpRoute => {
          if (httpRoute.match) {
            httpRoute.match.forEach(match => {
              const conditions = []
              if (match.uri) {
                Object.keys(match.uri).forEach(key => {
                  conditions.push(`URI ${key}: ${match.uri[key]}`)
                })
              }
              if (match.headers) {
                Object.keys(match.headers).forEach(headerName => {
                  const headerValue = match.headers[headerName]
                  if (typeof headerValue === 'object') {
                    Object.keys(headerValue).forEach(key => {
                      conditions.push(`Header ${headerName} ${key}: ${headerValue[key]}`)
                    })
                  } else {
                    conditions.push(`Header ${headerName}: ${headerValue}`)
                  }
                })
              }
              if (conditions.length > 0) {
                matchContents.push(conditions.join(', '))
              }
            })
          }
        })
      }
      return matchContents.join('; ')
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>
.match-content-text {
  display: inline-block;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}
</style>

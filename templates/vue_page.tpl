<template>
  <el-card class="box-card">
    <el-form :inline="true" :model="search">
      <el-form-item>
        <el-input
          v-model="search.keyword"
          placeholder="编号/名称/联系人/联系电话"
          clearable
          style="width:260px;"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="fetch${modelName}s" icon="el-icon-search">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="handleAdd" type="primary" icon="el-icon-plus">新增</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="${lowerModelName}s" border>
      <el-table-column type="index" align="center"></el-table-column>${tableColumnContents}
      <el-table-column prop="account" label="管理员账号" align="center"></el-table-column>
      <el-table-column prop="created_at" label="创建时间" align="center" :formatter="dateFormatter"></el-table-column>
      <el-table-column label="操作" align="center">
        <template slot-scope="scope">
          <el-button @click="handleEdit(scope.row)" type="text">编辑</el-button>
          <el-button @click="handleResetPassword(scope.row)" type="text">重置密码</el-button>
          <el-button @click="handleDel(scope.row)" type="text" style="color:red">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      style="margin-top:10px"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page.sync="search.page"
      :page-size="search.limit"
      :page-sizes="page.sizes"
      :total="page.total"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>
    <el-dialog
      :title="dialog.title"
      :visible.sync="dialog.show"
      width="40%"
      @close="closeDialog">
      <el-form label-width="100px" :model="form" :rules="rules" ref="form">${formContents}
        <el-form-item label="机构类型" prop="org_type_id">
            <el-select v-model="form.org_type_id" placeholder="请选择机构类型" style="width:60%;">
                <el-option v-for="item in org_types" :key="item.id" :label="item.name" :value="item.id"></el-option>
            </el-select>
        </el-form-item>
        <div>
          <el-form-item label="管理员账号" prop="account" v-if="!form.id">
            <el-input v-model="form.account" placeholder="请输入管理员账号" maxlength="20" class="form-item"></el-input>
          </el-form-item>
          <el-form-item label="密码" prop="password" v-if="!form.id">
            <el-input v-model="form.password" placeholder="请输入密码" class="form-item" type="password"></el-input>
          </el-form-item>
        </div>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialog.show = false">取 消</el-button>
        <el-button type="primary" @click="handleSubmit">保 存</el-button>
      </span>
    </el-dialog>
  </el-card>
</template>
<script>
import {
  findOrgTypesSelect,
  find${modelName}s,
  save${modelName},
  delete${modelName},
  saveEmployee,
  resetPasswordEmployee
} from "@/api/api"
import { LocalAccount } from "@/api/local-account"
import _ from "lodash"
import moment from "moment"

export default {
  data() {
    return {
      ${lowerModelName}s: [],
      page: { total: 0, sizes: [10, 20, 30, 50] },
      search: {
        page: 1,
        limit: 10
      },
	  params:{},
      form: {},
      org_types: [],
      dialog: {
        show: false,
        title: ""
      },
      rules: {
        ${ruleContents}
        org_type_id: [{ required: true, message: "机构类型不能为空", trigger: "blur" }],
        account: [{ required: true, message: "管理员账号不能为空", trigger: "blur" }],
        password: [{ required: true, message: "密码不能为空", trigger: "blur" }],
      }
    }
  },
  created() {
    this.fetchOrgTypes()
    this.fetch${modelName}s()
  },
  methods: {
    fetch${modelName}s() {
	  var storage = JSON.parse(LocalAccount.get("webadmin_account"))
	  this.params.search = this.search
	  this.params.parent_id = storage.user.OrgId
	  find${modelName}s(this.params).then(result => {
	  this.${lowerModelName}s = result.data
	  this.page.total = result.total
      })
    },
    fetchOrgTypes() {
        findOrgTypesSelect().then(result => {
            this.org_types = result.data
        })
    },
    closeDialog() {
      this.$refs.form.clearValidate()
    },
    handleAdd() {
      this.dialog.show = true
      this.dialog.title = "新增"
      this.form = {}
    },
    handleEdit(item) {
      this.form = _.cloneDeep(item)
      this.dialog.show = true
      this.dialog.title = "编辑"
    },
    handleSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
            let org_type = _.find(this.org_types, {id: this.form.org_type_id})
            var storage = JSON.parse(LocalAccount.get("webadmin_account"))
            this.form.parent_id = storage.user.OrgId
            if (!org_type) return
            this.form.org_type_name = org_type.name
            save${modelName}(this.form).then(res => {
            if (res.code == 0) {
              this.$message.success("已保存")
              this.fetch${modelName}s()
              this.dialog.show = false
            }
          })
        } else {
          console.log("error submit!!")
          return false
        }
      })
    },
    handleResetPassword(row) {
      this.$confirm("您确定重置"+row.account+"的密码为111111?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(() => {
        resetPasswordEmployee({ id: row.employee_id, password: "111111" }).then(
          res => {
            if (res.code == 0) {
              this.$message.success("已重置密码")
            }
          }
        )
      }).catch(() => {
        //取消
        console.log("n")
      })
    },
    handleDel(row) {
      this.$confirm("您确定要删除"+row.name+"吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(() => {
        delete${modelName}([row.id]).then(result => {
          if (result.code == 0) {
            this.fetch${modelName}s()
            this.$message.success("已删除")
          }
        })
      }).catch(() => {
        //取消
        console.log("n")
      })
    },
    handleSizeChange(val) {
      this.search.limit = val
      this.fetch${modelName}s()
    },
    handleCurrentChange(val) {
      this.fetch${modelName}s()
    },
    dateFormatter(row, col) {
      return moment(row[col.property]).format("YYYY-MM-DD HH:mm:ss")
    }
  }
}
</script>
<style scoped>
.pagination {
  padding: 10px;
  float: right;
}
.form-item {
  width: 60%;
}
</style>

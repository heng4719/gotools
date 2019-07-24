<template>
  <el-menu class="navbar" mode="horizontal">
    <!-- <div style="height:10px;"></div> -->
      <hamburger class="hamburger-container" :toggleClick="toggleSideBar" :isActive="sidebar.opened"></hamburger>
      <!-- <breadcrumb class="breadcrumb-container"></breadcrumb> -->
      <div class="right-menu">
        <!-- <error-log class="errLog-container right-menu-item"></error-log> -->
        <div class="shop-title"><span class="title">机构名称：</span>{{employee.mechanism_name}}</div>
        <div class="shop-title"><span class="title">机构类型：</span>{{brandName}}</div>
        <!-- <div class="shop-title">{{brandCode}}</div> -->

        <!--<el-tooltip effect="light" :content="$t('navbar.screenfull')" placement="bottom">
          <screenfull class="screenfull right-menu-item"></screenfull>
        </el-tooltip>-->

        <div class="shop-title"><span class="title">登录账号：</span>{{username}}</div>
        <el-dropdown class="avatar-container right-menu-item" trigger="click">
          <div class="avatar-wrapper">
            <img class="user-avatar" :src="avatar+'?imageView2/1/w/80/h/80'">
            <i class="el-icon-caret-bottom"></i>
          </div>
          <el-dropdown-menu slot="dropdown">
            <!-- <router-link to="/">
              <el-dropdown-item>
                首 页
              </el-dropdown-item>
            </router-link> -->
            <el-dropdown-item >
              <span @click="resetPassword" style="display:block;">修改密码</span>
            </el-dropdown-item>
            <el-dropdown-item divided>
              <span @click="logout" style="display:block;">退 出</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>

      <el-dialog
        title="提示"
        :visible.sync="dialogVisible"
        width="40%"
        >
          <el-form
            :model="employeeForm"
            label-width="90px"
            :rules="editFormRules"
            ref="employeeForm"
            class="form"
          >
            <el-form-item label="用户" prop="account" hidden>
              <el-input type="text" v-model="employeeForm.account" placeholder="用户" ></el-input>
            </el-form-item>
            <el-form-item label="原密码" prop="old_password">
              <el-input type="password" v-model="employeeForm.old_password" placeholder="请输入原密码"></el-input>
            </el-form-item>
            <el-form-item label="新密码" prop="password">
              <el-input type="password" v-model="employeeForm.password" placeholder="请输入新密码"></el-input>
            </el-form-item>
            <el-form-item label="确认密码" prop="password2">
              <el-input type="password" v-model="employeeForm.password2" placeholder="请输入确认密码"></el-input>
            </el-form-item>
          </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="submit">确 定</el-button>
        </span>
      </el-dialog>
  </el-menu>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import ErrorLog from '@/components/ErrorLog'
import Screenfull from '@/components/Screenfull'
import LangSelect from '@/components/LangSelect'
import ThemePicker from '@/components/ThemePicker'
import {LocalAccount} from '@/api/local-account'
import { updatePasswordEmployee } from "@/api/api";

export default {
  components: {
    Breadcrumb,
    Hamburger,
    ErrorLog,
    Screenfull,
    LangSelect,
    ThemePicker
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'name',
    ]),
    brandName() {
      return ''
    },
    brandCode() {
      if (LocalAccount.getBrandId() !== 0) {
        return ''
      }
    },
    username() {
      return LocalAccount.getUserInfo().account
    }
  },
  data() {
    var validatePass = (rule, value, callback) => {
      if (value == "") {
        callback(new Error("新密码不能为空"))
      } else {
        if (this.employeeForm.password2 !== "") {
          this.$refs.employeeForm.validateField("password2")
        }
        callback()
      }
    }
    var validatePass2 = (rule, value, callback) => {
      if (value == "") {
        callback(new Error("确认密码不能为空"))
      } else if (value !== this.employeeForm.password) {
        callback(new Error("两次输入密码不一致!"))
      } else {
        callback()
      }
    }
    return {
      avatar: 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif',
      dialogVisible:false,
      employeeForm: {},
      editFormRules: {
        old_password: [
          { required: true, message: "原密码不能为空", trigger: "blur" },
          { min: 6, message: '密码长度不能低于6位', trigger: 'blur' }
        ],
        password: [
          { validator: validatePass, trigger: "blur" },
          { min: 6, message: '密码长度不能低于6位', trigger: 'blur' }
        ],
        password2: [{ validator: validatePass2, trigger: "blur" }]
      },
      employee: LocalAccount.getUserInfo(),
    }
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('toggleSideBar')
    },
    logout() {
      this.$store.dispatch('LogOut').then(() => {
        this.$store.dispatch('setIsLoadedRoutes', false)
        location.reload()
      })
    },
    resetPassword(){
      this.dialogVisible=true
      this.employeeForm = {account: LocalAccount.getUserInfo().account}
    },
    submit() {
      this.$refs.employeeForm.validate(valid => {
        if (valid) {
          updatePasswordEmployee(this.employeeForm).then(result => {
            if (result.code == 0) {
              this.$message.success('密码已修改,请重新登录！')
              setTimeout(() => {
                this.$store.dispatch('LogOut').then(() => {
                  this.$store.dispatch('setIsLoadedRoutes', false)
                  location.reload()
                })
              }, 2000)
            }
          })
        }
      })
    },
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.navbar {
  padding-top:5px;
  padding-bottom:5px;
  background-color:rgb(48,65,85);
  height: 60px;
  line-height: 60px;
  vertical-align: middle;
  border-radius: 0px !important;
  .shop-title {
    height: 50px;
    line-height: 50px;
    display: inline-block;
    vertical-align: top;
    margin-left: 10px;
    font-size: 15px;
    color: #fff;
  }
  .hamburger-container {
    margin-top:5px;
    line-height: 50px;
    height: 50px;
    vertical-align: middle;
    float: left;
    padding: 0 10px;
  }
  .breadcrumb-container{
    float: left;
    font-size: 15px;
  }
  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }
  .right-menu {
    float: right;
    height: 50px;
    &:focus{
     outline: none;
    }
    .right-menu-item {
      display: inline-block;
      margin: 0 8px;
    }
    .screenfull {
      height: 20px;
    }
    .international{
      vertical-align: top;
    }
    .theme-switch {
      vertical-align: 15px;
    }
    .avatar-container {
      height: 50px;
      margin-right: 30px;
      .avatar-wrapper {
        cursor: pointer;
        margin-top: 5px;
        position: relative;
        .user-avatar {
          width: 40px;
          height: 40px;
          border-radius: 10px;
        }
        .el-icon-caret-bottom {
          position: absolute;
          right: -20px;
          top: 25px;
          font-size: 12px;
        }
      }
    }
  }
}

.title{
  font-size: 14px;
  color:rgb(151, 168, 190);
}
</style>

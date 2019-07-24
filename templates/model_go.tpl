package model

import (
	"fmt"
<<<<<<< HEAD:templates/model_go.tpl
	"${project}/internal/common/DB"
	"${project}/internal/common/casbins"
	"${project}/internal/common/enum/employee_enum"
	"${project}/internal/model/base"
	"${project}/internal/model/sys"
	"${project}/pkg/security"
=======
	"goadmin/internal/common/DB"
	"goadmin/internal/common/casbins"
	"goadmin/internal/common/enum/employee_enum"
	"goadmin/internal/model/sys"
	"goadmin/pkg/security"
>>>>>>> cd68f6bb8145f7b0dbdc0161f4331a570cd0bc9c:templates/project/internal/model/model.go
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/gogf/gf/g/util/gconv"
	"github.com/xormplus/xorm"
)

func Init() {
	//需要同步的表结构
	if err := DB.Engine.Sync2(
		new(sys.Employee),
		new(sys.Role),
		new(sys.OrgType),
		new(sys.SysMenu),
		new(sys.RoleMenu),
		${newModels}
	); err != nil {
		panic(err)
	}

	// 初始化数据
	if err := initData(); err != nil {
		panic(err)
	}
}

func initData() (err error) {
	if err = initRole(); err != nil {
		return err
	}

	if err = initAccount(); err != nil {
		return err
	}

	return
}

func initRole() (err error) {
	count, err := DB.Where("id = 1").Count(&sys.Role{})
	if err != nil {
		return fmt.Errorf("init superadmin role err: %v\n", err)
	}

	if count > 0 {
		return nil
	}

	role := &sys.Role{Id: 1, Code: "1001", IsAdmin: 1, Name: "超级管理员", Buildin: 1}
	_, err = DB.Insert(role)
	return err
}

// 初始化超级管理员账号
func initAccount() (err error) {
	count, err := DB.Where("account=?", "super_admin").Count(&sys.Employee{})
	if err != nil {
		glog.Fatalf("init superadmin account err: %v\n", err)
		panic(err)
	}
	if count > 0 {
		return
	}

	session := DB.Engine.NewSession()
	defer session.Close()

	if err = session.Begin(); err != nil {
		return fmt.Errorf("session begin err: %s", err)
	}

	password := security.MD5Password("111111")
	employee := &sys.Employee{
		Id:       1,
		Name:     "超级管理员",
		Account:  "super_admin",
		Password: password,
		Phone:    "18601694368",
		RoleId:   1,
		Code:     "1000",
		RoleName: "超级管理员",
		Type:     employee_enum.ADMIN,
	}
	if _, err = DB.InsertTx(session, employee); err != nil {
		session.Rollback()
		return err
	}

	apiUrls := make([]*sys.ApiUrl, 0)
	if err = initSuperAdminPermsission(session, apiUrls); err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	return err
}

// 超级管理员初始权限
func initSuperAdminPermsission(session *xorm.Session, apiUrls []*sys.ApiUrl) (err error) {
	// 超级管理员初始权限
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/sys_menu/query")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/sys_menu/save")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/role/query")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/role/save")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/role_menu/query")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/role_menu/save")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/org_type/query_by_tree")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/org_type/query_by_select")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/org_type/save")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/employee/query")
	apiUrls = appendApiUrl(apiUrls, "/v1/admin_api/employee/save")

	// 先删除角色所有权限，然后重新将角色权限持久化到casbin
	//casbins.RemoveRolePolicy("p", gconv.String(1))
	for _, apiUrl := range apiUrls {
		casbin := casbins.CasbinModel{
			Ptype:  "p",
			Role:   gconv.String(1),
			Path:   apiUrl.ApiUrl,
			Method: "*",
		}

		if err = casbins.AddCasbin(&casbin); err != nil {
			return err
		}
	}

	//在role_menu表中给超管添加默认记录
	//判断是否已经存在数据
	count, err := DB.Where("id = 1").Count(sys.RoleMenu{})
	if count > 0 {
		return err
	}

	sysMenus := []sys.SysMenu{}
	fmt.Println(apiUrls[0].ApiUrl)

	FindSysMenu:
	for i := 0; i < len(apiUrls); i++{
		sysMenu := sys.SysMenu{}
		has, err := DB.Where("api_urls like ?", "%"+apiUrls[i].ApiUrl+"%").Get(&sysMenu)
		if err != nil {
			return err
		}

		if !has {
			continue FindSysMenu
		}
		//去重
		for i := 0; i < len(sysMenus); i++{
			if sysMenu.Id == sysMenus[i].Id {
					continue FindSysMenu
			}
		}
		sysMenus = append(sysMenus, sysMenu)
	}

	roleMenu := new(sys.RoleMenu)
	roleMenu.Id = 0
	roleMenu.RoleId = 1
	roleMenu.Checked = 1
	roleMenu.CreatedAt = time.Now()
	roleMenu.UpdatedAt = time.Now()

	for i := 0; i < len(sysMenus); i++{
		roleMenu.MenuId = sysMenus[i].Id
		roleMenu.Id++

		if _, err = DB.InsertTx(session, roleMenu); err != nil {
			session.Rollback()
			return err
		}
	}
	return err
}

func UpdateSuperAdminPermsission(session *xorm.Session, apiUrls []*sys.ApiUrl) (err error) {
	if err = casbins.RemoveRolePolicy(gconv.String(1)); err != nil {
		return err
	}

	err = initSuperAdminPermsission(session, apiUrls)
	return err
}

func existsApiUrl(url string, apiUrls []*sys.ApiUrl) bool {
	for _, item := range apiUrls {
		if url == strings.TrimSpace(item.ApiUrl) {
			fmt.Printf("url: %s\n", url)
			fmt.Printf("item.ApiUrl: %s\n", item.ApiUrl)
			return true
		}
	}

	return false
}

func appendApiUrl(apiUrls []*sys.ApiUrl, url string) []*sys.ApiUrl {
	if !existsApiUrl(url, apiUrls) {
		fmt.Printf("append url: %s\n", url)
		apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: url})
	}

	return apiUrls
}

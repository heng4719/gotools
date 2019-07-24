package model

import (
	"fmt"
	"gotools/internal/DB"
	"gotools/internal/model/sys"
	"gotools/util"

	"github.com/xormplus/xorm"
)

var ()

func syncTable() (err error) {
	employee := new(sys.Employee)
	role := new(sys.Role)
	orgType := new(sys.OrgType)
	sysMenu := new(sys.SysMenu)
	roleMenu := new(sys.RoleMenu)

	DB.Engine.DropTables(employee, role, orgType, sysMenu, roleMenu)
	err = DB.Engine.Sync2(
		employee,
		role,
		orgType,
		sysMenu,
		roleMenu,
	)

	return err
}

func Init() {
	// 同步数据
	if err := syncTable(); err != nil {
		panic(err)
	}

	if err := initData(); err != nil {
		panic(err)
	}
}

// 初始化数据
func initData() (err error) {
	session := DB.Engine.NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		return err
	}

	if err = initOrgTypes(session); err != nil {
		session.Rollback()
		return err
	}

	if err = initSysMenu(session); err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	return err
}

// 读取org_type.json文件的数据
func readOrgTypes() (orgTypes []*sys.OrgType) {
	orgTypes = make([]*sys.OrgType, 0)
	util.ReadJSON("../configs/org_type.json", &orgTypes)

	return orgTypes
}

// 保存机构类型数据到数据库中
func initOrgTypes(session *xorm.Session) (err error) {
	orgTypes := readOrgTypes()
	for _, orgType := range orgTypes {
		if _, err = DB.InsertTx(session, orgType); err != nil {
			return err
		}

		role := sys.Role{}
		role.Name = fmt.Sprintf("%s管理员", orgType.Name)
		role.IsAdmin = 1
		role.OrgTypeId = orgType.Id
		role.OrgTypeName = orgType.Name
		role.CreateById = 1
		role.CreateBy = "super_admin"
		if _, err = DB.InsertTx(session, role); err != nil {
			return err
		}
	}

	return err
}

// 初始化菜单权限等数据
func initSysMenu(session *xorm.Session) (err error) {
	orgTypes := readOrgTypes()
	var orgTypeIds []int64
	for _, orgType := range orgTypes {
		orgTypeIds = append(orgTypeIds, orgType.Id)
	}

	baseMenu := sys.SysMenu{}
	baseMenu.Name = "基础设置"
	baseMenu.Level = "level1"
	baseMenu.Index = 1
	baseMenu.NodeType = "menu"
	baseMenu.OrgTypeIds = orgTypeIds

	sysMenu := sys.SysMenu{}
	sysMenu.Name = "系统管理"
	sysMenu.Level = "level1"
	sysMenu.Index = 2
	sysMenu.NodeType = "menu"
	sysMenu.OrgTypeIds = orgTypeIds

	if _, err = DB.InsertTx(session, &baseMenu); err != nil {
		return err
	}

	if _, err = DB.InsertTx(session, &sysMenu); err != nil {
		return err
	}

	if err = insertRoleMenus(session, &sysMenu); err != nil {
		return err
	}

	if err = insertEmployeeMenus(session, &sysMenu); err != nil {
		return err
	}

	if err = insertBaseMenus(session, &baseMenu); err != nil {
		return err
	}

	return err
}

// 新增角色菜单权限
func insertRoleMenus(session *xorm.Session, parent *sys.SysMenu) (err error) {
	// 角色管理菜单
	roleMenu := sys.SysMenu{}
	roleMenu.ParentId = parent.Id
	roleMenu.Name = "角色管理"
	roleMenu.RouteLink = "/sys/roles"
	roleMenu.Level = "level2"
	roleMenu.Index = 2
	roleMenu.NodeType = "menu"
	roleMenu.OrgTypeIds = parent.OrgTypeIds
	if _, err = DB.InsertTx(session, &roleMenu); err != nil {
		return err
	}

	// 角色查询权限
	roleMenus := make([]*sys.SysMenu, 0)
	roleQuery := sys.SysMenu{}
	roleQuery.ParentId = roleMenu.Id
	roleQuery.Name = "角色查询"
	roleQuery.Level = "level3"
	roleQuery.Index = 1
	roleQuery.NodeType = "permission"
	roleQuery.OrgTypeIds = parent.OrgTypeIds

	apiUrls := make([]*sys.ApiUrl, 0)
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/role/query"})
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/role/get"})
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/org_type/query_by_select"})
	roleQuery.ApiUrls = apiUrls
	roleMenus = append(roleMenus, &roleQuery)

	// 角色编辑权限
	roleSave := sys.SysMenu{}
	roleSave.ParentId = roleMenu.Id
	roleSave.Name = "角色编辑"
	roleSave.Level = "level3"
	roleSave.Index = 2
	roleSave.NodeType = "permission"
	roleSave.OrgTypeIds = parent.OrgTypeIds

	apiUrls = make([]*sys.ApiUrl, 0)
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/role/save"})
	roleSave.ApiUrls = apiUrls
	roleMenus = append(roleMenus, &roleSave)

	// 角色删除权限
	roleDelete := sys.SysMenu{}
	roleDelete.ParentId = roleMenu.Id
	roleDelete.Name = "角色删除"
	roleDelete.Level = "level3"
	roleDelete.Index = 2
	roleDelete.NodeType = "permission"
	roleDelete.OrgTypeIds = parent.OrgTypeIds

	apiUrls = make([]*sys.ApiUrl, 0)
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/role/del"})
	roleDelete.ApiUrls = apiUrls
	roleMenus = append(roleMenus, &roleDelete)
	if _, err = DB.InsertTx(session, &roleMenus); err != nil {
		return err
	}

	return err
}

// 新增员工菜单权限
func insertEmployeeMenus(session *xorm.Session, parent *sys.SysMenu) (err error) {
	// 员工管理菜单
	employeeMenu := sys.SysMenu{}
	employeeMenu.ParentId = parent.Id
	employeeMenu.Name = "员工管理"
	employeeMenu.RouteLink = "/sys/employees"
	employeeMenu.Level = "level2"
	employeeMenu.Index = 2
	employeeMenu.NodeType = "menu"
	employeeMenu.OrgTypeIds = parent.OrgTypeIds
	if _, err = DB.InsertTx(session, &employeeMenu); err != nil {
		return err
	}

	// 员工查询权限
	employeeMenus := make([]*sys.SysMenu, 0)
	employeeQuery := sys.SysMenu{}
	employeeQuery.ParentId = employeeMenu.Id
	employeeQuery.Name = "员工查询"
	employeeQuery.Level = "level3"
	employeeQuery.Index = 1
	employeeQuery.NodeType = "permission"
	employeeQuery.OrgTypeIds = parent.OrgTypeIds

	apiUrls := make([]*sys.ApiUrl, 0)
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/employee/query"})
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/employee/get"})
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/employee/update_password"})
	employeeQuery.ApiUrls = apiUrls
	employeeMenus = append(employeeMenus, &employeeQuery)

	// 员工编辑权限
	employeeSave := sys.SysMenu{}
	employeeSave.ParentId = employeeMenu.Id
	employeeSave.Name = "员工编辑"
	employeeSave.Level = "level3"
	employeeSave.Index = 2
	employeeSave.NodeType = "permission"
	employeeSave.OrgTypeIds = parent.OrgTypeIds

	apiUrls = make([]*sys.ApiUrl, 0)
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/employee/save"})
	employeeSave.ApiUrls = apiUrls
	employeeMenus = append(employeeMenus, &employeeSave)

	// 员工权限分配
	employeePerms := sys.SysMenu{}
	employeePerms.ParentId = employeeMenu.Id
	employeePerms.Name = "权限分配"
	employeePerms.Level = "level3"
	employeePerms.Index = 3
	employeePerms.NodeType = "permission"
	employeePerms.OrgTypeIds = parent.OrgTypeIds

	apiUrls = make([]*sys.ApiUrl, 0)
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/role_menu/save"})
	apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/role_menu/query"})
	employeePerms.ApiUrls = apiUrls
	employeeMenus = append(employeeMenus, &employeePerms)

	if _, err = DB.InsertTx(session, &employeeMenus); err != nil {
		return err
	}

	return err
}

// 根据机构类型，添加机构管理的菜单权限
func insertBaseMenus(session *xorm.Session, parent *sys.SysMenu) (err error) {
	orgTypes := readOrgTypes()
	for index, orgType := range orgTypes {
		if orgType.Code == "platform" {
			continue
		}

		// xxx管理菜单
		baseMenu := sys.SysMenu{}
		baseMenu.ParentId = parent.Id
		baseMenu.Name = fmt.Sprintf("%s管理", orgType.Name)
		baseMenu.RouteLink = fmt.Sprintf("/base/%ss", orgType.Code)
		baseMenu.Level = "level2"
		baseMenu.Index = index
		baseMenu.NodeType = "menu"
		baseMenu.OrgTypeIds = getOrgParentIds(orgType)
		if _, err = DB.InsertTx(session, &baseMenu); err != nil {
			return err
		}

		// xxx查询权限
		childMenus := make([]*sys.SysMenu, 0)
		baseQuery := sys.SysMenu{}
		baseQuery.ParentId = baseMenu.Id
		baseQuery.Name = fmt.Sprintf("%s查询", orgType.Name)
		baseQuery.Level = "level3"
		baseQuery.Index = 1
		baseQuery.NodeType = "permission"
		baseQuery.OrgTypeIds = getOrgParentIds(orgType)

		apiUrls := make([]*sys.ApiUrl, 0)
		apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: fmt.Sprintf("/v1/admin_api/%s/query", orgType.Code)})
		apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: fmt.Sprintf("/v1/admin_api/%s/get", orgType.Code)})
		apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: "/v1/admin_api/org_type/query_by_select"})
		baseQuery.ApiUrls = apiUrls
		childMenus = append(childMenus, &baseQuery)

		// xxx编辑权限
		baseSave := sys.SysMenu{}
		baseSave.ParentId = baseMenu.Id
		baseSave.Name = fmt.Sprintf("%s编辑", orgType.Name)
		baseSave.Level = "level3"
		baseSave.Index = 2
		baseSave.NodeType = "permission"
		baseSave.OrgTypeIds = getOrgParentIds(orgType)

		apiUrls = make([]*sys.ApiUrl, 0)
		apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: fmt.Sprintf("/v1/admin_api/%s/save", orgType.Code)})
		baseSave.ApiUrls = apiUrls
		childMenus = append(childMenus, &baseSave)

		// xxx编辑权限
		baseDelete := sys.SysMenu{}
		baseDelete.ParentId = baseMenu.Id
		baseDelete.Name = fmt.Sprintf("%s删除", orgType.Name)
		baseDelete.Level = "level3"
		baseDelete.Index = 3
		baseDelete.NodeType = "permission"
		baseDelete.OrgTypeIds = getOrgParentIds(orgType)

		apiUrls = make([]*sys.ApiUrl, 0)
		apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: fmt.Sprintf("/v1/admin_api/%s/del", orgType.Code)})
		baseDelete.ApiUrls = apiUrls
		childMenus = append(childMenus, &baseDelete)

		basePassWordReset := sys.SysMenu{}
		basePassWordReset.ParentId = baseMenu.Id
		basePassWordReset.Name = "密码重置"
		basePassWordReset.Level = "level3"
		basePassWordReset.Index = 4
		basePassWordReset.NodeType = "permission"
		basePassWordReset.OrgTypeIds = getOrgParentIds(orgType)

		apiUrls = make([]*sys.ApiUrl, 0)
		//apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: fmt.Sprintf("/v1/admin_api/%s/reset_password",orgType.Code)})
		apiUrls = append(apiUrls, &sys.ApiUrl{ApiUrl: fmt.Sprintf("/v1/admin_api/employee/reset_password")})
		basePassWordReset.ApiUrls = apiUrls
		childMenus = append(childMenus, &basePassWordReset)

		if _, err = DB.InsertTx(session, &childMenus); err != nil {
			return err
		}
	}

	return
}

func getOrgParentIds(orgType *sys.OrgType) []int64 {
	orgTypeIds := []int64{orgType.ParentId}
	return orgTypeIds
}

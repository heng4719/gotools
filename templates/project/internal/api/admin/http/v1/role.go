package v1

import (
	"goadmin/internal/model/sys"
	service "goadmin/internal/service/sys"

	"github.com/gin-gonic/gin"
)

var (
	roleService = &service.RoleService{}
)

type RoleController struct {
	Router gin.IRouter
}

func (*RoleController) Find(c *gin.Context) {
	page, limit := GetPageParams(c)
	claims := GetEmployeeClaims(c)

	roles, err := roleService.Find(claims.RoleId, page, limit)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, roles)
}

func (*RoleController) Save(c *gin.Context) {
	role := new(sys.Role)
	if err := BindJSON(c, role); err != nil {
		ResponseErrorf(c, "请求的数据不合法: %s", err)
		return
	}

	claims := GetEmployeeClaims(c)
	if claims.RoleId == 1 {
		role.IsAdmin = 1
	} else {
		orgType, err := orgTypeService.GetByRoleId(claims.RoleId)
		if err != nil {
			ResponseError(c, "请求的数据不合法: "+err.Error())
			return
		}

		role.OrgTypeId = orgType.Id
		role.OrgTypeName = orgType.Name
	}

	if err := roleService.Save(role); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (*RoleController) Delete(c *gin.Context) {
	id := GetId(c)
	if err := roleService.Delete(id); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (self *RoleController) Setup() {
	self.Router.GET("/role/query", self.Find)
	self.Router.POST("/role/save", self.Save)
	self.Router.GET("/role/del", self.Delete)
}

/**
 * Created by Wangwei on 2019-06-03 20:09.
 */

package v1

import (
	"goadmin/internal/model/sys"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	Router gin.IRouter
}

func (self *EmployeeController) Setup() {
	self.Router.GET("/employee/get", self.Get)
	self.Router.GET("/employee/query", self.Find)
	self.Router.POST("/employee/save", self.Save)
	self.Router.POST("/employee/reset_password", self.ResetPassword)
	self.Router.POST("/employee/update_password", self.UpdatePassword)
}

func (self *EmployeeController) Find(c *gin.Context) {
	claims := GetEmployeeClaims(c)
	page, limit := GetPageParams(c)
	keyword := c.Query("keyword")

	var isAdmin int
	if claims.RoleId == 1 {
		isAdmin = 1
	}

	employees, err := employeeService.FindByPage(page, limit, keyword, claims.OrgTypeId, claims.OrgId, isAdmin)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, employees)
}

func (self *EmployeeController) Get(c *gin.Context) {
	employeeId := GetId(c)
	employee, err := employeeService.FindById(employeeId)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, employee)
}

func (self *EmployeeController) Save(c *gin.Context) {
	employee := new(sys.Employee)
	if err := BindJSON(c, employee); err != nil {
		ResponseError(c, err)
		return
	}

	claims := GetEmployeeClaims(c)
	employee.OrgTypeId = claims.OrgTypeId
	employee.OrgId = claims.OrgId
	employee.OrgName = claims.OrgName
	if err := employeeService.Save(employee); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

type EmployeeReq struct {
	account string
}

func (self *EmployeeController) ResetPassword(c *gin.Context) {

	employeeReq := new(EmployeeReq)
	if err := BindJSON(c, employeeReq); err != nil {
		ResponseError(c, err)
		return
	}

	if err := employeeService.ResetPassword(employeeReq.account); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (self *EmployeeController) UpdatePassword(c *gin.Context) {

	passwordReq := new(sys.PasswordReq)
	if err := BindJSON(c, passwordReq); err != nil {
		ResponseError(c, err)
		return
	}

	if err := employeeService.UpdatePassword(passwordReq); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

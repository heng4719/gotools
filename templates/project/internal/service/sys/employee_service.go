/**
 * Created by Wangwei on 2019-06-04 17:21.
 */

package sys

import (
	"errors"
	"goadmin/internal/common/DB"
	"goadmin/internal/model/dto"
	"goadmin/internal/model/sys"
	"goadmin/pkg/security"

	"github.com/xormplus/xorm"
)

var (
	SELECT_EMPLOYEE_TOTAL = "select_employee_total.stpl"
	SELECT_EMPLOYEE       = "select_employee.stpl"
)

type EmployeeService struct{}

func (self *EmployeeService) Get(id int64) (*sys.Employee, error) {
	employee := sys.Employee{}
	has, err := DB.GetById(id, &employee)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("账户不存在")
	}

	return &employee, err
}

//分页查询
func (self *EmployeeService) FindByPage(page, limit int, keyword string, orgTypeId, orgId int64, isAdmin int) (pages *dto.Pages, err error) {
	params := map[string]interface{}{
		"offset":      GetOffset(page, limit),
		"limit":       limit,
		"keyword":     LikeStr(keyword),
		"org_type_id": orgTypeId,
		"org_id":      orgId,
		"is_admin":    isAdmin,
	}

	var total int64
	employees := make([]*sys.Employee, 0)
	err = DB.PageBySqlTemplateClient(SELECT_EMPLOYEE, &params, &employees, SELECT_EMPLOYEE_TOTAL, &total)

	pages = &dto.Pages{Total: total, Data: employees}
	return pages, err
}

func (EmployeeService) GetByAccount(account, password string) (employee *sys.Employee, err error) {
	var ok bool
	employee = &sys.Employee{}
	ok, err = DB.Where("`account`=? and `password`=?", account, password).Omit("password").Get(employee)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("用户名/密码错误")
	}
	return
}

//通过id查询
func (self *EmployeeService) FindById(id int64) (*sys.Employee, error) {
	employee := new(sys.Employee)
	ok, err := DB.GetById(id, employee)
	if !ok {
		return nil, errors.New("员工ID:" + string(employee.Id) + "不存在")
	}

	return employee, err
}

func (self *EmployeeService) Save(employee *sys.Employee) (err error) {
	session := NewSession()
	defer session.Close()

	return self.SaveTx(session, employee)
}

func (EmployeeService) SaveTx(session *xorm.Session, employee *sys.Employee) (err error) {
	count, err := DB.Where("account = ?", employee.Account).Count(&sys.Employee{})
	if err != nil {
		return err
	}

	if employee.Id == 0 && count > 0 {
		return errors.New("账号已存在，请检查！")
	}

	if employee.Id == 0 {
		employee.Password = security.MD5Password(employee.Password)
		_, err = DB.InsertTx(session, employee)
	} else {
		_, err = DB.UpdateByIdWithOmitTx(session, employee.Id, employee, "account", "password")
	}

	return
}

func (self *EmployeeService) ResetPassword(account string) (err error) {

	employee := new(sys.Employee)
	ok, err := DB.Where("`account`=? ", account).Get(employee)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("无此用户")
	}

	employee.Password = security.MD5Password("111111")
	_, err = DB.UpdateById(employee.Id, employee)

	return err
}

func (self *EmployeeService) UpdatePassword(req *sys.PasswordReq) (err error) {

	employee := new(sys.Employee)
	ok, err := DB.Where("`account`=? and `password`=?", req.Account, security.MD5Password(req.OldPassword)).Get(employee)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("密码错误")
	}

	employee.Password = security.MD5Password(req.Password)
	_, err = DB.UpdateById(employee.Id, employee)

	return err
}

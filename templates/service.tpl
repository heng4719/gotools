package base

import (
	"${project}/internal/common/DB"
	"${project}/internal/model/base"
	"${project}/internal/model/dto"
	"${project}/internal/model/sys"

	"github.com/jinzhu/copier"
)

type ${modelName}Service struct{}

func (*${modelName}Service) Find(page, limit int, superId int) (pages *dto.Pages, err error) {
	var total int64
	total, err = DB.Count(&base.${modelName}{})
	if err != nil {
		return nil, err
	}

	offset := GetOffset(page, limit)
	${lowerModelName}sTemp := make([]*base.${modelName}, 0)
	${lowerModelName}s := make([]*base.${modelName}, 0)
	if err = DB.Limit(offset, limit).Find(&${lowerModelName}sTemp); err != nil {
		return nil, err
	}
	
	//筛选出employees.superId == superId 的记录，即被请求搜索的人所创建出来的记录
	for i := 0; i < len(${lowerModelName}sTemp); i++ {
		count, err := DB.Where("account = ? and super_id = ?", ${lowerModelName}sTemp[i].Account, superId).Count(sys.Employee{})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			${lowerModelName}s = append(${lowerModelName}s, ${lowerModelName}sTemp[i])
		}
	}

	pages = &dto.Pages{total, &${lowerModelName}s}
	return
}

func (*${modelName}Service) Get(id int64) (${lowerModelName} *base.${modelName}, err error) {
	${lowerModelName} = &base.${modelName}{}
	_, err = DB.GetById(id, ${lowerModelName})
	return
}

// 新增服务商
func (*${modelName}Service) Save(req *base.${modelName}Req) (err error) {
	${lowerModelName} := base.${modelName}{}
	copier.Copy(&${lowerModelName}, req)
	if req.Id > 0 {
		_, err = DB.UpdateById(${lowerModelName}.Id, &${lowerModelName})
		return
	}

	session := NewSession()
	defer session.Close()
	session.Begin()

	role, err := roleService.GetAdminRoleByOrgId(${lowerModelName}.OrgTypeId)
	if err != nil {
		return err
	}

	if _, err = DB.InsertTx(session, &${lowerModelName}); err != nil {
		return err
	}

	employee := sys.Employee{
		Name:      ${lowerModelName}.ContactName,
		Account:   req.Account,
		Password:  req.Password,
		RoleId:    role.Id,
		RoleName:  role.Name,
		OrgTypeId: role.OrgTypeId,
		OrgId:     ${lowerModelName}.Id,
		OrgName:   ${lowerModelName}.Name,
		Phone:     req.ServicePhone,
		SuperId:   req.SuperId,
	}

	if err = employeeService.SaveTx(session, &employee); err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return
}
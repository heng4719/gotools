package base

import (
	"${project}/internal/common/DB"
	"${project}/internal/model/dto"
)

type ${modelName}Service struct{}

// 根据id获取数据
func (*${modelName}Service) Get(id int64) (${lowerModelName} *base.${modelName}, err error) {
	${lowerModelName} = &base.${modelName}{}
	_, err = DB.GetById(id, ${lowerModelName})
	return
}

// 按条件查询
func (*${modelName}Service) Find(page, limit int) (pages *dto.Pages, err error) {
	var total int64
	total, err = DB.Count(&base.${modelName}{})
	if err != nil {
		return nil, err
	}

	offset := GetOffset(page, limit)
	${lowerModelName}s := make([]*base.${modelName}, 0)
	if err = DB.Limit(offset, limit).Find(&${lowerModelName}s); err != nil {
		return nil, err
	}

	pages = &dto.Pages{total, &${lowerModelName}s}
	return
}

// 新增或修改
func (*${modelName}Service) Save(${lowerModelName} *base.${modelName}) (err error) {
	if ${lowerModelName}.Id > 0 {
		_, err = DB.UpdateById(${lowerModelName}.Id, &${lowerModelName})
		return
	}

	err = DB.Insert(${lowerModelName})
	return err
}

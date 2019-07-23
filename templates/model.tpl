/**
 * Created by Wangwei on 2019-06-28 11:27.
 */

package base

import "time"

type ${modelName} struct {
	Id           int64     `json:"id"`
	${modelFields}
	ParentId     int64	   `xorm:"comment('上级机构ID')"`
	OrgTypeId    int64     `xorm:"comment('所属组织类型ID')"`
	OrgTypeName  string    `xorm:"varchar(40) comment('所属组织类型名称')"`
	Account      string    `xorm:"varchar(40) comment('管理员账号')"`
	CreatedAt    time.Time `xorm:"created" json:"created_at"`
	UpdatedAt    time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt    time.Time `xorm:"deleted" json:"deleted_at"`
}

type ${modelName}Req struct {
	Id           int64
	${justFields}
	ParentId     int64
	OrgTypeId    int64
	OrgTypeName  string
	Account      string
	Password     string
}

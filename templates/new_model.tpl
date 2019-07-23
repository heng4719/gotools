/**
 * Created by Wangwei on 2019-06-28 11:27.
 */

package base

import "time"

type ${modelName} struct {
	Id           int64     `json:"id"`
	${modelFields}
	CreatedAt    time.Time `xorm:"created" json:"created_at"`
	UpdatedAt    time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt    time.Time `xorm:"deleted" json:"deleted_at"`
}

type ${modelName}Req struct {
	Id           int64
	${justFields}
}
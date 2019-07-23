/**
 * Created by Wangwei on 2019-06-20 23:14.
 */

package v1

import (
	"${project}/internal/model/base"
	service "${project}/internal/service/base"
	"strconv"
		
	"github.com/gin-gonic/gin"
)

var (
	${lowerModelName}Service = &service.${modelName}Service{}
)

type ${modelName}Controller struct {
	Router gin.IRouter
}

func (*${modelName}Controller) Find(c *gin.Context) {
	page, limit := GetPageParams(c)
	parentId, _ := strconv.Atoi(c.Query("parent_id"))
	${lowerModelName}s, err := ${lowerModelName}Service.Find(page, limit, parentId)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, ${lowerModelName}s)
}

func (*${modelName}Controller) Save(c *gin.Context) {
	req := base.${modelName}Req{}
	if err := BindJSON(c, &req); err != nil {
		ResponseErrorf(c, "请求的数据不合法: %s", err)
		return
	}

	if err := ${lowerModelName}Service.Save(&req); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (self *${modelName}Controller) Setup() {
	self.Router.GET("/${snakeModelName}/query", self.Find)
	self.Router.POST("/${snakeModelName}/save", self.Save)
}

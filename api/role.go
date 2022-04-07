package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"trojan/module/constant"
	"trojan/module/dto"
	"trojan/module/vo"
	"trojan/service"
)

func SelectRoleList(c *gin.Context) {
	var roleDto dto.RoleDto
	_ = c.ShouldBind(&roleDto)
	if err := validate.Struct(&roleDto); err != nil {
		log.Println(err)
		vo.Fail(constant.ValidateFailed, c)
		return
	}
	roleListVos, err := service.SelectRoleList(roleDto)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(roleListVos, c)
}

package model

import "gin-template/common"

func GetRight(uid int, iid int) int {
	// 传入用户id和实例id，返回用户对实例的权限
	var container UserContainer
	result := DB.Table("user_containers").Where("user_id = ? AND id = ?", uid, iid).First(&container)
	if result.Error != nil {
		common.SysLog("No right: " + result.Error.Error())
		return 0 // 可以根据实际情况选择更合适的错误处理方式
	} else {

		return 1
	}
}

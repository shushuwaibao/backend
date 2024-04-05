package model

func GetRight(uid int, iid int) int {
	// 传入用户id和实例id，返回用户对实例的权限
	DB.Table("user_containers").Find("user_id = ? and id = ?", uid, iid)
	return 0
}

package model

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetAvailableArchive() ([]ImageConfig, error) {
	var results []ImageConfig
	// 使用已经初始化的 DB 进行查询
	err := DB.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func DeleteImage(imageList []ImageConfig) (int, error) {
	database, err := sqlx.Open("mysql", "root:shushuwaibao@tcp(172.16.13.73:13306)/wb2?parseTime=true")
	if err != nil {
		return 0, err
	}
	defer database.Close() // 确保在函数结束时关闭数据库连接

	// 验证连接是否有效
	err = database.Ping()
	if err != nil {
		log.Fatalf("ping mysql failed: %v", err)
	}

	// 准备DELETE语句，这里以删除某个具体ID的记录为例
	stmt, err := database.Prepare("DELETE FROM image_configs WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	//删除
	for _, image := range imageList {

		_, err = stmt.Exec(image.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
	return 1, err
}

func UpdateImagePermission(imageList []ImageConfig, newValue []string) (int, error) {
	database, err := sqlx.Open("mysql", "root:shushuwaibao@tcp(172.16.13.73:13306)/wb2?parseTime=true")
	if err != nil {
		return 0, err
	}
	defer database.Close() // 确保在函数结束时关闭数据库连接

	// 验证连接是否有效
	err = database.Ping()
	if err != nil {
		log.Fatalf("ping mysql failed: %v", err)
		return 0, err
	}

	// 准备UPDATE语句，修改特定字段的值
	stmt, err := database.Prepare("UPDATE image_configs SET permission = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer stmt.Close()

	for i := 0; i < len(imageList); i++ {
		// 执行UPDATE语句
		_, err = stmt.Exec(newValue[i], imageList[i].ID)
		if err != nil {
			log.Fatal(err)
			return 0, err
		}
	}

	return 1, err
}

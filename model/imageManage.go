package model

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetAvailableArchive() ([]ImageConfig, error) {
	var results []ImageConfig
	database, err := sqlx.Open("mysql", "root:shushuwaibao@tcp(172.16.13.73:13306)/wb2?parseTime=true")
	if err != nil {
		return results, err
	}
	defer database.Close() // 确保在函数结束时关闭数据库连接

	// 验证连接是否有效
	err = database.Ping()
	if err != nil {
		log.Fatalf("ping mysql failed: %v", err)
		return results, err
	}

	//从表中获取镜像数据
	rows, err := database.Query("SELECT * FROM image_configs")
	if err != nil {
		log.Fatal(err)
		return results, err
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var data ImageConfig
		err := rows.Scan(&data.ID, &data.Nickname, &data.Name, &data.Registry, &data.Version, &data.Description, &data.Size, &data.BelongsToWho, &data.BelongsTo, &data.Permission) // ... 根据你的表结构扫描相应的列
		if err != nil {
			log.Fatal(err)
			return results, err
		}
		results = append(results, data)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return results, err
	}

	return results, err
}

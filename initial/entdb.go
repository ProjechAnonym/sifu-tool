package initial

import (
	"context"
	"fmt"
	"sifu-tool/ent"

	"entgo.io/ent/dialect"
	"go.uber.org/zap"
)

func InitEntdb(workDir string, logger *zap.Logger) *ent.Client {
	entClient, err := ent.Open(dialect.SQLite, fmt.Sprintf("file:%s/sifu-tool.db?cache=shared&_fk=1", workDir))
	if err != nil {
		logger.Error(fmt.Sprintf("连接Ent数据库失败: [%s]", err.Error()))
		panic(err)
	}
	logger.Info("连接Ent数据库完成")
	if err = entClient.Schema.Create(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("创建表资源失败: [%s]", err.Error()))
		panic(err)
	}
	logger.Info("自动迁移Ent数据库完成")
	return entClient
}
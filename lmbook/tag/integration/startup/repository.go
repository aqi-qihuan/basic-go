package startup

import (
	"basic-go/lmbook/pkg/logger"
	"basic-go/lmbook/tag/repository"
	"basic-go/lmbook/tag/repository/cache"
	"basic-go/lmbook/tag/repository/dao"
)

func InitRepository(d dao.TagDAO, c cache.TagCache, l logger.LoggerV1) repository.TagRepository {
	return repository.NewTagRepository(d, c, l)
}

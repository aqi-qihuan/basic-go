.PHONY: mock
mock:
	@mockgen -source=./lmbook/internal/service/user.go -package=svcmocks -destination=./lmbook/internal/service/mocks/user.mock.go
	@mockgen -source=./lmbook/internal/service/code.go -package=svcmocks -destination=./lmbook/internal/service/mocks/code.mock.go
	@mockgen -source=./lmbook/internal/service/article.go -package=svcmocks -destination=./lmbook/internal/service/mocks/article.mock.go
	@mockgen -source=./lmbook/internal/service/sms/types.go -package=smsmocks -destination=./lmbook/internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./lmbook/internal/repository/code.go -package=repomocks -destination=./lmbook/internal/repository/mocks/code.mock.go
	@mockgen -source=./lmbook/internal/repository/user.go -package=repomocks -destination=./lmbook/internal/repository/mocks/user.mock.go
	@mockgen -source=./lmbook/internal/repository/article.go -package=repomocks -destination=./lmbook/internal/repository/mocks/article.mock.go
	@mockgen -source=./lmbook/internal/repository/article_author.go -package=repomocks -destination=./lmbook/internal/repository/mocks/article_author.mock.go
	@mockgen -source=./lmbook/internal/repository/article_reader.go -package=repomocks -destination=./lmbook/internal/repository/mocks/article_reader.mock.go
	@mockgen -source=./lmbook/internal/repository/dao/user.go -package=daomocks -destination=./lmbook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./lmbook/internal/repository/dao/article_reader.go -package=daomocks -destination=./lmbook/internal/repository/dao/mocks/article_reader.mock.go
	@mockgen -source=./lmbook/internal/repository/dao/article_author.go -package=daomocks -destination=./lmbook/internal/repository/dao/mocks/article_author.mock.go
	@mockgen -source=./lmbook/internal/repository/cache/user.go -package=cachemocks -destination=./lmbook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./lmbook/internal/repository/cache/code.go -package=cachemocks -destination=./lmbook/internal/repository/cache/mocks/code.mock.go
	@mockgen -source=./lmbook/pkg/limiter/types.go -package=limitermocks -destination=./lmbook/pkg/limiter/mocks/limiter.mock.go
	@mockgen -package=redismocks -destination=lmbook/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy
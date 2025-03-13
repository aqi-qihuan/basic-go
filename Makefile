.PHONY: mock
mock:
	@mockgen -source=./lanmengbook/internal/service/user.go -package=svcmocks -destination=./lanmengbook/internal/service/mocks/user.mock.go
	@mockgen -source=./lanmengbook/internal/service/code.go -package=svcmocks -destination=./lanmengbook/internal/service/mocks/code.mock.go
	@mockgen -source=./lanmengbook/internal/service/sms/types.go -package=smsmocks -destination=./lanmengbook/internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./lanmengbook/internal/repository/code.go -package=repomocks -destination=./lanmengbook/internal/repository/mocks/code.mock.go
	@mockgen -source=./lanmengbook/internal/repository/user.go -package=repomocks -destination=./lanmengbook/internal/repository/mocks/user.mock.go
	@mockgen -source=./lanmengbook/internal/repository/dao/user.go -package=daomocks -destination=./lanmengbook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./lanmengbook/internal/repository/cache/user.go -package=cachemocks -destination=./lanmengbook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./lanmengbook/internal/repository/cache/code.go -package=cachemocks -destination=./lanmengbook/internal/repository/cache/mocks/code.mock.go
	@mockgen -source=./lanmengbook/pkg/limiter/types.go -package=limitermocks -destination=./lanmengbook/pkg/limiter/mocks/limiter.mock.go
	@mockgen -package=redismocks -destination=lanmengbook/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy
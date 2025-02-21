.PHONY: mock
mock:
	@mockgen -source=./lanmengbook/internal/service/user.go \
                 -package=svcmocks \
                 -destination=./lanmengbook/internal/service/mocks/user.mock.go
	@mockgen -source=./lanmengbook/internal/service/code.go \
                -package=svcmocks \
                -destination=./lanmengbook/internal/service/mocks/code.mock.go

	@mockgen -source=./lanmengbook/internal/repository/user.go \
                     -package=repomocks \
                     -destination=./lanmengbook/internal/repository/mocks/user.mock.go
	@mockgen -source=./lanmengbook/internal/repository/code.go \
                    -package=repomocks \
                    -destination=./lanmengbook/internal/repository/mocks/code.mock.go


	@go mod tidy
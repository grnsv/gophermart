#!/bin/bash

export PATH=$PATH:$GOPATH/bin

mockgen -destination=internal/mocks/mock_user_service.go -package=mocks github.com/grnsv/gophermart/internal/services UserService
mockgen -destination=internal/mocks/mock_order_service.go -package=mocks github.com/grnsv/gophermart/internal/services OrderService
mockgen -destination=internal/mocks/mock_jwt_service.go -package=mocks github.com/grnsv/gophermart/internal/services JWTService
mockgen -destination=internal/mocks/mock_validator.go -package=mocks github.com/grnsv/gophermart/internal/services Validator
mockgen -destination=internal/mocks/mock_accrual_service.go -package=mocks github.com/grnsv/gophermart/internal/services AccrualService
mockgen -destination=internal/mocks/mock_user_repository.go -package=mocks github.com/grnsv/gophermart/internal/storage UserRepository
mockgen -destination=internal/mocks/mock_order_repository.go -package=mocks github.com/grnsv/gophermart/internal/storage OrderRepository
mockgen -destination=internal/mocks/mock_withdrawal_repository.go -package=mocks github.com/grnsv/gophermart/internal/storage WithdrawalRepository

go mod tidy
go vet $(go list ./... | grep -v /vendor/)
go fmt $(go list ./... | grep -v /vendor/)
go test -race $(go list ./... | grep -v /vendor/)

cd cmd/gophermart
go build
cd ../..

docker compose build gophermart accrual
docker compose up -d db accrual

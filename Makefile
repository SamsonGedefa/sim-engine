.PHONY: migrate-create migrate-down migrate-force sqlc init

PWD := ${CURDIR}
ACCTPATH = $(PWD)/simulator
PORT = 5432
N = 1

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(PWD)/db/migrations -seq -digits 5 $(NAME)

migrate-up:
	migrate -source file://$(PWD)/db/migrations -database postgres://root:password@localhost:$(PORT)/simulator-dev?sslmode=disable up $(N)

migrate-down:
	migrate -source file://$(PWD)/db/migrations -database postgres://root:password@localhost:$(PORT)/simulator-dev?sslmode=disable down $(N)

migrate-force:
	migrate -source file://$(PWD)/db/migrations -database postgres://root:password@localhost:$(PORT)/simulator-dev?sslmode=disable force $(VERSION)

sqlc:
	docker run --rm -v "$(PWD):/src" -w /src kjconroy/sqlc generate

# mock:
# 	mockgen -source=kazna-api/model/interfaces.go -destination kazna-api/db/mock/interfaces_mock.go -package mockdb 

mock:
	cd $(ACCTPATH) && mockery --all --dir model --outpkg mocks --output ./model/mocks --with-expecter

	
init: 
	docker-compose up -d postgres-account && \
	$(MAKE) migrate-down APPPATH=account N= && \
	$(MAKE) migrate-up APPPATH=account N= && \
	docker-compose down


default:
	@echo "=============building Local API============="
	docker build -f Dockerfile -t api .

up: default
	@echo "=============starting api locally============="
	docker-compose up -d

logs:
	docker-compose logs -f

down:
	docker-compose down

test:
	go test -v -cover ./...

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f
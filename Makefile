build:
	go build -o dist/app cmd/service.go
run:
	docker compose up -d
stop:
	docker compose down
# database migrations
create_migration:
	migrate create -dir ./migrations -ext sql $(NAME)
migrate_up:
	migrate -path ./migrations -database postgres://root:pass@localhost:5432/ryde_database?sslmode=disable up $(N)
migrate_down:
	migrate -path ./migrations -database postgres://root:pass@localhost:5432/ryde_database?sslmode=disable down $(N)
migrate_version:
	migrate -path ./migrations -database postgres://root:pass@localhost:5432/ryde_database?sslmode=disable version

# tests
mock_interface:
	mockgen -source $(FILE_PATH) -destination=$(DESTINATION_FILE) -package=$(PACKAGE) -build_flags=mock
test:
	@go test -race -v -cover ./... -coverprofile=coverage.out && \
	echo "\nTotal coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"
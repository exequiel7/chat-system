start:
	docker-compose up -d
destroy:
	docker-compose down
build-doc:
	cd cmd/api && swag init --parseDependency --parseInternal
coverage:
	sh coverage.sh
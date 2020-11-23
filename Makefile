db-create:
	- ./typicalw pg create

db-drop:
	- ./typicalw pg drop

db-migrate:
	- ./typicalw pg migrate

db-rollback:
	- ./typicalw pg rollback

run:
	- ./typicalw run

run-test-controller:
	@go test -v -coverprofile=coverage.out ./app/book/controller/... -bench=.

run-test-service:
	@go test -v -coverprofile=coverage.out ./app/book/service/... -bench=.

run-test-repository:
	@go test -v -coverprofile=coverage.out ./app/book/repository/... -bench=.
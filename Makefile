# Run all tests in memory mode (no external dependencies required).
# Matches the default local development workflow.
test:
	go test ./... -v -race

# Run all tests including Redis integration tests.
# Starts an ephemeral Redis container on port 16379, runs all tests with the race
# detector, then removes the container regardless of test outcome.
# Mirrors what the CI workflow does (quality-and-security.yml).
test-distributed:
	@echo "Starting ephemeral Redis container..."
	@docker run -d --name wacraft-core-test-redis -p 16379:6379 redis:7-alpine > /dev/null
	@echo "Waiting for Redis to be ready..."
	@until docker exec wacraft-core-test-redis redis-cli ping 2>/dev/null | grep -q PONG; do sleep 0.1; done
	@echo "Running tests..."
	@REDIS_URL=redis://localhost:16379 go test ./... -v -race; \
		EXIT=$$?; \
		echo "Removing Redis container..."; \
		docker rm -f wacraft-core-test-redis > /dev/null; \
		exit $$EXIT

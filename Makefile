# Makefile for Go testing workflow with Docker.
# This file provides shortcuts for running tests in different configurations.

# .PHONY tells Make these don't represent files.
.PHONY: test test-all test-clients test-repositories test-clean

# Main test target - configurable via TEST_DIRS variable
test:
	@if [ -z "$(ALPHA_VANTAGE_API_KEY)" ] && [ -n "$$(echo $(TEST_DIRS) | grep clients)" ]; then \
		echo "Warning: ALPHA_VANTAGE_API_KEY environment variable not set"; \
		echo "Integration tests will be skipped"; \
	fi

	@echo "Running tests for: $(TEST_DIRS)"

	docker-compose -f docker-compose.yml up --build
	
	@echo "Coverage reports available in ./test-reports/"

# Test all packages
test-all:
	# Set TEST_DIRS to all packages and call the main test target
	TEST_DIRS="./..." $(MAKE) test

# Test only the clients package
test-clients:
	TEST_DIRS="./clients" $(MAKE) test

# Test only the repositories package
test-repositories:
	TEST_DIRS="./repositories" $(MAKE) test

# Test both clients and repositories packages together
test-core:
	# Comma-separated list for multiple directories
	TEST_DIRS="./clients,./repositories" $(MAKE) test


# Clean up test artifacts and Docker resources
test-clean:
	# Remove generated test reports
	rm -rf test-reports
	# Remove the Docker container and image
	# --rmi local: Remove the image built for this test
	docker-compose -f docker-compose.yml down --rmi local

#!/bin/bash

# Exit immediately if a command exits with non-zero status.
set -e 

# Initialize exit code variable to track if any tests fail.
EXIT_CODE=0

# Ensure the test-reports directory exists.
mkdir -p /app/test-reports

# Convert a comma separated directory list to array for separate processing.
IFS=',' read -ra DIRS <<< "$TEST_DIRS"

# Process each directory in the array
for dir in "${DIRS[@]}"; do
	# Trim leading/trailing whitespace from directory name
	dir=$(echo $dir | xargs)

	# Visual separator for test output
	echo "==============================================="
	echo "Running tests for: $dir"
	echo "==============================================="

	if [ "$COVERAGE" = "true" ]; then
		# Create a unique name for the coverage profile based on the directory name.
		PROFILE_NAME=$(basename $dir | sed 's/\.\///g')
		if [ "$PROFILE_NAME" = "..." ]; then
			PROFILE_NAME="all"
		fi

		echo "Generating coverage report for $dir"

		# -v: Verbose output.
		# -cover: Enable coverage analysis.
		# -coverprofile: Save coverage data to html file for later.
		if ! go test -v -cover -coverprofile=/app/test-reports/coverage-$PROFILE_NAME.out $dir; then 
			# If tests fail, set exit code to 1 but continue running other tests
			EXIT_CODE=1
		fi

		# Generate HTML coverage report from the coverage data
		# This creates a visual representation of coverage with color-coded source code
		go tool cover -html=/app/test-reports/coverage-$PROFILE_NAME.out -o=/app/test-reports/coverage-$PROFILE_NAME.html

		# Generate function-level coverage summary
		# Shows coverage support percentage for each function
		go tool cover -func=/app/test-reports/coverage-$PROFILE_NAME.out > /app/test-reports/coverage-$PROFILE_NAME-func.txt
	else 
		# When coverage is disabled, run tests without coverage for faster execution
		if ! go test -v $dir; then
			EXIT_CODE=1
		fi
	fi
done

# If coverage is enabled and we tested multiple directories,
# create a combined coverage report
if [ "$COVERAGE" = "true" ] && [ ${#DIRS[@]} -gt 1 ]; then
  echo "Generating combined coverage report"
  
  # Start with coverage mode line
  echo "mode: set" > /app/test-reports/coverage-combined.out
  
  # Append all coverage data, skipping the mode line from each file
  grep -h -v "mode: set" /app/test-reports/coverage-*.out >> /app/test-reports/coverage-combined.out
  
  # Generate combined HTML report
  go tool cover -html=/app/test-reports/coverage-combined.out -o=/app/test-reports/coverage-combined.html
  
  # Generate combined function coverage report
  go tool cover -func=/app/test-reports/coverage-combined.out > /app/test-reports/coverage-combined-func.txt
  
  # Print overall coverage percentage for quick reference
  echo "Overall coverage:"
  cat /app/test-reports/coverage-combined-func.txt | grep total:
fi

# Print summary information
echo ""
echo "==============================================="
echo "Test Summary"
echo "==============================================="
echo "Test reports saved to ./test-reports/"

# Print final status message based on test results
if [ $EXIT_CODE -eq 0 ]; then
  echo "All tests passed!"
else
  echo "Some tests failed! Check the logs above for details."
fi

# Return the exit code
# This ensures the container exits with non-zero status if tests failed
exit $EXIT_CODE

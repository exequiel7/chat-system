# This script measures the coverage and compares it to the minimum specified.
# If the minimum coverage is greater than the current coverage, it exits with an error.
# If the minimum coverage is equal to or less than the current coverage, it exits successfully.

# File containing paths to ignore
ignore_file=".coverignore"

# Build the ignore pattern for grep
if [ -f "$ignore_file" ] && [ -s "$ignore_file" ]; then
    ignore_pattern=$(sed 's/^/-e /' "$ignore_file" | tr '\n' ' ')
    echo "Ignoring the following patterns:"
    cat "$ignore_file"
else
    ignore_pattern=""
    echo "No patterns to ignore."
fi

# Get the list of packages, excluding those listed in the ignore file
packages=$(go list ./... | grep -v $ignore_pattern)

# If the packages variable is empty, it means everything was ignored, so handle it
if [ -z "$packages" ]; then
    echo "[ ERROR ] - All packages were excluded, please check your .coverignore"
    exit 1
fi

# Print the packages that will be tested
echo "Testing the following packages:"
echo "$packages"

# Run tests and create coverage report
go test $packages -coverprofile=cover.out.tmp -coverpkg=./... -covermode=atomic

# Filter out ignored files from the coverage report
grep -vFf "$ignore_file" cover.out.tmp > cover.out

# Generate HTML report
go tool cover -html=cover.out

min_coverage=85.0 # Minimum Coverage Allowed
coverage=$(go tool cover -func cover.out | grep total | awk '{print $3}' | sed -e 's/[%]//g') # Get the coverage obtained from running unit tests
_result=$(echo "$coverage >= $min_coverage" | bc)

if [ $_result == "1" ]; then
    rm cover.out cover.out.tmp
    echo "Coverage OK!, min coverage is $min_coverage% and the current is $coverage%";
    exit 0
else
    rm cover.out cover.out.tmp
    echo "[ ERROR ] - Minimum coverage error, min coverage is $min_coverage% and the current is $coverage%";
    exit 1
fi;

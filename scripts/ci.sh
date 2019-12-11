# This script contains commands to be executed by the CI tool.
echo "Running tests..."
echo "" > coverage.txt;
for d in $(go list ./... | grep -v 'vendor\|cmd'); do
	GOCACHE=off
	go test -v -race -tags test -coverprofile=profile.out -covermode=atomic $d
	if [ -f profile.out ]; then
		cat profile.out >> coverage.txt
		rm profile.out
	fi
done

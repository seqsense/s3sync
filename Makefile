.PHONY: test
test:
	go test .

.PHONY: cover
cover:
	go test -coverprofile=cover.out .
	go tool cover -html=cover.out -o report.html

.PHONY: s3
s3:
	docker run -p 4572:4572 -e SERVICES=s3 localstack/localstack

.PHONY: fixture
fixture:
	aws s3 --endpoint-url http://localhost:4572 mb s3://example-bucket
	aws s3 --endpoint-url http://localhost:4572 cp README.md s3://example-bucket
	aws s3 --endpoint-url http://localhost:4572 cp README.md s3://example-bucket/foo/
	aws s3 --endpoint-url http://localhost:4572 cp README.md s3://example-bucket/bar/baz/


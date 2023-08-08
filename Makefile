
testpostgres:
	go test ./postgres

testall:
	go test ./...

deploy_hello:
	gcloud functions deploy hello \
		--runtime go120 \
		--entry-point Hello \
		--trigger-http \
		--source /home/ahansson/github.com/andreashanson/golang-postgres-test/internal/hello/
MY_SQL_DATA_DIR := $(shell realpath "./testdata/mnt/")
TEST_USER := root
TEST_PASSWORD := example

stand-up:
	sudo rm -rf $(MY_SQL_DATA_DIR)
	mkdir -p $(MY_SQL_DATA_DIR)
	sudo MY_SQL_DATA_DIR=$(MY_SQL_DATA_DIR) docker compose -f ./testdata/compose.yaml up -d
	@echo "Waiting for MySQL to start..."
	sleep 10

stand-down:
	sudo MY_SQL_DATA_DIR=$(MY_SQL_DATA_DIR) docker compose -f ./testdata/compose.yaml down || true
	sudo MY_SQL_DATA_DIR=$(MY_SQL_DATA_DIR) docker compose -f ./testdata/compose.yaml rm || true
	sudo rm -rf $(MY_SQL_DATA_DIR) || true

test:
	MYSQL_SKIP_TEST_MAIN=1 go test ./... -count 1 -race -v -shuffle=on

login:
	sudo -E docker exec                                                          \
		--env "MYSQL_PWD=$(TEST_PASSWORD)"                                       \
		-it $$(sudo docker container ls -a -q --filter "NAME=testdata-mysql-1")  \
		mysql -u $(TEST_USER)

full-test: stand-up test stand-down


dockerList:
	sudo docker ps -a
dockerStart:
	sudo docker start d4773580d079
dbMigrate:
	goose -dir migrations postgres "host=127.0.0.1 port=5432 user=postgres password=12345 dbname=postgres sslmode=disable" up
deps:
	go get -u gorm.io/gorm \
    go get -u gorm.io/driver/postgres \
    go get github.com/golang/mock/gomock
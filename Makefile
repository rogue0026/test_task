.PHONY build_app:
build_app:
	go build -o ./cmd/bin/application ./cmd/application/main.go;

.PHONY run_app:
run_app:
	SMTP_HOST=sandbox.smtp.mailtrap.io SMTP_PORT=2525 SMTP_LOGIN=ec8ab8fc8c02a3 SMTP_PASSWORD=89f91f31e9ad80 ./cmd/bin/application;

.PHONY deploy_database:
deploy_database:
	docker run -d --name test_container -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=test -p 5432:5432 postgres;
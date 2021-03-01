migrate-up:
	migrate -database postgres://root:root@localhost:5432/wills-svc?sslmode=disable -path db/migrations up

migrate-down:
	migrate -database postgres://root:root@localhost:5432/wills-svc?sslmode=disable -path db/migrations down

migrate-force:
	migrate -database postgres://root:root@localhost:5432/wills-svc?sslmode=disable -path db/migrations force 000001
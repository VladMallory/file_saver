run: stop
	docker compose up -d --build

log:
	docker compose logs -f

stop:
	docker compose down 

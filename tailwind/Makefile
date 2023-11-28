up:
	docker compose up -d
down:
	docker compose down
build-dev:
	docker compose run tailwind npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /dst/styles.css
build-prod:
	docker compose run tailwind npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /dst/styles.css --minify
.PHONY: build run dev clean

build:
	docker build -t voicescribe-pro .

run:
	docker run -p 8000:8000 voicescribe-pro

# Dev commands
dev:
	docker compose up --build

dev-python:
	cd ml && uvicorn main:app --host 0.0.0.0 --port 9090 --reload

clean:
	docker compose down
	docker rmi voicescribe-pro

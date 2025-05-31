.PHONY: build run dev clean

build:
	docker build -t voicescribe-pro .

run:
	docker run -p 8000:8000 voicescribe-pro

dev:
	docker compose up --build

clean:
	docker compose down
	docker rmi voicescribe-pro

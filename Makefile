# Declare phony targets to avoid conflicts with files of the same name
.PHONY: build run dev dev-python clean

# Builds the Docker image for the VoiceScribe-Pro project using the Dockerfile in the current directory.
# The resulting image is tagged as 'voicescribe-pro'.
build:
	docker build -t voicescribe-pro .

run:
	docker run -p 8000:8000 voicescribe-pro

# Dev commands - run GO server
dev:
	docker compose up --build

# Run Python server
dev-python:
	cd ml && uvicorn main:app --host 0.0.0.0 --port 9090 --reload

# Run client. 
dev-client:
	cd web/ && dum dev


clean:
	docker compose down
	docker rmi voicescribe-pro

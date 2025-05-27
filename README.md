# VoiceScribe-Pro

## Backend Structure

```
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   ├── models/
│   ├── services/
│   │   ├── transcription/
│   │   ├── summarization/
│   │   └── search/
│   └── storage/
├── pkg/
│   ├── huggingface/
│   ├── websocket/
│   └── utils/
└── configs/
```
# Engineering Assistant

A real-time speech-to-text assistant powered by whisper.cpp and Claude AI. This application transcribes your speech and provides AI-powered responses through a clean web interface.

## Features

- Real-time speech transcription using whisper.cpp
- AI responses powered by Claude 3 Opus
- Clean web interface with cost monitoring
- Pause/Resume functionality
- Conversation history management

## Prerequisites

- Go 1.21 or later
- whisper.cpp installed and compiled
- Claude API key from Anthropic

## Installation

1. Clone the repository:
```bash
git clone https://github.com/dimitarkovachev/eng-assist.git
cd eng-assist
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the project root:
```env
WHISPER_CPP_PATH=/path/to/whisper/executable
ANTHROPIC_API_KEY=your_api_key_here
```

## Usage

1. Start the application:
```bash
go run cmd/assistant/main.go
```

Optional flags:
- `--buffer-timeout`: Time to wait before processing buffered text (default: 1.0s)
- `--debug`: Enable debug mode

2. Open your browser and navigate to `http://localhost:5000`

3. Start speaking! The application will:
   - Transcribe your speech in real-time
   - Process the transcription through Claude AI
   - Display responses in the web interface
   - Track API usage costs

## Project Structure

```
eng-assist/
├── cmd/
│   └── assistant/
│       └── main.go           # Application entry point
├── pkg/
│   ├── ai/
│   │   └── client.go        # Claude AI client
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── transcription/
│   │   └── transcription.go # Speech transcription
│   └── ui/
│       └── ui.go            # Web interface
└── ui/
    └── static/              # Frontend assets
```

## License

MIT License # eng-assist

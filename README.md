# RAG LangChain

This simple RAG (Retrieval Augmented Generation) golang based example using google gemini.
The application allows you to add documents and query them using natural language, powered by Gemini's language model.

## Prerequisites
* Go 1.23 or later
* Docker and Docker Compose
* Google Cloud Project with Gemini API access

## Setup
1. Get a Google Gemini API key from Google AI Studio

2. set the api key as an environment variable
```bash
export GEMINI_API_KEY=your_api_key_here
```

3. Start the Weaviate vector database:
```bash
docker-compose up -d
```

4. Build run the application
```bash
make run
```

## Usage
The server provides the following endpoints:

### Add Documents
```bash
curl -X POST http://localhost:8000/add/ \
  -H "Content-Type: application/json" \
  -d '{"documents": [{"text": "Your document text here"}]}'
```

### Query Documents
```bash
curl -X POST http://localhost:8000/query/ \
  -H "Content-Type: application/json" \
  -d '{"content": "Your question here"}'
```

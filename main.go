package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

const generativeModelName = "gemini-1.5-flash"
const embeddingModelName = "text-embedding-004"

func main() {
	fmt.Println("Running RAG model", generativeModelName)

	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	geminiClient, err := googleai.New(ctx,
		googleai.WithAPIKey(apiKey),
		googleai.WithDefaultEmbeddingModel(embeddingModelName))
	if err != nil {
		log.Fatal(err)
	}

	emb, err := embeddings.NewEmbedder(geminiClient)
	if err != nil {
		log.Fatal(err)
	}

	wvStore, err := weaviate.New(
		weaviate.WithEmbedder(emb),
		weaviate.WithScheme("http"),
		weaviate.WithHost("localhost:"+cmp.Or(os.Getenv("WVPORT"), "9035")),
		weaviate.WithIndexName("Document"),
	)

	server := &ragServer{
		ctx:          ctx,
		wvStore:      wvStore,
		geminiClient: geminiClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /healthz/", server.HealthEndpoint)

	port := cmp.Or(os.Getenv("SERVERPORT"), "9020")
	address := "localhost:" + port
	log.Println("listening on", address)
	log.Fatal(http.ListenAndServe(address, mux))
}

type ragServer struct {
	ctx          context.Context
	wvStore      weaviate.Store
	geminiClient *googleai.GoogleAI
}

// HealthEndpoint
func (h *ragServer) HealthEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service Healthy"))
}

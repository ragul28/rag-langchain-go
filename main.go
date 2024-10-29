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
	"github.com/tmc/langchaingo/schema"
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
	mux.HandleFunc("POST /add/", server.addDocumentsHandler)

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

func (rs *ragServer) addDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type document struct {
		Text string
	}
	type addRequest struct {
		Documents []document
	}
	ar := &addRequest{}

	err := readRequestJSON(req, ar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store documents and their embeddings in weaviate
	var wvDocs []schema.Document
	for _, doc := range ar.Documents {
		wvDocs = append(wvDocs, schema.Document{PageContent: doc.Text})
	}
	_, err = rs.wvStore.AddDocuments(rs.ctx, wvDocs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

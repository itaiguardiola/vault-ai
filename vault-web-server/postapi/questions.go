package postapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pashpashpash/vault/form"
)

type Context struct {
	Text  string `json:"text"`
	Title string `json:"title"`
}

type Answer struct {
	Answer  string    `json:"answer"`
	Context []Context `json:"context"`
	Tokens  int       `json:"tokens"`
}

// Handle Requests For Question
func (ctx *HandlerContext) QuestionHandler(w http.ResponseWriter, r *http.Request) {
	form := new(form.QuestionForm)

	if errs := FormParseVerify(form, "QuestionForm", w, r); errs != nil {
		return
	}

	log.Println("[QuestionHandler] Question:", form.Question)
	log.Println("[QuestionHandler] Model:", form.Model)
	log.Println("[QuestionHandler] UUID:", form.UUID)

	clientToUse := ctx.llmClient

	// step 1: Feed question to local LLM embeddings api to get an embedding back
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	questionEmbedding, err := clientToUse.GetEmbedding(ctxWithTimeout, form.Question)
	if err != nil {
		log.Println("[QuestionHandler ERR] OpenAI get embedding request error\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("[QuestionHandler] Question Embedding Length:", len(questionEmbedding))

	// step 2: Query vector db using questionEmbedding to get context matches
	matches, err := ctx.vectorDB.Retrieve(questionEmbedding, 4, form.UUID)
	if err != nil {
		log.Println("[QuestionHandler ERR] Vector DB query error\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("[QuestionHandler] Got matches from vector DB:", matches)

	// Extract context text and titles from the matches
	contexts := make([]Context, len(matches))
	for i, match := range matches {
		contexts[i].Text = match.Metadata["text"]
		contexts[i].Title = match.Metadata["title"]
	}
	log.Println("[QuestionHandler] Retrieved context from vector DB:\n", contexts)

	// step 3: Structure the prompt with a context section + question, using top x results from vector DB as the context
	contextTexts := make([]string, len(contexts))
	for i, context := range contexts {
		contextTexts[i] = context.Text
	}
	prompt, err := buildPrompt(contextTexts, form.Question)
	if prompt == "" {
		prompt = form.Question
	}
	if err != nil {
		log.Println("[QuestionHandler ERR] Error building prompt\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[QuestionHandler] Sending local LLM api request...\nPrompt:%s\n", prompt)

	llmCtx, llmCancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer llmCancel()

	llmResponse, err := clientToUse.CreateChatCompletionSimple(
		llmCtx,
		prompt,
		"You are a helpful assistant answering questions based on the context provided.",
		512,
		0.7,
	)

	if err != nil {
		log.Println("[QuestionHandler ERR] LLM answer questions request error\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("[QuestionHandler] LLM response:\n", llmResponse)
	response := OpenAIResponse{llmResponse, 0} // tokens not tracked for local models

	answer := Answer{response.Response, contexts, response.Tokens}
	jsonResponse, err := json.Marshal(answer)
	if err != nil {
		log.Println("[QuestionHandler ERR] OpenAI response marshalling error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

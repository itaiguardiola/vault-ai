# OP Vault - Local Edition

OP Vault is a document-based question-answering system that allows you to upload your own custom knowledgebase files and ask questions about their contents. This version runs **completely locally** without any internet connection or cloud services using Ollama (local LLM) + Qdrant (local vector database).

**Original version**: [vault.pash.city](https://vault.pash.city) - used OpenAI + Pinecone (cloud-based)
**This version**: Fully local - uses Ollama + Qdrant (no internet required)

## 🚀 Quick Start with Docker (Recommended!)

The easiest way to run OP Vault is with Docker:

```bash
docker-compose up
```

Then visit **http://localhost:8100** and you're ready to go!

**First-time setup** takes 5-10 minutes to download models. See [DOCKER_SETUP.md](DOCKER_SETUP.md) for full Docker documentation.

---

<img width="512" alt="Screen Shot 2023-04-09 at 1 53 33 AM" src="/static/img/common/vault_library.png">

With quick setup, you can launch your own version of this Golang server along with a user-friendly React frontend that allows users to ask OpenAI questions about the specific knowledge base provided. The primary focus is on human-readable content like books, letters, and other documents, making it a practical and valuable tool for knowledge extraction and question-answering. You can upload an entire library's worth of books and documents and recieve pointed answers along with the name of the file and specific section within the file that the answer is based on!

<img width="1498" alt="Screen Shot 2023-04-17 at 6 23 00 PM" src="https://user-images.githubusercontent.com/20898225/232645187-fff56d2b-f654-4c92-b061-4670734b2764.png">

## What can you do with OP Vault?

With The Vault, you can:

-   Upload a variety of popular document types via a simple react frontend to create a custom knowledge base
-   **OCR support** for scanned PDFs and images - automatically extracts text using Tesseract
-   Retrieve accurate and relevant answers based on the content of your uploaded documents
-   See the filenames and specific context snippets that inform the answer
-   Explore the power of local AI models in a user-friendly interface
-   Load entire libraries' worth of books into The Vault
-   **Manage documents** through the built-in Document Library at http://localhost:8100/documents
-   **Configure settings** through the built-in web UI at http://localhost:8100/config

## Setup Options

Choose your preferred setup method:

- **🐳 [Docker Setup](DOCKER_SETUP.md)** (Recommended) - One command, everything containerized
- **💻 [Local Setup](SETUP_LOCAL.md)** - Manual installation for more control
- Both methods work completely offline after initial setup!

## Manual Dependencies (for Local Setup)

-   node: v19+ (tested with v22)
-   go: v1.18.9+ (tested with v1.24.7)
-   poppler (for PDF text extraction)
-   tesseract-ocr (for scanned PDF and image text extraction)
-   Docker (for running Qdrant vector database)
-   Ollama (for local LLM and embeddings)

## Setup - Local Installation (No Internet Required)

### 1. Install manual dependencies

**Install Go:**
Follow the go docs [here](https://go.dev/doc/install)

**Install Node.js:**
I recommend [installing nvm and using it to install node v19+](https://medium.com/@iam_vinojan/how-to-install-node-js-and-npm-using-node-version-manager-nvm-143165b16ce1)

**Install Poppler:**
- Ubuntu/Debian: `sudo apt-get install -y poppler-utils`
- Mac: `brew install poppler`
- Fedora/RHEL: `sudo dnf install poppler-utils`

**Install Tesseract OCR:**
- Ubuntu/Debian: `sudo apt-get install -y tesseract-ocr tesseract-ocr-eng`
- Mac: `brew install tesseract`
- Fedora/RHEL: `sudo dnf install tesseract tesseract-langpack-eng`

**Install Docker:**
Follow the Docker docs [here](https://docs.docker.com/get-docker/)

**Install Ollama:**
```bash
# Linux
curl -fsSL https://ollama.com/install.sh | sh

# Mac
brew install ollama

# Windows
# Download from https://ollama.com/download
```

### 2. Start Qdrant Vector Database

Qdrant is a local vector database that stores document embeddings:

```bash
docker run -d -p 6333:6333 -v $(pwd)/qdrant_storage:/qdrant/storage qdrant/qdrant
```

This will:
- Run Qdrant on port 6333
- Persist data in `./qdrant_storage` directory
- Run in detached mode (background)

### 3. Start Ollama and Download Models

Start the Ollama service:
```bash
# Mac/Linux
ollama serve
```

In another terminal, download the required models:
```bash
# Download embedding model (768 dimensions)
ollama pull nomic-embed-text

# Download chat model (choose one)
ollama pull llama3       # Recommended: 8B parameters, good balance
# OR
ollama pull llama3:70b   # Better quality, requires more RAM (40GB+)
# OR
ollama pull mistral      # Alternative 7B model
```

**Model Recommendations:**
- **nomic-embed-text**: Required for embeddings (768-dimensional vectors)
- **llama3**: Best general-purpose model (requires ~8GB RAM)
- **llama3:70b**: Higher quality (requires ~40GB RAM)
- **mistral**: Lightweight alternative (requires ~4GB RAM)

### 4. (Optional) Configure Environment Variables

The app uses sensible defaults, but you can customize:

```bash
export OLLAMA_ENDPOINT="http://localhost:11434"           # Default
export OLLAMA_EMBED_MODEL="nomic-embed-text"              # Default
export OLLAMA_CHAT_MODEL="llama3"                         # Default
export QDRANT_API_ENDPOINT="http://localhost:6333"        # Default
```

**Note:** You don't need to set these if using default values!

### 5. Build and Run the Application

**Install dependencies:**
```bash
npm install
```

This will automatically compile the Go server via the postinstall script.

**Run the application:**

In one terminal, start the backend server (default port `:8100`):
```bash
npm start
```

In another terminal, run webpack to compile the React frontend:
```bash
npm run dev
```

**Access the application:**
Visit http://localhost:8100 in your browser

### Quick Start Summary

Once everything is installed:
```bash
# Terminal 1: Start Qdrant (only needed once, runs in background)
docker run -d -p 6333:6333 -v $(pwd)/qdrant_storage:/qdrant/storage qdrant/qdrant

# Terminal 2: Start Ollama (if not running as service)
ollama serve

# Terminal 3: Start the backend
npm start

# Terminal 4: Start the frontend dev server
npm run dev

# Open browser
open http://localhost:8100
```

## Screenshots:

In the example screenshots, I uploaded a couple of books by Plato and some letters by Alexander Hamilton, showcasing the ability of OP Vault to answer questions based on the uploaded content.

### Uploading files

<img width="1483" alt="Screen Shot 2023-04-17 at 6 16 40 PM" src="https://user-images.githubusercontent.com/20898225/232645162-e89dc752-ad69-40d3-9eda-8c9075ddeeda.png">
<img width="1509" alt="Screen Shot 2023-04-17 at 6 17 29 PM" src="https://user-images.githubusercontent.com/20898225/232645171-b8eb56d5-8797-4970-b163-e17ed76b5b97.png">

### Asking questions

<img width="1498" alt="Screen Shot 2023-04-17 at 6 20 25 PM" src="https://user-images.githubusercontent.com/20898225/232645180-f41b3ebc-e050-4df5-bd47-b2819f480081.png">
<img width="1500" alt="Screen Shot 2023-04-17 at 6 20 58 PM" src="https://user-images.githubusercontent.com/20898225/232645183-e28bc0fa-3545-48f3-9374-29529c513fe2.png">
<img width="1498" alt="Screen Shot 2023-04-17 at 6 23 00 PM" src="https://user-images.githubusercontent.com/20898225/232645187-fff56d2b-f654-4c92-b061-4670734b2764.png">

## Under the hood

The golang server uses POST APIs to process incoming uploads and respond to questions:

1.  `/upload` for uploading files

2.  `/api/questions` for answering questions

All api endpoints are declared in the [vault-web-server/main.go](vault-web-server/main.go) file.

### Uploading files and processing them into embeddings

The [vault-web-server/postapi/fileupload.go](vault-web-server/postapi/fileupload.go) file contains the `UploadHandler` logic for handling incoming uploads on the backend.
The UploadHandler function in the postapi package is responsible for handling file uploads (with a maximum total upload size of 50 MB) and processing them into embeddings to store in Qdrant. It accepts PDF, epub, .docx, and plain text files, extracts text from them, and divides the content into chunks. Using **Ollama's local embedding model (nomic-embed-text)**, it obtains embeddings for each chunk and upserts (inserts or updates) the embeddings into **Qdrant (local vector database)**. The function returns a JSON response containing information about the uploaded files and their processing status.

1. Limit the size of the request body to MAX_TOTAL_UPLOAD_SIZE (300 MB).
2. Parse the incoming multipart form data with a maximum allowed size of 300 MB.
3. Initialize response data with fields for successful and failed file uploads.
4. Iterate over the uploaded files, and for each file:
   a. Check if the file size is within the allowed limit (MAX_FILE_SIZE, 300 MB).
   b. Read the file into memory.
   c. If the file is a PDF, extract the text from it; otherwise, read the contents as plain text.
   d. Divide the file contents into chunks.
   e. Use Ollama local embedding API to obtain embeddings for each chunk (768-dimensional vectors).
   f. Upsert (insert or update) the embeddings into Qdrant local vector database.
   g. Update the response data with information about successful and failed uploads.
5. Return a JSON response containing information about the uploaded files and their processing status.

### Storing embeddings into Qdrant db

After getting embeddings from Ollama for each chunk of an uploaded file, the server stores all of the embeddings, along with metadata associated for each embedding in Qdrant DB. The metadata for each embedding is created in the [UpsertEmbeddings](vectordb/qdrant/qdrant.go) function, with the following keys and values:

-   `file_name`: The name of the file from which the text chunk was extracted.
-   `start`: The starting character position of the text chunk in the original file.
-   `end`: The ending character position of the text chunk in the original file.
-   `title`: The title of the chunk, which is also the file name in this case.
-   `text`: The text of the chunk.

This metadata is useful for providing context to the embeddings and is used to display additional information about the matched embeddings when retrieving results from the Qdrant database.

### Answering questions

The `QuestionHandler` function in [vault-web-server/postapi/questions.go](vault-web-server/postapi/questions.go) is responsible for handling all incoming questions. When a question is entered on the frontend and the user presses "search" (or enter), the server uses the **Ollama embedding API** once again to get an embedding for the question (a.k.a. query vector). This query vector is used to query **Qdrant local database** to get the most relevant context for the question. Finally, a prompt is built by packing the most relevant context + the question in a prompt string that adheres to LLM token limits (the go tiktoken library is used to estimate token count). The answer is generated using **Ollama's local LLM (llama3 or mistral)** instead of cloud-based APIs.

### Frontend info

The frontend is built using `React.js` and `less` for styling.

### Generative question-answering with long-term memory

This local implementation uses the same RAG (Retrieval-Augmented Generation) architecture as the cloud version:
1. Documents are chunked and embedded using local models
2. Embeddings are stored in a local vector database (Qdrant)
3. Questions are embedded and matched against stored documents
4. Relevant context is retrieved and passed to a local LLM
5. The LLM generates answers based on your documents

**Advantages of running locally:**
- ✅ **Complete privacy**: Your documents never leave your machine
- ✅ **No API costs**: No per-token charges from cloud providers
- ✅ **No internet required**: Works completely offline
- ✅ **Unlimited usage**: Query as much as you want
- ✅ **Customizable models**: Swap models based on your needs
- ✅ **Data sovereignty**: Full control over your data

If you'd like to read more about RAG (Retrieval-Augmented Generation), see:
-   https://www.pinecone.io/learn/retrieval-augmented-generation/
-   https://ollama.com/blog/embedding-models

I hope you enjoy it (:

## Uploading larger files

The max individual file size is set to 25MB and total upload size to 50MB. If you want to increase this limit, edit the `MAX_FILE_SIZE` and `MAX_TOTAL_UPLOAD_SIZE` constants in [fileupload.go](vault-web-server/postapi/fileupload.go#L21-L22).

### Supported Filetypes

**Document Files:**
- PDFs (text-based and scanned with automatic OCR fallback)
- Text files (.txt, .rtf)
- Word documents (.doc, .docx)
- EPUB files (.epub)
- Pages files (.pages)
- Plaintext

**Image Files (with OCR):**
- PNG (.png)
- JPEG (.jpg, .jpeg)
- GIF (.gif)
- TIFF (.tiff, .tif)
- BMP (.bmp)
- WebP (.webp)

**Note:** The system automatically detects scanned PDFs and image-based documents, falling back to OCR extraction when needed. See [OCR_FEATURE.md](OCR_FEATURE.md) for details.

## Troubleshooting

### Ollama connection issues
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Verify models are installed
ollama list
```

### Qdrant connection issues
```bash
# Check if Qdrant is running
curl http://localhost:6333/collections

# Check Docker container status
docker ps | grep qdrant
```

### Build errors
```bash
# Clean and rebuild
rm -rf node_modules bin
npm install
```

### Performance tips
- Use `llama3` (8B) for best balance of speed and quality
- Use `mistral` (7B) for faster responses on lower-end hardware
- Reduce batch size in embedding calls if running out of memory
- Increase Qdrant memory limit in Docker if handling large document sets

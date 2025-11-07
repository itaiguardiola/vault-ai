package chunk

import (
	"bytes"
	"log"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strings"

	"github.com/neurosnap/sentences/english"
	tke "github.com/pkoukk/tiktoken-go"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"github.com/gabriel-vasile/mimetype"


	"github.com/saintfish/chardet"


	"io"

	"code.sajari.com/docconv"
	"github.com/gen2brain/go-fitz"
	gosseract "github.com/otiai10/gosseract/v2"
)

type Chunk struct {
	Start int
	End   int
	Title string
	Text  string
}

// MaxTokensPerChunk is the maximum number of tokens allowed in a single chunk for OpenAI embeddings
// MaxTokensPerChunk is the maximum number of tokens allowed in a single chunk for OpenAI embeddings
const MaxTokensPerChunk = 1500
const EmbeddingModel = "text-embedding-ada-002"

func CreateChunks(fileContent string, title string) ([]Chunk, error) {
	chunks := []Chunk{}

	// Initialize sentence tokenizer
	tokenizer, _ := english.NewSentenceTokenizer(nil)
	sentences := tokenizer.Tokenize(fileContent)

	// Get tiktoken encoding for the model
	tiktoken, err := tke.EncodingForModel(EmbeddingModel)
	if err != nil {
		return []Chunk{}, fmt.Errorf("getEncoding: %v", err)
	}

	chunkStart := 0

	for chunkStart < len(sentences) {
		tokenCount := 0
		chunkText := ""
		chunkSentences := 0

		for i := chunkStart; i < len(sentences) && tokenCount < MaxTokensPerChunk; i++ {
			sentence := sentences[i].Text
			tiktokens := tiktoken.Encode(sentence, nil, nil)
			sentenceTokenCount := len(tiktokens)

			if sentenceTokenCount > MaxTokensPerChunk {
				continue // Skip sentence if longer than MaxTokensPerChunk
			}

			if tokenCount+sentenceTokenCount <= MaxTokensPerChunk {
				tokenCount += sentenceTokenCount
				chunkText += " " + sentence
				chunkSentences++
			} else {
				break
			}
		}

		trimmedText := strings.TrimSpace(chunkText)
		if len(trimmedText) > 0 {
			chunks = append(chunks, Chunk{
				Start: chunkStart,
				End:   chunkStart + tokenCount,
				Title: title,
				Text:  trimmedText,
			})
		}

		// Calculate stride dynamically based on chunk sentences
		sentenceStride := chunkSentences / 5
		if sentenceStride == 0 {
			sentenceStride = 1
		}

		// Move chunkStart forward by sentenceStride
		chunkStart += sentenceStride
	}

	if len(chunks) == 0 {
		return nil, errors.New("no chunks created")
	}

	return chunks, nil
}


func GetTextFromFile(f multipart.File) (string, error) {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	mime := mimetype.Detect(content)
	contentType := mime.String()

	var text string

	log.Println("[GetTextfromFile] ContentType:", contentType)
	switch contentType {
	case "application/msword": // .doc
		log.Println("[GetTextfromFile] .doc file encountered...")
		text, _, err = docconv.ConvertDoc(bytes.NewReader(content))
		if err != nil {
			return "", fmt.Errorf("error converting .doc file")
		}
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document": // .docx
		log.Println("[GetTextfromFile] .docx file encountered...")
		text, _, err = docconv.ConvertDocx(bytes.NewReader(content))
		if err != nil {
			return "", fmt.Errorf("error converting .docx file: %v", err)
		}
	case "application/zip": // .pages
		log.Println("[GetTextfromFile] .pages file encountered...")
		text, _, err = docconv.ConvertPages(bytes.NewReader(content))
		if err != nil {
			return "", fmt.Errorf("error converting .pages file: %v", err)
		}
		text = strings.TrimSpace(text)
	case "application/epub+zip": // .epub
		log.Println("[GetTextfromFile] .epub file encountered...")
		fitzDoc, err := fitz.NewFromReader(bytes.NewReader(content))
		if err != nil {
			return "", fmt.Errorf("error reading .epub file: %v", err)
		}
		defer fitzDoc.Close()
		for i := 0; i < fitzDoc.NumPage(); i++ {
			pageText, err := fitzDoc.Text(i)
			if err != nil {
				return "", fmt.Errorf("error getting text from page %d: %v", i, err)
			}
			// Preprocess the text by replacing newline characters with spaces
			pageText = strings.ReplaceAll(pageText, "\n", " ")
			text += pageText
		}
	case "image/png", "image/jpeg", "image/jpg", "image/gif", "image/tiff", "image/bmp", "image/webp": // Image files
		log.Printf("[GetTextfromFile] Image file encountered (%s), using OCR...", contentType)
		text, err = extractTextFromImage(content, contentType)
		if err != nil {
			return "", fmt.Errorf("error extracting text from image: %v", err)
		}
	default: // Assume plain text
		detector := chardet.NewTextDetector()
		result, err := detector.DetectBest(content)
		if err != nil {
			return "", fmt.Errorf("error detecting encoding: %v", err)
		}

		if strings.ToLower(result.Charset) == "utf-8" {
			text = string(content)
		} else {
			var enc encoding.Encoding
			switch strings.ToLower(result.Charset) {
			case "iso-8859-1":
				enc = charmap.ISO8859_1
			case "windows-1252":
				enc = charmap.Windows1252
			// Add more encodings here as needed
			default:
				return "", fmt.Errorf("unsupported encoding: %s", result.Charset)
			}

			text, _, err = transform.String(enc.NewDecoder(), string(content))
			if err != nil {
				return "", fmt.Errorf("error decoding content: %v", err)
			}
		}
	}

	return text, nil
}

// isTextSufficientForProcessing checks if extracted text is meaningful
// Returns false if text appears to be from a scanned/image-based PDF
func isTextSufficientForProcessing(text string) bool {
	trimmed := strings.TrimSpace(text)

	// Check if text is too short
	if len(trimmed) < 50 {
		return false
	}

	// Count words (split by whitespace)
	words := strings.Fields(trimmed)
	if len(words) < 10 {
		return false
	}

	// Check for meaningful text (not just numbers or special characters)
	alphaCount := 0
	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			alphaCount++
		}
	}

	// If less than 20% alphabetic characters, likely not meaningful text
	if float64(alphaCount)/float64(len(trimmed)) < 0.2 {
		return false
	}

	return true
}

// extractTextFromPDFWithOCR uses Tesseract OCR to extract text from image-based PDFs
func extractTextFromPDFWithOCR(f multipart.File) (string, error) {
	log.Println("[OCR] Attempting OCR extraction for scanned PDF...")

	// Reset file position
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to reset file position: %v", err)
	}

	// Read file content
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Open PDF with go-fitz to render pages as images
	fitzDoc, err := fitz.NewFromMemory(content)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF with fitz: %v", err)
	}
	defer fitzDoc.Close()

	var fullText strings.Builder

	// Process each page
	for pageNum := 0; pageNum < fitzDoc.NumPage(); pageNum++ {
		log.Printf("[OCR] Processing page %d/%d...", pageNum+1, fitzDoc.NumPage())

		// Render page as image (PNG format, 300 DPI for good OCR quality)
		img, err := fitzDoc.Image(pageNum)
		if err != nil {
			log.Printf("[OCR] Warning: failed to render page %d: %v", pageNum, err)
			continue
		}

		// Convert image to bytes
		var imgBuffer bytes.Buffer
		err = fitz.ImagePNG(&imgBuffer, img)
		if err != nil {
			log.Printf("[OCR] Warning: failed to encode page %d as PNG: %v", pageNum, err)
			continue
		}

		// Use Tesseract to extract text from image
		client := gosseract.NewClient()
		defer client.Close()

		err = client.SetImageFromBytes(imgBuffer.Bytes())
		if err != nil {
			log.Printf("[OCR] Warning: failed to set image for page %d: %v", pageNum, err)
			continue
		}

		// Set language to English (can be extended to support more languages)
		client.SetLanguage("eng")

		// Extract text
		pageText, err := client.Text()
		if err != nil {
			log.Printf("[OCR] Warning: failed to extract text from page %d: %v", pageNum, err)
			continue
		}

		fullText.WriteString(pageText)
		fullText.WriteString("\n\n")
	}

	result := strings.TrimSpace(fullText.String())
	log.Printf("[OCR] Successfully extracted %d characters from %d pages", len(result), fitzDoc.NumPage())

	return result, nil
}

// extractTextFromImage uses Tesseract OCR to extract text from image files
func extractTextFromImage(content []byte, contentType string) (string, error) {
	log.Printf("[OCR] Extracting text from image (%s)...", contentType)

	// Use Tesseract to extract text from image
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImageFromBytes(content)
	if err != nil {
		return "", fmt.Errorf("failed to set image: %v", err)
	}

	// Set language to English (can be extended to support more languages)
	client.SetLanguage("eng")

	// Extract text
	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %v", err)
	}

	result := strings.TrimSpace(text)
	log.Printf("[OCR] Successfully extracted %d characters from image", len(result))

	if !isTextSufficientForProcessing(result) {
		return "", fmt.Errorf("OCR extraction yielded insufficient text (%d characters)", len(result))
	}

	return result, nil
}

// extract human-readable text from a given pdf with support for spaces/whitespace.
// Automatically falls back to OCR if regular text extraction yields insufficient text.
func ExtractTextFromPDF(f multipart.File, fileSize int64) (string, error) {
	// Reset the file reader's position
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	// Try regular text extraction first (fastest method)
	log.Println("[PDF] Attempting regular text extraction...")
	bodyResult, _, err := docconv.ConvertPDF(f)

	// Remove extra whitespace and newlines
	text := strings.TrimSpace(bodyResult)

	// Check if extraction was successful and yielded sufficient text
	if err == nil && isTextSufficientForProcessing(text) {
		log.Printf("[PDF] Successfully extracted %d characters using regular method", len(text))
		return text, nil
	}

	// Regular extraction failed or yielded insufficient text
	// This likely means it's a scanned/image-based PDF
	if err != nil {
		log.Printf("[PDF] Regular extraction failed: %v", err)
	} else {
		log.Printf("[PDF] Regular extraction yielded insufficient text (%d characters, %d words)", len(text), len(strings.Fields(text)))
	}

	// Fall back to OCR
	log.Println("[PDF] Falling back to OCR extraction...")
	ocrText, ocrErr := extractTextFromPDFWithOCR(f)
	if ocrErr != nil {
		return "", fmt.Errorf("both regular extraction and OCR failed: regular=%v, ocr=%v", err, ocrErr)
	}

	// Check if OCR yielded sufficient text
	if !isTextSufficientForProcessing(ocrText) {
		return "", fmt.Errorf("OCR extraction yielded insufficient text (%d characters)", len(ocrText))
	}

	log.Printf("[PDF] OCR extraction successful: %d characters extracted", len(ocrText))
	return ocrText, nil
}

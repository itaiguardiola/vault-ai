# OCR (Optical Character Recognition) Feature

## Overview

OP Vault now includes comprehensive OCR support using **Tesseract OCR** to extract text from scanned PDFs and image files. This feature automatically detects when documents need OCR processing and falls back seamlessly when regular text extraction fails or yields insufficient text.

## Supported File Types

### PDF Files
- **Regular PDFs**: Text-based PDFs are processed using standard extraction (fastest method)
- **Scanned PDFs**: Image-based or scanned PDFs automatically fall back to OCR extraction
- **Hybrid PDFs**: PDFs with both text and images use OCR when needed

### Image Files
All common image formats are now supported with OCR:
- PNG (`.png`)
- JPEG (`.jpg`, `.jpeg`)
- GIF (`.gif`)
- TIFF (`.tiff`, `.tif`)
- BMP (`.bmp`)
- WebP (`.webp`)

## How It Works

### Intelligent Fallback System

The OCR system uses a smart fallback mechanism:

1. **For PDFs:**
   ```
   1. Attempt regular text extraction (docconv)
   2. Check if extracted text is sufficient
   3. If insufficient → Fall back to OCR
   4. Extract text from rendered PDF pages
   5. Return OCR results
   ```

2. **For Images:**
   ```
   1. Detect image file type
   2. Process directly with Tesseract OCR
   3. Extract and return text
   ```

### Text Quality Detection

The system automatically determines if extracted text is meaningful:

```go
// Text is considered sufficient if:
- Length > 50 characters
- Word count > 10 words
- Alphabetic characters > 20% of total
```

This prevents false positives where extraction succeeds but yields only formatting characters or numbers.

## Technical Implementation

### Architecture

```
┌─────────────────────┐
│  File Upload        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  File Type Check    │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
┌─────────┐   ┌─────────┐
│   PDF   │   │  Image  │
└────┬────┘   └────┬────┘
     │             │
     ▼             │
┌─────────────┐   │
│ Try Regular │   │
│ Extraction  │   │
└────┬────────┘   │
     │            │
     ▼            │
┌─────────────┐  │
│ Check Text  │  │
│ Quality     │  │
└────┬────────┘  │
     │           │
     ▼           ▼
┌────────────────────┐
│   Tesseract OCR    │
│  - Render pages    │
│  - Extract text    │
│  - Validate output │
└──────────┬─────────┘
           │
           ▼
┌─────────────────────┐
│  Chunk & Embed      │
└─────────────────────┘
```

### Key Components

#### 1. Text Quality Checker (`isTextSufficientForProcessing`)
Located in `chunk/fileprocessing.go:189-217`

Analyzes extracted text to determine if OCR is needed.

#### 2. PDF OCR Extractor (`extractTextFromPDFWithOCR`)
Located in `chunk/fileprocessing.go:220-291`

- Renders each PDF page as a high-quality image (300 DPI)
- Processes each page through Tesseract
- Aggregates text from all pages
- Handles errors gracefully (skips problematic pages)

#### 3. Image OCR Extractor (`extractTextFromImage`)
Located in `chunk/fileprocessing.go:294-323`

- Processes image files directly with Tesseract
- Validates text quality
- Returns clean, processed text

#### 4. Enhanced PDF Extraction (`ExtractTextFromPDF`)
Located in `chunk/fileprocessing.go:326-368`

- Main entry point for PDF processing
- Implements smart fallback logic
- Comprehensive logging for debugging

### Dependencies

```go
// Added to chunk/fileprocessing.go
import (
    gosseract "github.com/otiai10/gosseract/v2"
    "github.com/gen2brain/go-fitz"  // Already present for EPUB
)
```

### Docker Integration

Tesseract OCR is automatically installed in the Docker container:

```dockerfile
# Builder stage
RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-data-eng

# Runtime stage
RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-data-eng
```

## Usage

### Uploading Files

The OCR feature is completely transparent to users. Simply upload files as usual:

1. **Via Web Interface:**
   - Drag and drop PDFs or images
   - Click to select files
   - Upload button

2. **Via API:**
   ```bash
   curl -X POST http://localhost:8100/api/upload \
     -F "files=@scanned_document.pdf" \
     -F "uuid=your-collection-uuid"
   ```

### Monitoring OCR Activity

Check the server logs to see OCR activity:

```bash
# For scanned PDFs
[PDF] Attempting regular text extraction...
[PDF] Regular extraction yielded insufficient text (23 characters, 3 words)
[PDF] Falling back to OCR extraction...
[OCR] Attempting OCR extraction for scanned PDF...
[OCR] Processing page 1/5...
[OCR] Processing page 2/5...
...
[OCR] Successfully extracted 5432 characters from 5 pages
[PDF] OCR extraction successful: 5432 characters extracted

# For images
[GetTextfromFile] Image file encountered (image/png), using OCR...
[OCR] Extracting text from image (image/png)...
[OCR] Successfully extracted 234 characters from image
```

## Performance Considerations

### Speed Comparison

| Method | Speed | Use Case |
|--------|-------|----------|
| Regular PDF extraction | ~100ms | Text-based PDFs |
| OCR PDF extraction | ~2-5s per page | Scanned PDFs |
| Image OCR | ~1-3s | Image files |

### Optimization Tips

1. **Pre-process images**: Use high-quality scans (300 DPI recommended)
2. **Image format**: PNG and TIFF generally work better than JPEG
3. **Page count**: Large documents with many pages will take longer
4. **Concurrent uploads**: Process multiple files in parallel

## Configuration

### OCR Language

Currently set to English. To add more languages, modify the Docker image:

```dockerfile
# Add more language packs
RUN apk add --no-cache \
    tesseract-ocr-data-spa \  # Spanish
    tesseract-ocr-data-fra \  # French
    tesseract-ocr-data-deu    # German
```

Then update the code:

```go
// In extractTextFromPDFWithOCR and extractTextFromImage
client.SetLanguage("eng+spa+fra")  // Multiple languages
```

### Image Quality

Adjust rendering DPI in `extractTextFromPDFWithOCR`:

```go
// Higher DPI = better quality but slower
img, err := fitzDoc.ImageDPI(pageNum, 300)  // Current: 300 DPI
img, err := fitzDoc.ImageDPI(pageNum, 600)  // Better quality
```

## Troubleshooting

### OCR Not Working

**Symptom**: Files fail to process or extraction yields empty text

**Solutions**:
1. Check Tesseract installation:
   ```bash
   docker exec -it vault-app tesseract --version
   ```

2. Verify language data:
   ```bash
   docker exec -it vault-app ls /usr/share/tessdata/
   ```

3. Check file quality:
   - Ensure images are clear and readable
   - Verify PDF isn't corrupted
   - Try converting to PNG first

### Poor OCR Quality

**Symptom**: Extracted text has many errors

**Solutions**:
1. Improve source document quality
2. Increase rendering DPI (see Configuration)
3. Pre-process images:
   - Increase contrast
   - Remove noise
   - Straighten text

### Slow Processing

**Symptom**: File uploads take too long

**Solutions**:
1. Reduce image size before upload
2. Split large documents into smaller chunks
3. Reduce DPI for acceptable quality
4. Increase server resources (CPU)

## API Response

### Successful Upload with OCR

```json
{
  "message": "All files uploaded and processed successfully",
  "num_files_succeeded": 1,
  "num_files_failed": 0,
  "successful_file_names": ["scanned_document.pdf"],
  "failed_file_names": {}
}
```

### Failed OCR

```json
{
  "message": "Some files failed to upload and process",
  "num_files_succeeded": 0,
  "num_files_failed": 1,
  "successful_file_names": [],
  "failed_file_names": {
    "bad_scan.pdf": "Error extracting text from PDF"
  }
}
```

## Code Examples

### Custom OCR Processing

If you need custom OCR processing:

```go
import gosseract "github.com/otiai10/gosseract/v2"

func customOCR(imageBytes []byte) (string, error) {
    client := gosseract.NewClient()
    defer client.Close()

    client.SetImageFromBytes(imageBytes)
    client.SetLanguage("eng")

    // Custom configuration
    client.SetPageSegMode(gosseract.PSM_AUTO)
    client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 ")

    return client.Text()
}
```

## Future Enhancements

Potential improvements:

1. **Multi-language support**: Add more Tesseract language packs
2. **OCR confidence scores**: Return quality metrics
3. **Image preprocessing**: Auto-enhance images before OCR
4. **Parallel page processing**: Process PDF pages concurrently
5. **Caching**: Cache OCR results for duplicate files
6. **Alternative engines**: Support for other OCR engines (e.g., Cloud Vision API)

## Related Files

### Modified Files

- `chunk/fileprocessing.go` - Core OCR implementation
  - `isTextSufficientForProcessing()` (line 189)
  - `extractTextFromPDFWithOCR()` (line 220)
  - `extractTextFromImage()` (line 294)
  - `ExtractTextFromPDF()` (line 326)

- `Dockerfile` - Tesseract installation
  - Builder stage: lines 11-12
  - Runtime stage: lines 46-47

### Dependencies

- `github.com/otiai10/gosseract/v2` - Go bindings for Tesseract
- `github.com/gen2brain/go-fitz` - PDF rendering (already present)

## Summary

The OCR feature provides:

✅ **Automatic detection** - No configuration needed
✅ **Smart fallback** - Only uses OCR when necessary
✅ **Wide format support** - PDFs and all common images
✅ **Quality validation** - Ensures meaningful text extraction
✅ **Comprehensive logging** - Easy debugging and monitoring
✅ **Production ready** - Integrated with Docker deployment

Users can now upload scanned documents and images with confidence that text will be extracted accurately and efficiently.

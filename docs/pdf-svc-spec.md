# PDF Service (pdf-svc) Specification

## Overview
The PDF Service provides comprehensive PDF document processing capabilities including generation, manipulation, conversion, and analysis. It supports both in-process and HTTP-based communication patterns within the hybrid API ecosystem.

## Port Assignment
- **Default Port**: 8083
- **Health Check**: 8083/health

## Core Features

### 1. PDF Generation
- HTML to PDF conversion
- Template-based PDF generation
- Dynamic content insertion
- Custom styling and layouts
- Multi-page document support

### 2. PDF Manipulation
- Page extraction and insertion
- Document merging and splitting
- Watermark addition
- Password protection and encryption
- Metadata modification

### 3. PDF Conversion
- PDF to image conversion (PNG, JPEG)
- PDF to text extraction
- PDF to HTML conversion
- Image to PDF conversion

### 4. PDF Analysis
- Text extraction and searching
- Form field detection
- Document structure analysis
- Accessibility compliance checking

## API Endpoints

### PDF Generation

#### POST /pdf/generate/html
Generate PDF from HTML content.

**Request Body:**
```json
{
  "html": "<html><body><h1>Hello World</h1></body></html>",
  "options": {
    "pageSize": "A4",
    "orientation": "portrait",
    "margin": {
      "top": "1cm",
      "right": "1cm",
      "bottom": "1cm",
      "left": "1cm"
    },
    "displayHeaderFooter": true,
    "headerTemplate": "<div style='text-align: center; font-size: 10px;'>Page <span class='pageNumber'></span></div>",
    "footerTemplate": "<div style='text-align: center; font-size: 10px;'>Generated on <span class='date'></span></div>",
    "printBackground": true,
    "preferCSSPageSize": false
  },
  "filename": "document.pdf"
}
```

**Response (200):**
```json
{
  "success": true,
  "filename": "document.pdf",
  "size": 1024,
  "downloadUrl": "/pdf/download/doc-uuid",
  "pages": 1,
  "generatedAt": "2025-09-22T10:00:00Z"
}
```

#### POST /pdf/generate/template
Generate PDF from template with data.

**Request Body:**
```json
{
  "templateId": "invoice-template",
  "data": {
    "invoiceNumber": "INV-001",
    "customerName": "John Doe",
    "items": [
      {
        "description": "Product A",
        "quantity": 2,
        "price": 50.00
      }
    ],
    "total": 100.00
  },
  "options": {
    "pageSize": "A4",
    "orientation": "portrait"
  }
}
```

#### POST /pdf/generate/url
Generate PDF from URL.

**Request Body:**
```json
{
  "url": "https://example.com/page",
  "options": {
    "pageSize": "A4",
    "waitTime": 2000,
    "viewportWidth": 1920,
    "viewportHeight": 1080
  }
}
```

### PDF Manipulation

#### POST /pdf/merge
Merge multiple PDF documents.

**Request Body (multipart/form-data):**
```
files: [file1.pdf, file2.pdf, file3.pdf]
options: {
  "outputFilename": "merged-document.pdf",
  "bookmarks": true
}
```

#### POST /pdf/split
Split PDF document into separate pages or ranges.

**Request Body (multipart/form-data):**
```
file: document.pdf
options: {
  "splitType": "pages", // "pages", "ranges", "size"
  "ranges": ["1-3", "5-7"],
  "outputPrefix": "split-doc"
}
```

#### POST /pdf/extract-pages
Extract specific pages from PDF.

**Request Body (multipart/form-data):**
```
file: document.pdf
options: {
  "pages": [1, 3, 5],
  "outputFilename": "extracted-pages.pdf"
}
```

#### POST /pdf/add-watermark
Add watermark to PDF document.

**Request Body (multipart/form-data):**
```
file: document.pdf
watermark: watermark.png
options: {
  "position": "center", // "center", "top-left", "top-right", "bottom-left", "bottom-right"
  "opacity": 0.5,
  "scale": 0.3,
  "rotation": 45
}
```

#### POST /pdf/protect
Add password protection to PDF.

**Request Body (multipart/form-data):**
```
file: document.pdf
options: {
  "userPassword": "user123",
  "ownerPassword": "owner123",
  "permissions": {
    "allowPrinting": true,
    "allowCopyContent": false,
    "allowEditContent": false,
    "allowEditAnnotations": false
  },
  "encryptionLevel": 128
}
```

### PDF Conversion

#### POST /pdf/convert/to-images
Convert PDF pages to images.

**Request Body (multipart/form-data):**
```
file: document.pdf
options: {
  "format": "png", // "png", "jpeg", "webp"
  "quality": 90,
  "dpi": 300,
  "pages": [1, 2, 3], // empty array for all pages
  "outputPrefix": "page"
}
```

**Response (200):**
```json
{
  "success": true,
  "images": [
    {
      "page": 1,
      "filename": "page-1.png",
      "downloadUrl": "/pdf/download/img-uuid-1",
      "width": 2480,
      "height": 3508
    }
  ]
}
```

#### POST /pdf/convert/to-text
Extract text content from PDF.

**Request Body (multipart/form-data):**
```
file: document.pdf
options: {
  "pages": [1, 2], // empty for all pages
  "preserveLayout": true,
  "includeMetadata": true
}
```

**Response (200):**
```json
{
  "success": true,
  "text": "Extracted text content...",
  "pages": [
    {
      "page": 1,
      "text": "Page 1 content..."
    }
  ],
  "metadata": {
    "title": "Document Title",
    "author": "John Doe",
    "subject": "Document Subject",
    "keywords": "pdf, extraction",
    "creator": "PDF Creator",
    "producer": "PDF Producer",
    "creationDate": "2025-09-22T10:00:00Z",
    "modificationDate": "2025-09-22T10:00:00Z"
  }
}
```

#### POST /pdf/convert/images-to-pdf
Convert multiple images to PDF.

**Request Body (multipart/form-data):**
```
files: [image1.jpg, image2.png, image3.jpg]
options: {
  "pageSize": "A4",
  "orientation": "portrait",
  "fitToPage": true,
  "outputFilename": "images-combined.pdf"
}
```

### PDF Analysis

#### POST /pdf/analyze
Analyze PDF document structure and content.

**Request Body (multipart/form-data):**
```
file: document.pdf
options: {
  "includeText": true,
  "includeForms": true,
  "includeImages": true,
  "includeMetadata": true
}
```

**Response (200):**
```json
{
  "success": true,
  "analysis": {
    "pageCount": 10,
    "hasText": true,
    "hasForms": false,
    "hasImages": true,
    "isPasswordProtected": false,
    "fileSize": 2048000,
    "version": "1.4",
    "pages": [
      {
        "page": 1,
        "width": 595.28,
        "height": 841.89,
        "rotation": 0,
        "textBlocks": 5,
        "imageCount": 2,
        "formFields": 0
      }
    ],
    "fonts": [
      {
        "name": "Arial",
        "type": "TrueType",
        "embedded": true
      }
    ],
    "accessibility": {
      "hasStructure": false,
      "hasAlternativeText": false,
      "isTagged": false
    }
  }
}
```

#### POST /pdf/search
Search for text within PDF document.

**Request Body (multipart/form-data):**
```
file: document.pdf
options: {
  "query": "search term",
  "caseSensitive": false,
  "wholeWords": false,
  "includeContext": true
}
```

### Template Management

#### GET /pdf/templates
List available PDF templates.

**Response (200):**
```json
{
  "templates": [
    {
      "id": "invoice-template",
      "name": "Invoice Template",
      "description": "Standard invoice template with company branding",
      "version": "1.0",
      "fields": [
        {
          "name": "invoiceNumber",
          "type": "string",
          "required": true
        },
        {
          "name": "customerName",
          "type": "string",
          "required": true
        }
      ],
      "createdAt": "2025-09-22T10:00:00Z"
    }
  ]
}
```

#### POST /pdf/templates
Create new PDF template.

#### PUT /pdf/templates/:id
Update existing template.

#### DELETE /pdf/templates/:id
Delete template.

### File Management

#### GET /pdf/download/:id
Download generated PDF or converted file.

#### DELETE /pdf/files/:id
Delete temporary file from server.

#### GET /pdf/files
List temporary files (admin only).

## Data Models

### PDF Generation Options
```go
type GenerationOptions struct {
    PageSize             string            `json:"pageSize"` // A4, Letter, Legal, A3, A5
    Orientation          string            `json:"orientation"` // portrait, landscape
    Margin               MarginOptions     `json:"margin"`
    DisplayHeaderFooter  bool              `json:"displayHeaderFooter"`
    HeaderTemplate       string            `json:"headerTemplate"`
    FooterTemplate       string            `json:"footerTemplate"`
    PrintBackground      bool              `json:"printBackground"`
    PreferCSSPageSize    bool              `json:"preferCSSPageSize"`
    WaitTime            int               `json:"waitTime"` // milliseconds
    ViewportWidth       int               `json:"viewportWidth"`
    ViewportHeight      int               `json:"viewportHeight"`
    Scale               float64           `json:"scale"`
    CustomCSS           string            `json:"customCSS"`
    JavaScript          string            `json:"javascript"`
}

type MarginOptions struct {
    Top    string `json:"top"`
    Right  string `json:"right"`
    Bottom string `json:"bottom"`
    Left   string `json:"left"`
}
```

### Template Model
```go
type Template struct {
    ID          string           `json:"id" db:"id"`
    Name        string           `json:"name" db:"name"`
    Description string           `json:"description" db:"description"`
    Content     string           `json:"content" db:"content"` // HTML template
    Version     string           `json:"version" db:"version"`
    Fields      []TemplateField  `json:"fields" db:"fields"`
    Options     GenerationOptions `json:"options" db:"options"`
    CreatedAt   time.Time        `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time        `json:"updatedAt" db:"updated_at"`
    CreatedBy   string           `json:"createdBy" db:"created_by"`
}

type TemplateField struct {
    Name        string `json:"name"`
    Type        string `json:"type"` // string, number, boolean, array, object
    Required    bool   `json:"required"`
    Description string `json:"description"`
    Default     interface{} `json:"default"`
}
```

### File Model
```go
type GeneratedFile struct {
    ID          string    `json:"id"`
    Filename    string    `json:"filename"`
    OriginalName string   `json:"originalName"`
    MimeType    string    `json:"mimeType"`
    Size        int64     `json:"size"`
    Path        string    `json:"path"`
    DownloadURL string    `json:"downloadUrl"`
    ExpiresAt   time.Time `json:"expiresAt"`
    CreatedAt   time.Time `json:"createdAt"`
    CreatedBy   string    `json:"createdBy"`
    Metadata    map[string]interface{} `json:"metadata"`
}
```

## SDK Interface

### Direct Mode (In-Process)
```go
type PDFSdk struct {
    GenerateFromHTML     func(html string, options GenerationOptions) (*GeneratedFile, error)
    GenerateFromTemplate func(templateID string, data interface{}, options GenerationOptions) (*GeneratedFile, error)
    GenerateFromURL      func(url string, options GenerationOptions) (*GeneratedFile, error)
    MergePDFs           func(files [][]byte, options MergeOptions) (*GeneratedFile, error)
    SplitPDF            func(file []byte, options SplitOptions) ([]*GeneratedFile, error)
    ConvertToImages     func(file []byte, options ConversionOptions) ([]*GeneratedFile, error)
    ExtractText         func(file []byte, options ExtractionOptions) (*TextExtractionResult, error)
    AnalyzePDF          func(file []byte, options AnalysisOptions) (*PDFAnalysis, error)
    AddWatermark        func(file []byte, watermark []byte, options WatermarkOptions) (*GeneratedFile, error)
    ProtectPDF          func(file []byte, options ProtectionOptions) (*GeneratedFile, error)
}

func NewPDFSdk(mode string) *PDFSdk
```

### HTTP Mode (Network)
Standard HTTP multipart file upload and JSON response handling.

## Environment Variables

```env
# Server Configuration
PDF_SERVICE_PORT=8083
PDF_SERVICE_HOST=0.0.0.0

# Storage Configuration
PDF_STORAGE_PATH=/tmp/pdf-service
PDF_TEMP_DIR=/tmp/pdf-temp
PDF_MAX_FILE_SIZE=100MB
PDF_FILE_RETENTION=24h

# Processing Configuration
PDF_MAX_PAGES=1000
PDF_MAX_CONCURRENT_JOBS=10
PDF_PROCESSING_TIMEOUT=300s
PDF_MEMORY_LIMIT=2GB

# Template Configuration
PDF_TEMPLATE_DIR=/app/templates
PDF_DEFAULT_FONT=Arial
PDF_FONT_DIR=/app/fonts

# Chrome/Chromium Configuration (for HTML to PDF)
CHROME_EXECUTABLE_PATH=/usr/bin/chromium
CHROME_USER_DATA_DIR=/tmp/chrome-data
CHROME_HEADLESS=true
CHROME_NO_SANDBOX=true

# Database (for templates and metadata)
PDF_DB_HOST=localhost
PDF_DB_PORT=5432
PDF_DB_NAME=pdf_db
PDF_DB_USER=pdf_user
PDF_DB_PASSWORD=pdf_password

# Security
PDF_MAX_UPLOAD_SIZE=50MB
PDF_ALLOWED_EXTENSIONS=pdf,html,png,jpg,jpeg,gif
PDF_SCAN_UPLOADS=true

# Monitoring
PDF_METRICS_ENABLED=true
PDF_LOG_LEVEL=info
```

## Dependencies

```go
// Add to go.mod
require (
    github.com/chromedp/chromedp v0.9.2
    github.com/pdfcpu/pdfcpu v0.5.0
    github.com/gen2brain/go-fitz v1.23.7
    github.com/jung-kurt/gofpdf v1.16.2
    github.com/signintech/gopdf v0.19.0
    github.com/google/uuid v1.3.1
    github.com/disintegration/imaging v1.6.2
    gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)
```

## Performance Requirements

- PDF generation time: < 5s for typical documents
- File processing: < 30s for documents up to 100 pages
- Concurrent processing: 10 simultaneous jobs
- Memory usage: < 2GB per worker process
- Supported file sizes: up to 100MB

## Security Considerations

### File Validation
- MIME type verification
- File size limits
- Malicious content scanning
- Extension whitelist

### Input Sanitization
- HTML content sanitization
- JavaScript execution control
- Template injection prevention
- Path traversal protection

### Output Security
- Temporary file cleanup
- Secure file storage
- Access control for downloads
- File expiration policies

## Error Handling

### HTTP Status Codes
- 200: Success
- 400: Bad Request (invalid parameters)
- 413: Payload Too Large (file size exceeded)
- 415: Unsupported Media Type
- 422: Unprocessable Entity (processing failed)
- 500: Internal Server Error

### Error Response Format
```json
{
  "error": {
    "code": "PDF_GENERATION_FAILED",
    "message": "Failed to generate PDF from HTML",
    "details": {
      "reason": "Chrome process crashed",
      "html_length": 50000,
      "processing_time": 120
    }
  }
}
```

## Testing Requirements

### Unit Tests
- HTML to PDF conversion
- PDF manipulation functions
- Template rendering
- File validation

### Integration Tests
- End-to-end PDF generation
- Multi-file processing
- Template system
- Error handling

### Performance Tests
- Large file processing
- Concurrent request handling
- Memory usage optimization
- Processing time benchmarks

## Monitoring & Logging

### Metrics to Track
```
# Processing metrics
pdf_generation_duration_seconds{type}
pdf_processing_jobs_total{type, status}
pdf_file_size_bytes{type}
pdf_pages_processed_total

# Resource metrics
pdf_memory_usage_bytes
pdf_temp_files_count
pdf_active_jobs_count

# Error metrics
pdf_errors_total{type, reason}
pdf_timeouts_total{type}
```

### Audit Logging
- File upload events
- PDF generation requests
- Template usage
- Error occurrences
- Performance metrics

## Migration Path

### Phase 1: Basic PDF Generation
- HTML to PDF conversion
- Simple template system
- File download API

### Phase 2: Advanced Processing
- PDF manipulation (merge, split)
- Image conversion
- Text extraction

### Phase 3: Enterprise Features
- Advanced templates
- Batch processing
- Webhook notifications
- Advanced security features
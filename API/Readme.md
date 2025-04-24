# Motion-Index Document Processing System

## Overview

Motion-Index is a comprehensive system for processing, analyzing, and indexing legal documents. It automates the extraction of text from various document formats, classifies them using AI, and makes them searchable through Elasticsearch, with optional storage in AWS S3.

Key features:

1. **Text Extraction**: Extracts text from multiple document formats using the textract library
2. **Document Classification**: Uses OpenAI's LLM to classify legal documents and extract metadata
3. **S3 Storage**: Uploads processed documents to AWS S3 with organized folder structure (optional)
4. **Elasticsearch Indexing**: Creates searchable indices of document content and metadata
5. **Parallel Processing**: Efficiently processes large document collections using multithreading

## Project Structure

```
API/
├── Dockerfile                # Container definition for the application
├── docker-compose.yml        # Multi-container Docker setup
├── docker-run.sh             # Helper script to run the Docker environment
├── requirements.txt          # Python dependencies
├── .env                      # Environment variables (create from .env.template)
├── src/
│   ├── main.py               # Entry point and CLI argument handling
│   ├── core/
│   │   └── document_processor.py  # Main processing pipeline
│   ├── handlers/
│   │   ├── document_classifier.py  # LLM-based document classification
│   │   ├── elasticsearch_handler.py # Elasticsearch operations
│   │   ├── file_processor.py       # Text extraction from documents
│   │   └── s3_handler.py           # AWS S3 operations
│   ├── models/
│   │   ├── document.py             # Document data model
│   │   └── metadata.py             # Metadata data model
│   └── utils/
│       └── constants.py            # Application constants and configuration
```

## Detailed Component Descriptions

### Core Components

- **document_processor.py**: Orchestrates the entire document processing workflow, managing parallel execution, and coordinating between the various handlers.

### Handlers

- **file_processor.py**: Handles the extraction of text from various document formats using the textract library. Includes file validation, size checking, and format support verification.

- **document_classifier.py**: Uses OpenAI's GPT model to analyze document content and classify it into appropriate legal document types. Extracts metadata such as case numbers, parties, and dates.

- **elasticsearch_handler.py**: Manages connections to Elasticsearch, creates and updates indices, and handles bulk document indexing operations.

- **s3_handler.py**: Handles uploading documents to AWS S3, organizing them in a structured manner, and generating accessible URIs.

### Models

- **document.py**: Defines the Document class that represents a processed document with its text content, metadata, and classification.

- **metadata.py**: Defines the Metadata class that stores extracted information about documents such as title, author, case information, etc.

### Utilities

- **constants.py**: Contains application-wide constants including supported file formats, document types, and default configuration values.

## Processing Flow

1. **Document Discovery**:
   - Recursively scans the specified directory for documents
   - Filters files based on extension and size constraints
   - Skips unsupported formats and oversized files

2. **Text Extraction**:
   - Validates each file can be processed (correct format, size within limits)
   - Extracts text content using the textract library
   - Handles encoding issues and format-specific extraction challenges

3. **Document Classification**:
   - Analyzes extracted text using OpenAI's GPT model
   - Classifies the document type (motion, brief, order, etc.)
   - Extracts relevant metadata (case numbers, parties, dates, etc.)

4. **S3 Upload** (Optional):
   - Uploads the original document to AWS S3
   - Organizes files in a structured folder hierarchy
   - Generates and stores S3 URIs for future reference

5. **Elasticsearch Indexing**:
   - Creates document records with extracted text, classification, and metadata
   - Indexes documents in Elasticsearch for efficient searching
   - Handles batch processing for performance optimization

## Setup and Installation

### Prerequisites

- Python 3.9+
- Docker and Docker Compose (for containerized deployment)
- Elasticsearch server (local or cloud-based like Elastic Cloud)
- AWS S3 bucket (optional, for document storage)
- OpenAI API key (for document classification)

### Installation

#### Local Installation

1. Clone the repository

2. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

3. Configure environment variables by copying the template:
   ```bash
   cp .env.template .env
   ```
   Then edit `.env` with your API keys and configuration.

#### Docker Installation (Recommended)

1. Clone the repository

2. Configure environment variables in `.env`

3. Run the Docker setup:
   ```bash
   ./docker-run.sh
   ```

## Usage

### Basic Usage

#### Local Execution

Process a directory of documents:

```bash
python -m src.main /path/to/documents
```

#### Docker Execution

```bash
./docker-run.sh
```

This will process documents in the `./data` directory (created automatically if it doesn't exist).

### Command Line Options

#### Elasticsearch Options

- `--es-host`: Elasticsearch host (default: localhost)
- `--es-port`: Elasticsearch port (default: 9200)
- `--es-index`: Elasticsearch index name (default: documents)
- `--es-username`: Elasticsearch username
- `--es-password`: Elasticsearch password
- `--es-api-key`: Elasticsearch API key (for Elastic Cloud)
- `--es-cloud-id`: Elasticsearch Cloud ID (for Elastic Cloud)
- `--es-ssl`: Use SSL for Elasticsearch connection

#### S3 Options

- `--s3-bucket`: S3 bucket name for document storage
- `--s3-region`: AWS region for S3
- `--s3-access-key`: AWS access key ID
- `--s3-secret-key`: AWS secret access key

#### Processing Options

- `--workers`: Number of worker threads (default: 4)
- `--batch-size`: Batch size for indexing (default: 100)
- `--max-file-size`: Maximum file size in MB (default: 50)
- `--extensions`: List of file extensions to process (e.g., pdf docx txt)
- `--process-all`: Process all files regardless of extension
- `--force`: Process files even if they already exist in the index
- `--no-llm`: Disable LLM-based document classification

## Supported Document Types

The system currently supports the following document formats (based on textract capabilities):

- PDF (.pdf)
- Microsoft Word (.docx, .doc)
- Text (.txt)
- Rich Text Format (.rtf)
- OpenDocument Text (.odt)
- HTML (.html, .htm)
- PowerPoint (.pptx)
- Excel (.xlsx, .xls)
- CSV (.csv)
- Images (.png, .jpg, .jpeg, .tiff, .tif)
- Email (.eml, .msg)
- Audio transcription (.mp3, .wav, .ogg)
- And more formats supported by textract

**Note:** Some formats like WordPerfect (.wpd) are not supported by the textract library and will be skipped during processing.

## Legal Document Classification

The system classifies documents into the following categories:

- Motion
- Petition
- Order
- Brief
- Report
- Exhibit
- Memorandum
- Response
- Opposition
- Complaint
- Answer
- Discovery Request
- Discovery Response
- Notice
- Declaration
- Affidavit
- Judgment
- Transcript
- Settlement Agreement

## Error Handling and Logging

The system includes comprehensive error handling and logging:

- Files that are too large are skipped with a warning
- Unsupported file formats are identified and skipped
- Processing errors are logged with detailed information
- Statistics are collected and reported at the end of processing

## Performance Considerations

- **File Size Limits**: Default maximum file size is 50MB, configurable via `--max-file-size`
- **Parallel Processing**: Multiple worker threads process documents concurrently
- **Batch Indexing**: Documents are indexed in batches for better performance
- **Skip Existing**: By default, documents already in the index are skipped (override with `--force`)

## Environment Variables

All command-line options can also be set via environment variables in the `.env` file:

```
# Elasticsearch settings
ES_HOST=localhost
ES_PORT=9200
ES_INDEX=documents
ES_USERNAME=
ES_PASSWORD=
ES_API_KEY=
ES_CLOUD_ID=
ES_USE_SSL=True

# AWS S3 settings
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=us-east-1
S3_BUCKET_NAME=

# OpenAI settings
OPENAI_API_KEY=

# Processing settings
MAX_WORKERS=4
BATCH_SIZE=100
MAX_FILE_SIZE=50
```

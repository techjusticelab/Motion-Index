# Motion-Index Document Processing System

Note: doc strings in this repo and this readme itself have been edited/written by LLMs.

A system for processing legal documents that:
1. Extracts text from various document formats
2. Classifies documents using AI/LLM
3. Uploads documents to AWS S3
4. Indexes document metadata in Elasticsearch

## Setup and Installation

### Prerequisites
- Python 3.8+
- Elasticsearch server (local or remote)
- AWS S3 bucket (optional, for document storage)
- OpenAI API key (for document classification)

### Installation

1. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Configure environment variables by copying the template:
   ```bash
   cp .env.template .env
   ```
   Then edit `.env` with your API keys and configuration.

## Usage

### Basic Usage

Process a directory of documents:

```bash
python -m src.main /path/to/documents
```

### Command Line Options

#### Elasticsearch Options
- `--es-host`: Elasticsearch host (default: localhost)
- `--es-port`: Elasticsearch port (default: 9200)
- `--es-index`: Elasticsearch index name (default: documents)
- `--es-username`: Elasticsearch username
- `--es-password`: Elasticsearch password
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

## Project Structure

### Core Components

- `src/core/document_processor.py`: Main document processing pipeline

### Handlers

- `src/handlers/file_processor.py`: Extracts text from document files
- `src/handlers/document_classifier.py`: Classifies documents using OpenAI
- `src/handlers/s3_handler.py`: Manages uploads to AWS S3
- `src/handlers/elasticsearch_handler.py`: Indexes documents in Elasticsearch

### Models

- `src/models/document.py`: Document data model
- `src/models/metadata.py`: Metadata data model

### Utilities

- `src/utils/constants.py`: Application constants and configuration

## Processing Flow

1. **Document Discovery**: Find all files in the specified directory
2. **Text Extraction**: Extract text content from each document
3. **Document Classification**: Use LLM to classify documents and extract metadata
4. **S3 Upload**: Upload documents to S3 with organized folder structure
5. **Elasticsearch Indexing**: Index document metadata and references in Elasticsearch

## Supported Document Types

- PDF (.pdf)
- Microsoft Word (.docx, .doc)
- Text (.txt)
- Rich Text Format (.rtf)
- OpenDocument Text (.odt)
- HTML (.html, .htm)
- WordPerfect (.wpd, .wp, .wp5)
- PowerPoint (.pptx, .ppt)
- Excel (.xlsx, .xls)

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

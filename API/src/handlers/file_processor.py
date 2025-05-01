"""
Advanced file processor for extracting text from various document formats.
"""
import os
import logging
import hashlib
import mimetypes
import tempfile
import subprocess
import concurrent.futures
import shutil
from pathlib import Path
from typing import Optional, Dict, List, Tuple, Any, Union
import time
import re
import threading

# Optional dependencies for advanced features
try:
    import pytesseract
    from pdf2image import convert_from_path
    ENABLE_OCR = True
except ImportError:
    ENABLE_OCR = False

try:
    import tabula
    import pandas as pd
    ENABLE_TABLE_EXTRACTION = True
except ImportError:
    ENABLE_TABLE_EXTRACTION = False

try:
    from pdfminer.high_level import extract_text as pdfminer_extract_text
    ENABLE_PDFMINER = True
except ImportError:
    ENABLE_PDFMINER = False

import textract

from src.utils.constants import MAX_FILE_SIZE_DEFAULT, SUPPORTED_FORMATS, MIME_TYPES

# Configure logging
logger = logging.getLogger("file_processor")

# Add missing mimetypes
for ext, mime_type in MIME_TYPES.items():
    mimetypes.add_type(mime_type, ext)

# Add specialized mimetypes
mimetypes.add_type('application/rtf', '.rtf')
mimetypes.add_type('application/vnd.ms-outlook', '.msg')
mimetypes.add_type('message/rfc822', '.eml')


class FileProcessor:
    """
    Advanced processor for extracting text and metadata from various document formats.
    """
    
    def __init__(
        self, 
        max_file_size: int = MAX_FILE_SIZE_DEFAULT,
        enable_ocr: bool = ENABLE_OCR,
        ocr_language: str = "eng",
        enable_table_extraction: bool = ENABLE_TABLE_EXTRACTION,
        max_workers: int = 4,
        temp_dir: Optional[str] = None,
        extraction_timeout: int = 300,
        ocr_dpi: int = 300,
        detect_rotation: bool = True,
        extraction_fallbacks: bool = True
    ):
        """
        Initialize the processor with enhanced options.
        
        Args:
            max_file_size: Maximum file size in bytes (default: 50MB)
            enable_ocr: Whether to enable OCR for images and scanned PDFs
            ocr_language: Language for OCR (default: English)
            enable_table_extraction: Whether to extract tables from documents
            max_workers: Maximum number of worker threads for parallel processing
            temp_dir: Custom temporary directory for file processing
            extraction_timeout: Timeout for text extraction in seconds
            ocr_dpi: DPI for OCR image processing (higher is more accurate but slower)
            detect_rotation: Whether to detect and correct page rotation during OCR
            extraction_fallbacks: Whether to try multiple extraction methods
        """
        self.max_file_size = max_file_size
        self.enable_ocr = enable_ocr and ENABLE_OCR
        self.ocr_language = ocr_language
        self.enable_table_extraction = enable_table_extraction and ENABLE_TABLE_EXTRACTION
        self.max_workers = max_workers
        self.temp_dir = temp_dir
        self.extraction_timeout = extraction_timeout
        self.ocr_dpi = ocr_dpi
        self.detect_rotation = detect_rotation
        self.extraction_fallbacks = extraction_fallbacks
        
        # Create a thread pool for parallel processing
        self.thread_pool = concurrent.futures.ThreadPoolExecutor(max_workers=max_workers)
        
        # Thread lock for OCR operations
        self.ocr_lock = threading.Lock()
        
        # Validate OCR setup if enabled
        if self.enable_ocr and not ENABLE_OCR:
            logger.warning("OCR requested but dependencies not available. Install pytesseract and pdf2image.")
            self.enable_ocr = False
            
        # Validate table extraction setup if enabled
        if self.enable_table_extraction and not ENABLE_TABLE_EXTRACTION:
            logger.warning("Table extraction requested but dependencies not available. Install tabula-py and pandas.")
            self.enable_table_extraction = False
            
        # Set up the textract methods to use
        self.textract_methods = {
            '.pdf': ['pdfminer', 'tesseract'] if self.enable_ocr else ['pdfminer'],
            '.docx': ['docx2txt'],
            '.doc': ['antiword', 'textract'],
            '.rtf': ['unrtf'],
            '.txt': ['plaintext'],
            '.html': ['html2text'],
            '.htm': ['html2text'],
            '.eml': ['email_message'],
            '.msg': ['msg_extractor'],
            '.xls': ['xlrd'],
            '.xlsx': ['xlrd'],
            '.csv': ['plaintext'],
            '.json': ['plaintext'],
            '.jpg': ['tesseract'] if self.enable_ocr else None,
            '.jpeg': ['tesseract'] if self.enable_ocr else None,
            '.png': ['tesseract'] if self.enable_ocr else None,
            '.tif': ['tesseract'] if self.enable_ocr else None,
            '.tiff': ['tesseract'] if self.enable_ocr else None,
        }
        
        # Log initialization information
        logger.info(f"Initialized FileProcessor with max_file_size={max_file_size}, OCR={'enabled' if self.enable_ocr else 'disabled'}")
        if self.enable_ocr:
            try:
                ocr_version = subprocess.run(["tesseract", "--version"], capture_output=True, text=True).stdout.split("\n")[0]
                logger.info(f"Using Tesseract OCR: {ocr_version}")
            except:
                logger.warning("Could not determine Tesseract version")
    
    def __del__(self):
        """Clean up resources when the processor is destroyed."""
        try:
            self.thread_pool.shutdown(wait=False)
        except:
            pass
    
    def get_file_hash(self, file_path: str) -> str:
        """
        Generate a hash of file contents to use as a unique identifier.
        
        Args:
            file_path: Path to the file
            
        Returns:
            MD5 hash of the file contents
        """
        try:
            # For efficiency with large files, read in chunks
            md5_hash = hashlib.md5()
            with open(file_path, 'rb') as f:
                for chunk in iter(lambda: f.read(4096), b""):
                    md5_hash.update(chunk)
            return md5_hash.hexdigest()
        except Exception as e:
            logger.error(f"Failed to hash file {file_path}: {e}")
            return ""
    
    def get_file_metadata(self, file_path: str) -> Dict[str, Any]:
        """
        Extract basic metadata from a file.
        
        Args:
            file_path: Path to the file
            
        Returns:
            Dictionary of file metadata
        """
        path = Path(file_path)
        
        metadata = {
            "filename": path.name,
            "extension": path.suffix.lower(),
            "size_bytes": path.stat().st_size if path.exists() else 0,
            "mime_type": mimetypes.guess_type(file_path)[0] or "application/octet-stream",
            "last_modified": time.ctime(path.stat().st_mtime) if path.exists() else None,
            "created": time.ctime(path.stat().st_ctime) if path.exists() else None,
        }
        
        # Calculate hash for smaller files (under 20MB) upfront
        if metadata["size_bytes"] < 20 * 1024 * 1024:
            metadata["hash"] = self.get_file_hash(file_path)
        
        # Try to extract more specific metadata based on file type
        try:
            if metadata["extension"] == ".pdf":
                pdf_metadata = self._extract_pdf_metadata(file_path)
                if pdf_metadata:
                    metadata.update(pdf_metadata)
                    
        except Exception as e:
            logger.warning(f"Error extracting extended metadata for {file_path}: {e}")
        
        return metadata
    
    def _extract_pdf_metadata(self, file_path: str) -> Dict[str, Any]:
        """Extract metadata from PDF documents."""
        metadata = {}
        try:
            # Try using PyPDF2 if available
            try:
                from PyPDF2 import PdfReader
                reader = PdfReader(file_path)
                if reader.metadata:
                    # Extract standard metadata fields
                    for key in ['Author', 'Creator', 'Producer', 'Subject', 'Title']:
                        if hasattr(reader.metadata, key) and getattr(reader.metadata, key):
                            metadata[key.lower()] = getattr(reader.metadata, key)
                            
                    # Extract creation and modification dates
                    if hasattr(reader.metadata, 'creation_date'):
                        metadata['pdf_creation_date'] = str(reader.metadata.creation_date)
                    if hasattr(reader.metadata, 'modification_date'):
                        metadata['pdf_modification_date'] = str(reader.metadata.modification_date)
                            
                # Count pages
                metadata['page_count'] = len(reader.pages)
                
                # Check if document is scanned
                metadata['is_scanned'] = self._is_scanned_pdf(reader)
                
                return metadata
            except ImportError:
                logger.debug("PyPDF2 not available, skipping detailed PDF metadata extraction")
                
        except Exception as e:
            logger.warning(f"Error extracting PDF metadata: {e}")
        
        return metadata
    
    def _is_scanned_pdf(self, pdf_reader) -> bool:
        """
        Determine if a PDF appears to be a scanned document.
        
        Args:
            pdf_reader: The PyPDF2 PdfReader object
            
        Returns:
            True if the PDF appears to be scanned, False otherwise
        """
        try:
            # Check first few pages for text
            pages_to_check = min(3, len(pdf_reader.pages))
            text_count = 0
            
            for i in range(pages_to_check):
                page = pdf_reader.pages[i]
                text = page.extract_text()
                
                if text and len(text.strip()) > 100:  # At least 100 chars of real text
                    text_count += 1
            
            # If less than half of checked pages have text, likely scanned
            return text_count < (pages_to_check / 2)
        except Exception as e:
            logger.debug(f"Error checking if PDF is scanned: {e}")
            # If we can't determine, default to False
            return False
    
    def can_process(self, file_path: str) -> bool:
        """
        Check if a file can be processed based on extension and size.
        
        Args:
            file_path: Path to the file
            
        Returns:
            True if the file can be processed, False otherwise
        """
        path = Path(file_path)
        
        # Check file size
        try:
            if path.stat().st_size > self.max_file_size:
                logger.warning(f"File too large to process: {file_path} ({path.stat().st_size} bytes)")
                return False
        except Exception as e:
            logger.error(f"Error checking file size {file_path}: {e}")
            return False
            
        # Check file extension
        ext = path.suffix.lower()
        
        # Skip WPD files as requested
        if ext in ['.wpd', '.wp', '.wp5', '.wps', '.wpt', '.wri']:
            logger.info(f"Skipping WordPerfect file: {file_path}")
            return False
        
        # Add support for image files if OCR is enabled
        if self.enable_ocr and ext.lower() in ['.jpg', '.jpeg', '.png', '.tiff', '.tif', '.bmp', '.gif']:
            return True
        
        # Check textract supported formats
        textract_supported = [
            '.csv','.eml', '.epub', '.gif', '.htm', '.html',
            '.jpeg', '.jpg', '.json', '.log', '.mp3', '.msg', '.odt', '.ogg',
            '.pdf', '.png', '.pptx', '.ps', '.psv', '.rtf', '.tab', '.tff',
            '.tif', '.tiff', '.tsv', '.txt', ".docx", '.wav', '.xls', '.xlsx'
        ]
        
        if ext in textract_supported or ext in SUPPORTED_FORMATS:
            return True
            
        # Try to determine MIME type
        mime_type, _ = mimetypes.guess_type(file_path)
        if mime_type and any(mime_type.startswith(t) for t in [
            'text/', 'application/pdf', 'application/msword', 
            'application/vnd.ms-', 'application/vnd.openxmlformats-',
            'application/rtf'
        ]):
            return True
        
        # Try more sophisticated detection for files without extensions
        if not ext:
            try:
                # Use file command to detect file type
                file_output = subprocess.run(
                    ["file", "--mime-type", "--brief", file_path], 
                    capture_output=True, 
                    text=True, 
                    timeout=5
                ).stdout.strip()
                
                return any(t in file_output for t in [
                    'text/', 'pdf', 'word', 'excel', 'powerpoint', 'rtf', 'xml', 'json'
                ])
            except:
                pass
                
        return False
    
    def _create_temp_dir(self) -> str:
        """Create a temporary directory for file processing."""
        if self.temp_dir and os.path.isdir(self.temp_dir):
            # Create a subdirectory within the specified temp directory
            temp_dir = tempfile.mkdtemp(dir=self.temp_dir)
        else:
            # Use system default temp directory
            temp_dir = tempfile.mkdtemp()
        
        return temp_dir
    
    def _clean_temp_dir(self, temp_dir: str):
        """Clean up a temporary directory."""
        try:
            if os.path.exists(temp_dir) and os.path.isdir(temp_dir):
                shutil.rmtree(temp_dir)
        except Exception as e:
            logger.warning(f"Error cleaning up temporary directory {temp_dir}: {e}")
    
    def _extract_text_with_ocr(self, file_path: str) -> str:
        """
        Extract text from an image or scanned document using OCR.
        
        Args:
            file_path: Path to the image or PDF file
            
        Returns:
            Extracted text
        """
        if not self.enable_ocr:
            return ""
            
        try:
            # Use a lock to prevent multiple OCR processes from running simultaneously
            # as Tesseract can be memory-intensive
            with self.ocr_lock:
                if file_path.lower().endswith('.pdf'):
                    # For PDFs, convert to images first
                    temp_dir = self._create_temp_dir()
                    try:
                        images = convert_from_path(
                            file_path, 
                            dpi=self.ocr_dpi, 
                            output_folder=temp_dir,
                            fmt="jpeg",
                            grayscale=True,
                            use_pdftocairo=True
                        )
                        
                        text = ""
                        for i, img in enumerate(images):
                            logger.debug(f"Running OCR on page {i+1}/{len(images)} of {file_path}")
                            
                            # Configure OCR options
                            ocr_config = f"--oem 1 --psm 3 -l {self.ocr_language}"
                            if self.detect_rotation:
                                ocr_config += " --rotate-pages"
                                
                            # Run OCR
                            page_text = pytesseract.image_to_string(
                                img, 
                                config=ocr_config
                            )
                            text += page_text + "\n\n"
                            
                        return text
                    finally:
                        self._clean_temp_dir(temp_dir)
                else:
                    # For direct image files
                    import PIL.Image
                    img = PIL.Image.open(file_path)
                    
                    # Configure OCR options
                    ocr_config = f"--oem 1 --psm 3 -l {self.ocr_language}"
                    if self.detect_rotation:
                        ocr_config += " --rotate-pages"
                        
                    # Run OCR
                    return pytesseract.image_to_string(
                        img, 
                        config=ocr_config
                    )
        except Exception as e:
            logger.error(f"OCR extraction failed for {file_path}: {e}")
            return ""
    
    def _extract_tables(self, file_path: str) -> List[Dict[str, Any]]:
        """
        Extract tables from a document.
        
        Args:
            file_path: Path to the document
            
        Returns:
            List of tables as dictionaries
        """
        if not self.enable_table_extraction:
            return []
            
        tables = []
        try:
            if file_path.lower().endswith('.pdf'):
                # Extract tables from PDF using tabula-py
                raw_tables = tabula.read_pdf(
                    file_path, 
                    pages='all', 
                    multiple_tables=True,
                    pandas_options={'header': None}
                )
                
                # Convert tables to dictionaries
                for i, table in enumerate(raw_tables):
                    if not table.empty:
                        # Convert DataFrame to dictionary
                        table_dict = {
                            'table_id': f"table_{i+1}",
                            'page_number': None,  # tabula doesn't provide page numbers easily
                            'rows': len(table),
                            'columns': len(table.columns),
                            'data': table.fillna('').to_dict('records')
                        }
                        tables.append(table_dict)
            
            # Could add support for other formats like XLSX, DOCX, etc.
                
        except Exception as e:
            logger.error(f"Table extraction failed for {file_path}: {e}")
            
        return tables
    
    def _extract_text_from_pdf_with_pdfminer(self, file_path: str) -> str:
        """
        Extract text from a PDF using pdfminer.
        
        Args:
            file_path: Path to the PDF file
            
        Returns:
            Extracted text
        """
        if not ENABLE_PDFMINER:
            return ""
            
        try:
            return pdfminer_extract_text(file_path)
        except Exception as e:
            logger.error(f"pdfminer extraction failed for {file_path}: {e}")
            return ""
    
    def _try_textract_extraction(self, file_path: str, method: Optional[str] = None) -> str:
        """
        Extract text using textract with specified method.
        
        Args:
            file_path: Path to the file
            method: Optional specific textract method to use
            
        Returns:
            Extracted text
        """
        try:
            # Build the extraction options
            options = {}
            if method:
                options['method'] = method
                
            # For PDFs, add specific options
            if file_path.lower().endswith('.pdf') and method == 'tesseract':
                options['language'] = self.ocr_language
                
            # Extract text
            text = textract.process(
                file_path,
                encoding='utf-8',
                **options
            ).decode('utf-8')
            
            return text
        except UnicodeDecodeError:
            # Handle encoding issues
            try:
                text = textract.process(
                    file_path,
                    encoding='utf-8',
                    **options
                ).decode('utf-8', errors='replace')
                
                logger.warning(f"Processed {file_path} with character replacement due to encoding issues")
                return text
            except Exception as e:
                logger.error(f"Textract extraction with method {method} failed for {file_path} after encoding retry: {e}")
                return ""
        except Exception as e:
            logger.error(f"Textract extraction with method {method} failed for {file_path}: {e}")
            return ""
    
    def extract_text(self, file_path: str) -> str:
        """
        Extract text from a document file using the most appropriate method.
        
        Args:
            file_path: Path to the document file
            
        Returns:
            Extracted text content
        """
        if not self.can_process(file_path):
            logger.warning(f"Cannot process file: {file_path}")
            return ""
        
        path = Path(file_path)
        ext = path.suffix.lower()
        
        # Check if file has a valid extension
        if not ext and self.enable_ocr:
            # Try to infer type for files without extension
            try:
                mime_type = subprocess.run(
                    ["file", "--mime-type", "--brief", file_path], 
                    capture_output=True, 
                    text=True, 
                    timeout=5
                ).stdout.strip()
                
                if "image/" in mime_type:
                    # It's an image, use OCR
                    logger.info(f"File without extension detected as image, using OCR: {file_path}")
                    return self._extract_text_with_ocr(file_path)
                elif "pdf" in mime_type:
                    # It's a PDF
                    ext = ".pdf"
                elif "text/" in mime_type:
                    # It's a text file
                    ext = ".txt"
                else:
                    return ext
            except:
                
                logger.warning(f"Could not determine MIME type for file without extension: {file_path}")
                pass
        
        print(f"Processing file: {file_path} with extension: {ext}")
        # Get methods to try for this file type
        methods_to_try = self.textract_methods.get(ext, [])
        
        # For image files, use OCR directly
        if ext.lower() in ['.jpg', '.jpeg', '.png', '.tiff', '.tif', '.bmp', '.gif'] and self.enable_ocr:
            return self._extract_text_with_ocr(file_path)
        
        # For PDFs, try both pdfminer and OCR if the file appears to be scanned
        if ext.lower() == '.pdf':
            # First try pdfminer for native text extraction
            text = self._extract_text_from_pdf_with_pdfminer(file_path)
            
            # If we got reasonable text, return it
            if text and len(text.strip()) > 100:
                return text
                
            # Otherwise, if OCR is enabled, it might be a scanned document
            if self.enable_ocr:
                logger.info(f"Falling back to OCR for possible scanned PDF: {file_path}")
                return self._extract_text_with_ocr(file_path)
        
        # For all other files, try textract with each method in sequence
        extracted_text = ""
        if methods_to_try:
            for method in methods_to_try:
                extracted_text = self._try_textract_extraction(file_path, method)
                if extracted_text and len(extracted_text.strip()) > 0:
                    return extracted_text
        
        # Fall back to default textract without specifying method
        if not extracted_text:
            try:
                extracted_text = textract.process(file_path, encoding='utf-8').decode('utf-8')
            except UnicodeDecodeError:
                try:
                    extracted_text = textract.process(file_path).decode('utf-8', errors='replace')
                except Exception as e:
                    logger.error(f"Default textract extraction failed for {file_path}: {e}")
            except Exception as e:
                logger.error(f"Default textract extraction failed for {file_path}: {e}")
        
        # Return whatever we've got, even if it's empty
        return extracted_text
    
    def process_file(self, file_path: str, extract_tables: bool = False) -> Dict[str, Any]:
        """
        Process a file to extract text, metadata, and optionally tables.
        
        Args:
            file_path: Path to the file
            extract_tables: Whether to extract tables (if supported)
            
        Returns:
            Dictionary with extracted text, metadata, and tables
        """
        result = {
            "success": False,
            "text": "",
            "metadata": {},
            "tables": [],
            "error": None
        }
        
        start_time = time.time()
        
        try:
            # Check if we can process this file
            if not self.can_process(file_path):
                result["error"] = f"Cannot process file: {file_path}"
                return result
            
            # Extract metadata
            result["metadata"] = self.get_file_metadata(file_path)
            
            # Extract text
            result["text"] = self.extract_text(file_path)
            
            # Extract tables if requested
            if extract_tables and self.enable_table_extraction:
                result["tables"] = self._extract_tables(file_path)
            
            # Set success flag
            result["success"] = True
            
            # Add processing stats
            processing_time = time.time() - start_time
            result["processing_time"] = processing_time
            
            logger.info(f"Successfully processed {file_path} in {processing_time:.2f} seconds")
            
            return result
            
        except Exception as e:
            logger.error(f"Error processing file {file_path}: {e}")
            result["error"] = str(e)
            return result
    
    def batch_process(self, file_paths: List[str], extract_tables: bool = False) -> List[Dict[str, Any]]:
        """
        Process multiple files in parallel.
        
        Args:
            file_paths: List of paths to files
            extract_tables: Whether to extract tables
            
        Returns:
            List of results for each file
        """
        start_time = time.time()
        logger.info(f"Starting batch processing of {len(file_paths)} files")
        
        # Filter out files we can't process
        processable_files = []
        results = []
        
        for file_path in file_paths:
            if not os.path.exists(file_path):
                results.append({
                    "file_path": file_path,
                    "success": False,
                    "error": "File not found"
                })
                continue
                
            if not self.can_process(file_path):
                results.append({
                    "file_path": file_path,
                    "success": False,
                    "error": "File type not supported or too large"
                })
                continue
                
            processable_files.append(file_path)
        
        # Process files in parallel using the thread pool
        futures = []
        for file_path in processable_files:
            future = self.thread_pool.submit(self.process_file, file_path, extract_tables)
            futures.append((file_path, future))
        
        # Collect results
        for file_path, future in futures:
            try:
                result = future.result(timeout=self.extraction_timeout)
                result["file_path"] = file_path
                results.append(result)
            except concurrent.futures.TimeoutError:
                logger.error(f"Processing timed out for {file_path}")
                results.append({
                    "file_path": file_path,
                    "success": False,
                    "error": "Processing timed out"
                })
            except Exception as e:
                logger.error(f"Error getting result for {file_path}: {e}")
                results.append({
                    "file_path": file_path,
                    "success": False,
                    "error": str(e)
                })
        
        # Log summary
        elapsed = time.time() - start_time
        success_count = sum(1 for r in results if r.get("success", False))
        logger.info(f"Batch processing complete: {success_count}/{len(file_paths)} files successful in {elapsed:.2f} seconds")
        
        return results
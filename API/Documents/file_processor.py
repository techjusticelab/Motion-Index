import os
from pathlib import Path
import textract
import mimetypes
import hashlib
import logging
from typing import Optional, Dict, Any, Tuple
import re
from datetime import datetime

# Import constants
from constants import (
    MAX_FILE_SIZE_DEFAULT, 
    SUPPORTED_FORMATS, 
    METADATA_PATTERNS,
    MIME_TYPES
)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("document_processing.log"),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger("file_processor")

# Add missing mimetypes
for ext, mime_type in MIME_TYPES.items():
    mimetypes.add_type(mime_type, ext)

class FileProcessor:
    """Handles reading and extracting text from various document formats"""
    
    def __init__(self, max_file_size: int = MAX_FILE_SIZE_DEFAULT):
        """Initialize the processor with a maximum file size (default 50MB)"""
        self.max_file_size = max_file_size
    
    def get_file_hash(self, file_path: str) -> str:
        """Generate a hash of file contents to use as a unique identifier"""
        try:
            with open(file_path, 'rb') as f:
                file_hash = hashlib.md5(f.read()).hexdigest()
            return file_hash
        except Exception as e:
            logger.error(f"Failed to hash file {file_path}: {e}")
            return ""
    
    def can_process(self, file_path: str) -> bool:
        """Check if a file can be processed based on extension and size"""
        path = Path(file_path)
        
        # Check file size
        try:
            if path.stat().st_size > self.max_file_size:
                logger.warning(f"File too large to process: {file_path}")
                return False
        except Exception as e:
            logger.error(f"Error checking file size {file_path}: {e}")
            return False
            
        # Check file extension
        ext = path.suffix.lower()
        if ext in SUPPORTED_FORMATS:
            return True
            
        # Try to determine MIME type
        mime_type, _ = mimetypes.guess_type(file_path)
        if mime_type and any(mime_type.startswith(t) for t in ['text/', 'application/pdf', 'application/msword', 'application/vnd.ms-', 'application/vnd.openxmlformats-']):
            return True
            
        return False
    
    def extract_text(self, file_path: str) -> Tuple[str, Dict[str, Any]]:
        """Extract text and basic metadata from a document file"""
        if not self.can_process(file_path):
            logger.warning(f"Cannot process file: {file_path}")
            return "", {}
            
        text = ""
        try:
            text = textract.process(file_path, encoding='utf-8').decode('utf-8')
        except Exception as e:
            logger.error(f"Error extracting text from {file_path}: {e}")
            return "", {}
            
        # Extract potential metadata from text
        metadata = self.extract_metadata_from_text(text, file_path)
        
        return text, metadata
    
    def extract_metadata_from_text(self, text: str, file_path: str) -> Dict[str, Any]:
        """Extract metadata from document text using regex patterns"""
        metadata = {
            'document_name': Path(file_path).name,
            'status': 'processed',
            'timestamp': datetime.now()
        }
        
        # Apply regex patterns to extract basic metadata
        for key, pattern in METADATA_PATTERNS.items():
            match = re.search(pattern, text, re.IGNORECASE)
            if match:
                metadata[key] = match.group(1).strip()
        
        # Determine subject from filename if not found in text
        if 'subject' not in metadata:
            filename = Path(file_path).stem
            # Remove numbers and special chars for a cleaner subject
            clean_name = re.sub(r'[\d_\-\.]+', ' ', filename).strip()
            if clean_name:
                metadata['subject'] = clean_name
        
        return metadata
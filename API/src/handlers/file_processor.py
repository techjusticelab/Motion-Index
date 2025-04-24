"""
File processor for extracting text from various document formats.
"""
import os
import logging
import hashlib
import mimetypes
import textract
from pathlib import Path
from typing import Optional

from src.utils.constants import MAX_FILE_SIZE_DEFAULT, SUPPORTED_FORMATS, MIME_TYPES

# Configure logging
logger = logging.getLogger("file_processor")

# Add missing mimetypes
for ext, mime_type in MIME_TYPES.items():
    mimetypes.add_type(mime_type, ext)


class FileProcessor:
    """
    Handles reading and extracting text from various document formats.
    """
    
    def __init__(self, max_file_size: int = MAX_FILE_SIZE_DEFAULT):
        """
        Initialize the processor with a maximum file size.
        
        Args:
            max_file_size: Maximum file size in bytes (default: 50MB)
        """
        self.max_file_size = max_file_size
    
    def get_file_hash(self, file_path: str) -> str:
        """
        Generate a hash of file contents to use as a unique identifier.
        
        Args:
            file_path: Path to the file
            
        Returns:
            MD5 hash of the file contents
        """
        try:
            with open(file_path, 'rb') as f:
                file_hash = hashlib.md5(f.read()).hexdigest()
            return file_hash
        except Exception as e:
            logger.error(f"Failed to hash file {file_path}: {e}")
            return ""
    
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
                logger.warning(f"File too large to process: {file_path}")
                return False
        except Exception as e:
            logger.error(f"Error checking file size {file_path}: {e}")
            return False
            
        # Check file extension
        ext = path.suffix.lower()
        
        # Check if the extension is in the list of supported formats by textract
        # This avoids attempting to process files that textract doesn't support
        textract_supported = [
            '.csv', '.doc', '.docx', '.eml', '.epub', '.gif', '.htm', '.html',
            '.jpeg', '.jpg', '.json', '.log', '.mp3', '.msg', '.odt', '.ogg',
            '.pdf', '.png', '.pptx', '.ps', '.psv', '.rtf', '.tab', '.tff',
            '.tif', '.tiff', '.tsv', '.txt', '.wav', '.xls', '.xlsx'
        ]
        
        if ext not in textract_supported:
            logger.warning(f"File extension {ext} not supported by textract: {file_path}")
            return False
        
        if ext in SUPPORTED_FORMATS:
            return True
            
        # Try to determine MIME type
        mime_type, _ = mimetypes.guess_type(file_path)
        if mime_type and any(mime_type.startswith(t) for t in [
            'text/', 'application/pdf', 'application/msword', 
            'application/vnd.ms-', 'application/vnd.openxmlformats-'
        ]):
            return True
            
        return False
    
    def extract_text(self, file_path: str) -> str:
        """
        Extract text from a document file.
        
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
        
        text = ""
        try:
            # Process the file using textract
            text = textract.process(file_path, encoding='utf-8').decode('utf-8')
        except UnicodeDecodeError:
            # Handle encoding issues
            try:
                text = textract.process(file_path).decode('utf-8', errors='replace')
                logger.warning(f"Processed {file_path} with character replacement due to encoding issues")
            except Exception as e:
                logger.error(f"Error extracting text from {file_path} after encoding retry: {e}")
                return ""
        except Exception as e:
            logger.error(f"Error extracting text from {file_path}: {e}")
            return ""
            
        return text

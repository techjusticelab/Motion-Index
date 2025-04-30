"""
File processor for extracting text from various document formats.
"""
import os
import logging
import hashlib
import mimetypes
import textract
import tempfile
import subprocess
from pathlib import Path
from typing import Optional

from src.utils.constants import MAX_FILE_SIZE_DEFAULT, SUPPORTED_FORMATS, MIME_TYPES

# Configure logging
logger = logging.getLogger("file_processor")

# Add missing mimetypes
for ext, mime_type in MIME_TYPES.items():
    mimetypes.add_type(mime_type, ext)

# Add WordPerfect Document mimetype
mimetypes.add_type('application/wordperfect', '.wpd')


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
        
        # WPD files are handled separately by converting to PDF first
        if ext == '.wpd':
            return True
        
        if ext not in textract_supported:
            logger.warning(f"File extension {ext} not supported by textract: {file_path}")
            return False
        
        if ext in SUPPORTED_FORMATS:
            return True
            
        # Try to determine MIME type
        mime_type, _ = mimetypes.guess_type(file_path)
        if mime_type and any(mime_type.startswith(t) for t in [
            'text/', 'application/pdf', 'application/msword', 
            'application/vnd.ms-', 'application/vnd.openxmlformats-',
            'application/wordperfect'
        ]):
            return True
            
        return False
    
    def _convert_wpd_to_pdf(self, wpd_path: str) -> Optional[str]:
        """
        Convert a WPD file to PDF using LibreOffice or another converter.
        
        Args:
            wpd_path: Path to the WPD file
            
        Returns:
            Path to the converted PDF file or None if conversion failed
        """
        try:
            # Create a temporary directory for the output
            with tempfile.TemporaryDirectory() as temp_dir:
                pdf_path = os.path.join(temp_dir, "converted.pdf")
                
                # Use LibreOffice to convert WPD to PDF
                # Adjust the command based on your system's configuration
                cmd = ["libreoffice", "--headless", "--convert-to", "pdf", 
                       "--outdir", temp_dir, wpd_path]
                
                process = subprocess.run(cmd, capture_output=True, text=True)
                
                if process.returncode != 0:
                    logger.error(f"Failed to convert WPD to PDF: {process.stderr}")
                    return None
                
                # Create a temporary file that will persist after this function returns
                temp_pdf = tempfile.NamedTemporaryFile(suffix='.pdf', delete=False)
                temp_pdf.close()
                
                # Copy the converted PDF to our persistent temp file
                with open(pdf_path, 'rb') as src, open(temp_pdf.name, 'wb') as dst:
                    dst.write(src.read())
                
                return temp_pdf.name
        except Exception as e:
            logger.error(f"Error converting WPD to PDF: {e}")
            return None
    
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
        
        # Special handling for WPD files - convert to PDF first
        temp_pdf_path = None
        if ext == '.wpd':
            logger.info(f"Converting WPD file to PDF: {file_path}")
            temp_pdf_path = self._convert_wpd_to_pdf(file_path)
            if temp_pdf_path:
                file_path = temp_pdf_path
            else:
                logger.error(f"Failed to convert WPD file to PDF: {file_path}")
                return ""
        
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
        finally:
            # Clean up the temporary PDF file if one was created
            if temp_pdf_path and os.path.exists(temp_pdf_path):
                try:
                    os.unlink(temp_pdf_path)
                except Exception as e:
                    logger.warning(f"Failed to delete temporary PDF file {temp_pdf_path}: {e}")
            
        return text
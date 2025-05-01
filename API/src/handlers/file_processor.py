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
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')

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
            '.tif', '.tiff', '.tsv', '.txt', '.wav', 'wdp', '.xls', '.xlsx'
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
                # Get the base filename without extension
                base_filename = os.path.basename(wpd_path)
                base_name_no_ext = os.path.splitext(base_filename)[0]

                # LibreOffice may create the PDF with this name pattern
                expected_pdf_name = f"{base_name_no_ext}.pdf"
                expected_pdf_path = os.path.join(temp_dir, expected_pdf_name)

                logger.info(f"Using temp directory: {temp_dir}")
                logger.info(f"Base filename: {base_filename}")
                logger.info(f"Expected PDF path: {expected_pdf_path}")

                # Search for LibreOffice executable in various locations
                libreoffice_paths = [
                    "libreoffice",  # Regular PATH
                    "/usr/bin/libreoffice",
                    "/usr/lib/libreoffice/program/soffice",
                    "/opt/libreoffice/program/soffice",
                    "/usr/lib/libreoffice/program/libreoffice",
                ]

                libreoffice_exec = None
                for path in libreoffice_paths:
                    try:
                        subprocess.run([path, "--version"], capture_output=True, text=True)
                        libreoffice_exec = path
                        logger.info(f"Found LibreOffice at: {path}")
                        break
                    except (FileNotFoundError, subprocess.SubprocessError):
                        continue
                    
                if not libreoffice_exec:
                    logger.error("LibreOffice executable not found in any of the expected locations")
                    return None

                # Use LibreOffice to convert WPD to PDF
                cmd = [libreoffice_exec, "--headless", "--convert-to", "pdf", 
                       "--outdir", temp_dir, wpd_path]

                logger.info(f"Running conversion command: {' '.join(cmd)}")

                process = subprocess.run(cmd, capture_output=True, text=True)

                # Log the output regardless of success
                logger.info(f"Command stdout: {process.stdout}")
                logger.info(f"Command stderr: {process.stderr}")

                if process.returncode != 0 and not process.stdout.strip():
                    logger.error(f"Failed to convert WPD to PDF: {process.stderr}")
                    return None

                # List files in the output directory to find the converted file
                logger.info(f"Checking output directory contents: {temp_dir}")
                files_in_dir = os.listdir(temp_dir)
                logger.info(f"Files in output directory: {files_in_dir}")

                pdf_files = [f for f in files_in_dir if f.endswith('.pdf')]

                if not pdf_files:
                    logger.error(f"No PDF files found in output directory")
                    return None

                # If we found PDF files, use the first one
                actual_pdf_path = os.path.join(temp_dir, pdf_files[0])
                logger.info(f"Found PDF file: {actual_pdf_path}")

                # Create a temporary file that will persist after this function returns
                temp_pdf = tempfile.NamedTemporaryFile(suffix='.pdf', delete=False)
                temp_pdf.close()

                # Copy the converted PDF to our persistent temp file
                logger.info(f"Copying PDF to persistent temp file: {temp_pdf.name}")
                with open(actual_pdf_path, 'rb') as src, open(temp_pdf.name, 'wb') as dst:
                    dst.write(src.read())

                logger.info(f"Successfully created persistent PDF at: {temp_pdf.name}")
                return temp_pdf.name
        except Exception as e:
            logger.error(f"Error converting WPD to PDF: {str(e)}", exc_info=True)
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


if __name__ == "__main__":
    # Define mock constants if they're not available (for standalone testing)
    try:
        from src.utils.constants import MAX_FILE_SIZE_DEFAULT, SUPPORTED_FORMATS, MIME_TYPES
    except ImportError:
        MAX_FILE_SIZE_DEFAULT = 50 * 1024 * 1024  # 50MB
        SUPPORTED_FORMATS = ['.txt', '.pdf', '.doc', '.docx', '.wpd']
        MIME_TYPES = {
            '.wpd': 'application/wordperfect',
            '.pdf': 'application/pdf',
            '.docx': 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
            '.doc': 'application/msword',
            '.txt': 'text/plain'
        }
    
    # Test if LibreOffice is installed
    print("Checking for LibreOffice installation...")
    try:
        libreoffice_check = subprocess.run(["libreoffice", "--version"], 
                                        capture_output=True, text=True)
        if libreoffice_check.returncode == 0:
            print(f"LibreOffice is available: {libreoffice_check.stdout.strip()}")
        else:
            print("Warning: LibreOffice check returned non-zero exit code.")
            print(f"Error: {libreoffice_check.stderr}")
    except FileNotFoundError:
        print("ERROR: LibreOffice not found! Please install LibreOffice.")
        print("On Ubuntu/Debian: sudo apt-get install libreoffice")
        print("On macOS: brew install --cask libreoffice")
        print("On Windows: Download from libreoffice.org")
        exit(1)
    except Exception as e:
        print(f"Error checking LibreOffice: {e}")
    
    # Test with the specific file
    test_file_path = "cpra.sheriff.1.wpd"
    
    # Ensure file exists
    if not os.path.exists(test_file_path):
        print(f"Error: Test file '{test_file_path}' not found.")
        exit(1)
    
    print(f"Testing WordPerfect conversion with file: {test_file_path}")
    
    # Create processor instance
    processor = FileProcessor()
    
    # Test WPD to PDF conversion
    pdf_path = processor._convert_wpd_to_pdf(test_file_path)
    if pdf_path:
        print(f"Successfully converted WPD to PDF: {pdf_path}")
        
        # Try to extract text from the PDF directly as a test
        try:
            pdf_text = textract.process(pdf_path, method='pdfminer').decode('utf-8')
            print(f"Successfully extracted text from PDF. Length: {len(pdf_text)}")
            preview = pdf_text[:200] + ("..." if len(pdf_text) > 200 else "")
            print(f"PDF text preview: {preview}")
        except Exception as e:
            print(f"Error extracting text from PDF: {e}")
        
        # Clean up the temp PDF file
        try:
            os.unlink(pdf_path)
            print(f"Cleaned up temporary PDF file.")
        except Exception as e:
            print(f"Note: Could not clean up temporary PDF file: {e}")
    else:
        print("Failed to convert WPD to PDF.")
    
    # Test text extraction
    print("\nTesting full text extraction process...")
    text = processor.extract_text(test_file_path)
    
    print(f"\nExtracted text length: {len(text)} characters")
    if text:
        # Print a preview of the extracted text (first 500 chars)
        preview = text[:500] + ("..." if len(text) > 500 else "")
        print(f"\nText preview:\n{preview}")
        
        # Save the extracted text to a file
        output_text_file = test_file_path + ".txt"
        try:
            with open(output_text_file, "w", encoding="utf-8") as f:
                f.write(text)
            print(f"\nExtracted text saved to: {output_text_file}")
        except Exception as e:
            print(f"Error saving text to file: {e}")
    else:
        print("No text was extracted from the file.")
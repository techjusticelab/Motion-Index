import os
import subprocess
from pathlib import Path
import time
import shutil
import zipfile
import gzip
import re

# ANSI color codes for colorful terminal output
class Colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    CYAN = '\033[96m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def colored_print(text, color):
    """Print text with color"""
    print(f"{color}{text}{Colors.ENDC}")

def check_command_exists(command):
    """Check if a command exists on the system"""
    try:
        subprocess.run(["which", command], check=True, capture_output=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False

def get_file_type(file_path):
    """Determine file type using the 'file' command"""
    try:
        result = subprocess.run(["file", "--brief", str(file_path)], capture_output=True, text=True)
        return result.stdout.strip()
    except Exception:
        return "unknown"

def is_binary_file(file_path):
    """Check if a file is binary (not text)"""
    try:
        with open(file_path, 'rb') as f:
            chunk = f.read(1024)
            return b'\0' in chunk  # Binary files often contain null bytes
    except Exception:
        return True  # If can't read, assume binary

def is_macosx_metadata(file_path):
    """Check if the file is a macOS metadata file"""
    path_parts = str(file_path).split(os.sep)
    # Check for __MACOSX directory or files starting with ._
    return "__MACOSX" in path_parts or os.path.basename(file_path).startswith("._")

def extract_zip(zip_path, extract_dir):
    """Extract a ZIP file"""
    try:
        # First check if it's really a ZIP file
        if not zipfile.is_zipfile(zip_path):
            colored_print(f"Not a valid ZIP file: {zip_path}", Colors.YELLOW)
            return False
            
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            # Filter out __MACOSX and ._ files
            files_to_extract = [f for f in zip_ref.namelist() 
                                if not f.startswith('__MACOSX/') and not os.path.basename(f).startswith('._')]
            
            # If no valid files to extract, report it
            if not files_to_extract:
                colored_print(f"ZIP file contains only macOS metadata, skipping: {zip_path}", Colors.YELLOW)
                return False
                
            # Extract only valid files
            for file in files_to_extract:
                try:
                    zip_ref.extract(file, extract_dir)
                except Exception as e:
                    colored_print(f"Error extracting {file} from ZIP: {str(e)}", Colors.RED)
                    
        colored_print(f"Extracted ZIP file: {zip_path}", Colors.GREEN)
        return True
    except Exception as e:
        colored_print(f"Error extracting ZIP file {zip_path}: {str(e)}", Colors.RED)
        return False

def extract_gzip(gz_path, extract_dir):
    """Extract a GZIP file"""
    try:
        # Check if it's really a gzip file
        with open(gz_path, 'rb') as f:
            magic = f.read(2)
            if magic != b'\x1f\x8b':  # gzip magic number
                colored_print(f"Not a valid GZIP file: {gz_path}", Colors.YELLOW)
                return False
        
        output_path = os.path.join(extract_dir, os.path.basename(gz_path)[:-3])
        with gzip.open(gz_path, 'rb') as f_in:
            with open(output_path, 'wb') as f_out:
                shutil.copyfileobj(f_in, f_out)
        colored_print(f"Extracted GZIP file: {gz_path}", Colors.GREEN)
        return True
    except Exception as e:
        colored_print(f"Error extracting GZIP file {gz_path}: {str(e)}", Colors.RED)
        return False

def convert_to_txt(input_file, output_txt):
    """Convert unknown file to text if possible"""
    try:
        # Check if it's a binary file
        if is_binary_file(input_file):
            # Try to use 'strings' command to extract text
            with open(output_txt, 'w') as f:
                subprocess.run(["strings", str(input_file)], stdout=f, check=True)
        else:
            # If it's already text-like, just copy it
            shutil.copy(input_file, output_txt)
        
        # Check if the output is empty or too small
        if os.path.getsize(output_txt) < 10:  # Very small files might be useless
            return False
            
        return True
    except Exception as e:
        colored_print(f"Error converting to text: {str(e)}", Colors.RED)
        return False

def convert_text_to_pdf(text_file, pdf_file):
    """Convert text file to PDF using enscript and ps2pdf"""
    try:
        ps_file = text_file.with_suffix(".ps")
        subprocess.run(["enscript", "-p", str(ps_file), str(text_file)], check=True, capture_output=True)
        subprocess.run(["ps2pdf", str(ps_file), str(pdf_file)], check=True, capture_output=True)
        if os.path.exists(ps_file):
            os.remove(ps_file)  # Clean up the PS file
        return True
    except Exception as e:
        colored_print(f"Error converting text to PDF: {str(e)}", Colors.RED)
        return False

def convert_html_to_pdf_alternative(html_file, pdf_file):
    """Convert HTML to PDF using LibreOffice as an alternative to wkhtmltopdf"""
    try:
        parent_dir = os.path.dirname(html_file)
        cmd = [
            "libreoffice",
            "--headless",
            "--convert-to", "pdf",
            "--outdir", parent_dir,
            str(html_file)
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        
        # LibreOffice outputs to a file with .pdf extension
        expected_pdf = os.path.splitext(html_file)[0] + ".pdf"
        if os.path.exists(expected_pdf) and expected_pdf != pdf_file:
            shutil.move(expected_pdf, pdf_file)
            
        return result.returncode == 0 and os.path.exists(pdf_file)
    except Exception as e:
        colored_print(f"Error converting HTML to PDF: {str(e)}", Colors.RED)
        return False

def convert_to_pdf_via_libreoffice(input_file, output_dir):
    """Try to convert file to PDF using LibreOffice"""
    try:
        cmd = [
            "libreoffice",
            "--headless",
            "--convert-to", "pdf",
            "--outdir", output_dir,
            str(input_file)
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        output_pdf = os.path.join(output_dir, os.path.basename(input_file).rsplit('.', 1)[0] + ".pdf")
        
        return result.returncode == 0 and os.path.exists(output_pdf)
    except Exception:
        return False

def is_junk_file(file_path):
    """Check if file is a system junk file that should be deleted"""
    junk_patterns = [
        r'\.DS_Store$',
        r'Thumbs\.db$',
        r'desktop\.ini$',
        r'\.\.?$'  # . or ..
    ]
    
    filename = os.path.basename(file_path).lower()
    
    # Check against patterns
    for pattern in junk_patterns:
        if re.match(pattern, filename, re.IGNORECASE):
            return True
    
    # Check specific names
    junk_names = [
        'ds_store',
        '.ds_store',
        'thumbs.db',
        'desktop.ini',
        '.directory'
    ]
    
    if filename in junk_names:
        return True
    
    # Check if it's a macOS metadata file
    if is_macosx_metadata(file_path):
        return True
    
    return False

def safe_makedirs(directory):
    """Safely create a directory, handling spaces and special characters"""
    try:
        os.makedirs(directory, exist_ok=True)
        return True
    except Exception as e:
        colored_print(f"Error creating directory {directory}: {str(e)}", Colors.RED)
        return False

def process_all_files():
    """
    Process all files in the data directory:
    1. Extract archives (zip, gz)
    2. Convert files to PDF using appropriate methods
    3. Delete junk files
    4. Convert unknown files to text and then to PDF
    """
    root_directory = "./data"
    
    # Ensure the root directory exists
    if not os.path.isdir(root_directory):
        colored_print(f"Error: Directory '{root_directory}' does not exist.", Colors.RED)
        return
    
    # List of file formats that textract already supports (to keep as is)
    textract_supported = [
        '.csv','.eml', '.epub', '.gif', '.htm', '.html',
        '.jpeg', '.jpg', '.json', '.log', '.mp3', '.msg', '.odt', '.ogg',
        '.pdf', '.png', '.pptx', '.ps', '.psv', '.rtf', '.tab', '.tff',
        '.tif', '.tiff', '.tsv', '.txt', '.wav', '.xls', '.xlsx'
    ]
    
    # File types that should be processed/handled specially
    archive_extensions = ['.zip', '.gz', '.tar', '.rar', '.7z', '.bz2', '.xz']
    document_extensions = ['.doc', '.docx', '.wpd', '.wp', '.wp5', '.wps', '.wpt', '.wri', '.pages']
    presentation_extensions = ['.ppt', '.pptx', '.odp']
    spreadsheet_extensions = ['.xls', '.xlsx', '.ods']
    
    colored_print("\n🔍 SCANNING AND PROCESSING FILES...", Colors.HEADER + Colors.BOLD)
    
    # First pass: Delete all macOS metadata and junk files
    colored_print("\n🧹 PHASE 1: CLEANING UP JUNK FILES...", Colors.HEADER)
    
    junk_files_deleted = 0
    
    for root, dirs, files in os.walk(root_directory):
        # Skip processing the __MACOSX directories
        dirs[:] = [d for d in dirs if d != "__MACOSX"]
        
        for file in files:
            file_path = os.path.join(root, file)
            if is_junk_file(file_path):
                try:
                    os.remove(file_path)
                    junk_files_deleted += 1
                    colored_print(f"Deleted junk file: {file_path}", Colors.YELLOW)
                except Exception as e:
                    colored_print(f"Error deleting junk file {file_path}: {str(e)}", Colors.RED)
    
    colored_print(f"Junk files deleted: {junk_files_deleted}", Colors.CYAN)
    
    # Second pass: Extract all archives
    colored_print("\n📦 PHASE 2: EXTRACTING ARCHIVES...", Colors.HEADER)
    
    archives_extracted = 0
    archives_failed = 0
    
    for root, dirs, files in os.walk(root_directory):
        # Skip processing the __MACOSX directories
        dirs[:] = [d for d in dirs if d != "__MACOSX"]
        
        for file in files:
            file_path = os.path.join(root, file)
            file_ext = os.path.splitext(file)[1].lower()
            
            # Skip if it's a macOS metadata file
            if is_macosx_metadata(file_path):
                continue
                
            if file_ext in archive_extensions:
                colored_print(f"Found archive: {file_path}", Colors.CYAN)
                
                # Create extraction directory (named after the archive without extension)
                extract_dir = os.path.join(root, os.path.splitext(file)[0])
                
                # Skip if the directory can't be created
                if not safe_makedirs(extract_dir):
                    archives_failed += 1
                    continue
                
                # Extract based on file type
                success = False
                if file_ext == '.zip':
                    success = extract_zip(file_path, extract_dir)
                elif file_ext == '.gz':
                    success = extract_gzip(file_path, extract_dir)
                # Add additional archive formats here
                
                if success:
                    archives_extracted += 1
                    # Delete the archive after extraction
                    try:
                        os.remove(file_path)
                        colored_print(f"Deleted archive after extraction: {file_path}", Colors.BLUE)
                    except Exception as e:
                        colored_print(f"Error deleting archive {file_path}: {str(e)}", Colors.RED)
                else:
                    archives_failed += 1
    
    colored_print(f"Archives extracted: {archives_extracted}, Failed: {archives_failed}", Colors.CYAN)
    
    # Third pass: Process all files
    colored_print("\n🔄 PHASE 3: CONVERTING FILES...", Colors.HEADER)
    
    files_kept = 0
    files_converted = 0
    files_deleted = 0
    files_failed = 0
    unknown_converted = 0
    
    # Get all files after extraction
    all_files = []
    for root, dirs, files in os.walk(root_directory):
        # Skip processing the __MACOSX directories
        dirs[:] = [d for d in dirs if d != "__MACOSX"]
        
        for file in files:
            file_path = os.path.join(root, file)
            # Skip macOS metadata files
            if not is_macosx_metadata(file_path):
                all_files.append(file_path)
    
    total_files = len(all_files)
    colored_print(f"Total files to process: {total_files}", Colors.CYAN)
    
    for i, file_path in enumerate(all_files):
        try:
            progress = f"[{i+1}/{total_files}]"
            filename = os.path.basename(file_path)
            file_ext = os.path.splitext(filename)[1].lower()
            parent_dir = os.path.dirname(file_path)
            
            # Skip files that already have a PDF version
            pdf_version = os.path.join(parent_dir, os.path.splitext(filename)[0] + ".pdf")
            if os.path.exists(pdf_version):
                colored_print(f"{progress} Skipping {filename} - PDF version exists", Colors.BLUE)
                files_kept += 1
                continue
                
            # Check if it's a junk file (again to catch any new files from extraction)
            if is_junk_file(file_path):
                colored_print(f"{progress} Deleting junk file: {filename}", Colors.YELLOW)
                os.remove(file_path)
                files_deleted += 1
                continue
                
            # Handle files based on type
            if file_ext.lower() in textract_supported:
                # Keep textract-supported files
                colored_print(f"{progress} Keeping textract-supported file: {filename}", Colors.BLUE)
                files_kept += 1
                
            elif file_ext.lower() in document_extensions + presentation_extensions + spreadsheet_extensions:
                # Try to convert to PDF using LibreOffice
                colored_print(f"{progress} Converting to PDF: {filename}", Colors.CYAN)
                if convert_to_pdf_via_libreoffice(file_path, parent_dir):
                    files_converted += 1
                    colored_print(f"{progress} ✅ Successfully converted: {filename}", Colors.GREEN)
                    os.remove(file_path)
                    colored_print(f"{progress} 🗑️ Deleted original: {filename}", Colors.BLUE)
                else:
                    # If LibreOffice fails, try alternative methods
                    colored_print(f"{progress} LibreOffice conversion failed, trying alternatives...", Colors.YELLOW)
                    
                    # For wpd files, try specialized converters
                    if file_ext.lower() in ['.wpd', '.wp', '.wp5']:
                        # Try wpd2text
                        if check_command_exists("wpd2text"):
                            temp_txt = os.path.join(parent_dir, os.path.splitext(filename)[0] + ".txt")
                            try:
                                with open(temp_txt, 'w') as f:
                                    subprocess.run(["wpd2text", file_path], stdout=f, check=True)
                                
                                # Convert text to PDF
                                if convert_text_to_pdf(Path(temp_txt), Path(pdf_version)):
                                    files_converted += 1
                                    colored_print(f"{progress} ✅ Successfully converted via wpd2text: {filename}", Colors.GREEN)
                                    os.remove(file_path)
                                    if os.path.exists(temp_txt):
                                        os.remove(temp_txt)
                                    continue
                            except Exception:
                                if os.path.exists(temp_txt):
                                    os.remove(temp_txt)
                    
                    # For HTML files, use LibreOffice as alternative to wkhtmltopdf
                    if file_ext.lower() in ['.htm', '.html']:
                        if convert_html_to_pdf_alternative(file_path, pdf_version):
                            files_converted += 1
                            colored_print(f"{progress} ✅ Successfully converted HTML: {filename}", Colors.GREEN)
                            os.remove(file_path)
                            continue
                    
                    # Fallback for all files: convert to text then PDF
                    temp_txt = os.path.join(parent_dir, os.path.splitext(filename)[0] + ".txt")
                    if convert_to_txt(file_path, temp_txt):
                        if convert_text_to_pdf(Path(temp_txt), Path(pdf_version)):
                            files_converted += 1
                            colored_print(f"{progress} ✅ Successfully converted via text: {filename}", Colors.GREEN)
                            os.remove(file_path)
                            if os.path.exists(temp_txt):
                                os.remove(temp_txt)
                            continue
                        if os.path.exists(temp_txt):
                            os.remove(temp_txt)
                    
                    files_failed += 1
                    colored_print(f"{progress} ❌ Failed to convert: {filename}", Colors.RED)
                    
            elif file_ext == '':
                # No extension - try to determine file type
                colored_print(f"{progress} Processing file with no extension: {filename}", Colors.YELLOW)
                file_type = get_file_type(file_path)
                
                # Try to convert based on file type
                if "text" in file_type.lower():
                    # It's text, convert to PDF
                    if convert_text_to_pdf(Path(file_path), Path(pdf_version)):
                        files_converted += 1
                        colored_print(f"{progress} ✅ Successfully converted text file: {filename}", Colors.GREEN)
                        os.remove(file_path)
                        continue
                elif "word" in file_type.lower() or "document" in file_type.lower():
                    # Try LibreOffice
                    if convert_to_pdf_via_libreoffice(file_path, parent_dir):
                        files_converted += 1
                        colored_print(f"{progress} ✅ Successfully converted document: {filename}", Colors.GREEN)
                        os.remove(file_path)
                        continue
                
                # Fallback: convert to text
                colored_print(f"{progress} Converting unknown file to text: {filename}", Colors.YELLOW)
                temp_txt = os.path.join(parent_dir, filename + ".txt")
                if convert_to_txt(file_path, temp_txt):
                    if convert_text_to_pdf(Path(temp_txt), Path(pdf_version)):
                        unknown_converted += 1
                        colored_print(f"{progress} ✅ Successfully converted unknown file: {filename}", Colors.GREEN)
                        os.remove(file_path)
                        if os.path.exists(temp_txt):
                            os.remove(temp_txt)
                        continue
                    if os.path.exists(temp_txt):
                        os.remove(temp_txt)
                
                files_failed += 1
                colored_print(f"{progress} ❌ Could not process file: {filename}", Colors.RED)
                
            else:
                # Unknown extension - try to convert to text
                colored_print(f"{progress} Converting unknown format to text: {filename}", Colors.YELLOW)
                temp_txt = os.path.join(parent_dir, os.path.splitext(filename)[0] + ".txt")
                if convert_to_txt(file_path, temp_txt):
                    if convert_text_to_pdf(Path(temp_txt), Path(pdf_version)):
                        unknown_converted += 1
                        colored_print(f"{progress} ✅ Successfully converted to PDF via text: {filename}", Colors.GREEN)
                        os.remove(file_path)
                        if os.path.exists(temp_txt):
                            os.remove(temp_txt)
                        continue
                    if os.path.exists(temp_txt):
                        os.remove(temp_txt)
                
                files_failed += 1
                colored_print(f"{progress} ❌ Failed to convert unknown format: {filename}", Colors.RED)
                
        except Exception as e:
            files_failed += 1
            colored_print(f"Error processing {file_path}: {str(e)}", Colors.RED)
    
    # Fourth pass: Clean up empty directories
    colored_print("\n🧹 PHASE 4: CLEANING UP EMPTY DIRECTORIES...", Colors.HEADER)
    
    empty_dirs_removed = 0
    for root, dirs, files in os.walk(root_directory, topdown=False):
        for dir_name in dirs:
            dir_path = os.path.join(root, dir_name)
            try:
                # Check if the directory is empty
                if not os.listdir(dir_path):
                    os.rmdir(dir_path)
                    empty_dirs_removed += 1
                    colored_print(f"Removed empty directory: {dir_path}", Colors.YELLOW)
            except Exception as e:
                colored_print(f"Error removing empty directory {dir_path}: {str(e)}", Colors.RED)
    
    # Print summary
    colored_print("\n📊 PROCESSING SUMMARY", Colors.HEADER + Colors.BOLD)
    colored_print(f"Total files processed: {total_files}", Colors.CYAN)
    colored_print(f"✅ Files kept (textract supported): {files_kept}", Colors.BLUE)
    colored_print(f"✅ Files converted to PDF: {files_converted}", Colors.GREEN)
    colored_print(f"✅ Unknown files converted: {unknown_converted}", Colors.GREEN)
    colored_print(f"🗑️ Junk files deleted: {junk_files_deleted + files_deleted}", Colors.YELLOW)
    colored_print(f"❌ Files that couldn't be processed: {files_failed}", Colors.RED)
    colored_print(f"🧹 Empty directories removed: {empty_dirs_removed}", Colors.YELLOW)
    
    if archives_extracted > 0:
        colored_print(f"📦 Archives extracted: {archives_extracted}", Colors.CYAN)

if __name__ == "__main__":
    colored_print("\n📁 COMPREHENSIVE FILE PROCESSING UTILITY", Colors.HEADER + Colors.BOLD)
    colored_print("Converting all files to formats compatible with textract...\n", Colors.CYAN)
    process_all_files()
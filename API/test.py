#!/usr/bin/env python3
import os
import subprocess
import sys

def convert_wpd_to_pdf(input_file, output_dir=None):
    """
    Convert a WordPerfect document to PDF using LibreOffice.
    
    Args:
        input_file (str): Path to the WPD file
        output_dir (str, optional): Directory to save the PDF. Defaults to same directory as input.
    
    Returns:
        str: Path to the output PDF file or None if conversion failed
    """
    # Validate input file
    if not os.path.isfile(input_file):
        print(f"Error: Input file '{input_file}' does not exist.")
        return None
    
    # Get absolute path
    input_file = os.path.abspath(input_file)
    
    # Set output directory
    if output_dir is None:
        output_dir = os.path.dirname(input_file)
    else:
        if not os.path.isdir(output_dir):
            os.makedirs(output_dir)
    
    # Get base filename without extension
    base_name = os.path.splitext(os.path.basename(input_file))[0]
    
    # Define output PDF path
    output_file = os.path.join(output_dir, f"{base_name}.pdf")
    
    # Command to execute LibreOffice in headless mode
    # Different commands based on operating system
    if sys.platform.startswith('win'):
        # Windows
        libreoffice_paths = [
            r'C:\Program Files\LibreOffice\program\soffice.exe',
            r'C:\Program Files (x86)\LibreOffice\program\soffice.exe',
        ]
        
        libreoffice_exec = None
        for path in libreoffice_paths:
            if os.path.isfile(path):
                libreoffice_exec = f'"{path}"'
                break
                
        if not libreoffice_exec:
            print("Error: LibreOffice not found. Please install LibreOffice.")
            return None
    else:
        # Linux/Mac
        libreoffice_exec = 'libreoffice'
    
    # Build the conversion command
    cmd = [
        libreoffice_exec, 
        '--headless', 
        '--convert-to', 
        'pdf', 
        '--outdir', 
        output_dir,
        input_file
    ]
    
    try:
        # Execute the command
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout, stderr = process.communicate()
        
        if process.returncode != 0:
            print(f"Error during conversion: {stderr.decode('utf-8')}")
            return None
        
        # Check if the file was created
        if os.path.isfile(output_file):
            print(f"Successfully converted to: {output_file}")
            return output_file
        else:
            print("Conversion seemed to succeed, but output file not found.")
            return None
            
    except Exception as e:
        print(f"Exception during conversion: {str(e)}")
        return None

if __name__ == "__main__":
    # The specific file mentioned in the prompt
    wpd_file = "./cpra.sheriff.1.wpd"
    
    # Verify file exists
    if not os.path.exists(wpd_file):
        print(f"Error: The file '{wpd_file}' doesn't exist.")
        print("Please check the file path and try again.")
        sys.exit(1)
    
    # Convert the file
    convert_wpd_to_pdf(wpd_file)
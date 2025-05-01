import os
import subprocess
from pathlib import Path

def convert_all_to_pdf():
    """
    Converts all compatible files in ./data directory to PDF format using LibreOffice,
    then deletes the original files after successful conversion.
    """
    directory = "./data"
    
    # Ensure the directory exists
    if not os.path.isdir(directory):
        print(f"Error: Directory '{directory}' does not exist.")
        return
    
    # Find all files in the directory (excluding existing PDFs)
    input_files = [f for f in Path(directory).glob("*") if f.is_file() and f.suffix.lower() != '.pdf']
    
    if not input_files:
        print(f"No convertible files found in '{directory}'.")
        return
    
    print(f"Found {len(input_files)} file(s) to convert.")
    
    # Try to use LibreOffice for conversion
    try:
        # Attempt to find LibreOffice
        libreoffice_path = "libreoffice"  # Use default path for Linux
        
        # Convert each file to PDF
        success_count = 0
        deleted_count = 0
        
        for input_file in input_files:
            output_pdf = input_file.with_suffix(".pdf")
            print(f"Converting {input_file.name} to PDF...")
            
            try:
                # Use LibreOffice to convert the file
                cmd = [
                    libreoffice_path,
                    "--headless",
                    "--convert-to", "pdf",
                    "--outdir", directory,
                    str(input_file)
                ]
                
                result = subprocess.run(cmd, capture_output=True, text=True)
                
                if result.returncode == 0 and os.path.exists(output_pdf):
                    success_count += 1
                    print(f"Converted: {input_file.name} → {output_pdf.name}")
                    
                    # Delete the original file
                    os.remove(input_file)
                    deleted_count += 1
                else:
                    print(f"Failed to convert {input_file.name}")
                    
            except Exception as e:
                print(f"Error processing {input_file.name}: {str(e)}")
        
        print(f"Conversion complete. Successfully converted {success_count} of {len(input_files)} files.")
        print(f"Deleted {deleted_count} original files.")
        
    except Exception as e:
        print(f"Error: {str(e)}")
        print("Make sure LibreOffice is installed on your system.")

if __name__ == "__main__":
    convert_all_to_pdf()
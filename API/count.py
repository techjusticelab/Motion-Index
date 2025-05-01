import os
from collections import Counter
import sys
from pathlib import Path

def count_file_extensions(root_dir):
    if not os.path.exists(root_dir):
        print(f"Error: Directory '{root_dir}' does not exist.")
        return
    
    extension_counter = Counter()
    
    # Walk through all subdirectories
    for dirpath, dirnames, filenames in os.walk(root_dir):
        for filename in filenames:
            # Get the file extension (or empty string if none)
            _, extension = os.path.splitext(filename)
            # Convert to lowercase for consistency and remove the dot
            extension = extension.lower()[1:] if extension else "no_extension"
            extension_counter[extension] += 1
    
    return extension_counter

def main():
    root_directory = "./data"
    
    print(f"Scanning directory: {Path(root_directory).resolve()}")
    extensions = count_file_extensions(root_directory)
    
    if not extensions:
        print("No files found or directory doesn't exist.")
        return
    
    print("\nFile types found:")
    print("-" * 40)
    print(f"{'Extension':<20} {'Count':<10}")
    print("-" * 40)
    
    # Sort by count (descending)
    for ext, count in sorted(extensions.items(), key=lambda x: x[1], reverse=True):
        print(f"{ext:<20} {count:<10}")
    
    print("-" * 40)
    print(f"Total file types: {len(extensions)}")
    print(f"Total files: {sum(extensions.values())}")

if __name__ == "__main__":
    main()
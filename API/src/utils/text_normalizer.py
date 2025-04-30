"""
Text normalization utilities for standardizing text fields.
"""
import re
from typing import Dict, List, Optional


def normalize_court_name(court_name: str) -> str:
    """
    Normalize court names to a standard format to reduce duplication in dropdowns.
    
    Args:
        court_name: The original court name
        
    Returns:
        Normalized court name
    """
    if not court_name:
        return ""
        
    # Convert to title case first (handles all-caps cases)
    normalized = court_name.strip().title()
    
    # Standardize common variations
    patterns = [
        # Standardize "Superior Court of the State of California, County of X"
        (r"Superior Court Of (The )?State Of California,? County Of (.+)", 
         r"Superior Court of California, County of \2"),
         
        # Standardize "SUPERIOR COURT OF CALIFORNIA, COUNTY OF X"
        (r"Superior Court Of California,? County Of (.+)", 
         r"Superior Court of California, County of \1"),
         
        # Standardize "Supreme Court of the State of California"
        (r"Supreme Court Of (The )?State Of California", 
         r"Supreme Court of California"),
         
        # Standardize "Court of Appeal of the State of California"
        (r"Court Of Appeal Of (The )?State Of California,? (.+)", 
         r"Court of Appeal of California, \2"),
    ]
    
    # Apply all patterns
    for pattern, replacement in patterns:
        normalized = re.sub(pattern, replacement, normalized, flags=re.IGNORECASE)
    
    # Fix specific court names with Justice Centers
    if "Justice Center" in normalized:
        # Extract county name and justice center
        match = re.search(r"County Of ([^,]+),?\s+(.+?)\s*Justice Center", normalized, re.IGNORECASE)
        if match:
            county, center = match.groups()
            normalized = f"Superior Court of California, County of {county}, {center} Justice Center"
    
    # Fix specific court names with divisions
    if "Division" in normalized:
        # Standardize division format
        normalized = re.sub(r"([^,]+) Division", r"Division \1", normalized)
    
    return normalized


def group_similar_court_names(court_names: List[str]) -> List[str]:
    """
    Group similar court names and return a deduplicated list.
    
    Args:
        court_names: List of court names from Elasticsearch
        
    Returns:
        Deduplicated list of normalized court names
    """
    # Normalize all court names
    normalized_map: Dict[str, List[str]] = {}
    
    for name in court_names:
        normalized = normalize_court_name(name)
        if normalized not in normalized_map:
            normalized_map[normalized] = []
        normalized_map[normalized].append(name)
    
    # Return the normalized names
    return list(normalized_map.keys())

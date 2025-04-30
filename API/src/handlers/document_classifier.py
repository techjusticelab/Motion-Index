"""
Document classifier for legal document classification and metadata extraction.
"""
import re
import json
import logging
from typing import Dict, Any, Optional, Union
import openai

from src.utils.constants import DOCUMENT_TYPES, OPENAI_MODEL, LEGAL_TAGS_LIST

# Configure logging
logger = logging.getLogger("document_classifier")

# List of allowed document types for classification
ALLOWED_TYPES_LIST = list(DOCUMENT_TYPES.values())


def classify_and_extract_by_llm(text: str, client: Any) -> Dict[str, Optional[Any]]:
    """
    Classifies the document AND extracts information using the OpenAI API.
    
    Args:
        text: The full document text (or a significant portion).
        client: An initialized OpenAI client instance.
        
    Returns:
        A dictionary containing the classified document type, legal tags, and other extracted fields,
        or a dictionary with 'document_type' as 'Unknown' if parsing fails.
    """
    logger.info("Attempting LLM classification and extraction...")

    # Define the desired JSON structure
    json_format_instructions = f"""
    {{
        "document_type": "...", // MUST be one of [{', '.join(ALLOWED_TYPES_LIST)}]
        "subject": "...",       // A concise summary (5-10 words) of the document's main topic or purpose
        "status": "...",        // e.g., Granted, Denied, Filed, Served, Proposed (null if not applicable/found)
        "timestamp": "...",     // IMPORTANT: The document's filing date, signature date, or publication date (not just any date mentioned)
                                // Format as YYYY-MM-DD if possible, otherwise use the original format found in the document
                                // Set to null if no filing/signature/publication date can be found
        "case_name": "...",     // e.g., "Plaintiff Corp. v. Defendant Inc." (null if not applicable/found)
        "case_number": "...",   // e.g., "3:24-cv-01234-ABC", "INF-xxxxxxx" (null if not found)
        "author": "...",        // Authoring attorney, law firm, or entity (null if not found)
        "judge": "...",         // Presiding judge's name (null if not applicable/found)
        "court": "...",         // Name of the court, including county/district if specified (null if not applicable/found)
        "legal_tags": [...]     // Array of relevant legal tags from the California Criminal Law Guide 
                                // (e.g., ["Search and Seizure", "Warrantless Search", "Miranda Rights"])
    }}
    """

    # Define a comprehensive list of legal tags based on Witkin California Criminal Law Guide
    

    prompt = f"""
    Analyze the following legal document text. Perform these tasks:
    1. Classify the document into ONE of the following categories: {', '.join(ALLOWED_TYPES_LIST)}
    2. Identify 2-5 MOST RELEVANT legal tags from the comprehensive list provided - only select tags that are directly applicable
    3. Extract the specified fields from the text with high precision.

    Pay close attention to these fields:
    - "subject": Provide a concise summary of the main topic or purpose (5-10 words)
    - "timestamp": CRITICALLY IMPORTANT - Look for the document's filing date, signature date, or publication date
      * Search for dates near terms like "Filed on", "Dated", "Signed this", "Entered", "Filed", "Submitted", etc.
      * Look at the document header, footer, signature blocks, and certificate of service
      * If multiple dates appear, prioritize the filing date over other dates
      * Format as YYYY-MM-DD only
      * Only use null if absolutely no filing/signature/publication date can be found
    - "legal_tags": Identify 2-5 most relevant legal tags from the provided list that apply to this document
      * Look for topic matter, procedural aspects, legal doctrines, and key issues
      * Tags should reflect substantive legal issues addressed and touched on in document
      * If document discusses search and seizure without a warrant, include BOTH "Search and Seizure" and "Warrantless Search" for example
      * Choose the most specific applicable tags (e.g., "Consent Searches" rather than just "Search and Seizure")
      * Include tag on document type ie "dismissal motion" or "motion to dismiss" if neccessary

    Respond ONLY with a single, valid JSON object matching this exact structure:
    {json_format_instructions}

    If a field cannot be found or is not applicable, use the JSON value null.
    Do not include any explanations, apologies, or text outside the JSON object.

    Available legal tags (select only from this comprehensive list):
    {', '.join(LEGAL_TAGS_LIST)}

    Document Text:
    ---
    {text[:8000]} # Reduced limit to accommodate all legal tags while staying within token limits
    ---

    JSON Output:
    """

    default_result = {"document_type": DOCUMENT_TYPES['unknown'], "legal_tags": []}
    try:
        instruction_keys = json.loads(json_format_instructions).keys()
    except json.JSONDecodeError:
        # Fallback keys if instruction string is somehow invalid
        instruction_keys = ["document_type", "subject", "status", "timestamp", "case_name", 
                           "case_number", "author", "judge", "court", "legal_tags"]

    for key in instruction_keys:
        if key != "document_type" and key != "legal_tags":
            default_result[key] = None

    try:
        # For OpenAI v0.28.0
        response = client.ChatCompletion.create(
            model=OPENAI_MODEL,
            messages=[
                {"role": "system", "content": "You are an expert legal document analyzer. You meticulously classify documents, identify relevant legal tags, and extract key information based on the provided text, responding only in the specified valid JSON format."},
                {"role": "user", "content": prompt}
            ],
            temperature=0.1,
            max_tokens=800,  # Increased token limit to accommodate tags
            n=1,
            stop=None,
        )
        llm_response_content = response.choices[0]['message']['content'].strip()
        logger.debug(f"Raw LLM Response:\n{llm_response_content}")
        
        try:
            # Clean up response if needed
            if llm_response_content.startswith("```json"):
                llm_response_content = llm_response_content[7:]
            if llm_response_content.endswith("```"):
                llm_response_content = llm_response_content[:-3]
            llm_response_content = llm_response_content.strip()
            
            # Remove JavaScript-style comments before parsing
            # First, remove single-line comments (// comment)
            cleaned_json = re.sub(r'\s*//.*?$', '', llm_response_content, flags=re.MULTILINE)
            # Then ensure the JSON is valid by removing any trailing commas
            cleaned_json = re.sub(r',\s*([\]\}])', r'\1', cleaned_json)
            
            parsed_json = json.loads(cleaned_json)
            
            # Create final result with all expected fields
            final_result = {}
            for key in default_result:
                if key == "legal_tags" and key not in parsed_json:
                    final_result[key] = []
                else:
                    final_result[key] = parsed_json.get(key, default_result[key])

            # Validate the document_type
            if final_result.get("document_type") not in ALLOWED_TYPES_LIST:
                logger.warning(f"Invalid or missing 'document_type' in response: {final_result.get('document_type')}")
                # Try to extract a valid type from the raw response
                found_type_in_raw = next(
                    (t for t in ALLOWED_TYPES_LIST if re.search(r'\b' + re.escape(t) + r'\b', llm_response_content, re.IGNORECASE)), 
                    DOCUMENT_TYPES['unknown']
                )
                final_result['document_type'] = found_type_in_raw
            
            # Validate legal tags
            if not isinstance(final_result.get("legal_tags"), list):
                if final_result.get("legal_tags") is None:
                    final_result["legal_tags"] = []
                else:
                    # Convert to list if it's a string
                    try:
                        tag_text = str(final_result.get("legal_tags"))
                        if "," in tag_text:
                            final_result["legal_tags"] = [tag.strip() for tag in tag_text.split(",")]
                        else:
                            final_result["legal_tags"] = [tag_text]
                    except:
                        final_result["legal_tags"] = []
            
            # Validate that tags are from the approved list
            if final_result.get("legal_tags"):
                validated_tags = []
                for tag in final_result["legal_tags"]:
                    # Exact match to approved tag
                    if tag in LEGAL_TAGS_LIST:
                        validated_tags.append(tag)
                    else:
                        # Try to find closest matching tag using more sophisticated matching
                        # First, try direct substring match
                        closest_matches = [t for t in LEGAL_TAGS_LIST 
                                          if re.search(r'\b' + re.escape(tag) + r'\b', t, re.IGNORECASE) or 
                                             re.search(r'\b' + re.escape(t) + r'\b', tag, re.IGNORECASE)]
                        
                        if closest_matches:
                            # If multiple matches, prefer the shortest one as it's likely more specific
                            closest_tag = min(closest_matches, key=len)
                            validated_tags.append(closest_tag)
                        else:
                            # If no direct match, try fuzzy matching based on word overlap
                            tag_words = set(tag.lower().split())
                            best_match = None
                            best_overlap = 0
                            
                            for t in LEGAL_TAGS_LIST:
                                t_words = set(t.lower().split())
                                overlap = len(tag_words.intersection(t_words))
                                if overlap > best_overlap:
                                    best_overlap = overlap
                                    best_match = t
                            
                            # Only use the match if there's meaningful word overlap
                            if best_overlap > 0 and best_match:
                                validated_tags.append(best_match)
                
                final_result["legal_tags"] = validated_tags

            # Clean up values
            for key, value in final_result.items():
                if key != "legal_tags" and (value == "null" or value == "N/A" or value == "" or (isinstance(value, str) and not value.strip())):
                    final_result[key] = None

            logger.info(f"LLM classification/extraction successful: Type '{final_result.get('document_type')}', Tags: {final_result.get('legal_tags')}, Timestamp: '{final_result.get('timestamp')}'")
            return final_result

        except json.JSONDecodeError as json_e:
            logger.error(f"Failed to parse JSON response: {json_e}")
            logger.debug(f"Response that failed parsing:\n---\n{llm_response_content}\n---")
            # Try to extract document type from raw response with more context awareness
            # Look for phrases like "document_type": "Motion" or "type is Motion"
            type_patterns = [
                r'"document_type"\s*:\s*"([^"]+)"',  # JSON format
                r'document_type\s*:\s*([\w\s]+)',      # Relaxed format
                r'type\s+is\s+([\w\s]+)',              # Natural language
                r'classified\s+as\s+([\w\s]+)'         # Natural language
            ]
            
            found_type = DOCUMENT_TYPES['unknown']
            for pattern in type_patterns:
                matches = re.findall(pattern, llm_response_content, re.IGNORECASE)
                if matches:
                    # Clean up the match and check if it's in our allowed types
                    for match in matches:
                        clean_match = match.strip().rstrip(',').rstrip('"').strip()
                        if clean_match in ALLOWED_TYPES_LIST:
                            found_type = clean_match
                            break
                    if found_type != DOCUMENT_TYPES['unknown']:
                        break
            
            # If still not found, fall back to simple word matching
            if found_type == DOCUMENT_TYPES['unknown']:
                found_type = next(
                    (t for t in ALLOWED_TYPES_LIST if re.search(r'\b' + re.escape(t) + r'\b', llm_response_content, re.IGNORECASE)), 
                    DOCUMENT_TYPES['unknown']
                )
                
            default_result['document_type'] = found_type
            
            # Try to extract tags from raw response with improved context awareness
            # Look for patterns like "legal_tags": ["Search and Seizure", "Warrantless Search"]
            found_tags = []
            tag_pattern = r'"legal_tags"\s*:\s*\[([^\]]+)\]'
            tag_matches = re.findall(tag_pattern, llm_response_content, re.IGNORECASE)
            
            if tag_matches:
                # Process the content inside the brackets
                for match in tag_matches:
                    # Split by commas and clean up each tag
                    potential_tags = [t.strip().strip('"').strip('\'') for t in match.split(',')]
                    for potential_tag in potential_tags:
                        # Check if this is a valid tag
                        if potential_tag in LEGAL_TAGS_LIST:
                            found_tags.append(potential_tag)
            
            # If no tags found through JSON pattern, fall back to direct matching
            if not found_tags:
                for tag in LEGAL_TAGS_LIST:
                    if re.search(r'\b' + re.escape(tag) + r'\b', llm_response_content, re.IGNORECASE):
                        found_tags.append(tag)
                        if len(found_tags) >= 5:  # Limit to 5 tags
                            break
            
            default_result['legal_tags'] = found_tags
            
            logger.info(f"Fallback: Extracted type '{found_type}' and tags {found_tags} from raw response.")
            return default_result

    except openai.APIError as e:
        logger.error(f"OpenAI API Error: {e}")
        return default_result
    except Exception as e:
        logger.error(f"An unexpected error occurred during LLM processing: {e}")
        return default_result

def process_document_llm(
    document_name: str,
    document_text: str,
    openai_client: Any
) -> Dict[str, Optional[Any]]:
    """
    Classifies and extracts information from a legal document using the LLM API.
    
    Args:
        document_name: The name/identifier of the document (e.g., filename).
        document_text: The full text content of the document.
        openai_client: An initialized OpenAI client instance.
        
    Returns:
        A dictionary containing the classification and extracted fields.
    """
    result_structure = {
        "document_name": document_name,
        "document_type": DOCUMENT_TYPES['unknown'],
        "subject": None,
        "status": None,
        "timestamp": None,
        "case_name": None,
        "case_number": None,
        "legal_tags": None,
        "author": None,
        "judge": None,
        "court": None,
    }
    

    if not document_text or not document_text.strip():
        result_structure["document_type"] = "Unclassified (Empty Input)"
        return result_structure

    llm_extracted_data = classify_and_extract_by_llm(document_text, openai_client)

    final_result = result_structure.copy()  # Start with the base structure
    for key, value in llm_extracted_data.items():
        if key in final_result:  # Only update keys defined in our initial structure
            # Ensure value is not just whitespace before assigning
            if value is not None and (not isinstance(value, str) or value.strip()):
                final_result[key] = value
        else:
            logger.warning(f"LLM returned unexpected key '{key}' - ignoring.")

    # Final check - if type is still Unknown after LLM, log it
    if final_result["document_type"] == DOCUMENT_TYPES['unknown']:
        logger.warning("LLM returned 'Unknown' or classification failed.")
    
    return final_result

"""
Document classifier for legal document classification and metadata extraction.
"""
import re
import json
import logging
from typing import Dict, Any, Optional, Union
import openai

from src.utils.constants import DOCUMENT_TYPES, OPENAI_MODEL

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
        A dictionary containing the classified document type and other extracted fields,
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
        "court": "..."          // Name of the court, including county/district if specified (null if not applicable/found)
    }}
    """

    prompt = f"""
    Analyze the following legal document text. Perform two tasks:
    1. Classify the document into ONE of the following categories: {', '.join(ALLOWED_TYPES_LIST)}
    2. Extract the specified fields from the text.

    Pay close attention to these fields:
    - "subject": Provide a concise summary of the main topic or purpose (5-10 words)
    - "timestamp": CRITICALLY IMPORTANT - Look for the document's filing date, signature date, or publication date
      * Search for dates near terms like "Filed on", "Dated", "Signed this", "Entered", "Filed", "Submitted", etc.
      * Look at the document header, footer, signature blocks, and certificate of service
      * If multiple dates appear, prioritize the filing date over other dates
      * Format as YYYY-MM-DD only
      * Only use null if absolutely no filing/signature/publication date can be found

    Respond ONLY with a single, valid JSON object matching this exact structure:
    {json_format_instructions}

    If a field cannot be found or is not applicable, use the JSON value null.
    Do not include any explanations, apologies, or text outside the JSON object.

    Document Text:
    ---
    {text[:8000]} # Limit text length to avoid token limits
    ---

    JSON Output:
    """

    default_result = {"document_type": DOCUMENT_TYPES['unknown']}
    try:
        instruction_keys = json.loads(json_format_instructions).keys()
    except json.JSONDecodeError:
        # Fallback keys if instruction string is somehow invalid
        instruction_keys = ["document_type", "subject", "status", "timestamp", "case_name", 
                           "case_number", "author", "judge", "court"]

    for key in instruction_keys:
        if key != "document_type":
            default_result[key] = None

    try:
        # For OpenAI v0.28.0
        response = client.ChatCompletion.create(
            model=OPENAI_MODEL,
            messages=[
                {"role": "system", "content": "You are an expert legal document analyzer. You meticulously classify documents and extract key information based on the provided text, responding only in the specified valid JSON format."},
                {"role": "user", "content": prompt}
            ],
            temperature=0.1,
            max_tokens=600,
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

            parsed_json = json.loads(llm_response_content)
            
            # Create final result with all expected fields
            final_result = {}
            for key in default_result:
                final_result[key] = parsed_json.get(key, None)

            # Validate the document_type
            if final_result.get("document_type") not in ALLOWED_TYPES_LIST:
                logger.warning(f"Invalid or missing 'document_type' in response: {final_result.get('document_type')}")
                # Try to extract a valid type from the raw response
                found_type_in_raw = next(
                    (t for t in ALLOWED_TYPES_LIST if re.search(r'\b' + re.escape(t) + r'\b', llm_response_content, re.IGNORECASE)), 
                    DOCUMENT_TYPES['unknown']
                )
                final_result['document_type'] = found_type_in_raw

            # Clean up values
            for key, value in final_result.items():
                if value == "null" or value == "N/A" or value == "" or (isinstance(value, str) and not value.strip()):
                    final_result[key] = None

            logger.info(f"LLM classification/extraction successful: Type '{final_result.get('document_type')}', Timestamp: '{final_result.get('timestamp')}'")
            return final_result

        except json.JSONDecodeError as json_e:
            logger.error(f"Failed to parse JSON response: {json_e}")
            logger.debug(f"Response that failed parsing:\n---\n{llm_response_content}\n---")
            # Try to extract document type from raw response
            found_type = next(
                (t for t in ALLOWED_TYPES_LIST if re.search(r'\b' + re.escape(t) + r'\b', llm_response_content, re.IGNORECASE)), 
                DOCUMENT_TYPES['unknown']
            )
            default_result['document_type'] = found_type
            logger.info(f"Fallback: Extracted type '{found_type}' from raw response.")
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

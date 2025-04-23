"""
Document model for representing processed documents.
"""
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional

from src.models.metadata import Metadata


class Document:
    """
    Represents a processed document with text content, metadata, and storage information.
    """
    def __init__(
        self,
        file_path: str,
        text: str,
        metadata: Metadata,
        doc_type: str = "unknown",
        category: Optional[str] = None,
        chunk_id: int = 0,
        embedding: Optional[List[float]] = None,
        hash_value: Optional[str] = None,
        created_at: Optional[datetime] = None,
        s3_uri: Optional[str] = None
    ):
        """
        Initialize a Document object.
        
        Args:
            file_path: Path to the document file (local or S3)
            text: Extracted text content from the document
            metadata: Metadata object containing document information
            doc_type: Type of document (e.g., "Motion", "Order", "Brief")
            category: Category of the document based on file type or content
            chunk_id: ID for document chunking (for large documents)
            embedding: Vector embedding of the document for semantic search
            hash_value: Hash of the document content for deduplication
            created_at: Timestamp when the document was processed
            s3_uri: S3 URI if the document is stored in S3
        """
        self.file_path = file_path
        self.file_name = Path(file_path).name
        self.text = text
        self.metadata = metadata
        self.doc_type = doc_type
        self.category = category
        self.chunk_id = chunk_id
        self.embedding = embedding
        self.hash = hash_value
        self.created_at = created_at or datetime.now()
        self.s3_uri = s3_uri
    
    def __str__(self) -> str:
        """String representation of the document."""
        return f"Document: {self.file_name} ({self.doc_type})"
    
    def __repr__(self) -> str:
        """Detailed string representation of the document."""
        return f"Document('{self.file_name}', '{self.doc_type}')"
    
    def to_dict(self) -> Dict[str, Any]:
        """
        Convert document to dictionary format for storage/indexing.
        
        Returns:
            Dictionary representation of the document
        """
        doc_dict = {
            "file_path": self.file_path,
            "file_name": self.file_name,
            "text": self.text,
            "doc_type": self.doc_type,
            "chunk_id": self.chunk_id,
            "metadata": self.metadata.to_dict(),
            "hash": self.hash,
            "created_at": self.created_at.isoformat() if self.created_at else None
        }
        
        if self.category:
            doc_dict["category"] = self.category
            
        if self.embedding:
            doc_dict["embedding"] = self.embedding
            
        if self.s3_uri:
            doc_dict["s3_uri"] = self.s3_uri
            
        return doc_dict
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Document":
        """
        Create a Document instance from a dictionary.
        
        Args:
            data: Dictionary containing document fields
            
        Returns:
            New Document instance
        """
        # Handle datetime conversion
        if "created_at" in data and isinstance(data["created_at"], str):
            data["created_at"] = datetime.fromisoformat(data["created_at"])
            
        # Extract and convert metadata
        metadata_dict = data.pop("metadata", {})
        metadata = Metadata.from_dict(metadata_dict)
        
        # Rename hash field if necessary
        if "hash" in data:
            data["hash_value"] = data.pop("hash")
            
        return cls(metadata=metadata, **data)

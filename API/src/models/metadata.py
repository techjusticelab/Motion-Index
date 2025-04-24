"""
Metadata model for document information.
"""
from datetime import datetime
from typing import Dict, Any, Optional


class Metadata:
    """
    Represents metadata extracted from a document.
    Contains information such as document name, subject, status, case information, etc.
    """
    def __init__(
        self,
        document_name: str,
        subject: str,
        status: Optional[str] = None,
        timestamp: Optional[datetime] = None,
        case_name: Optional[str] = None,
        case_number: Optional[str] = None,
        author: Optional[str] = None,
        judge: Optional[str] = None,
        court: Optional[str] = None
    ):
        """
        Initialize a Metadata object.
        
        Args:
            document_name: Name of the document
            subject: Subject or brief description of the document
            status: Status of the document (e.g., "Granted", "Denied", "Filed")
            timestamp: Date/time associated with the document
            case_name: Name of the legal case (e.g., "Plaintiff v. Defendant")
            case_number: Case identification number
            author: Author of the document
            judge: Judge associated with the document
            court: Court where the document was filed
        """
        self.document_name = document_name
        self.subject = subject
        self.status = status
        
        # Optional fields
        self.timestamp = timestamp  # Keep as None if not provided
        self.case_name = case_name
        self.case_number = case_number
        self.author = author
        self.judge = judge
        self.court = court
        
    def __str__(self) -> str:
        """String representation of the metadata."""
        return f"Metadata: {self.document_name} ({self.subject})"
    
    def __repr__(self) -> str:
        """Detailed string representation of the metadata."""
        return f"Metadata('{self.document_name}', '{self.subject}', status='{self.status}')"
    
    def to_dict(self) -> Dict[str, Any]:
        """
        Convert metadata to dictionary format.
        
        Returns:
            Dictionary representation of the metadata
        """
        metadata_dict = {
            "document_name": self.document_name,
            "subject": self.subject,
            "status": self.status,
            "case_name": self.case_name,
            "case_number": self.case_number,
            "author": self.author,
            "judge": self.judge,
            "court": self.court
        }
        
        # Properly format timestamp as ISO format string if it exists
        if self.timestamp:
            metadata_dict["timestamp"] = self.timestamp.isoformat()
        else:
            metadata_dict["timestamp"] = None
            
        return metadata_dict
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Metadata":
        """
        Create a Metadata instance from a dictionary.
        
        Args:
            data: Dictionary containing metadata fields
            
        Returns:
            New Metadata instance
        """
        # Convert timestamp string to datetime if needed
        if "timestamp" in data and isinstance(data["timestamp"], str):
            data["timestamp"] = datetime.fromisoformat(data["timestamp"])
        return cls(**data)

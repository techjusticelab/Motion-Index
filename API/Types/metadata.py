from datetime import datetime
from typing import Dict, Any, Optional

class Metadata:
    def __init__(
        self,
        document_name: str,
        subject: str,
        timestamp: Optional[datetime] = None,
        case_name: Optional[str] = None,
        author: Optional[str] = None,
        judge: Optional[str] = None, # WARNING: This probso needs indexing #
        court: Optional[str] = None # WARNING: This probso needs indexing
    ):
        self.document_name = document_name
        self.subject = subject
        
        # Optional fields with defaults
        self.timestamp = timestamp or datetime.now()
        self.case_name = case_name
        self.author = author
        self.judge = judge
        self.court = court
        
    def __str__(self) -> str:
        return f"Metadata: {self.document_name} ({self.subject})"
    
    def __repr__(self) -> str:
        return f"Metadata('{self.document_name}', '{self.subject}', '{self.status}')"
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert metadata to dictionary format."""
        return {
            "document_name": self.document_name,
            "subject": self.subject,
            "status": self.status,
            "timestamp": self.timestamp,
            "case_name": self.case_name,
            "author": self.author,
            "judge": self.judge,
            "court": self.court
        }
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Metadata":
        """Create a Metadata instance from a dictionary."""
        if "timestamp" in data and isinstance(data["timestamp"], str):
            data["timestamp"] = datetime.fromisoformat(data["timestamp"])
        return cls(**data)

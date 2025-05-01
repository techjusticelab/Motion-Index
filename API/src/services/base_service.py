"""
Base service class for dependency injection and common functionality.
"""

from abc import ABC, abstractmethod
from typing import Any, Dict, Optional
from elasticsearch import Elasticsearch


class BaseService(ABC):
    """
    Base class for all services, providing common functionality and dependency injection.
    """
    
    def __init__(
        self, 
        es_client: Optional[Elasticsearch] = None,
        config: Optional[Dict[str, Any]] = None
    ):
        """
        Initialize the base service.
        
        Args:
            es_client: Optional Elasticsearch client
            config: Optional configuration dictionary
        """
        self.es_client = es_client
        self.config = config or {}
        
    @abstractmethod
    def initialize(self) -> None:
        """Initialize service-specific resources."""
        pass
    
    def get_config(self, key: str, default: Any = None) -> Any:
        """Get a configuration value with optional default."""
        return self.config.get(key, default)
    
    def set_config(self, key: str, value: Any) -> None:
        """Set a configuration value."""
        self.config[key] = value
    
    def close(self) -> None:
        """Clean up resources."""
        if self.es_client:
            self.es_client.close()
            self.es_client = None

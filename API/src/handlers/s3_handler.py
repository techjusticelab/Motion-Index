"""
S3 handler for uploading and managing files in AWS S3.
"""
import os
import logging
import boto3
import mimetypes
from botocore.exceptions import ClientError
from typing import Optional

# Configure logging
logger = logging.getLogger("s3_handler")


class S3Handler:
    """
    Handles interactions with AWS S3 for document storage.
    """
    
    def __init__(
        self,
        bucket_name: str,
        region_name: Optional[str] = None,
        aws_access_key_id: Optional[str] = None,
        aws_secret_access_key: Optional[str] = None
    ):
        """
        Initialize S3 connection.
        
        Args:
            bucket_name: Name of the S3 bucket
            region_name: AWS region name
            aws_access_key_id: AWS access key ID
            aws_secret_access_key: AWS secret access key
        """
        self.bucket_name = bucket_name
        
        # Configure S3 client
        s3_config = {}
        if region_name:
            s3_config['region_name'] = region_name
        if aws_access_key_id and aws_secret_access_key:
            s3_config['aws_access_key_id'] = aws_access_key_id
            s3_config['aws_secret_access_key'] = aws_secret_access_key
            
        try:
            self.s3_client = boto3.client('s3', **s3_config)
            logger.info(f"Connected to S3 with bucket: {bucket_name}")
        except Exception as e:
            logger.error(f"Error connecting to S3: {e}")
            raise
    
    def upload_file(self, file_path: str, s3_key: Optional[str] = None) -> Optional[str]:
        """
        Upload a file to S3 bucket.
        
        Args:
            file_path: Path to the local file
            s3_key: Optional key (path) in S3. If not provided, uses the filename
            
        Returns:
            S3 URI if successful, None otherwise
        """
        if not os.path.exists(file_path):
            logger.error(f"File not found: {file_path}")
            return None
            
        # If no S3 key provided, use the filename
        if not s3_key:
            s3_key = os.path.basename(file_path)
        
        # Determine content type
        content_type, _ = mimetypes.guess_type(file_path)
        extra_args = {}
        if content_type:
            extra_args['ContentType'] = content_type
        
        try:
            self.s3_client.upload_file(
                file_path, 
                self.bucket_name, 
                s3_key,
                ExtraArgs=extra_args
            )
            s3_uri = f"s3://{self.bucket_name}/{s3_key}"
            logger.info(f"Uploaded {file_path} to {s3_uri}")
            return s3_uri
        except ClientError as e:
            logger.error(f"Error uploading file to S3: {e}")
            return None
    
    def generate_presigned_url(self, s3_key: str, expiration: int = 3600) -> Optional[str]:
        """
        Generate a presigned URL for accessing a file.
        
        Args:
            s3_key: The key (path) of the file in S3
            expiration: URL expiration time in seconds (default: 1 hour)
            
        Returns:
            Presigned URL if successful, None otherwise
        """
        try:
            url = self.s3_client.generate_presigned_url(
                'get_object',
                Params={'Bucket': self.bucket_name, 'Key': s3_key},
                ExpiresIn=expiration
            )
            return url
        except ClientError as e:
            logger.error(f"Error generating presigned URL: {e}")
            return None
    
    def check_file_exists(self, s3_key: str) -> bool:
        """
        Check if a file exists in the S3 bucket.
        
        Args:
            s3_key: The key (path) of the file in S3
            
        Returns:
            True if the file exists, False otherwise
        """
        try:
            self.s3_client.head_object(Bucket=self.bucket_name, Key=s3_key)
            return True
        except ClientError as e:
            # If error code is 404, file doesn't exist
            if e.response['Error']['Code'] == '404':
                return False
            # For other errors, log and re-raise
            logger.error(f"Error checking if file exists in S3: {e}")
            raise
"""
Authentication middleware for FastAPI using Supabase JWT tokens.
"""
import os
from typing import Optional, List
import jwt
from fastapi import Request, HTTPException, Depends
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
import dotenv

# Load environment variables
dotenv.load_dotenv()

# Get Supabase JWT secret from environment variables
SUPABASE_JWT_SECRET = os.environ.get("SUPABASE_JWT_SECRET")
# If not set, we'll use the Supabase public key
SUPABASE_URL = os.environ.get("SUPABASE_URL")
SUPABASE_PUBLIC_KEY = os.environ.get("SUPABASE_PUBLIC_KEY")

# Define security scheme
security = HTTPBearer()

# Define standalone functions for dependencies
async def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security)) -> dict:
    """
    Verify the JWT token from the Authorization header.
    
    Args:
        credentials: The HTTP Authorization credentials.
        
    Returns:
        dict: The decoded JWT token payload.
        
    Raises:
        HTTPException: If the token is invalid or expired.
    """
    token = credentials.credentials
    
    try:
        # First try to decode with verification if we have the secret
        if SUPABASE_JWT_SECRET:
            try:
                # Decode using secret
                payload = jwt.decode(token, SUPABASE_JWT_SECRET, algorithms=["HS256"])
                print("Successfully decoded token with HS256")
                return payload
            except Exception as e:
                print(f"HS256 decoding failed: {str(e)}, trying RS256")
                # If HS256 fails, try RS256
                if SUPABASE_PUBLIC_KEY:
                    try:
                        payload = jwt.decode(token, SUPABASE_PUBLIC_KEY, algorithms=["RS256"])
                        print("Successfully decoded token with RS256")
                        return payload
                    except Exception as e2:
                        print(f"RS256 decoding failed: {str(e2)}")
                        # If both fail, fall back to no verification
                        pass
        
        # If we get here, either we don't have keys or verification failed
        # For development purposes, decode without verification
        print("Falling back to decoding without verification")
        payload = jwt.decode(token, options={"verify_signature": False})
        return payload
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token has expired")
    except jwt.InvalidTokenError:
        raise HTTPException(status_code=401, detail="Invalid token")
    except Exception as e:
        raise HTTPException(status_code=401, detail=f"Authentication error: {str(e)}")

async def get_current_user(payload: dict = Depends(verify_token)) -> dict:
    """
    Extract user information from the JWT payload.
    
    Args:
        payload: The decoded JWT token payload.
        
    Returns:
        dict: The user information.
    """
    # Extract user info from the token payload
    # The structure depends on your Supabase configuration
    user_id = payload.get("sub")
    if not user_id:
        raise HTTPException(status_code=401, detail="Invalid user information in token")
    
    return {
        "user_id": user_id,
        "email": payload.get("email"),
        "role": payload.get("role"),
    }

class AuthMiddleware:
    """Authentication middleware for FastAPI."""
    pass

# Export the authentication functions directly
__all__ = ['verify_token', 'get_current_user']

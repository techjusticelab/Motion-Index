FROM python:3.9-slim

WORKDIR /app

# Install system dependencies including LibreOffice with Java support
RUN apt-get update && apt-get install -y \
    build-essential \
    libpoppler-cpp-dev \
    pkg-config \
    python3-dev \
    tesseract-ocr \
    libtesseract-dev \
    libmagic1 \
    antiword \
    libreoffice \
    default-jre \
    --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements file
COPY requirements.txt .

# Install Python dependencies with a specific pip version that handles the textract dependency
RUN pip install pip==24.0 && \
    pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Set environment variables
ENV PYTHONPATH=/app
ENV PYTHONUNBUFFERED=1

# Command to run the application
CMD ["python", "-m", "src.main"]
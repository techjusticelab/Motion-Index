FROM python:3.9-slim

WORKDIR /app

# Install system dependencies in multiple steps to better manage disk space
# First, essential build tools and SSL
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    pkg-config \
    python3-dev \
    openssl \
    libssl-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Second, document processing tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpoppler-cpp-dev \
    tesseract-ocr \
    libtesseract-dev \
    libmagic1 \
    antiword \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Finally, LibreOffice and Java (these are the largest)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libreoffice-writer \
    libreoffice-calc \
    default-jre-headless \
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
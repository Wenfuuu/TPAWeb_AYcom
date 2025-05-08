FROM python:3.10-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    libsndfile1 \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements first for better layer caching
COPY service-ai/requirements.txt .
RUN pip install --upgrade pip
RUN pip install --no-cache-dir --timeout=120 -r requirements.txt

# Copy the rest of the application
COPY service-ai/. .

# Create directory for model if it doesn't exist
RUN mkdir -p /app/model

# Expose Flask ports
EXPOSE 5000

CMD ["python", "app.py"]
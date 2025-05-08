from flask import Flask, request, jsonify
from flask_cors import CORS
import tempfile, base64, os
import librosa
from dotenv import load_dotenv
from transformers import pipeline
import torch

# Load environment variables
load_dotenv()

app = Flask(__name__)
CORS(app)

# Global variables
asr_model = None
model_path = os.environ.get('MODEL_PATH', './model')

def load_model():
    """Initialize the ASR model"""
    global asr_model
    print(f"Loading speech recognition model from {model_path}")
    
    try:
        # Load the speech recognition model
        asr_model = pipeline(
            "automatic-speech-recognition",
            model=model_path,
            device=0 if torch.cuda.is_available() else -1
        )
        print("Model loaded successfully")
    except Exception as e:
        print(f"Error loading model: {str(e)}")
        asr_model = None

@app.route('/transcribe', methods=['POST'])
def transcribe_audio():
    """Endpoint to transcribe audio to text"""
    if asr_model is None:
        return jsonify({"error": "Model not loaded"}), 503

    try:
        if 'file' in request.files:
            # Handle file upload
            audio_file = request.files['file']
            with tempfile.NamedTemporaryFile(delete=False, suffix='.wav') as temp_file:
                audio_file.save(temp_file.name)
                temp_filename = temp_file.name
        elif request.is_json and 'audio_data' in request.json:
            # Handle base64 encoded audio
            audio_data = base64.b64decode(request.json['audio_data'])
            with tempfile.NamedTemporaryFile(delete=False, suffix='.wav') as temp_file:
                temp_file.write(audio_data)
                temp_filename = temp_file.name
        else:
            return jsonify({"error": "No audio data provided"}), 400

        # Load and preprocess audio
        audio, sample_rate = librosa.load(temp_filename, sr=16000)
        os.unlink(temp_filename)

        # Perform transcription
        result = asr_model(audio)

        return jsonify({
            "text": result["text"],
            "success": True
        })

    except Exception as e:
        return jsonify({"error": str(e), "success": False}), 500


# Make sure the app is reachable externally on port 5000
if __name__ == "__main__":
    load_model()

    app.run(host="0.0.0.0", port=5000)
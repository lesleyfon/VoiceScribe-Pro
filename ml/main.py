from fastapi import FastAPI, File
from transformers import WhisperProcessor, WhisperForConditionalGeneration
import io
import librosa
import time
import numpy as np
from audio_chunking import router as audio_chunking_router
from libs.chunk_audio import process_audio_array_to_chunks


app = FastAPI()
app.include_router(audio_chunking_router)

processor = WhisperProcessor.from_pretrained("openai/whisper-base")
model = WhisperForConditionalGeneration.from_pretrained("openai/whisper-base")  # fastest for some reason


CHUNK_LENGTH_SAMPLES = 30 * 16000  # 30 seconds
CHUNK_OVERLAP_SAMPLES = 5 * 16000  # 5 seconds - the amo




@app.get('/')
def read_root():
  return {"Well": "Come"}

@app.post('/process-audio')
async def process_audio(audio_file: bytes = File()):
    start_time = time.time()
    audio_array, sample_rate = librosa.load(
        io.BytesIO(audio_file), 
        sr=16000,
        mono=True
    )
    print("Array length:", len(audio_array), "Sample rate:", sample_rate)
    print("Duration (s):", len(audio_array) / sample_rate)
    if len(audio_array) == 0:
        return "No Audio detected"

    # Chunk the audio
    chunks = process_audio_array_to_chunks(audio_array, CHUNK_LENGTH_SAMPLES, CHUNK_OVERLAP_SAMPLES)
    transcriptions = []

    for i, chunk in enumerate(chunks):
        print(f"Transcribing chunk {i+1}/{len(chunks)}")
        if len(chunk) == 0:
            continue
        try:
            inputs = processor(
                chunk,
                sampling_rate=sample_rate,
                return_tensors="pt"
            )
            input_features = inputs.input_features
            attention_mask = inputs.get("attention_mask", None)
            
            if attention_mask is not None:
                predicted_ids = model.generate(
                    input_features,
                    attention_mask=attention_mask,
                    task="transcribe",
                    language="en"
                )
            else:
                predicted_ids = model.generate(
                    input_features,
                    task="transcribe",
                    language="en"
                )
            transcription = processor.batch_decode(predicted_ids, skip_special_tokens=True)[0]
            # Stream the response here
            transcriptions.append(transcription)
        except Exception e:
            # Decide how to handle failed chunks, e.g., append a placeholder or skip
            # For now, skipping failed chunks from the list to merge
            print(f"Error transcribing chunk {i+1}: {str(e)}")
            
        

    full_transcription = " ".join(transcriptions)
    end_time = time.time()
    request_duration = end_time - start_time

    return {"transcription": full_transcription, "request_duration": request_duration}


# Instead of waiting for all the response to fullfil, try streaming the response chunks.

# uvicorn main:app --host 0.0.0.0 --port 9090




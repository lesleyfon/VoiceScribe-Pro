from fastapi import APIRouter, File
import librosa
import time
import io

CHUNK_LENGTH_SAMPLES = 30 * 16000  # 30 seconds
CHUNK_OVERLAP_SAMPLES = 5 * 16000  # 5 seconds

def process_audio_array_to_chunks(audio_array, chunk_size, overlap):
    chunks = []
    start = 0
    while start < len(audio_array):
        end = min(start + chunk_size, len(audio_array))
        chunks.append(audio_array[start:end])
        if end == len(audio_array):
            break
        start += chunk_size - overlap
    return chunks



def process_audio_file_audio_chunks(audio_file: bytes = File(), request_start_time=time.time()):
    try:
        audio_array, sample_rate = librosa.load(
            io.BytesIO(audio_file), 
            sr=16000, 
            mono=True
        )
    except Exception as e:
        return {
            "error": f"Error loading audio: {str(e)}",
            "request_duration": time.time() - request_start_time,
        }

    print("Array length:", len(audio_array), "Sample rate:", sample_rate)
    duration_s = len(audio_array) / sample_rate
    print("Duration (s):", duration_s)
    
    if len(audio_array) == 0:
        return [], None
    
    audio_chunks = process_audio_array_to_chunks( 
        audio_array, 
        CHUNK_LENGTH_SAMPLES, 
        CHUNK_OVERLAP_SAMPLES
    )
    
    return audio_chunks, sample_rate
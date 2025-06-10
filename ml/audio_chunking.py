import time # For tracking processing duration (optional)
import re # For regex checking during audio chunk merging
from fastapi import APIRouter, File
from transformers import WhisperProcessor, WhisperForConditionalGeneration
import io
import librosa
import numpy as np
from libs.chunk_audio import process_audio_array_to_chunks, process_audio_file_audio_chunks

CHUNK_LENGTH_SAMPLES = 30 * 16000  # 30 seconds
CHUNK_OVERLAP_SAMPLES = 5 * 16000  # 5 seconds


# This function is adapted from the Groq audio chunking tutorial's principles
def find_best_overlap_and_merge(
    text1: str, text2: str, match_by_words: bool = True
) -> tuple[str, float]:
    """
    Merges two texts by finding the best overlap.
    Returns the merged text and the best overlap score found.
    Score is -1.0 if no overlap met the criteria.
    """
    if not text1:
        return text2, 1.0  # Perfect merge if text1 is empty
    if not text2:
        return text1, 1.0  # Perfect merge if text2 is empty

    # Tokenize
    if match_by_words:
        seq1 = [word for word in re.split(r"(\s+\S+)", text1) if word.strip()]
        seq2 = [word for word in re.split(r"(\s+\S+)", text2) if word.strip()]
    else:
        seq1 = list(text1)
        seq2 = list(text2)

    if not seq1:
        # text1 was whitespace only
        return ("".join(seq2) if match_by_words else "".join(seq2)), 1.0
    if not seq2:
        # text2 was whitespace only
        return ("".join(seq1) if match_by_words else "".join(seq1)), 1.0

    max_overlap_len = 0
    best_overlap_score = -1.0  # Initialize to indicate no good overlap found yet

    # Iterate over possible overlap lengths, from min(len(seq1), len(seq2)) down to 1
    for overlap_len in range(min(len(seq1), len(seq2)), 0, -1):
        suffix_seq1 = seq1[-overlap_len:]
        prefix_seq2 = seq2[:overlap_len]

        matches = sum(
            s1_token == s2_token
            for s1_token, s2_token in zip(suffix_seq1, prefix_seq2)
        )

        current_score = matches / overlap_len

        # Prefer longer overlaps with high match rate
        if current_score > 0.75:  # Threshold for a "good" overlap
            # We are iterating from longest possible overlap downwards,
            # so the first one that meets the criteria is the longest good one.
            max_overlap_len = overlap_len
            best_overlap_score = current_score
            break  # Found the longest good overlap

    if max_overlap_len > 0:
        # Merge: text1 up to the overlap + text2 (which starts with the overlap)
        merged_sequence = seq1[:-max_overlap_len] + seq2
    else:
        # No significant overlap found, concatenate with a space
        merged_sequence = seq1 + (
            [" "] if match_by_words and seq1 and seq2 else []
        ) + seq2
        # best_overlap_score remains -1.0 or its last value if no overlap > 0.75 was found
        # If we want to assign a score for simple concatenation, we could do it here.
        # For now, -1.0 signifies no "good" overlap based on the threshold.

    return "".join(merged_sequence), best_overlap_score


def merge_transcription_list(
    transcriptions: list[str],
) -> tuple[str, list[float]]:
    """
    Merges a list of transcriptions.
    Returns the final merged text and a list of overlap scores for each merge.
    """
    if not transcriptions:
        return "", []

    # Filter out empty or whitespace-only transcriptions initially
    valid_transcriptions = [t for t in transcriptions if t.strip()]
    if not valid_transcriptions:
        return "", []

    merged_text = valid_transcriptions[0]
    overlap_scores = []

    for i in range(1, len(valid_transcriptions)):
        current_chunk_text = valid_transcriptions[i]
        # No need to check .strip() here again as we pre-filtered
        merged_text, score = find_best_overlap_and_merge(
            merged_text, current_chunk_text
        )
        overlap_scores.append(score)
    return merged_text, overlap_scores


router = APIRouter()

processor = WhisperProcessor.from_pretrained("openai/whisper-base")
model = WhisperForConditionalGeneration.from_pretrained(
    "openai/whisper-base"
)








@router.post('/process-audio/transcribe')
async def process_audio_endpoint(audio_file: bytes = File()):
    request_start_time = time.time()

    # Chunk the audio
    audio_chunks, sample_rate = process_audio_file_audio_chunks(audio_file, request_start_time)

    if len(audio_chunks) == 0:
        return {
            "transcription": "No Audio detected",
            "request_duration": time.time() - request_start_time,
        }

    transcriptions = []
    print(f"Processing {len(audio_chunks)} chunks...")

    for i, chunk_data in enumerate(audio_chunks): # Renamed variable for clarity
        print(f"Transcribing chunk {i+1}/{len(audio_chunks)}")
        if len(chunk_data) == 0:
            continue # Skip empty chunks for transcription list
        try:
            inputs = processor(
                chunk_data, sampling_rate=sample_rate, return_tensors="pt"
            )
            input_features = inputs.input_features
            generate_kwargs = {"task": "transcribe", "language": "en"}


            predicted_ids = model.generate(input_features, **generate_kwargs)
            transcription_text = processor.batch_decode( # Renamed variable
                predicted_ids, skip_special_tokens=True
            )[0]
            # Ideally we would want to stream results from here. 
            transcriptions.append(transcription_text.strip())
        except Exception as e:
            # Decide how to handle failed chunks, e.g., append a placeholder or skip
            # For now, skipping failed chunks from the list to merge
            print(f"Error transcribing chunk {i+1}: {str(e)}")


    for i, t in enumerate(transcriptions):
        print(f"Chunk {i+1}: {t}")

    full_transcription, overlap_scores = merge_transcription_list(
        transcriptions
    )
    
    avg_overlap_score = (
        sum(overlap_scores) / len(overlap_scores) if overlap_scores else -1.0
    )


    processing_end_time = time.time() 
    request_duration = processing_end_time - request_start_time

    print(f"\nFinal merged transcription: {full_transcription}")
    print(f"Overlap scores for merges: {overlap_scores}")
    print(f"Average overlap score: {avg_overlap_score:.2f}")
    print(f"Total request duration: {request_duration:.2f}s")

    return {
        "transcription": full_transcription,
        "overlap_scores": overlap_scores,
        "average_overlap_score": avg_overlap_score,
        "request_duration": request_duration,
    }


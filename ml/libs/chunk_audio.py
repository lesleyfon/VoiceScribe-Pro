
def process_chunk(audio_array, chunk_size, overlap):
    chunks = []
    start = 0
    while start < len(audio_array):
        end = min(start + chunk_size, len(audio_array))
        chunks.append(audio_array[start:end])
        if end == len(audio_array):
            break
        start += chunk_size - overlap
    return chunks
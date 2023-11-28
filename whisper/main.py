from pydantic import BaseModel

from fastapi import FastAPI

from faster_whisper import WhisperModel

model_size = "medium"
model = WhisperModel(model_size, device="cpu", compute_type="int8")

app = FastAPI()

class TranscriptionRequest(BaseModel):
    audio_path: str

@app.post("/api/v1/transcribe")
def read_root(req: TranscriptionRequest):
    print("recebido audio", req.audio_path)
    segments, info = model.transcribe(req.audio_path, vad_filter=True, language="pt", beam_size=5)
    transcription = []
    for segment in segments:
        transcription.append("%s" % (segment.text))
    return {"result": transcription, "language": info.language, "probability": info.language_probability}

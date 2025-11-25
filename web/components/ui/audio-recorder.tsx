"use client";
import { cn } from "@/lib/utils";
import { useEffect, useRef, useState } from "react";
export default function AudioRecorder() {
	const [mediaRecorder, setMediaRecorder] = useState<null | MediaRecorder>();
	const [isRecording, setIsRecording] = useState<boolean>(false);
	const [error, setError] = useState<string | null>(null);

	const audioRef = useRef(null);
	const streamRef = useRef<MediaStream | null>(null);
	const chunksRef = useRef<Blob[]>([]);
	const audioUrlRef = useRef<string | null>(null);

	useEffect(() => {
		if (window === undefined) return;

		(async () => {
			if (!navigator?.mediaDevices) {
				setError("getUserMedia not supported on your browser!");
				return;
			}

			try {
				const stream = await navigator.mediaDevices.getUserMedia({
					audio: true,
				});

				streamRef.current = stream;

				const mediaRecorderInstance = new MediaRecorder(stream);
				setMediaRecorder(mediaRecorderInstance);

				mediaRecorderInstance.ondataavailable = (e) => {
					if (e.data.size > 0) {
						chunksRef.current.push(e.data);
					}
				};
				mediaRecorderInstance.onstop = () => {
					console.log("stopping");
					if (!audioRef.current) {
						console.error(`
            ${"=".repeat(20)}
            = NO AUDIO REF
            ${"=".repeat(20)}
            `);
						return;
					}

					if (audioRef.current) {
						URL.revokeObjectURL(audioRef.current);
					}
					const audio = audioRef.current as HTMLAudioElement;

					const blob = new Blob(chunksRef.current, {
						type: "audio/ogg; codecs=opus",
					});

					chunksRef.current = [];

					const audioURL = URL.createObjectURL(blob);
					audioUrlRef.current = audioURL;
					audio.src = audioURL;
				};
			} catch (err) {
				setError(
					`Error accessing microphone: ${
						err instanceof Error ? err.message : String(err)
					}`
				);
			}
		})();

		return () => {
			if (streamRef.current) {
				const tracks = streamRef.current.getTracks();
				for (const track of tracks) {
					track.stop();
				}
				if (audioUrlRef.current) {
					URL.revokeObjectURL(audioUrlRef.current);
				}
			}
		};
	}, []);

	const handleRecord = () => {
		if (!mediaRecorder) return;
		mediaRecorder.start();
		setIsRecording(true);
	};

	const handleStop = () => {
		if (!mediaRecorder) return;
		mediaRecorder.stop();
		setIsRecording(false);
	};

	if (error) {
		return (
			<article className="text-red-600">
				<p>Error: {error}</p>
			</article>
		);
	}
	return (
		<>
			<article>
				<audio controls ref={audioRef} className="mb-4">
					<track
						kind="captions"
						src="/captions.vtt"
						srcLang="en"
						label="English captions"
						default
					/>
				</audio>
				<p>Your Clip name</p>
				<div className=" flex w-[300px] justify-between">
					<button
						type="button"
						className={cn(
							"cursor-pointer bg-[#2288CC] px-10 py-3 transition-colors",
							"disabled:cursor-not-allowed disabled:opacity-50",
							isRecording && "bg-red-600"
						)}
						onClick={handleRecord}
					>
						{isRecording ? "Recording..." : "Record"}
					</button>
					<button
						type="button"
						className={cn(
							"cursor-pointer bg-[#2288CC] px-10 py-3 transition-colors",
							"disabled:cursor-not-allowed disabled:opacity-50"
						)}
						onClick={handleStop}
					>
						Stop
					</button>
				</div>
			</article>
		</>
	);
}

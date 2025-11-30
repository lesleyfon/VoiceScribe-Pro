"use client";
import { useSocket } from "@/hooks/use-sockets";
import { cn } from "@/lib/utils";
import { MicIcon, MicOff } from "lucide-react";
import { useEffect, useRef, useState } from "react";
export default function AudioRecorder() {
	const [mediaRecorder, setMediaRecorder] = useState<null | MediaRecorder>();
	const [isRecording, setIsRecording] = useState<boolean>(false);
	const [error, setError] = useState<string | null>(null);

	const audioRef = useRef(null);
	const streamRef = useRef<MediaStream | null>(null);
	const chunksRef = useRef<Blob[]>([]);
	const audioUrlRef = useRef<string | null>(null);

	const { error: socketError } = useSocket({
		url: "ws://127.0.0.1:8000/ws", // TODO: WE DONT NEED THIS. THIS SHOULD BE READ FROM AN ENV FILE
		onConnect: () => console.log("Connected!"),
		onDisconnect: (reason) => console.log("Disconnected:", reason),
		onError: (err) => console.error("Socket error:", err),
		onAuthSuccess: () => console.log("Authenticated!"),
		onAuthError: (err) => console.error("Auth failed:", err),
	});

	useEffect(() => {
		//TODO: Move this into a hook
		if (typeof window === "undefined") return;

		const setupMediaRecorder = async () => {
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
		};
		setupMediaRecorder();

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

	if (error || socketError) {
		return (
			<>
				<article className="text-red-600">
					<p>Error: {error ?? socketError}</p>
				</article>
			</>
		);
	}
	return (
		<>
			<article>
				{/* biome-ignore lint/a11y/useMediaCaption: <explanation> */}
				<audio controls ref={audioRef} className="mb-4" />
				<div className=" flex w-[300px] justify-between">
					<button
						type="button"
						className={cn(
							"cursor-pointer bg-[#2288CC]  py-3 transition-colors flex justify-center align-middle content-center flex-[0.4]",
							"disabled:cursor-not-allowed disabled:opacity-50",
							isRecording && "bg-red-600"
						)}
						disabled={isRecording}
						onClick={handleRecord}
					>
						<div className="mr-0.5">{isRecording ? "Recording..." : "Record"}</div>
						<MicIcon />
					</button>
					<button
						type="button"
						className={cn(
							"cursor-pointer bg-[#2288CC]  py-3 transition-colors flex justify-center align-middle content-center flex-[0.4]",
							"disabled:cursor-not-allowed disabled:opacity-50"
						)}
						onClick={handleStop}
					>
						<div className="mr-0.5">Stop</div> <MicOff />
					</button>
				</div>
			</article>
		</>
	);
}

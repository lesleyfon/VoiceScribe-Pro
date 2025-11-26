"use client";
import { cn } from "@/lib/utils";
import { useAuth } from "@clerk/nextjs";
import { MicIcon, MicOff } from "lucide-react";
import { useEffect, useRef, useState } from "react";
export default function AudioRecorder() {
	const { getToken } = useAuth();

	const [mediaRecorder, setMediaRecorder] = useState<null | MediaRecorder>();

	const [isRecording, setIsRecording] = useState<boolean>(false);
	const [error, setError] = useState<string | null>(null);
	const [isConnected, setIsConnected] = useState(false);

	const socketRef = useRef<WebSocket | null>(null);
	const audioRef = useRef(null);
	const streamRef = useRef<MediaStream | null>(null);
	const chunksRef = useRef<Blob[]>([]);
	const audioUrlRef = useRef<string | null>(null);

	useEffect(() => {
		const setupWebSocket = async () => {
			const token = await getToken({ template: "App-Template" });

			// Use native WebSocket with token in query param or subprotocol
			// Option 1: Query parameter
			socketRef.current = new WebSocket(`ws://127.0.0.1:8000/ws?token=${token}`);

			const socket = socketRef.current;

			socket.onopen = () => {
				console.log("WebSocket connected");
				setIsConnected(true);
			};

			socket.onmessage = (event) => {
				console.log("Message from server:", event.data);
				try {
					const data = JSON.parse(event.data);
					console.log("Parsed data:", data);
				} catch {
					console.log("Raw data:", event.data);
				}
			};

			socket.onerror = (error) => {
				console.error("WebSocket error:", error);
				setIsConnected(false);
			};

			socket.onclose = (event) => {
				console.log("WebSocket closed:", event.code, event.reason);
				setIsConnected(false);
			};
		};

		setupWebSocket();

		return () => {
			if (socketRef.current) {
				socketRef.current.close();
			}
		};
	}, [getToken]);

	useEffect(() => {
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
		if (socketRef.current?.readyState === WebSocket.OPEN) {
			socketRef.current.send(
				JSON.stringify({
					ping: "ping",
				})
			);
		}

		if (!mediaRecorder) return;
		mediaRecorder.stop();
		setIsRecording(false);
	};

	const sendMessage = (message: string) => {
		if (socketRef.current?.readyState === WebSocket.OPEN) {
			socketRef.current.send(message);
			console.log("Sent:", message);
		} else {
			console.error("WebSocket is not connected");
		}
	};

	if (error) {
		return (
			<>
				<article className="text-red-600">
					<p>Error: {error}</p>
				</article>
			</>
		);
	}
	return (
		<>
			<article>
				{/* biome-ignore lint/a11y/useMediaCaption: <explanation> */}
				<audio controls ref={audioRef} className="mb-4" />
				<p>Your Clip name</p>
				<div className=" flex w-[300px] justify-between">
					<button
						type="button"
						className={cn(
							"cursor-pointer bg-[#2288CC] px-10 py-3 transition-colors flex justify-center align-middle content-center ",
							"disabled:cursor-not-allowed disabled:opacity-50",
							isRecording && "bg-red-600"
						)}
						onClick={handleRecord}
					>
						<div className="mr-0.5">{isRecording ? "Recording..." : "Record"}</div>
						<MicIcon />
					</button>
					<button
						type="button"
						className={cn(
							"cursor-pointer bg-[#2288CC] px-10 py-3 transition-colors flex justify-center align-middle content-center ",
							"disabled:cursor-not-allowed disabled:opacity-50"
						)}
						onClick={handleStop}
					>
						<div className="mr-0.5">Stop</div> <MicOff />
					</button>
				</div>
				<button onClick={() => sendMessage("Message from click me button")} type="button">
					Click me
				</button>
			</article>
		</>
	);
}

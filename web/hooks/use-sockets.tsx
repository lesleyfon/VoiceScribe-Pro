import { useAuth } from "@clerk/nextjs";
import { useEffect, useRef, useState } from "react";

type MessageHandler = (...args: unknown[]) => void;

interface UseSocketOptions {
	url: string;
	onConnect?: () => void;
	onDisconnect?: (reason: string) => void;
	onError?: (error: Error) => void;
	onAuthSuccess?: () => void;
	onAuthError?: (error: string) => void;
}
export function useSocket({ onConnect, onDisconnect, onError }: UseSocketOptions) {
	const { getToken } = useAuth();

	const socketRef = useRef<WebSocket | null>(null);
	const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
	const eventHandlersRef = useRef<Map<string, Set<MessageHandler>>>(new Map());

	const [isConnected, setIsConnected] = useState<boolean>(false);
	const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
	const [error, setError] = useState<string | null>(null);
	const WS_URL = "ws://127.0.0.1:8000/ws";

	useEffect(() => {
		const connect = async () => {
			const token = await getToken();
			try {
				socketRef.current = new WebSocket(`${WS_URL}?token=${token}`);
				const socket = socketRef.current;

				socket.onopen = () => {
					console.log("WebSocket connected");
					setIsConnected(true);
					setError(null);
					onConnect?.();

					// Send Authentication token
					const payload = JSON.stringify({
						type: "auth",
						token: token,
					});

					socket.send(payload);
				};

				// socket.onmessage = (event) => {
				// 	try {
				// 		const data = JSON.parse(event.data);

				// 		// Handle auth response
				// 		if (data.type === "auth_success") {
				// 			setIsAuthenticated(true);
				// 			onAuthSuccess?.();
				// 			return;
				// 		}

				// 		if (data.type === "auth_error" || data.error) {
				// 			const errorMsg = data.error || data.message || "Authentication failed";
				// 			setError(errorMsg);
				// 			onAuthError?.(errorMsg);
				// 			socket.close();
				// 			return;
				// 		}

				// 		// Emit to registered event handlers
				// 		const eventType = data.type || data.event || "message";
				// 		const handlers = eventHandlersRef.current.get(eventType);

				// 		if (handlers) {
				// 			for (const handler of handlers) {
				// 				try {
				// 					handler(data);
				// 				} catch (err) {
				// 					console.error(`Error in event handler for ${eventType}:`, err);
				// 				}
				// 			}
				// 		}

				// 		// Also emit to wildcard handlers
				// 		const wildcardHandlers = eventHandlersRef.current.get("*");
				// 		if (wildcardHandlers) {
				// 			for (const handler of wildcardHandlers) {
				// 				handler(data);
				// 			}
				// 		}
				// 	} catch (err) {
				// 		console.error("Failed to parse message:", event.data, err);
				// 	}
				// };

				socket.onerror = (event) => {
					const error = new Error("WebSocket error occurred");
					console.error("WebSocket error:", event);
					setError(error.message);
					setIsConnected(false);
					onError?.(error);
				};

				socket.onclose = (event) => {
					const reason = event.reason || `Connection closed with code ${event.code}`;
					console.log(`Websocket closed: ${reason}`);

					setIsConnected(false);

					setIsAuthenticated(false); //?? TODO: Verify the reason, if it is an unauthenticated error, set this to false
					onDisconnect?.(reason);

					// Auto-reconnect after 3 seconds if not a normal closure
					if (event.code !== 1000 && event.code !== 1001) {
						reconnectTimeoutRef.current = setTimeout(() => {
							console.log("Attempting to reconnect...");
							connect();
						}, 3000);
					}
				};
			} catch (err) {
				const error = err instanceof Error ? err : new Error(String(err));
				setError(error.message);
				onError?.(error);
			}
		};

		//
		connect();
		return () => {
			if (reconnectTimeoutRef.current) {
				clearTimeout(reconnectTimeoutRef.current);
			}

			if (socketRef.current) {
				//TODO: This seems to be closing the connection on every mount.
				// socketRef.current.close(3000, "Component unmounting");
			}
			eventHandlersRef.current.clear();
		};
	}, [getToken, onConnect, onDisconnect, onError]);

	type DataType = string | ArrayBufferLike | Blob | ArrayBufferView<ArrayBufferLike>;
	const emit = (event: string, data?: DataType) => {
		if (socketRef.current?.readyState === WebSocket.OPEN) {
			const message: { type: string; data?: DataType } = {
				type: event,
			};
			if (data) {
				message.data = data;
			}
			const payload = JSON.stringify(message);

			// Emit message to the server
			socketRef.current.send(payload);
		} else {
			console.warn("WebSocket is not connected. Cannot send message.");
		}
	};

	const off = (event: string, handler?: MessageHandler) => {
		if (!handler) {
			// Remove all handlers for this event
			eventHandlersRef.current.delete(event);
		} else {
			// Remove specific handler
			const handlers = eventHandlersRef.current.get(event);
			if (handlers) {
				handlers.delete(handler);
				if (handlers.size === 0) {
					eventHandlersRef.current.delete(event);
				}
			}
		}
	};

	const on = (event: string, handler: MessageHandler) => {
		if (!eventHandlersRef.current.has(event)) {
			eventHandlersRef.current.set(event, new Set());
		}
		eventHandlersRef.current.get(event)?.add(handler);
	};

	const close = () => {
		if (reconnectTimeoutRef.current) {
			clearTimeout(reconnectTimeoutRef.current);
		}
		socketRef.current?.close(1000, "Manual close");
	};
	return {
		socket: socketRef.current,
		isConnected,
		isAuthenticated,
		error,
		emit,
		off,
		close,
		on,
	};
}

import { EventSourcePolyfill } from "event-source-polyfill";
import { useCallback,useEffect,useRef } from "react";

interface SSEOptions {
  path?: string;
  events: string[];
  appId?: string;
  accessToken?: string;
  onEvent?: (type: string, data: unknown) => void;
  onError?: (event: unknown) => void;
  shouldRun?: boolean;
  baseUrl: string;
}

export function useSSE(options: SSEOptions) {
  const {
  path,
  events,
    shouldRun,
    accessToken,
    appId,
    onError,
    onEvent,
    baseUrl,
  } = options;
  const eventSourceRef = useRef<EventSource | null>(null);

  const close = useCallback(() => {
    eventSourceRef.current?.close();
    eventSourceRef.current = null;
  }, []);

  useEffect(() => {
    if (
    shouldRun &&
    events.length > 0 &&
    accessToken &&
    accessToken.length > 0 &&
    appId &&
    appId.length > 0
  ) {
    const url = new URL(baseUrl);

      if (path) {
        url.pathname = path;
      }

      const options: {
        heartbeatTimeout: number;
        headers?: { [key: string]: string };
      } = {
        heartbeatTimeout: 999999999,
        headers: {
          Authorization: `Bearer ${accessToken}`,
          "x-app-id": appId ?? "",
        },
      };

      const eventSource = new EventSourcePolyfill(url.toString(), options);
      eventSourceRef.current = eventSource;

    const handleMessage = (event: MessageEvent) => {
    try {
      const parsedData = JSON.parse(event.data);
      if (onEvent) onEvent(event.type, parsedData);
    } catch {
      if (onEvent) onEvent(event.type, event.data);
    }
    };

    const handleError = (event: unknown) => {
    if (onError) onError(event);
    };

    for (const eventType of events) {
    eventSourceRef.current?.addEventListener(eventType, handleMessage);
    }

    eventSource.addEventListener("error", handleError);
  } else {
    close();
    }

    return () => {
      close();
    };
  }, [
    path,
    events,
    accessToken,
    shouldRun,
    onEvent,
    onError,
    appId,
    baseUrl,
    close,
  ]);

  return { close };
}

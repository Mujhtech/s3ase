import { useEffect, useRef, useState } from "react";
import { EventSourcePolyfill } from "event-source-polyfill";
import { isEqual, isEmpty } from "lodash-es";

interface SSEOptions {
  path?: string;
  events: string[];
  appId?: string;
  accessToken?: string;
  onEvent?: (type: string, data: any) => void;
  onError?: (event: Event) => void;
  shouldRun?: boolean;
  baseUrl: string;
}

export function useSSE<T = any>(options: SSEOptions) {
  const {
    path,
    events: _events,
    shouldRun,
    accessToken,
    appId,
    onError,
    onEvent,
    baseUrl,
  } = options;
  const [events, setEvents] = useState(_events);
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    if (!isEqual(events, _events)) {
      setEvents(_events);
    }
  }, [_events, setEvents, events]);

  useEffect(() => {
    if (
      shouldRun &&
      events.length > 0 &&
      accessToken &&
      isEmpty(accessToken) == false &&
      appId &&
      isEmpty(appId) == false
    ) {
      let url = new URL(baseUrl);

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

      eventSource.addEventListener("ping", (event: any) => {
        console.log(`Got ping from server!`);
      });

      const handleMessage = (event: MessageEvent) => {
        console.log(event);
        try {
          const parsedData = JSON.parse(event.data);
          if (onEvent) onEvent(event.type, parsedData);
        } catch (err) {
          if (onEvent) onEvent(event.type, event.data as T);
        }
      };

      const handleError = (event: Event) => {
        if (onError) onError(event);
        // setError(
        //   error instanceof Error ? error : new Error("SSE connection failed")
        // );
        eventSourceRef?.current?.close();
      };

      for (const i in events) {
        const eventType = events[i];
        eventSourceRef.current?.addEventListener(eventType, handleMessage);
      }

      eventSource.addEventListener("error", (error: any) => handleError);
    } else {
      close();
    }

    return () => {
      close();
    };
  }, [path, events, accessToken, shouldRun, onEvent, onError, appId]);

  // Method to manually close the connection
  const close = () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }
  };

  return { data, error, close };
}

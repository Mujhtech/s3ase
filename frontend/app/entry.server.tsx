import {
  createReadableStreamFromReadable,
  type EntryContext,
} from "@remix-run/node"; // or cloudflare/deno
import { RemixServer } from "@remix-run/react";
import { renderToPipeableStream } from "react-dom/server";
import { PassThrough } from "stream";
import { isbot } from "isbot";

const ABORT_DELAY = 30000;

export default function handleRequest(
  request: Request,
  responseStatusCode: number,
  responseHeaders: Headers,
  remixContext: EntryContext
) {
  // If the request is from a bot, we want to wait for the full
  // response to render before sending it to the client. This
  // ensures that bots can see the full page content.
  if (isbot(request.headers.get("user-agent"))) {
    return handleBotRequest(
      request,
      responseStatusCode,
      responseHeaders,
      remixContext
    );
  }

  return handleBrowserRequest(
    request,
    responseStatusCode,
    responseHeaders,
    remixContext
  );
}

function handleBotRequest(
  request: Request,
  responseStatusCode: number,
  responseHeaders: Headers,
  remixContext: EntryContext
) {
  return new Promise((resolve, reject) => {
    let shellRendered = false;
    const { pipe, abort } = renderToPipeableStream(
      <RemixServer
        context={remixContext}
        url={request.url}
        abortDelay={ABORT_DELAY}
      />,
      {
        onAllReady() {
          shellRendered = true;
          const body = new PassThrough();
          const stream = createReadableStreamFromReadable(body);

          responseHeaders.set("Content-Type", "text/html");

          resolve(
            new Response(stream, {
              headers: responseHeaders,
              status: responseStatusCode,
            })
          );

          pipe(body);
        },
        onShellError(error: unknown) {
          reject(error);
        },
        onError(error: unknown) {
          responseStatusCode = 500;
          // Log streaming rendering errors from inside the shell.  Don't log
          // errors encountered during initial shell rendering since they'll
          // reject and get logged in handleDocumentRequest.
          if (shellRendered) {
            console.error(error);
          }
        },
      }
    );

    setTimeout(abort, ABORT_DELAY);
  });
}

function handleBrowserRequest(
  request: Request,
  responseStatusCode: number,
  responseHeaders: Headers,
  remixContext: EntryContext
) {
  return new Promise((resolve, reject) => {
    let shellRendered = false;
    const { pipe, abort } = renderToPipeableStream(
      <RemixServer
        context={remixContext}
        url={request.url}
        abortDelay={ABORT_DELAY}
      />,
      {
        onShellReady() {
          shellRendered = true;
          const body = new PassThrough();
          const stream = createReadableStreamFromReadable(body);

          responseHeaders.set("Content-Type", "text/html");

          resolve(
            new Response(stream, {
              headers: responseHeaders,
              status: responseStatusCode,
            })
          );

          pipe(body);
        },
        onShellError(error: unknown) {
          reject(error);
        },
        onError(error: unknown) {
          responseStatusCode = 500;
          // Log streaming rendering errors from inside the shell.  Don't log
          // errors encountered during initial shell rendering since they'll
          // reject and get logged in handleDocumentRequest.
          if (shellRendered) {
            console.error(error);
          }
        },
      }
    );

    setTimeout(abort, ABORT_DELAY);
  });
}

// export function handleError(
//   error: unknown,
//   { request, params, context }: DataFunctionArgs
// ) {
//   logError(error, request);
// }

// Worker.init().catch((error) => {
//   logError(error);
// });

// function logError(error: unknown, request?: Request) {
//   console.error(error);

//   if (error instanceof Error && error.message.startsWith("There are locked jobs present")) {
//     console.log("⚠️  graphile-worker migration issue detected!");
//   }
// }

// process.on("uncaughtException", (error, origin) => {
//   if (
//     error instanceof Prisma.PrismaClientKnownRequestError ||
//     error instanceof Prisma.PrismaClientUnknownRequestError
//   ) {
//     // Don't exit the process if the error is a Prisma error
//     logger.error("uncaughtException prisma error", {
//       error,
//       prismaMessage: error.message,
//       code: "code" in error ? error.code : undefined,
//       meta: "meta" in error ? error.meta : undefined,
//       stack: error.stack,
//       origin,
//     });
//   } else {
//     logger.error("uncaughtException", {
//       error: { name: error.name, message: error.message, stack: error.stack },
//       origin,
//     });
//   }

//   process.exit(1);
// });

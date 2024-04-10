import { isRouteErrorResponse } from "react-router-dom";

// Coerce an unknown error into a string.
//   See [discussion](https://github.com/remix-run/react-router/discussions/9628#discussioncomment-5555901)
export function errorMessage(error: unknown): string {
  if (isRouteErrorResponse(error)) {
    return `${error.status} ${error.statusText}`;
  } else if (error instanceof Error) {
    return error.message;
  } else if (typeof error === "string") {
    return error;
  } else {
    console.error(error);
    return "Unknown error";
  }
}

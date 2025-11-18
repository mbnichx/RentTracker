import { Alert } from "react-native";

// Base URL for API requests. This is currently set to a local LAN IP used in
// development. Replace with an environment variable or production URL when
// building for release.
const BASE_URL = "http://10.0.0.67:8080"; // LAN IP

/**
 * Low-level helper for making HTTP requests to the backend.
 * It builds a fetch() call with JSON headers, serializes the body when
 * provided, and attempts to parse responses as JSON while falling back to
 * plain text when parsing fails.
 *
 * @param endpoint - Path appended to the BASE_URL (e.g. `/properties`)
 * @param method - HTTP method to use (GET, POST, PUT, DELETE, ...)
 * @param body - Optional payload object which will be JSON.stringified
 * @returns Parsed JSON response or raw text from the server
 * @throws Error when network fails or HTTP status is not ok
 */
async function apiRequest(
  endpoint: string,
  method: string = "GET",
  body?: object
) {
  const url = `${BASE_URL}${endpoint}`; // append endpoint to base URL
  console.log(url);

  const options: RequestInit = {
    method,
    headers: { "Content-Type": "application/json" },
  };

  // Attach JSON body when provided
  if (body) {
    options.body = JSON.stringify(body);
  }

  try {
    const res = await fetch(url, options);

    // Read response as text first so we can safely attempt JSON.parse
    const text = await res.text();
    if (!res.ok) {
      // Surface server-provided message when possible
      throw new Error(text || `HTTP ${res.status}`);
    }

    try {
      // If the response is JSON, return the parsed object
      return JSON.parse(text);
    } catch {
      // Otherwise return the plain text body
      return text;
    }
  } catch (err: any) {
    // Normalize and rethrow the error so callers can handle it consistently
    console.error("API request error:", err);
    throw new Error(err.message || "Network error");
  }
}

export default apiRequest;

/**
 * Convenience wrapper for calling API functions with built-in Alert UI for
 * success and a generic error alert on failure. This keeps UI code small when
 * performing actions from screens.
 *
 * Example: await safeCall(() => createLease(leasePayload), "Lease saved");
 *
 * @param fn - Asynchronous function that performs an API call
 * @param successMsg - Optional message shown in an Alert on success
 */
export const safeCall = async (fn: () => Promise<any>, successMsg?: string) => {
  try {
    const result = await fn();
    if (successMsg) Alert.alert("Success", successMsg);
    return result;
  } catch (err) {
    // Log the real error to console for debugging and present a generic
    // message to the user to avoid leaking internal details.
    console.error(err);
    Alert.alert("Error", "Something went wrong.");
  }
};

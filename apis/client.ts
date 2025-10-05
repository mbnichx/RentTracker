const BASE_URL = "http://10.0.0.67:8080"; // LAN IP

async function apiRequest(
  endpoint: string,
  method: string = "GET",
  body?: object
) {
  const url = `${BASE_URL}${endpoint}`; // append endpoint to base URL

  const options: RequestInit = {
    method,
    headers: { "Content-Type": "application/json" },
  };

  if (body) {
    options.body = JSON.stringify(body);
  }

  try {
    const res = await fetch(url, options);

    const text = await res.text(); // read as text first
    if (!res.ok) {
      throw new Error(text || `HTTP ${res.status}`);
    }

    try {
      return JSON.parse(text); // parse JSON only if valid
    } catch {
      return text; // fallback to plain text
    }
  } catch (err: any) {
    console.error("API request error:", err);
    throw new Error(err.message || "Network error");
  }
}

export default apiRequest;

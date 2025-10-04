const BASE_URL = "http://localhost:8080"; // change to your LAN IP when testing on mobile

async function apiRequest(endpoint: string, method: string = "GET", body?: object) {
  const options: RequestInit = {
    method,
    headers: { "Content-Type": "application/json" },
  };

  if (body) {
    options.body = JSON.stringify(body);
  }

  const res = await fetch(endpoint, options);
  const data = await res.json();
  return data;
}

export default apiRequest;

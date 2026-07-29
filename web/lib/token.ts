const STORAGE_KEY = "slw_token";

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(STORAGE_KEY);
}

export function setToken(token: string): void {
  window.localStorage.setItem(STORAGE_KEY, token);
}

export function clearToken(): void {
  if (typeof window === "undefined") return;
  try {
    window.localStorage.removeItem(STORAGE_KEY);
  } catch {
    // Blocked storage (extension/private mode) must not stop the 401 handler
    // in api.ts from still throwing UnauthorizedError and redirecting.
  }
}

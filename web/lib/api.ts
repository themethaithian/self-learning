import { clearToken, getToken } from "./token";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export class UnauthorizedError extends ApiError {
  constructor() {
    super(401, "unauthorized");
    this.name = "UnauthorizedError";
  }
}

export class NotFoundError extends ApiError {
  constructor() {
    super(404, "not found");
    this.name = "NotFoundError";
  }
}

export interface Concept {
  slug: string;
  title: string;
  position: number;
  has_lesson: boolean;
  est_minutes: number | null;
}

export interface Chapter {
  slug: string;
  title: string;
  position: number;
  concepts: Concept[];
}

export interface Topic {
  slug: string;
  title: string;
  position: number;
  chapters: Chapter[];
}

export interface Track {
  track: string;
  topics: Topic[];
}

export interface CurriculumResponse {
  tracks: Track[];
}

export interface Reference {
  title: string;
  source: string;
  why: string;
}

export interface RecallCheck {
  position: number;
  type: "short_answer" | "mcq";
  question: string;
  expected_answer: string;
  options?: string[];
}

export interface Lesson {
  topic: string;
  concept: string;
  title_en: string;
  est_minutes: number;
  body_md: string;
  references: Reference[];
  recall_checks: RecallCheck[];
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    ...(init?.headers as Record<string, string> | undefined),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...(init?.body ? { "Content-Type": "application/json" } : {}),
  };

  let res: Response;
  try {
    res = await fetch(`${API_BASE}${path}`, {
      ...init,
      headers,
      signal: init?.signal ?? AbortSignal.timeout(10_000),
    });
  } catch {
    throw new ApiError(0, "network request failed — is the API reachable?");
  }

  if (res.status === 401) {
    // Discard the rejected token here, not in each caller — otherwise a
    // stale token survives the redirect to /token and bounces straight
    // back to the page that redirected it, looping forever.
    clearToken();
    throw new UnauthorizedError();
  }
  if (res.status === 404) {
    throw new NotFoundError();
  }
  if (!res.ok) {
    const body = await res.text().catch(() => "");
    throw new ApiError(res.status, body || `request failed with status ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export function getCurriculum(): Promise<CurriculumResponse> {
  return apiFetch<CurriculumResponse>("/api/v1/curriculum");
}

export function getLesson(topicSlug: string, conceptSlug: string): Promise<Lesson> {
  return apiFetch<Lesson>(`/api/v1/lessons/${encodeURIComponent(topicSlug)}/${encodeURIComponent(conceptSlug)}`);
}

export interface FocusTrack {
  track: string | null;
}

// Focus track is an enhancement, not data the page needs to function — a
// failure here shouldn't turn a successful curriculum load into a full-page
// error. An expired token is the one failure the caller still needs to see.
export async function getFocusTrack(): Promise<FocusTrack> {
  try {
    return await apiFetch<FocusTrack>("/api/v1/prefs/focus-track");
  } catch (err) {
    if (err instanceof UnauthorizedError) throw err;
    return { track: null };
  }
}

export function setFocusTrack(track: string | null): Promise<FocusTrack> {
  return apiFetch<FocusTrack>("/api/v1/prefs/focus-track", {
    method: "PUT",
    body: JSON.stringify({ track }),
  });
}

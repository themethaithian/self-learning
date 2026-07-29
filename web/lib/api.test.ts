import { afterEach, describe, expect, it, vi } from "vitest";
import { getFocusTrack, getProgress, postAttempt, UnauthorizedError } from "./api";

function jsonResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
    text: async () => JSON.stringify(body),
  } as Response;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("getProgress", () => {
  it("reports ok with entries on a well-shaped response", async () => {
    const entry = { topic: "t", concept: "c", state: "passed", last_read_at: null, first_passed_at: null };
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(200, { concepts: [entry] })));
    expect(await getProgress()).toEqual({ kind: "ok", entries: [entry] });
  });

  it("reports error, not a false ok, when the body doesn't contain a concepts array (shape drift)", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(200, { entries: [] })));
    expect(await getProgress()).toEqual({ kind: "error" });
  });

  it("reports error, not a silent empty ok, on a 500", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(500, { error: "boom" })));
    expect(await getProgress()).toEqual({ kind: "error" });
  });

  it("propagates UnauthorizedError on 401 instead of swallowing it into {kind: error}", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(401, {})));
    await expect(getProgress()).rejects.toBeInstanceOf(UnauthorizedError);
  });
});

describe("postAttempt", () => {
  it("POSTs exactly the four documented keys — no kind field (server resolves it, and DisallowUnknownFields 400s on extras)", async () => {
    const record = {
      check_key: "k",
      question: "Q?",
      kind: "mcq",
      confidence: "confident",
      outcome: "correct",
      selected_option: "A",
      graded_by: "self",
      created_at: "2026-01-01T00:00:00Z",
    };
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(201, record));
    vi.stubGlobal("fetch", fetchMock);

    const result = await postAttempt("t1", "c1", {
      question: "Q?",
      confidence: "confident",
      outcome: "correct",
      selected_option: "A",
    });

    expect(result).toEqual({ kind: "ok", record });
    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toContain("/api/v1/progress/t1/c1/attempts");
    expect(init.method).toBe("POST");
    const body = JSON.parse(init.body as string);
    expect(Object.keys(body).sort()).toEqual(["confidence", "outcome", "question", "selected_option"]);
    expect(body).toEqual({ question: "Q?", confidence: "confident", outcome: "correct", selected_option: "A" });
  });

  it("sends selected_option: null verbatim for a short_answer attempt, not an omitted key", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(201, {
        check_key: "k",
        question: "Q?",
        kind: "short_answer",
        confidence: "guessed",
        outcome: "incorrect",
        selected_option: null,
        graded_by: "self",
        created_at: "2026-01-01T00:00:00Z",
      }),
    );
    vi.stubGlobal("fetch", fetchMock);

    await postAttempt("t1", "c1", { question: "Q?", confidence: "guessed", outcome: "incorrect", selected_option: null });

    const body = JSON.parse(fetchMock.mock.calls[0][1].body as string);
    expect(body.selected_option).toBeNull();
    expect("selected_option" in body).toBe(true);
  });

  it("reports {kind:'error'}, not a thrown exception, on a non-auth failure — the caller must keep the reader usable", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(400, { error: "invalid recall attempt" })));
    expect(await postAttempt("t1", "c1", { question: "Q", confidence: "guessed", outcome: "correct", selected_option: null })).toEqual({
      kind: "error",
    });
  });

  it("propagates UnauthorizedError on 401 instead of swallowing it into {kind: error}", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(401, {})));
    await expect(
      postAttempt("t1", "c1", { question: "Q", confidence: "guessed", outcome: "correct", selected_option: null }),
    ).rejects.toBeInstanceOf(UnauthorizedError);
  });

  it("percent-encodes topic and concept in the URL", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(201, {
        check_key: "k",
        question: "Q",
        kind: "mcq",
        confidence: "guessed",
        outcome: "correct",
        selected_option: null,
        graded_by: "self",
        created_at: "2026-01-01T00:00:00Z",
      }),
    );
    vi.stubGlobal("fetch", fetchMock);

    await postAttempt("a topic/slash", "a concept?", {
      question: "Q",
      confidence: "guessed",
      outcome: "correct",
      selected_option: null,
    });

    const [url] = fetchMock.mock.calls[0];
    expect(url).toContain(
      `/api/v1/progress/${encodeURIComponent("a topic/slash")}/${encodeURIComponent("a concept?")}/attempts`,
    );
  });
});

describe("getFocusTrack", () => {
  it("propagates UnauthorizedError on 401 instead of swallowing it into {track: null}", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(401, {})));
    await expect(getFocusTrack()).rejects.toBeInstanceOf(UnauthorizedError);
  });

  it("degrades a non-auth failure to {track: null}", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(500, {})));
    expect(await getFocusTrack()).toEqual({ track: null });
  });
});

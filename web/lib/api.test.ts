import { afterEach, describe, expect, it, vi } from "vitest";
import { getFocusTrack, getProgress, UnauthorizedError } from "./api";

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

// @vitest-environment jsdom
import { describe, expect, it, vi, beforeEach } from "vitest";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { UnauthorizedError, type Chapter, type Concept, type CurriculumResponse, type FocusTrack, type ProgressResult, type Topic } from "@/lib/api";

// vi.hoisted so these exist before vi.mock's factory runs — a plain `const`
// here would still hit a TDZ error once the file also imports a real (non
// type-only) value from "@/lib/api", since that import resolves the mocked
// module before this file's own top-level statements run.
const { routerMock, searchParamsMock, getCurriculum, getFocusTrack, getProgress } = vi.hoisted(() => ({
  routerMock: { replace: vi.fn(), push: vi.fn() },
  searchParamsMock: vi.fn(() => new URLSearchParams()),
  getCurriculum: vi.fn<() => Promise<CurriculumResponse>>(),
  getFocusTrack: vi.fn<() => Promise<FocusTrack>>(),
  getProgress: vi.fn<() => Promise<ProgressResult>>(),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => routerMock,
  useSearchParams: () => searchParamsMock(),
}));

vi.mock("@/lib/api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/lib/api")>();
  return { ...actual, getCurriculum, getFocusTrack, getProgress };
});

// Imported after the mocks above so the module under test picks them up.
const { default: LearnPage } = await import("./page");

function concept(slug: string, hasLesson: boolean, position: number): Concept {
  return { slug, title: `Concept ${slug}`, position, has_lesson: hasLesson, est_minutes: hasLesson ? 5 : null };
}
function chapter(slug: string, title: string, position: number, concepts: Concept[]): Chapter {
  return { slug, title, position, concepts };
}
function topic(slug: string, title: string, position: number, chapters: Chapter[]): Topic {
  return { slug, title, position, chapters };
}

// 2 tracks so a "which track did this stat belong to" mix-up can't hide
// behind there being only one candidate; ddd has 2 chapters mixing all
// three progress states plus a lesson-less concept, same as the other
// component-level fixtures.
const ddd = {
  track: "ddd",
  topics: [
    topic("t-ddd", "Domain-Driven Design", 1, [
      chapter("c1", "Model-Driven Foundations", 1, [
        concept("k1", true, 1),
        concept("k2", true, 2),
        concept("k3", true, 3),
      ]),
      chapter("c2", "Building Blocks", 2, [concept("k4", true, 1), concept("k5", false, 2)]),
    ]),
  ],
};
const go = {
  track: "go",
  topics: [topic("t-go", "Go", 1, [chapter("c-go", "Internals", 1, [concept("g1", true, 1)])])],
};

const curriculum: CurriculumResponse = { tracks: [ddd, go] };
const realProgress: ProgressResult = {
  kind: "ok",
  entries: [
    { topic: "t-ddd", concept: "k1", state: "passed", last_read_at: null, first_passed_at: null },
    { topic: "t-ddd", concept: "k2", state: "in_progress", last_read_at: null, first_passed_at: null },
  ],
};

beforeEach(() => {
  vi.clearAllMocks();
  searchParamsMock.mockReturnValue(new URLSearchParams());
  getCurriculum.mockResolvedValue(curriculum);
  getFocusTrack.mockResolvedValue({ track: null });
});

describe("LearnPage — list view", () => {
  it("shows real per-track read counts once progress loads", async () => {
    getProgress.mockResolvedValue(realProgress);
    render(<LearnPage />);

    await waitFor(() => expect(screen.getByText(/1\/4 lessons read/)).toBeTruthy());
    expect(screen.queryByRole("alert")).toBeNull();
    expect(screen.getAllByRole("progressbar")).toHaveLength(2);
  });

  it("never passes an empty progress map to a track card when real progress exists", async () => {
    getProgress.mockResolvedValue(realProgress);
    render(<LearnPage />);
    await waitFor(() => expect(screen.getAllByRole("progressbar").length).toBeGreaterThan(0));
    // If the wiring silently dropped the fetched entries, ddd would render
    // as 0/4 instead of 1/4.
    expect(screen.queryByText(/0\/4 lessons read/)).toBeNull();
  });

  it("shows the banner and hides every bar when the progress fetch fails", async () => {
    getProgress.mockResolvedValue({ kind: "error" });
    render(<LearnPage />);

    await waitFor(() => expect(screen.getByRole("alert")).toBeTruthy());
    expect(screen.queryAllByRole("progressbar")).toHaveLength(0);
    expect(screen.queryByText(/lessons read\b/)).toBeNull();
  });

  it("shows an empty state instead of a success view when the curriculum has no tracks", async () => {
    getCurriculum.mockResolvedValue({ tracks: [] });
    getProgress.mockResolvedValue({ kind: "ok", entries: [] });
    render(<LearnPage />);

    await waitFor(() => expect(screen.getByText(/No curriculum has been imported yet/)).toBeTruthy());
    expect(screen.queryAllByRole("progressbar")).toHaveLength(0);
    expect(routerMock.replace).not.toHaveBeenCalled();
  });
});

describe("LearnPage — auth failure", () => {
  it("redirects to /token when the curriculum fetch is unauthorized, not to any other route", async () => {
    getCurriculum.mockRejectedValue(new UnauthorizedError());
    getProgress.mockResolvedValue({ kind: "ok", entries: [] });
    render(<LearnPage />);

    await waitFor(() => expect(routerMock.replace).toHaveBeenCalledWith("/token"));
    expect(routerMock.replace).not.toHaveBeenCalledWith("/dashboard");
    expect(routerMock.replace).toHaveBeenCalledTimes(1);
  });

  it("shows the generic error state, not a redirect, for a non-401 curriculum failure", async () => {
    getCurriculum.mockRejectedValue(new Error("network request failed"));
    getProgress.mockResolvedValue({ kind: "ok", entries: [] });
    render(<LearnPage />);

    await waitFor(() => expect(screen.getByText(/Could not load the curriculum/)).toBeTruthy());
    expect(routerMock.replace).not.toHaveBeenCalled();
  });
});

describe("LearnPage — track detail view", () => {
  it("shows read-state markers on the detail view when progress loads successfully", async () => {
    searchParamsMock.mockReturnValue(new URLSearchParams("track=ddd"));
    getProgress.mockResolvedValue(realProgress);
    render(<LearnPage />);

    const button = await screen.findByRole("button", { name: /Model-Driven Foundations/ });
    fireEvent.click(button);
    expect(await screen.findByRole("img", { name: "Read" })).toBeTruthy();
    expect(screen.getByRole("img", { name: "In progress" })).toBeTruthy();
    expect(screen.queryByRole("alert")).toBeNull();
  });

  it("shows the banner and no markers on the detail view when the progress fetch fails", async () => {
    searchParamsMock.mockReturnValue(new URLSearchParams("track=ddd"));
    getProgress.mockResolvedValue({ kind: "error" });
    render(<LearnPage />);

    await waitFor(() => expect(screen.getByRole("alert")).toBeTruthy());
    const button = await screen.findByRole("button", { name: /Model-Driven Foundations/ });
    fireEvent.click(button);
    expect(screen.queryAllByRole("img")).toHaveLength(0);
    // The chapter header must fall back to "n/N ready", never a confident
    // "0/3 read" — that would happen if a failed fetch's null ever got
    // coalesced into an empty (but non-null, so "loaded") progress map
    // somewhere on the way down to TrackTopics.
    expect(within(button).getByText("3/3 ready")).toBeTruthy();
    expect(within(button).queryByText(/\bread\b/)).toBeNull();
  });
});

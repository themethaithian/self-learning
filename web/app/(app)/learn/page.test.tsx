// @vitest-environment jsdom
import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, waitFor, within } from "@testing-library/react";
import type { Chapter, Concept, CurriculumResponse, FocusTrack, ProgressResult, Topic } from "@/lib/api";

const routerMock = { replace: vi.fn(), push: vi.fn() };
const searchParamsMock = vi.fn(() => new URLSearchParams());

vi.mock("next/navigation", () => ({
  useRouter: () => routerMock,
  useSearchParams: () => searchParamsMock(),
}));

const getCurriculum = vi.fn<() => Promise<CurriculumResponse>>();
const getFocusTrack = vi.fn<() => Promise<FocusTrack>>();
const getProgress = vi.fn<() => Promise<ProgressResult>>();

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
});

describe("LearnPage — track detail view", () => {
  it("shows read-state markers on the detail view when progress loads successfully", async () => {
    searchParamsMock.mockReturnValue(new URLSearchParams("track=ddd"));
    getProgress.mockResolvedValue(realProgress);
    render(<LearnPage />);

    const button = await screen.findByRole("button", { name: /Model-Driven Foundations/ });
    button.click();
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
    button.click();
    expect(screen.queryAllByRole("img")).toHaveLength(0);
    // The chapter header must fall back to "n/N ready", never a confident
    // "0/3 read" — that would happen if a failed fetch's null ever got
    // coalesced into an empty (but non-null, so "loaded") progress map
    // somewhere on the way down to TrackTopics.
    expect(within(button).getByText("3/3 ready")).toBeTruthy();
    expect(within(button).queryByText(/\bread\b/)).toBeNull();
  });
});

// @vitest-environment jsdom
import { describe, expect, it, vi, beforeEach } from "vitest";
import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import type { CurriculumResponse, Lesson, ProgressEntry } from "@/lib/api";

const { routerMock, searchParamsMock, getLesson, getCurriculum, setProgress } = vi.hoisted(() => ({
  routerMock: { replace: vi.fn(), push: vi.fn() },
  searchParamsMock: vi.fn(() => new URLSearchParams()),
  getLesson: vi.fn<(topic: string, concept: string) => Promise<Lesson>>(),
  getCurriculum: vi.fn<() => Promise<CurriculumResponse>>(),
  setProgress: vi.fn<(topic: string, concept: string, state: "in_progress" | "passed") => Promise<ProgressEntry>>(),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => routerMock,
  useSearchParams: () => searchParamsMock(),
}));

vi.mock("@/lib/api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/lib/api")>();
  return { ...actual, getLesson, getCurriculum, setProgress };
});

const { default: LessonPage } = await import("./page");

function lesson(topic: string, concept: string): Lesson {
  return {
    topic,
    concept,
    title_en: `Lesson ${concept}`,
    est_minutes: 5,
    body_md: "Body text.",
    references: [],
    recall_checks: [
      {
        position: 0,
        type: "mcq",
        question: `Question for ${concept}`,
        expected_answer: `Answer for ${concept}`,
        options: [`Wrong for ${concept}`, `Answer for ${concept}`, `Other for ${concept}`],
      },
    ],
  };
}

const lessonA = lesson("t1", "c1");
const lessonB = lesson("t1", "c2");

beforeEach(() => {
  vi.clearAllMocks();
  getCurriculum.mockResolvedValue({ tracks: [] });
  setProgress.mockResolvedValue({ topic: "t1", concept: "c1", state: "in_progress", last_read_at: null, first_passed_at: null });
  getLesson.mockImplementation((topic: string, concept: string) =>
    Promise.resolve(concept === "c1" ? lessonA : lessonB),
  );
});

function setParams(topic: string, concept: string) {
  searchParamsMock.mockReturnValue(new URLSearchParams({ topic, concept }));
}

describe("LessonPage — per-check state resets across a query-param-only lesson change", () => {
  it("does not leak the previous lesson's revealed answer or finished state into the next lesson's same-position check", async () => {
    setParams("t1", "c1");
    const { rerender } = render(<LessonPage />);

    await screen.findByText("Question for c1");
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getByRole("radio", { name: "Answer for c1" }));
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    await waitFor(() => expect(screen.getByRole("status").textContent).toContain("Correct"));
    expect(screen.getByRole("button", { name: "Finish lesson" }).getAttribute("aria-disabled")).toBe("false");

    // Query-param-only navigation: same page component, new searchParams —
    // LessonPage does not remount, matching real Next.js App Router behaviour.
    setParams("t1", "c2");
    rerender(<LessonPage />);

    await screen.findByText("Question for c2");
    expect(screen.queryByText("Answer for c1")).toBeNull();
    expect(screen.queryByRole("status")).toBeNull();
    expect(screen.getByRole("button", { name: "Show options" })).toBeTruthy();
    expect(screen.getByRole("button", { name: "Finish lesson" }).getAttribute("aria-disabled")).toBe("true");
  });
});

describe("LessonPage — finish() lesson-identity guard (pre-existing since UX-5, pinned here)", () => {
  it("ignores a stale Finish PUT response that resolves after the user has already navigated to a different lesson", async () => {
    let resolvePassedPut: (entry: ProgressEntry) => void = () => {};
    const passedPutPromise = new Promise<ProgressEntry>((resolve) => {
      resolvePassedPut = resolve;
    });
    setProgress.mockImplementation((topic: string, concept: string, state: string) =>
      state === "passed"
        ? passedPutPromise
        : Promise.resolve({ topic, concept, state: "in_progress", last_read_at: null, first_passed_at: null }),
    );

    setParams("t1", "c1");
    const { rerender } = render(<LessonPage />);
    await screen.findByText("Question for c1");
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getByRole("radio", { name: "Answer for c1" }));
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    await waitFor(() =>
      expect(screen.getByRole("button", { name: "Finish lesson" }).getAttribute("aria-disabled")).toBe("false"),
    );

    fireEvent.click(screen.getByRole("button", { name: "Finish lesson" }));
    await screen.findByRole("button", { name: "Saving…" });

    // Navigate to a different lesson while lesson A's PUT is still in flight —
    // Next is visible (and clickable) before Finish resolves, by design.
    setParams("t1", "c2");
    rerender(<LessonPage />);
    await screen.findByText("Question for c2");

    // The stale PUT for lesson A finally comes back.
    await act(async () => {
      resolvePassedPut({ topic: "t1", concept: "c1", state: "passed", last_read_at: null, first_passed_at: null });
      await passedPutPromise;
      await Promise.resolve();
      await Promise.resolve();
    });

    // Lesson B must not be marked finished by a response meant for lesson A.
    expect(screen.queryByText("Lesson finished")).toBeNull();
    expect(screen.getByRole("button", { name: "Finish lesson" })).toBeTruthy();
  });
});

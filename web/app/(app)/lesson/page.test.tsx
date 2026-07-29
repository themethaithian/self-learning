// @vitest-environment jsdom
import { describe, expect, it, vi, beforeEach } from "vitest";
import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { UnauthorizedError, type AttemptInput, type CurriculumResponse, type Lesson, type PostAttemptResult, type ProgressEntry } from "@/lib/api";

const { routerMock, searchParamsMock, getLesson, getCurriculum, setProgress, postAttempt } = vi.hoisted(() => ({
  routerMock: { replace: vi.fn(), push: vi.fn() },
  searchParamsMock: vi.fn(() => new URLSearchParams()),
  getLesson: vi.fn<(topic: string, concept: string) => Promise<Lesson>>(),
  getCurriculum: vi.fn<() => Promise<CurriculumResponse>>(),
  setProgress: vi.fn<(topic: string, concept: string, state: "in_progress" | "passed") => Promise<ProgressEntry>>(),
  postAttempt: vi.fn<(topic: string, concept: string, attempt: AttemptInput) => Promise<PostAttemptResult>>(),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => routerMock,
  useSearchParams: () => searchParamsMock(),
}));

vi.mock("@/lib/api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/lib/api")>();
  return { ...actual, getLesson, getCurriculum, setProgress, postAttempt };
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

function shortAnswerLesson(topic: string, concept: string): Lesson {
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
        type: "short_answer",
        question: `Short question for ${concept}`,
        expected_answer: `Short answer for ${concept}`,
      },
    ],
  };
}

function okAttemptResult(): PostAttemptResult {
  return {
    kind: "ok",
    record: {
      check_key: "k",
      question: "q",
      kind: "mcq",
      confidence: "confident",
      outcome: "correct",
      selected_option: null,
      graded_by: "self",
      created_at: "2026-01-01T00:00:00Z",
    },
  };
}

beforeEach(() => {
  vi.clearAllMocks();
  getCurriculum.mockResolvedValue({ tracks: [] });
  setProgress.mockResolvedValue({ topic: "t1", concept: "c1", state: "in_progress", last_read_at: null, first_passed_at: null });
  getLesson.mockImplementation((topic: string, concept: string) =>
    Promise.resolve(concept === "c1" ? lessonA : lessonB),
  );
  postAttempt.mockResolvedValue(okAttemptResult());
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

async function completeMcq(correct: boolean) {
  await screen.findByText("Question for c1");
  fireEvent.click(screen.getByRole("button", { name: "Show options" }));
  fireEvent.click(screen.getByRole("radio", { name: correct ? "Answer for c1" : "Wrong for c1" }));
  fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
  fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
}

describe("LessonPage — attempt submission (Q-2b)", () => {
  it("POSTs the exact body (question/confidence/outcome/selected_option) exactly once when an mcq check completes correctly", async () => {
    setParams("t1", "c1");
    render(<LessonPage />);
    await completeMcq(true);

    await waitFor(() => expect(postAttempt).toHaveBeenCalledTimes(1));
    expect(postAttempt).toHaveBeenCalledWith("t1", "c1", {
      question: "Question for c1",
      confidence: "confident",
      outcome: "correct",
      selected_option: "Answer for c1",
    });
  });

  it("POSTs outcome=incorrect with the picked (wrong) option text, not a hardcoded outcome", async () => {
    setParams("t1", "c1");
    render(<LessonPage />);
    await completeMcq(false);

    await waitFor(() => expect(postAttempt).toHaveBeenCalledTimes(1));
    expect(postAttempt).toHaveBeenCalledWith(
      "t1",
      "c1",
      expect.objectContaining({ outcome: "incorrect", selected_option: "Wrong for c1" }),
    );
  });

  it("does not submit at stage 2 (commit) — only once Reveal is reached", async () => {
    setParams("t1", "c1");
    render(<LessonPage />);
    await screen.findByText("Question for c1");
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getByRole("radio", { name: "Answer for c1" }));
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));

    await new Promise((r) => setTimeout(r, 0));
    expect(postAttempt).not.toHaveBeenCalled();
  });

  it.each(["Guessed", "Unsure", "Confident"] as const)(
    "propagates confidence='%s' verbatim in the request body",
    async (label) => {
      setParams("t1", "c1");
      render(<LessonPage />);
      await screen.findByText("Question for c1");
      fireEvent.click(screen.getByRole("button", { name: "Show options" }));
      fireEvent.click(screen.getByRole("radio", { name: "Answer for c1" }));
      fireEvent.click(screen.getByRole("radio", { name: label }));
      fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

      await waitFor(() =>
        expect(postAttempt).toHaveBeenCalledWith(
          "t1",
          "c1",
          expect.objectContaining({ confidence: label.toLowerCase() }),
        ),
      );
    },
  );

  it("short_answer: submits selected_option=null, and re-submits as a correction when the Pass/Not yet rating changes", async () => {
    getLesson.mockImplementation(() => Promise.resolve(shortAnswerLesson("t1", "sa1")));
    setParams("t1", "sa1");
    render(<LessonPage />);
    await screen.findByText("Short question for sa1");
    fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
    fireEvent.click(screen.getByRole("radio", { name: "Guessed" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    fireEvent.click(screen.getByRole("button", { name: "Pass" }));
    await waitFor(() => expect(postAttempt).toHaveBeenCalledTimes(1));
    expect(postAttempt).toHaveBeenLastCalledWith("t1", "sa1", {
      question: "Short question for sa1",
      confidence: "guessed",
      outcome: "correct",
      selected_option: null,
    });

    fireEvent.click(screen.getByRole("button", { name: "Not yet" }));
    await waitFor(() => expect(postAttempt).toHaveBeenCalledTimes(2));
    expect(postAttempt).toHaveBeenLastCalledWith("t1", "sa1", {
      question: "Short question for sa1",
      confidence: "guessed",
      outcome: "incorrect",
      selected_option: null,
    });
  });

  it("shows a 'Not saved' indicator (never a false success) on a failed save, then clears it once Retry succeeds", async () => {
    postAttempt.mockResolvedValueOnce({ kind: "error" });
    setParams("t1", "c1");
    render(<LessonPage />);
    await completeMcq(true);

    await screen.findByText("Not saved");

    // A deferred promise (not mockResolvedValueOnce) so the assertions below
    // run strictly after the retry has actually settled — a bare `waitFor`
    // on "Not saved" disappearing would also pass transiently during the
    // "saving" state in between, before the retry's own result lands, and
    // that false pass would hide a regression where the retry succeeds but
    // the indicator wrongly stays (or the Retry affordance never clears).
    let resolveRetry: (result: PostAttemptResult) => void = () => {};
    const retryPromise = new Promise<PostAttemptResult>((resolve) => {
      resolveRetry = resolve;
    });
    postAttempt.mockReturnValueOnce(retryPromise);
    fireEvent.click(screen.getByRole("button", { name: "Retry" }));
    await waitFor(() => expect(postAttempt).toHaveBeenCalledTimes(2));

    await act(async () => {
      resolveRetry(okAttemptResult());
      await retryPromise;
    });

    expect(screen.queryByText("Not saved")).toBeNull();
    expect(screen.queryByRole("button", { name: "Retry" })).toBeNull();
  });

  it("surfaces unsaved attempts near Finish instead of silently implying everything recorded", async () => {
    postAttempt.mockResolvedValueOnce({ kind: "error" });
    setParams("t1", "c1");
    render(<LessonPage />);
    await completeMcq(true);

    await screen.findByText(/recall attempts didn.t save/);
  });

  it("redirects to /token on a 401 from the attempts endpoint, matching useFocusTrack's policy", async () => {
    postAttempt.mockRejectedValueOnce(new UnauthorizedError());
    setParams("t1", "c1");
    render(<LessonPage />);
    await completeMcq(true);

    await waitFor(() => expect(routerMock.replace).toHaveBeenCalledWith("/token"));
  });

  it("ignores a stale attempt response for a lesson the user has already navigated away from", async () => {
    let resolveAttempt: (result: PostAttemptResult) => void = () => {};
    const attemptPromise = new Promise<PostAttemptResult>((resolve) => {
      resolveAttempt = resolve;
    });
    postAttempt.mockReturnValueOnce(attemptPromise);

    setParams("t1", "c1");
    const { rerender } = render(<LessonPage />);
    await completeMcq(true);
    await waitFor(() => expect(postAttempt).toHaveBeenCalledTimes(1));

    // Navigate to a different lesson while lesson A's attempt POST is still
    // in flight, same as the finish() guard test above.
    setParams("t1", "c2");
    rerender(<LessonPage />);
    await screen.findByText("Question for c2");

    await act(async () => {
      resolveAttempt({ kind: "error" });
      await attemptPromise;
      await Promise.resolve();
      await Promise.resolve();
    });

    // Lesson B's own (untouched) check is still at stage 1 (recall) — its
    // per-check "Not saved" indicator only renders at stage 3, so asserting
    // on that text alone would pass even with the guard deleted. The
    // aggregate Finish-step banner is keyed by lesson-view-level state
    // (attemptStatus), not by a child's stage, so it is the assertion that
    // actually observes whether lesson A's stale response leaked in.
    expect(screen.queryByText(/recall attempts didn.t save/)).toBeNull();
    expect(screen.queryByText("Not saved")).toBeNull();
  });
});

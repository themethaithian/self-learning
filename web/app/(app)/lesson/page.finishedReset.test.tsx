// @vitest-environment jsdom
//
// RecallCheckCard's own mount effect immediately reports finished=false for a
// fresh "recall"-stage card, which means lesson/page.tsx's explicit
// `setFinishedChecks({})` on lesson change is defense-in-depth today (see
// docs/tickets/quiz.md) — the real component always self-corrects a split
// second later anyway. That self-correction would mask a regression in the
// parent's own reset line. This file replaces RecallCheckCard with a fake
// that does NOT self-report on mount, so the only thing that can clear a
// stale `finishedChecks` entry across a lesson change is the parent's own
// reset — isolating exactly what that line is responsible for.
import { forwardRef } from "react";
import { describe, expect, it, vi, beforeEach } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import type { CurriculumResponse, Lesson, ProgressEntry, RecallCheck } from "@/lib/api";

const { routerMock, searchParamsMock, getLesson, getCurriculum, setProgress } = vi.hoisted(() => ({
  routerMock: { replace: vi.fn(), push: vi.fn() },
  searchParamsMock: vi.fn(() => new URLSearchParams()),
  getLesson: vi.fn<(topic: string, concept: string) => Promise<Lesson>>(),
  getCurriculum: vi.fn<() => Promise<CurriculumResponse>>(),
  setProgress: vi.fn<() => Promise<ProgressEntry>>(),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => routerMock,
  useSearchParams: () => searchParamsMock(),
}));

vi.mock("@/lib/api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/lib/api")>();
  return { ...actual, getLesson, getCurriculum, setProgress };
});

vi.mock("@/components/RecallCheckCard", () => ({
  RecallCheckCard: forwardRef<HTMLDivElement, { check: RecallCheck; onFinishedChange: (p: number, f: boolean) => void }>(
    function FakeRecallCheckCard({ check, onFinishedChange }, ref) {
      return (
        <div ref={ref} tabIndex={-1}>
          <span>{check.question}</span>
          <button onClick={() => onFinishedChange(check.position, true)}>Mark finished {check.position}</button>
        </div>
      );
    },
  ),
}));

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
      { position: 0, type: "short_answer", question: `Question for ${concept}`, expected_answer: `Answer for ${concept}` },
    ],
  };
}

const lessonA = lesson("t1", "c1");
const lessonB = lesson("t1", "c2");

beforeEach(() => {
  vi.clearAllMocks();
  getCurriculum.mockResolvedValue({ tracks: [] });
  setProgress.mockResolvedValue({ topic: "t1", concept: "c1", state: "in_progress", last_read_at: null, first_passed_at: null });
  getLesson.mockImplementation((topic: string, concept: string) => Promise.resolve(concept === "c1" ? lessonA : lessonB));
});

function setParams(topic: string, concept: string) {
  searchParamsMock.mockReturnValue(new URLSearchParams({ topic, concept }));
}

describe("LessonPage — finishedChecks resets on lesson change even if the card never self-corrects", () => {
  it("re-gates Finish for the new lesson instead of trusting a stale finished=true from the previous one", async () => {
    setParams("t1", "c1");
    const { rerender } = render(<LessonPage />);
    await screen.findByText("Question for c1");

    fireEvent.click(screen.getByRole("button", { name: "Mark finished 0" }));
    expect(screen.getByRole("button", { name: "Finish lesson" }).getAttribute("aria-disabled")).toBe("false");

    setParams("t1", "c2");
    rerender(<LessonPage />);
    await screen.findByText("Question for c2");

    // The fake card never reports finished=false on its own — if the parent
    // did not clear finishedChecks itself, position 0 would still read
    // `true` from lesson A and Finish would be wrongly enabled here.
    expect(screen.getByRole("button", { name: "Finish lesson" }).getAttribute("aria-disabled")).toBe("true");
  });
});

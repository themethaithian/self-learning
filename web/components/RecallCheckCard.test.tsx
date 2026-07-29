// @vitest-environment jsdom
import { describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import { RecallCheckCard, isRecallCheckFinished } from "./RecallCheckCard";
import type { RecallCheck } from "@/lib/api";

const zeroRng = () => 0;

const mcqCheck: RecallCheck = {
  position: 0,
  type: "mcq",
  question: "Which one is B?",
  expected_answer: "Option B",
  options: ["Option A", "Option B", "Option C"],
};

const shortAnswerCheck: RecallCheck = {
  position: 1,
  type: "short_answer",
  question: "What is CAP?",
  expected_answer: "Consistency, Availability, Partition tolerance",
};

function renderCard(check: RecallCheck, opts?: { rng?: () => number; onFinishedChange?: (p: number, f: boolean) => void }) {
  const onFinishedChange = opts?.onFinishedChange ?? vi.fn();
  const utils = render(<RecallCheckCard check={check} index={0} onFinishedChange={onFinishedChange} rng={opts?.rng} />);
  return { ...utils, onFinishedChange };
}

describe("isRecallCheckFinished — the one place 'finished' is defined", () => {
  it.each([
    ["mcq", "recall", null, false],
    ["mcq", "commit", null, false],
    ["mcq", "reveal", null, true],
    ["short_answer", "recall", null, false],
    ["short_answer", "commit", "pass", false],
    ["short_answer", "reveal", null, false],
    ["short_answer", "reveal", "pass", true],
    ["short_answer", "reveal", "fail", true],
  ] as const)("type=%s stage=%s rating=%s -> %s", (type, stage, rating, expected) => {
    expect(isRecallCheckFinished(type, stage, rating)).toBe(expected);
  });
});

describe("RecallCheckCard — stage 1 (recall)", () => {
  it("shows only the question, no options and no expected_answer, for mcq", () => {
    const { container } = renderCard(mcqCheck);
    expect(screen.getByText(mcqCheck.question)).toBeTruthy();
    expect(container.innerHTML).not.toContain("Option A");
    expect(container.innerHTML).not.toContain(mcqCheck.expected_answer);
    expect(screen.getByRole("button", { name: "Show options" })).toBeTruthy();
  });

  it("labels the stage-1 control 'I've answered' for short_answer, not 'Show options'", () => {
    renderCard(shortAnswerCheck);
    expect(screen.getByRole("button", { name: "I've answered" })).toBeTruthy();
    expect(screen.queryByRole("button", { name: "Show options" })).toBeNull();
  });

  it("never puts expected_answer in the serialized DOM before reveal", () => {
    const { container } = renderCard(shortAnswerCheck);
    expect(container.innerHTML).not.toContain(shortAnswerCheck.expected_answer);
    fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
    expect(container.innerHTML).not.toContain(shortAnswerCheck.expected_answer);
  });
});

describe("RecallCheckCard — stage 2 (commit) gating", () => {
  it("uses real radio semantics for options and confidence (fieldset/legend + input radio)", () => {
    renderCard(mcqCheck);
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    expect(screen.getByText("Choose one").tagName).toBe("LEGEND");
    expect(screen.getByText("How confident are you?").tagName).toBe("LEGEND");
    const radios = screen.getAllByRole("radio");
    // 3 options + 3 confidence levels.
    expect(radios).toHaveLength(6);
  });

  it("does not reach reveal (mcq) without an option selected, even with confidence set", () => {
    // Note: mcqCheck's correct option's own label text equals expected_answer
    // ("Option B"), so asserting DOM-innerHTML-excludes-expected_answer would
    // false-fail here (the option list legitimately shows that string at
    // stage 2) — role=status only exists once stage 3 mounts, so its absence
    // is the unambiguous "still on stage 2" signal.
    renderCard(mcqCheck);
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getByRole("radio", { name: "Guessed" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(screen.queryByRole("status")).toBeNull();
    expect(screen.getByText("Choose one")).toBeTruthy();
  });

  it("does not reach reveal (mcq) without a confidence, even with an option selected", () => {
    renderCard(mcqCheck);
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(screen.queryByRole("status")).toBeNull();
  });

  it("does not reach reveal (short_answer) without a confidence", () => {
    const { container } = renderCard(shortAnswerCheck);
    fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(container.innerHTML).not.toContain(shortAnswerCheck.expected_answer);
  });

  it("reaches reveal once both an option and a confidence are committed", () => {
    const { container } = renderCard(mcqCheck, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(container.innerHTML).toContain(mcqCheck.expected_answer);
  });

  it("moves focus into the commit region when leaving stage 1", () => {
    renderCard(mcqCheck);
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    expect(document.activeElement).toBe(screen.getByTestId("stage-commit"));
  });
});

describe("RecallCheckCard — the shuffle", () => {
  it("is stable across re-renders (lazy useState, not recomputed every render)", () => {
    const rng = vi.fn(zeroRng);
    const onFinishedChange = vi.fn();
    const { rerender } = render(
      <RecallCheckCard check={mcqCheck} index={0} onFinishedChange={onFinishedChange} rng={rng} />,
    );
    const callsAfterMount = rng.mock.calls.length;
    expect(callsAfterMount).toBeGreaterThan(0);

    rerender(<RecallCheckCard check={mcqCheck} index={0} onFinishedChange={onFinishedChange} rng={rng} />);
    rerender(<RecallCheckCard check={mcqCheck} index={0} onFinishedChange={onFinishedChange} rng={rng} />);

    expect(rng.mock.calls.length).toBe(callsAfterMount);
  });

  it("actually reorders the source array (identity shuffle would render source order)", () => {
    renderCard(mcqCheck, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    // Fisher-Yates with rng()=0 on [A,B,C] yields [B,C,A] — verified directly
    // in shuffle.test.ts; asserting it here proves the component actually
    // calls shuffleOptions instead of passing options straight through.
    const labels = screen.getAllByRole("radio", { name: /^Option/ }).map((r) => r.closest("label")?.textContent);
    expect(labels).toEqual(["Option B", "Option C", "Option A"]);
  });

  it("marks the correct option by matching expected_answer, not by its original index", () => {
    // Source order deliberately puts the correct answer NOT first.
    const check: RecallCheck = {
      position: 2,
      type: "mcq",
      question: "Q",
      expected_answer: "Correct",
      options: ["Wrong1", "Correct", "Wrong2"],
    };
    // rng()=0 shuffles [Wrong1, Correct, Wrong2] -> [Correct, Wrong2, Wrong1]
    // (same transform as the A,B,C case above). The correct answer's source
    // index was 1; if matching were done by that original index instead of
    // by text, position 1 (Wrong2) would be mislabelled correct instead.
    const { container } = renderCard(check, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[1]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    const items = screen.getAllByRole("listitem");
    expect(items[0].textContent).toContain("Correct");
    expect(items[0].textContent).toContain("Correct answer");
    expect(items[1].textContent).not.toContain("Correct answer");
    void container;
  });

  it("keys options by shuffled position, not text — selecting one duplicate never checks its twin", () => {
    const check: RecallCheck = {
      position: 3,
      type: "mcq",
      question: "Q",
      expected_answer: "Different",
      options: ["Same text", "Same text", "Different"],
    };
    // rng()=0 shuffles [Same, Same, Different] -> [Same(idx1), Different, Same(idx0)]
    renderCard(check, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    const radios = screen.getAllByRole("radio").filter((r) => r.getAttribute("name")?.startsWith("mcq-"));
    expect(radios).toHaveLength(3);

    fireEvent.click(radios[2]);
    expect((radios[2] as HTMLInputElement).checked).toBe(true);
    expect((radios[0] as HTMLInputElement).checked).toBe(false);
  });
});

describe("RecallCheckCard — reveal stage (mcq)", () => {
  it("derives correctness instead of trusting a self-report, and announces it via role=status", () => {
    renderCard(mcqCheck, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    // Shuffled order is [B, C, A]; picking index 0 is the correct option.
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Guessed" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    const region = screen.getByRole("status");
    expect(region.textContent).toContain("Correct");
    expect(document.activeElement).toBe(region);
  });

  it("marks an incorrect pick as incorrect without colour as the only signal", () => {
    renderCard(mcqCheck, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    // Shuffled order is [B, C, A]; index 2 (A) is wrong.
    fireEvent.click(screen.getAllByRole("radio")[2]);
    fireEvent.click(screen.getByRole("radio", { name: "Unsure" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    const region = screen.getByRole("status");
    expect(region.textContent).toContain("Not quite");
    expect(screen.getByText(/Your answer — incorrect/)).toBeTruthy();
  });
});

describe("RecallCheckCard — finished callback", () => {
  it("reports mcq finished only once stage 3 is reached", () => {
    const onFinishedChange = vi.fn();
    renderCard(mcqCheck, { rng: zeroRng, onFinishedChange });
    expect(onFinishedChange).toHaveBeenLastCalledWith(0, false);

    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    expect(onFinishedChange).toHaveBeenLastCalledWith(0, false);

    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    expect(onFinishedChange).toHaveBeenLastCalledWith(0, false);

    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(onFinishedChange).toHaveBeenLastCalledWith(0, true);
  });

  it("does not report short_answer finished at stage 3 until Pass/Not yet is chosen", () => {
    const onFinishedChange = vi.fn();
    renderCard(shortAnswerCheck, { onFinishedChange });

    fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
    fireEvent.click(screen.getByRole("radio", { name: "Unsure" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(onFinishedChange).toHaveBeenLastCalledWith(1, false);

    fireEvent.click(screen.getByRole("button", { name: "Pass" }));
    expect(onFinishedChange).toHaveBeenLastCalledWith(1, true);
  });
});

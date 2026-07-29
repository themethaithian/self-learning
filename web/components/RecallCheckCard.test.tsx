// @vitest-environment jsdom
import { StrictMode, useState } from "react";
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

interface CompletedAttempt {
  confidence: "guessed" | "unsure" | "confident";
  outcome: "correct" | "incorrect";
  selectedOption: string | null;
}

function renderCard(
  check: RecallCheck,
  opts?: {
    rng?: () => number;
    onFinishedChange?: (p: number, f: boolean) => void;
    onAttemptReady?: (p: number, a: CompletedAttempt) => void;
    saveStatus?: "idle" | "saving" | "saved" | "error";
    onRetrySave?: () => void;
  },
) {
  const onFinishedChange = opts?.onFinishedChange ?? vi.fn();
  const onAttemptReady = opts?.onAttemptReady ?? vi.fn();
  const onRetrySave = opts?.onRetrySave ?? vi.fn();
  const utils = render(
    <RecallCheckCard
      check={check}
      index={0}
      onFinishedChange={onFinishedChange}
      onAttemptReady={onAttemptReady}
      saveStatus={opts?.saveStatus ?? "idle"}
      onRetrySave={onRetrySave}
      rng={opts?.rng}
    />,
  );
  return { ...utils, onFinishedChange, onAttemptReady, onRetrySave };
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

describe("RecallCheckCard — confidence selection", () => {
  it.each(["Guessed", "Unsure", "Confident"] as const)(
    "clicking '%s' checks exactly that radio, not one of the other two",
    (clicked) => {
      renderCard(shortAnswerCheck);
      fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
      fireEvent.click(screen.getByRole("radio", { name: clicked }));

      for (const label of ["Guessed", "Unsure", "Confident"] as const) {
        const radio = screen.getByRole("radio", { name: label }) as HTMLInputElement;
        expect(radio.checked).toBe(label === clicked);
      }
    },
  );
});

describe("RecallCheckCard — the shuffle", () => {
  it("is stable across re-renders (lazy useState, not recomputed every render)", () => {
    const rng = vi.fn(zeroRng);
    const onFinishedChange = vi.fn();
    const { rerender } = render(
      <RecallCheckCard
        check={mcqCheck}
        index={0}
        onFinishedChange={onFinishedChange}
        onAttemptReady={vi.fn()}
        saveStatus="idle"
        onRetrySave={vi.fn()}
        rng={rng}
      />,
    );
    const callsAfterMount = rng.mock.calls.length;
    expect(callsAfterMount).toBeGreaterThan(0);

    rerender(
      <RecallCheckCard
        check={mcqCheck}
        index={0}
        onFinishedChange={onFinishedChange}
        onAttemptReady={vi.fn()}
        saveStatus="idle"
        onRetrySave={vi.fn()}
        rng={rng}
      />,
    );
    rerender(
      <RecallCheckCard
        check={mcqCheck}
        index={0}
        onFinishedChange={onFinishedChange}
        onAttemptReady={vi.fn()}
        saveStatus="idle"
        onRetrySave={vi.fn()}
        rng={rng}
      />,
    );

    expect(rng.mock.calls.length).toBe(callsAfterMount);
  });

  it("uses real randomness on the production path (no rng prop) — repeated mounts are not pinned to one order", () => {
    // shuffle.test.ts proves shuffleOptions() is a good shuffle in isolation,
    // and every other test here injects a fixed rng to prove the plumbing —
    // neither exercises the one call the component itself makes without a
    // caller-supplied rng, which is the only path production traffic takes.
    const orders = new Set<string>();
    for (let i = 0; i < 40; i++) {
      const { unmount } = render(
        <RecallCheckCard
          check={mcqCheck}
          index={0}
          onFinishedChange={vi.fn()}
          onAttemptReady={vi.fn()}
          saveStatus="idle"
          onRetrySave={vi.fn()}
        />,
      );
      fireEvent.click(screen.getByRole("button", { name: "Show options" }));
      const order = screen
        .getAllByRole("radio", { name: /^Option/ })
        .map((r) => r.closest("label")?.textContent)
        .join("|");
      orders.add(order);
      unmount();
    }
    expect(orders.size).toBeGreaterThan(1);
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
    renderCard(check, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[1]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    const items = screen.getAllByRole("listitem");
    expect(items[0].textContent).toContain("Correct");
    expect(items[0].textContent).toContain("Correct answer");
    expect(items[1].textContent).not.toContain("Correct answer");
  });

  it("tracks selection by index, not option text — a duplicate-text option's twin is never checked in stage 2 or marked 'your answer' in stage 3", () => {
    const check: RecallCheck = {
      position: 3,
      type: "mcq",
      question: "Q",
      expected_answer: "Different",
      options: ["Same text", "Same text", "Different"],
    };
    // rng()=0 shuffles [Same(idx0), Same(idx1), Different(idx2)] into
    // [Same(idx1), Different(idx2), Same(idx0)].
    renderCard(check, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    const radios = screen.getAllByRole("radio").filter((r) => r.getAttribute("name")?.startsWith("mcq-"));
    expect(radios).toHaveLength(3);

    fireEvent.click(radios[2]);
    expect((radios[2] as HTMLInputElement).checked).toBe(true);
    expect((radios[0] as HTMLInputElement).checked).toBe(false);

    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    const items = screen.getAllByRole("listitem");
    // Text-based "your answer" matching would mark BOTH "Same text" rows —
    // only the one actually clicked (shuffled position 2) should carry the tag.
    expect(items[0].textContent).not.toContain("Your answer");
    expect(items[2].textContent).toContain("Your answer — incorrect");
  });
});

describe("RecallCheckCard — reveal stage (mcq)", () => {
  it("derives correctness instead of trusting a self-report, and announces the result via a scoped role=status", () => {
    renderCard(mcqCheck, { rng: zeroRng });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    // Shuffled order is [B, C, A]; picking index 0 is the correct option.
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Guessed" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    const region = screen.getByRole("status");
    expect(region.textContent).toContain("Correct");
    // Focus lands on the whole reveal panel, not the (narrower) status line —
    // a live region wrapping the Pass/Not yet buttons elsewhere in the panel
    // would re-announce on every aria-pressed toggle, so role=status here is
    // scoped to just this one line.
    expect(document.activeElement).toBe(screen.getByTestId("stage-reveal"));
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

  it("does not wrap the short_answer Pass/Not yet buttons in a live region", () => {
    renderCard(shortAnswerCheck);
    fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
    fireEvent.click(screen.getByRole("radio", { name: "Guessed" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(screen.queryByRole("status")).toBeNull();
    expect(document.activeElement).toBe(screen.getByTestId("stage-reveal"));
    expect(screen.getByRole("button", { name: "Pass" })).toBeTruthy();
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

describe("RecallCheckCard — onAttemptReady (Q-2b)", () => {
  it("fires once, at stage 3, not at stage 2 (commit)", () => {
    const onAttemptReady = vi.fn();
    renderCard(mcqCheck, { rng: zeroRng, onAttemptReady });

    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    expect(onAttemptReady).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(onAttemptReady).toHaveBeenCalledTimes(1);
  });

  it("reports outcome=correct with the picked option text when the mcq answer is right", () => {
    const onAttemptReady = vi.fn();
    renderCard(mcqCheck, { rng: zeroRng, onAttemptReady });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    // Shuffled order is [B, C, A]; index 0 (B) is the expected_answer.
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(onAttemptReady).toHaveBeenCalledWith(0, {
      confidence: "confident",
      outcome: "correct",
      selectedOption: "Option B",
    });
  });

  it("reports outcome=incorrect with the picked option text when the mcq answer is wrong", () => {
    const onAttemptReady = vi.fn();
    renderCard(mcqCheck, { rng: zeroRng, onAttemptReady });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    // Shuffled order is [B, C, A]; index 2 (A) is wrong.
    fireEvent.click(screen.getAllByRole("radio")[2]);
    fireEvent.click(screen.getByRole("radio", { name: "Unsure" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(onAttemptReady).toHaveBeenCalledWith(0, {
      confidence: "unsure",
      outcome: "incorrect",
      selectedOption: "Option A",
    });
  });

  it.each(["Guessed", "Unsure", "Confident"] as const)(
    "propagates confidence='%s' verbatim, not a hardcoded default",
    (label) => {
      const onAttemptReady = vi.fn();
      renderCard(mcqCheck, { rng: zeroRng, onAttemptReady });
      fireEvent.click(screen.getByRole("button", { name: "Show options" }));
      fireEvent.click(screen.getAllByRole("radio")[0]);
      fireEvent.click(screen.getByRole("radio", { name: label }));
      fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

      expect(onAttemptReady).toHaveBeenCalledWith(
        0,
        expect.objectContaining({ confidence: label.toLowerCase() }),
      );
    },
  );

  it("reports selectedOption=null for short_answer (never the mcq option shape)", () => {
    const onAttemptReady = vi.fn();
    renderCard(shortAnswerCheck, { onAttemptReady });
    fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
    fireEvent.click(screen.getByRole("radio", { name: "Guessed" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    fireEvent.click(screen.getByRole("button", { name: "Pass" }));

    expect(onAttemptReady).toHaveBeenCalledWith(1, {
      confidence: "guessed",
      outcome: "correct",
      selectedOption: null,
    });
  });

  it("re-fires as a correction when the short_answer rating changes (Pass -> Not yet), but not on a repeat click of the same rating", () => {
    const onAttemptReady = vi.fn();
    renderCard(shortAnswerCheck, { onAttemptReady });
    fireEvent.click(screen.getByRole("button", { name: "I've answered" }));
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    fireEvent.click(screen.getByRole("button", { name: "Pass" }));
    expect(onAttemptReady).toHaveBeenCalledTimes(1);
    expect(onAttemptReady).toHaveBeenLastCalledWith(1, { confidence: "confident", outcome: "correct", selectedOption: null });

    // Clicking the same rating again must not resubmit — React bails out of
    // a state update to an identical primitive, so the effect never re-runs.
    fireEvent.click(screen.getByRole("button", { name: "Pass" }));
    expect(onAttemptReady).toHaveBeenCalledTimes(1);

    // An actual correction (Pass -> Not yet) is new information for SRS —
    // decided as "resubmit", not "ignore" or "first one wins" (see quiz.md).
    fireEvent.click(screen.getByRole("button", { name: "Not yet" }));
    expect(onAttemptReady).toHaveBeenCalledTimes(2);
    expect(onAttemptReady).toHaveBeenLastCalledWith(1, { confidence: "confident", outcome: "incorrect", selectedOption: null });
  });

  it("does not re-fire on a re-render that leaves the completed attempt unchanged, even if onAttemptReady's own identity changes", () => {
    const calls: unknown[] = [];
    function Wrapper() {
      const [, setTick] = useState(0);
      return (
        <>
          <button onClick={() => setTick((t) => t + 1)}>bump</button>
          <RecallCheckCard
            check={mcqCheck}
            index={0}
            onFinishedChange={vi.fn()}
            onAttemptReady={(p, a) => calls.push([p, a])}
            saveStatus="idle"
            onRetrySave={vi.fn()}
            rng={zeroRng}
          />
        </>
      );
    }
    render(<Wrapper />);
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));
    expect(calls).toHaveLength(1);

    fireEvent.click(screen.getByRole("button", { name: "bump" }));
    expect(calls).toHaveLength(1);
  });

  it("does not call onAttemptReady twice under React StrictMode's double-invoked effects", () => {
    const onAttemptReady = vi.fn();
    render(
      <StrictMode>
        <RecallCheckCard
          check={mcqCheck}
          index={0}
          onFinishedChange={vi.fn()}
          onAttemptReady={onAttemptReady}
          saveStatus="idle"
          onRetrySave={vi.fn()}
          rng={zeroRng}
        />
      </StrictMode>,
    );
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(onAttemptReady).toHaveBeenCalledTimes(1);
  });
});

describe("RecallCheckCard — save-status indicator", () => {
  it("shows nothing extra while idle or saved", () => {
    renderCard(mcqCheck, { rng: zeroRng, saveStatus: "idle" });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(screen.queryByText("Not saved")).toBeNull();
    expect(screen.queryByRole("button", { name: "Retry" })).toBeNull();
  });

  it("shows a 'Not saved' alert with a Retry action on saveStatus='error' — never implying success", () => {
    const onRetrySave = vi.fn();
    renderCard(mcqCheck, { rng: zeroRng, saveStatus: "error", onRetrySave });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(screen.getByRole("alert").textContent).toContain("Not saved");
    fireEvent.click(screen.getByRole("button", { name: "Retry" }));
    expect(onRetrySave).toHaveBeenCalledTimes(1);
  });

  it("does not show the Retry affordance when saveStatus='saved' (the save actually succeeded)", () => {
    renderCard(mcqCheck, { rng: zeroRng, saveStatus: "saved" });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(screen.queryByRole("button", { name: "Retry" })).toBeNull();
    expect(screen.queryByText("Not saved")).toBeNull();
  });

  it("shows 'Saving…' on saveStatus='saving', with no Retry affordance", () => {
    renderCard(mcqCheck, { rng: zeroRng, saveStatus: "saving" });
    fireEvent.click(screen.getByRole("button", { name: "Show options" }));
    fireEvent.click(screen.getAllByRole("radio")[0]);
    fireEvent.click(screen.getByRole("radio", { name: "Confident" }));
    fireEvent.click(screen.getByRole("button", { name: "Reveal answer" }));

    expect(screen.getByText("Saving…")).toBeTruthy();
    expect(screen.queryByRole("button", { name: "Retry" })).toBeNull();
  });
});

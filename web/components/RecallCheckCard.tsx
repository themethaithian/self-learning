"use client";

import { forwardRef, useEffect, useRef, useState } from "react";
import type { RecallCheck } from "@/lib/api";
import { shuffleOptions } from "@/lib/shuffle";
import { Button } from "@/components/Button";
import { CheckIcon, XIcon } from "@/components/icons";

const RATE_OPTIONS = [
  { value: "pass", label: "Pass" },
  { value: "fail", label: "Not yet" },
] as const;

export type RecallRating = (typeof RATE_OPTIONS)[number]["value"];

const CONFIDENCE_OPTIONS = [
  { value: "guessed", label: "Guessed" },
  { value: "unsure", label: "Unsure" },
  { value: "confident", label: "Confident" },
] as const;

export type Confidence = (typeof CONFIDENCE_OPTIONS)[number]["value"];

export type RecallStage = "recall" | "commit" | "reveal";

/**
 * The one place "finished" is defined per check type: mcq has no self-report,
 * so reaching the reveal stage is the whole signal; short_answer can't be
 * graded automatically, so it also needs the Pass/Not yet rating.
 */
export function isRecallCheckFinished(
  checkType: RecallCheck["type"],
  stage: RecallStage,
  shortAnswerRating: RecallRating | null,
): boolean {
  if (stage !== "reveal") return false;
  return checkType === "mcq" || shortAnswerRating != null;
}

interface RecallCheckCardProps {
  check: RecallCheck;
  index: number;
  onFinishedChange: (position: number, finished: boolean) => void;
  rng?: () => number;
}

export const RecallCheckCard = forwardRef<HTMLDivElement, RecallCheckCardProps>(function RecallCheckCard(
  { check, index, onFinishedChange, rng },
  ref,
) {
  const hasOptions = check.type === "mcq" && !!check.options && check.options.length > 0;

  // Lazy useState initialiser: React guarantees this runs exactly once per
  // mount, unlike useMemo which is only a cache React is allowed to discard —
  // a mid-interaction reshuffle would move the option under the user's finger
  // between tapping and confirming.
  const [shuffledOptions] = useState(() => (hasOptions ? shuffleOptions(check.options as string[], rng) : []));

  const [stage, setStage] = useState<RecallStage>("recall");
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const [confidence, setConfidence] = useState<Confidence | null>(null);
  const [shortAnswerRating, setShortAnswerRating] = useState<RecallRating | null>(null);

  const commitRef = useRef<HTMLDivElement>(null);
  const revealRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (stage === "commit") commitRef.current?.focus({ preventScroll: true });
    if (stage === "reveal") revealRef.current?.focus({ preventScroll: true });
  }, [stage]);

  useEffect(() => {
    onFinishedChange(check.position, isRecallCheckFinished(check.type, stage, shortAnswerRating));
  }, [check.position, check.type, stage, shortAnswerRating, onFinishedChange]);

  const canReveal = hasOptions ? selectedIndex != null && confidence != null : confidence != null;

  function handleReveal() {
    if (!canReveal) return;
    setStage("reveal");
  }

  // selectedIndex indexes into shuffledOptions (this card's own stable
  // order), never the option text — two options with identical text would
  // otherwise both light up as "selected" together.
  const selectedOptionText = hasOptions && selectedIndex != null ? shuffledOptions[selectedIndex] : null;
  const isMcqCorrect = selectedOptionText != null && selectedOptionText === check.expected_answer;

  return (
    <div
      ref={ref}
      tabIndex={-1}
      className="rounded-2xl border border-subtle bg-surface p-5 shadow-sm focus:outline-none focus:ring-2 focus:ring-focus focus:ring-offset-2"
    >
      <div className="flex items-start justify-between gap-3">
        <p className="font-thai text-base leading-[1.8] text-body">{check.question}</p>
        <span className="shrink-0 rounded-full bg-page px-2 py-0.5 text-xs font-medium text-faint">#{index + 1}</span>
      </div>

      {stage === "recall" && (
        <Button variant="ghost" className="mt-4" onClick={() => setStage("commit")}>
          {check.type === "mcq" ? "Show options" : "I've answered"}
        </Button>
      )}

      {stage === "commit" && (
        <div
          ref={commitRef}
          tabIndex={-1}
          data-testid="stage-commit"
          className="mt-4 space-y-4 focus:outline-none focus:ring-2 focus:ring-focus focus:ring-offset-2 rounded-xl"
        >
          {hasOptions && (
            <fieldset>
              <legend className="text-xs font-medium uppercase tracking-wide text-faint">Choose one</legend>
              <div className="mt-2 space-y-2">
                {shuffledOptions.map((option, i) => (
                  <label
                    key={i}
                    className="flex cursor-pointer items-center gap-2 rounded-lg border border-subtle px-3 py-2 transition-colors duration-150 ease-out has-[:checked]:border-accent has-[:checked]:bg-accent-tint has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-focus has-[:focus-visible]:ring-offset-2"
                  >
                    <input
                      type="radio"
                      name={`mcq-${check.position}`}
                      checked={selectedIndex === i}
                      onChange={() => setSelectedIndex(i)}
                      className="h-4 w-4 shrink-0 accent-accent"
                    />
                    <span className="font-thai text-sm leading-[1.8] text-body">{option}</span>
                  </label>
                ))}
              </div>
            </fieldset>
          )}

          <fieldset>
            <legend className="text-xs font-medium uppercase tracking-wide text-faint">How confident are you?</legend>
            <div className="mt-2 flex flex-wrap gap-2">
              {CONFIDENCE_OPTIONS.map((opt) => (
                <label
                  key={opt.value}
                  className="cursor-pointer rounded-full border border-subtle px-3 py-1 text-xs font-medium text-muted transition-colors duration-150 ease-out has-[:checked]:border-accent has-[:checked]:bg-accent-tint has-[:checked]:text-accent-strong has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-focus has-[:focus-visible]:ring-offset-2"
                >
                  <input
                    type="radio"
                    name={`confidence-${check.position}`}
                    checked={confidence === opt.value}
                    onChange={() => setConfidence(opt.value)}
                    className="sr-only"
                  />
                  {opt.label}
                </label>
              ))}
            </div>
          </fieldset>

          <Button variant="ghost" aria-disabled={!canReveal} onClick={handleReveal}>
            Reveal answer
          </Button>
        </div>
      )}

      {stage === "reveal" && (
        <div
          ref={revealRef}
          tabIndex={-1}
          role="status"
          data-testid="stage-reveal"
          className="mt-4 space-y-3 focus:outline-none focus:ring-2 focus:ring-focus focus:ring-offset-2 rounded-xl"
        >
          {hasOptions && selectedOptionText != null && (
            <div
              className={`flex items-center gap-2 rounded-xl px-4 py-3 text-sm font-medium ${
                isMcqCorrect ? "bg-success/10 text-success-strong" : "bg-danger/10 text-danger-strong"
              }`}
            >
              {isMcqCorrect ? <CheckIcon /> : <XIcon />}
              <span>{isMcqCorrect ? "Correct" : "Not quite"}</span>
            </div>
          )}

          <div className="rounded-xl bg-accent-tint px-4 py-3">
            <p className="text-xs font-medium uppercase tracking-wide text-accent-strong">Answer</p>
            <p className="mt-1 font-thai text-sm leading-[1.8] text-body">{check.expected_answer}</p>
          </div>

          {hasOptions && (
            <ul className="space-y-1">
              {shuffledOptions.map((option, i) => {
                // Correctness is derived by matching expected_answer, never
                // by index — the index only survives the shuffle, meaning
                // stays fixed with the option text.
                const isCorrectOption = option === check.expected_answer;
                const isYourAnswer = i === selectedIndex;
                const tag = isCorrectOption && isYourAnswer
                  ? "Your answer — correct"
                  : isCorrectOption
                    ? "Correct answer"
                    : isYourAnswer
                      ? "Your answer — incorrect"
                      : null;
                return (
                  <li
                    key={i}
                    className={`flex items-center gap-2 rounded-lg border px-3 py-1.5 font-thai text-sm ${
                      isCorrectOption
                        ? "border-success-strong bg-success/10 text-success-strong"
                        : isYourAnswer
                          ? "border-danger-strong bg-danger/10 text-danger-strong"
                          : "border-subtle text-muted"
                    }`}
                  >
                    {isCorrectOption ? <CheckIcon /> : isYourAnswer ? <XIcon /> : <span className="h-4 w-4 shrink-0" />}
                    <span>{option}</span>
                    {tag && <span className="text-xs font-medium">({tag})</span>}
                  </li>
                );
              })}
            </ul>
          )}

          {check.type === "short_answer" && (
            <div className="flex flex-wrap items-center gap-2 pt-1">
              <span className="text-xs font-medium uppercase tracking-wide text-faint">How did you do?</span>
              {RATE_OPTIONS.map((opt) => (
                <button
                  key={opt.value}
                  type="button"
                  onClick={() => setShortAnswerRating(opt.value)}
                  aria-pressed={shortAnswerRating === opt.value}
                  className={`rounded-full px-3 py-1 text-xs font-medium transition-colors duration-150 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2 ${
                    shortAnswerRating === opt.value
                      ? opt.value === "pass"
                        ? "bg-success-strong text-white"
                        : "bg-danger text-white"
                      : "border border-subtle text-muted hover:bg-page"
                  }`}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
});

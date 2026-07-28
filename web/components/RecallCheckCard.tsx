"use client";

import { forwardRef, useState } from "react";
import type { RecallCheck } from "@/lib/api";

const RATE_OPTIONS = [
  { value: "pass", label: "Pass" },
  { value: "fail", label: "Not yet" },
] as const;

export type RecallRating = (typeof RATE_OPTIONS)[number]["value"];

interface RecallCheckCardProps {
  check: RecallCheck;
  index: number;
  rating: RecallRating | null;
  onRate: (value: RecallRating) => void;
}

export const RecallCheckCard = forwardRef<HTMLDivElement, RecallCheckCardProps>(function RecallCheckCard(
  { check, index, rating, onRate },
  ref,
) {
  const [revealed, setRevealed] = useState(false);

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

      {!revealed ? (
        <button
          type="button"
          onClick={() => setRevealed(true)}
          className="mt-4 rounded-xl border border-subtle px-4 py-2 text-sm font-medium text-accent transition-colors duration-150 ease-out hover:bg-accent-tint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
        >
          Reveal answer
        </button>
      ) : (
        <div className="mt-4 space-y-3">
          <div className="rounded-xl bg-accent-tint px-4 py-3">
            <p className="text-xs font-medium uppercase tracking-wide text-accent-strong">Answer</p>
            <p className="mt-1 font-thai text-sm leading-[1.8] text-body">{check.expected_answer}</p>
          </div>

          {check.type === "mcq" && check.options && check.options.length > 0 && (
            <ul className="space-y-1">
              {check.options.map((option) => {
                const isCorrect = option === check.expected_answer;
                return (
                  <li
                    key={option}
                    className={`rounded-lg border px-3 py-1.5 font-thai text-sm ${
                      isCorrect ? "border-success-strong bg-success/10 text-success-strong" : "border-subtle text-muted"
                    }`}
                  >
                    {option}
                    {isCorrect && <span className="ml-2 text-xs font-medium">(correct)</span>}
                  </li>
                );
              })}
            </ul>
          )}

          <div className="flex flex-wrap items-center gap-2 pt-1">
            <span className="text-xs font-medium uppercase tracking-wide text-faint">How did you do?</span>
            {RATE_OPTIONS.map((opt) => (
              <button
                key={opt.value}
                type="button"
                onClick={() => onRate(opt.value)}
                aria-pressed={rating === opt.value}
                className={`rounded-full px-3 py-1 text-xs font-medium transition-colors duration-150 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2 ${
                  rating === opt.value
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
        </div>
      )}
    </div>
  );
});

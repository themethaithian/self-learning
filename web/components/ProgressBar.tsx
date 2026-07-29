// value is always out of "lessons available right now" (has_lesson count),
// never out of the full planned-concept count — a full bar must only ever
// mean "read everything that currently exists" (UX-7), not "the curriculum
// is finished".
export function ProgressBar({ value }: { value: number }) {
  const pct = Math.max(0, Math.min(100, value));
  return (
    <div
      role="progressbar"
      aria-valuenow={Math.round(pct)}
      aria-valuemin={0}
      aria-valuemax={100}
      className="h-2 w-full overflow-hidden rounded-full bg-subtle"
    >
      <div
        className="h-full rounded-full bg-accent transition-[width] duration-150 ease-out"
        style={{ width: `${pct}%` }}
      />
    </div>
  );
}

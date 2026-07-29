interface ProgressBarProps {
  read: number;
  available: number;
  label: string;
}

// Percentage and the divide-by-zero guard are computed here, not by the
// caller — a caller can only ever pass the raw counts, never a pre-baked
// percentage that might secretly be out of the wrong denominator.
// `label` is required (not optional) so a role="progressbar" can never ship
// without an accessible name (WCAG 1.1.1 / axe aria-progressbar-name).
export function ProgressBar({ read, available, label }: ProgressBarProps) {
  const pct = available <= 0 ? 0 : Math.max(0, Math.min(100, (read / available) * 100));
  return (
    <div
      role="progressbar"
      aria-label={label}
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

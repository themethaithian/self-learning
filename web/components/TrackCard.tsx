import { Card } from "@/components/Card";
import { LinkButton } from "@/components/Button";

export interface TrackStats {
  track: string;
  label: string;
  totalConcepts: number;
  lessonsReady: number;
  chapterCount: number;
}

interface TrackCardProps {
  stats: TrackStats;
  isFocus: boolean;
  onSetFocus: (track: string | null) => void;
  saving: boolean;
  busy: boolean;
}

function plural(count: number, word: string): string {
  return `${count} ${word}${count === 1 ? "" : "s"}`;
}

export function TrackCard({ stats, isFocus, onSetFocus, saving, busy }: TrackCardProps) {
  return (
    <Card className="space-y-4">
      <div className="flex flex-wrap items-start justify-between gap-2">
        <h3 className="min-w-0 break-words text-base font-semibold text-heading">{stats.label}</h3>
        <button
          type="button"
          onClick={() => {
            // aria-disabled, not the disabled attribute: a native disabled
            // button is forced to blur, which would kick keyboard focus off
            // this button the instant its own click starts the request.
            if (busy) return;
            onSetFocus(isFocus ? null : stats.track);
          }}
          aria-disabled={busy}
          aria-pressed={isFocus}
          className={`inline-flex shrink-0 items-center gap-1 rounded-full px-3 py-1 text-xs font-medium transition-colors duration-150 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2 ${
            busy ? "opacity-50" : ""
          } ${
            isFocus
              ? "bg-accent text-white"
              : "border border-subtle text-muted hover:bg-accent-tint hover:text-accent-strong"
          }`}
        >
          <span aria-hidden>{isFocus ? "★" : "☆"}</span>
          {saving ? "Saving…" : isFocus ? "Focus" : "Set focus"}
        </button>
      </div>

      {/* No ProgressBar here yet: reading progress isn't tracked (UX-4), and a
          bar that's always full (n = lessonsReady = N) would read as "fully
          read" when really nothing has been read at all — worse than no bar. */}
      <div className="flex flex-wrap items-center justify-between gap-x-3 gap-y-0.5 text-xs text-muted">
        <span>{plural(stats.lessonsReady, "lesson")} ready</span>
        <span>
          {plural(stats.chapterCount, "chapter")} · {plural(stats.totalConcepts, "concept")} planned
        </span>
      </div>

      <LinkButton href={`/learn?track=${encodeURIComponent(stats.track)}`} variant="ghost">
        View track
      </LinkButton>
    </Card>
  );
}

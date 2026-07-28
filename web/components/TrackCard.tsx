import Link from "next/link";
import { Card } from "@/components/Card";

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
  onSetFocus: (track: string) => void;
  settingFocus: boolean;
}

export function TrackCard({ stats, isFocus, onSetFocus, settingFocus }: TrackCardProps) {
  return (
    <Card className="space-y-4">
      <div className="flex flex-wrap items-start justify-between gap-2">
        <h3 className="min-w-0 break-words text-base font-semibold text-heading">{stats.label}</h3>
        {isFocus ? (
          <span className="inline-flex shrink-0 items-center gap-1 rounded-full bg-accent px-3 py-1 text-xs font-medium text-white">
            <span aria-hidden>★</span> Focus
          </span>
        ) : (
          <button
            type="button"
            onClick={() => onSetFocus(stats.track)}
            disabled={settingFocus}
            className="inline-flex shrink-0 items-center gap-1 rounded-full border border-subtle px-3 py-1 text-xs font-medium text-muted transition-colors duration-150 ease-out hover:bg-accent-tint hover:text-accent-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          >
            <span aria-hidden>☆</span> Set focus
          </button>
        )}
      </div>

      {/* No ProgressBar here yet: reading progress isn't tracked (UX-4), and a
          bar that's always full (n = lessonsReady = N) would read as "fully
          read" when really nothing has been read at all — worse than no bar. */}
      <div className="flex flex-wrap items-center justify-between gap-x-3 gap-y-0.5 text-xs text-muted">
        <span>{stats.lessonsReady} lessons ready</span>
        <span>
          {stats.chapterCount} chapters · {stats.totalConcepts} concepts planned
        </span>
      </div>

      <Link
        href={`/learn?track=${encodeURIComponent(stats.track)}`}
        className="inline-flex items-center justify-center gap-2 rounded-xl border border-subtle px-4 py-2 text-sm font-medium text-body transition-colors duration-150 ease-out hover:bg-accent-tint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
      >
        View track
      </Link>
    </Card>
  );
}

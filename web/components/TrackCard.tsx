import { Card } from "@/components/Card";
import { LinkButton } from "@/components/Button";
import { FocusToggleButton } from "@/components/FocusToggleButton";
import { ProgressBar } from "@/components/ProgressBar";
import type { TrackReadStats } from "@/lib/curriculum";

export interface TrackStats {
  track: string;
  label: string;
  totalConcepts: number;
  lessonsReady: number;
  chapterCount: number;
}

interface TrackCardProps {
  stats: TrackStats;
  readStats: TrackReadStats | null;
  isFocus: boolean;
  onSetFocus: (track: string | null) => void;
  saving: boolean;
  busy: boolean;
}

export function plural(count: number, word: string): string {
  return `${count} ${word}${count === 1 ? "" : "s"}`;
}

export function TrackCard({ stats, readStats, isFocus, onSetFocus, saving, busy }: TrackCardProps) {
  return (
    <Card className="space-y-4">
      <div className="flex flex-wrap items-start justify-between gap-2">
        <h3 className="min-w-0 break-words text-base font-semibold text-heading">{stats.label}</h3>
        <FocusToggleButton
          isFocus={isFocus}
          saving={saving}
          busy={busy}
          onToggle={() => onSetFocus(isFocus ? null : stats.track)}
        />
      </div>

      {readStats ? (
        <div className="space-y-1.5">
          <ProgressBar value={readStats.available === 0 ? 0 : (readStats.read / readStats.available) * 100} />
          <p className="text-xs text-muted">
            {readStats.read}/{readStats.available} {readStats.available === 1 ? "lesson" : "lessons"} read
          </p>
        </div>
      ) : (
        <p className="text-xs text-muted">{plural(stats.lessonsReady, "lesson")} ready</p>
      )}

      <p className="text-xs text-muted">
        {plural(stats.chapterCount, "chapter")} · {plural(stats.totalConcepts, "concept")} planned
      </p>

      <LinkButton href={`/learn?track=${encodeURIComponent(stats.track)}`} variant="ghost">
        View track
      </LinkButton>
    </Card>
  );
}

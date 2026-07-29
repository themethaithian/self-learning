import { Card } from "@/components/Card";
import { LinkButton } from "@/components/Button";
import { FocusToggleButton } from "@/components/FocusToggleButton";

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

export function plural(count: number, word: string): string {
  return `${count} ${word}${count === 1 ? "" : "s"}`;
}

export function TrackCard({ stats, isFocus, onSetFocus, saving, busy }: TrackCardProps) {
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

      {/* No ProgressBar here: this card counts lesson availability
          (lessonsReady of totalConcepts), not reading completion, so a bar
          would always render full — completion numbers live on /today. */}
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

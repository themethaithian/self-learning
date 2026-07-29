"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  getCurriculum,
  getFocusTrack,
  getProgress,
  setFocusTrack as putFocusTrack,
  UnauthorizedError,
  type Track,
} from "@/lib/api";
import { computeTrackProgress, indexProgress, pickNextUp, type NextUpResult, type ProgressByKey } from "@/lib/curriculum";
import { sortTracksByDisplayOrder, trackLabel } from "@/lib/trackMeta";
import { Card } from "@/components/Card";
import { Button, LinkButton } from "@/components/Button";
import { EmptyState } from "@/components/EmptyState";
import { FocusToggleButton } from "@/components/FocusToggleButton";
import { TodaySkeleton } from "@/components/Skeleton";
import { BookIcon, WarningIcon } from "@/components/icons";

type ViewState =
  | { status: "loading" }
  | { status: "error"; message: string }
  | { status: "success"; tracks: Track[] };

function NextUpSection({ next }: { next: NextUpResult }) {
  if (next.kind === "no-lessons") {
    return <EmptyState icon={<BookIcon />} message="No lessons are imported yet — check back once a batch is written." />;
  }

  if (next.kind === "all-done") {
    return (
      <EmptyState
        icon={<BookIcon />}
        message="You've read everything available right now — nice work. Check back once a new batch is written."
      />
    );
  }

  return (
    <Card className="space-y-3">
      <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Next up</h2>
      <p className="text-sm text-muted">
        {trackLabel(next.track)} · {next.chapterTitle}
      </p>
      <h3 className="min-w-0 break-words text-lg font-semibold text-heading">{next.title}</h3>
      {next.estMinutes != null && next.estMinutes > 0 && (
        <span className="inline-flex items-center rounded-full bg-accent-tint px-3 py-1 text-xs font-medium text-accent-strong">
          ~{next.estMinutes} min
        </span>
      )}
      <div>
        <LinkButton
          variant="primary"
          href={`/lesson?topic=${encodeURIComponent(next.topic)}&concept=${encodeURIComponent(next.concept)}`}
        >
          {next.state === "in_progress" ? "Continue" : "Start"}
        </LinkButton>
      </div>
    </Card>
  );
}

function OtherTrackRow({
  label,
  done,
  total,
  track,
  isFocus,
  onSetFocus,
  saving,
  busy,
}: {
  label: string;
  done: number;
  total: number;
  track: string;
  isFocus: boolean;
  onSetFocus: (track: string | null) => void;
  saving: boolean;
  busy: boolean;
}) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-subtle bg-surface px-4 py-3">
      <div className="min-w-0">
        <p className="truncate text-sm font-medium text-heading">{label}</p>
        <p className="text-xs text-muted">
          {done}/{total} done
        </p>
      </div>
      <div className="flex shrink-0 items-center gap-2">
        <FocusToggleButton isFocus={isFocus} saving={saving} busy={busy} onToggle={() => onSetFocus(isFocus ? null : track)} />
        <LinkButton href={`/learn?track=${encodeURIComponent(track)}`} variant="ghost">
          View track
        </LinkButton>
      </div>
    </div>
  );
}

export default function TodayPage() {
  const router = useRouter();
  const [state, setState] = useState<ViewState>({ status: "loading" });
  const [progressByKey, setProgressByKeyValue] = useState<ProgressByKey>({});
  const [focusTrack, setFocusTrackValue] = useState<string | null>(null);
  const [pendingTrack, setPendingTrack] = useState<string | null>(null);
  const [focusError, setFocusError] = useState<string | null>(null);

  // getProgress/getFocusTrack already degrade non-auth failures to empty
  // defaults internally (see lib/api.ts) — a failing progress or focus-track
  // call resolves rather than rejects, so it can never take the curriculum
  // load down with it.
  const load = useCallback(async () => {
    setState({ status: "loading" });
    try {
      const [curriculum, progress, focus] = await Promise.all([getCurriculum(), getProgress(), getFocusTrack()]);
      setProgressByKeyValue(indexProgress(progress));
      setFocusTrackValue(focus.track);
      setState({ status: "success", tracks: curriculum.tracks });
    } catch (err) {
      if (err instanceof UnauthorizedError) {
        router.replace("/token");
        return;
      }
      setState({ status: "error", message: err instanceof Error ? err.message : "Something went wrong" });
    }
  }, [router]);

  useEffect(() => {
    load();
  }, [load]);

  const handleSetFocus = useCallback(
    async (track: string | null) => {
      const previous = focusTrack;
      const pendingCard = track ?? previous;
      setFocusTrackValue(track);
      setPendingTrack(pendingCard);
      setFocusError(null);
      try {
        await putFocusTrack(track);
      } catch {
        setFocusTrackValue(previous);
        setFocusError("Could not save focus track — please try again.");
      } finally {
        setPendingTrack(null);
      }
    },
    [focusTrack],
  );

  if (state.status === "loading") {
    return <TodaySkeleton />;
  }

  if (state.status === "error") {
    return (
      <EmptyState
        icon={<WarningIcon />}
        message={`Could not load today's plan: ${state.message}`}
        action={<Button onClick={load}>Retry</Button>}
      />
    );
  }

  const { tracks } = state;
  const nextUp = pickNextUp(tracks, progressByKey, focusTrack);
  const otherTracks = sortTracksByDisplayOrder(tracks)
    .map((track) => ({ label: trackLabel(track.track), ...computeTrackProgress(track, progressByKey) }))
    .filter((stats) => stats.total > 0);

  return (
    <div className="space-y-8">
      {focusError && (
        <div
          role="alert"
          className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger-strong"
        >
          <span>{focusError}</span>
          <button
            type="button"
            onClick={() => setFocusError(null)}
            className="rounded-md text-xs font-medium underline underline-offset-2 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
          >
            Dismiss
          </button>
        </div>
      )}

      <NextUpSection next={nextUp} />

      {otherTracks.length > 0 && (
        <section className="space-y-3">
          <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Other tracks</h2>
          <div className="space-y-2">
            {otherTracks.map((stats) => (
              <OtherTrackRow
                key={stats.track}
                label={stats.label}
                done={stats.done}
                total={stats.total}
                track={stats.track}
                isFocus={stats.track === focusTrack}
                onSetFocus={handleSetFocus}
                saving={pendingTrack === stats.track}
                busy={pendingTrack !== null}
              />
            ))}
          </div>
        </section>
      )}
    </div>
  );
}

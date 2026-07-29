"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getCurriculum, getFocusTrack, getProgress, UnauthorizedError, type Track } from "@/lib/api";
import { computeTrackProgress, indexProgress, pickNextUp, pinFocusFirst, type NextUpResult, type ProgressByKey } from "@/lib/curriculum";
import { sortTracksByDisplayOrder, trackLabel } from "@/lib/trackMeta";
import { useFocusTrack } from "@/lib/useFocusTrack";
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
  progressAvailable,
  isFocus,
  onSetFocus,
  saving,
  busy,
}: {
  label: string;
  done: number;
  total: number;
  track: string;
  progressAvailable: boolean;
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
          {progressAvailable ? `${done}/${total} done` : `${total} lesson${total === 1 ? "" : "s"}`}
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
  const [progressAvailable, setProgressAvailable] = useState(true);
  const { focusTrack, setFocusTrackValue, pendingTrack, error: focusError, setFocus: handleSetFocus, dismiss: dismissFocusError } =
    useFocusTrack();

  // getFocusTrack already degrades non-auth failures to {track: null}
  // internally (see lib/api.ts) — only getCurriculum() failing can reject
  // this Promise.all, so progress/focus-track never take the page down.
  const load = useCallback(async () => {
    setState({ status: "loading" });
    try {
      const [curriculum, progress, focus] = await Promise.all([getCurriculum(), getProgress(), getFocusTrack()]);
      if (progress.ok) {
        setProgressByKeyValue(indexProgress(progress.entries));
        setProgressAvailable(true);
      } else {
        setProgressByKeyValue({});
        setProgressAvailable(false);
      }
      setFocusTrackValue(focus.track);
      setState({ status: "success", tracks: curriculum.tracks });
    } catch (err) {
      if (err instanceof UnauthorizedError) {
        router.replace("/token");
        return;
      }
      setState({ status: "error", message: err instanceof Error ? err.message : "Something went wrong" });
    }
  }, [router, setFocusTrackValue]);

  useEffect(() => {
    load();
  }, [load]);

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
  const otherTracks = pinFocusFirst(sortTracksByDisplayOrder(tracks), focusTrack)
    .map((track) => ({ label: trackLabel(track.track), ...computeTrackProgress(track, progressByKey) }))
    .filter((stats) => stats.total > 0)
    .filter((stats) => !(nextUp.kind === "next" && stats.track === nextUp.track));

  return (
    <div className="space-y-8">
      {!progressAvailable && (
        <div role="alert" className="rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger-strong">
          Could not load your reading progress — counts below and the Continue/Start label above may be out of date.
        </div>
      )}

      {focusError && (
        <div
          role="alert"
          className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger-strong"
        >
          <span>{focusError}</span>
          <button
            type="button"
            onClick={dismissFocusError}
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
                progressAvailable={progressAvailable}
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

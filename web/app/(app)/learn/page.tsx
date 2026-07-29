"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { getCurriculum, getFocusTrack, getProgress, UnauthorizedError, type Track } from "@/lib/api";
import { TrackTopics } from "@/components/TrackTopics";
import { countAvailableLessons, indexProgress, pinFocusFirst, trackReadStats, type ProgressByKey } from "@/lib/curriculum";
import { useFocusTrack } from "@/lib/useFocusTrack";
import { TrackCard, type TrackStats } from "@/components/TrackCard";
import { TrackCardsSkeleton } from "@/components/Skeleton";
import { EmptyState } from "@/components/EmptyState";
import { Button, LinkButton } from "@/components/Button";
import { BookIcon, WarningIcon } from "@/components/icons";
import { sortTracksByDisplayOrder, trackLabel } from "@/lib/trackMeta";

type ViewState =
  | { status: "loading" }
  | { status: "error"; message: string }
  | { status: "empty" }
  | { status: "success"; tracks: Track[] };

function computeStats(track: Track): TrackStats {
  let totalConcepts = 0;
  let lessonsReady = 0;
  let chapterCount = 0;
  for (const topic of track.topics) {
    chapterCount += topic.chapters.length;
    for (const chapter of topic.chapters) {
      const { available, total } = countAvailableLessons(chapter.concepts);
      totalConcepts += total;
      lessonsReady += available;
    }
  }
  return { track: track.track, label: trackLabel(track.track), totalConcepts, lessonsReady, chapterCount };
}

function ComingSoonRow({ label }: { label: string }) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-subtle/60 px-4 py-3 text-sm text-faint">
      <span>{label}</span>
      <span className="text-xs">Coming soon</span>
    </div>
  );
}

function TrackDetail({
  track,
  progressByKey,
  progressAvailable,
}: {
  track: Track;
  progressByKey: ProgressByKey;
  progressAvailable: boolean;
}) {
  return (
    <div className="space-y-6">
      <LinkButton href="/learn" variant="ghost">
        ← Back to Learn
      </LinkButton>
      <h2 className="text-xl font-semibold text-heading">{trackLabel(track.track)}</h2>
      <TrackTopics topics={track.topics} progressByKey={progressByKey} progressAvailable={progressAvailable} />
    </div>
  );
}

function ProgressUnavailableBanner() {
  return (
    <div role="alert" className="rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger-strong">
      Could not load your reading progress — the read markers on each concept, the read count on each chapter, and
      the progress bar on each track card may be wrong until this loads.
    </div>
  );
}

function LearnView() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const trackParam = searchParams.get("track");

  const [state, setState] = useState<ViewState>({ status: "loading" });
  const [progressByKey, setProgressByKeyValue] = useState<ProgressByKey>({});
  const [progressAvailable, setProgressAvailable] = useState(true);
  const { focusTrack, setFocusTrackValue, pendingTrack, error: focusError, setFocus: handleSetFocus, dismiss: dismissFocusError } =
    useFocusTrack();

  // getFocusTrack/getProgress already degrade non-auth failures internally
  // (see lib/api.ts) — only getCurriculum() failing can reject this
  // Promise.all, so a progress-fetch failure never takes the whole page down.
  const load = useCallback(async () => {
    setState({ status: "loading" });
    try {
      const [curriculum, progress, focus] = await Promise.all([getCurriculum(), getProgress(), getFocusTrack()]);
      if (progress.kind === "ok") {
        setProgressByKeyValue(indexProgress(progress.entries));
        setProgressAvailable(true);
      } else {
        setProgressByKeyValue({});
        setProgressAvailable(false);
      }
      setFocusTrackValue(focus.track);
      setState(
        curriculum.tracks.length === 0 ? { status: "empty" } : { status: "success", tracks: curriculum.tracks },
      );
    } catch (err) {
      if (err instanceof UnauthorizedError) {
        router.replace("/token");
        return;
      }
      setState({
        status: "error",
        message: err instanceof Error ? err.message : "Something went wrong",
      });
    }
  }, [router, setFocusTrackValue]);

  useEffect(() => {
    load();
  }, [load]);

  if (state.status === "loading") {
    return <TrackCardsSkeleton />;
  }

  if (state.status === "error") {
    return (
      <EmptyState
        icon={<WarningIcon />}
        message={`Could not load the curriculum: ${state.message}`}
        action={<Button onClick={load}>Retry</Button>}
      />
    );
  }

  if (state.status === "empty") {
    return <EmptyState icon={<BookIcon />} message="No curriculum has been imported yet." />;
  }

  const { tracks } = state;

  if (trackParam) {
    const matched = tracks.find((t) => t.track === trackParam);
    if (!matched) {
      return (
        <EmptyState
          icon={<BookIcon />}
          message={`Track "${trackParam}" was not found.`}
          action={
            <Button variant="ghost" onClick={() => router.push("/learn")}>
              Back to Learn
            </Button>
          }
        />
      );
    }
    return (
      <div className="space-y-6">
        {!progressAvailable && <ProgressUnavailableBanner />}
        <TrackDetail track={matched} progressByKey={progressByKey} progressAvailable={progressAvailable} />
      </div>
    );
  }

  const allStats = sortTracksByDisplayOrder(tracks).map(computeStats);
  const availableStats = pinFocusFirst(
    allStats.filter((s) => s.lessonsReady > 0),
    focusTrack,
  );
  const comingSoon = allStats.filter((s) => s.lessonsReady === 0);

  return (
    <div className="space-y-8">
      {!progressAvailable && <ProgressUnavailableBanner />}

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

      {availableStats.length === 0 ? (
        <EmptyState
          icon={<BookIcon />}
          message="No lessons are imported yet — check back once a batch is written."
        />
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {availableStats.map((stats) => {
            const trackTree = tracks.find((t) => t.track === stats.track);
            return (
              <TrackCard
                key={stats.track}
                stats={stats}
                readStats={trackTree ? trackReadStats(trackTree, progressByKey, progressAvailable) : null}
                isFocus={stats.track === focusTrack}
                onSetFocus={handleSetFocus}
                saving={pendingTrack === stats.track}
                busy={pendingTrack !== null}
              />
            );
          })}
        </div>
      )}

      {comingSoon.length > 0 && (
        <section className="space-y-2">
          <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Coming soon</h2>
          <div className="space-y-2">
            {comingSoon.map((stats) => (
              <ComingSoonRow key={stats.track} label={stats.label} />
            ))}
          </div>
        </section>
      )}
    </div>
  );
}

export default function LearnPage() {
  return (
    <Suspense fallback={<TrackCardsSkeleton />}>
      <LearnView />
    </Suspense>
  );
}

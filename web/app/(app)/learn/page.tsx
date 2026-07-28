"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import {
  getCurriculum,
  getFocusTrack,
  setFocusTrack as putFocusTrack,
  UnauthorizedError,
  type Track,
} from "@/lib/api";
import { TrackTopics } from "@/components/TrackTopics";
import { TrackCard, type TrackStats } from "@/components/TrackCard";
import { CurriculumTreeSkeleton } from "@/components/Skeleton";
import { EmptyState } from "@/components/EmptyState";
import { Button } from "@/components/Button";
import { BookIcon, WarningIcon } from "@/components/icons";
import { TRACK_ORDER, trackLabel } from "@/lib/trackMeta";

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
      totalConcepts += chapter.concepts.length;
      lessonsReady += chapter.concepts.filter((concept) => concept.has_lesson).length;
    }
  }
  return { track: track.track, label: trackLabel(track.track), totalConcepts, lessonsReady, chapterCount };
}

// Focus is pinned first; everything else keeps the curriculum's canonical
// track order so the grid doesn't reshuffle as lessons get added over time.
function orderStats(stats: TrackStats[], focusTrack: string | null): TrackStats[] {
  if (!focusTrack) return stats;
  const focusIndex = stats.findIndex((s) => s.track === focusTrack);
  if (focusIndex === -1) return stats;
  const focus = stats[focusIndex];
  return [focus, ...stats.slice(0, focusIndex), ...stats.slice(focusIndex + 1)];
}

function ComingSoonRow({ label }: { label: string }) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-subtle/60 px-4 py-3 text-sm text-faint">
      <span>{label}</span>
      <span className="text-xs">Coming soon</span>
    </div>
  );
}

function TrackDetail({ track, onBack }: { track: Track; onBack: () => void }) {
  return (
    <div className="space-y-6">
      <button
        type="button"
        onClick={onBack}
        className="inline-flex items-center gap-1 rounded-lg text-sm font-medium text-muted transition-colors duration-150 ease-out hover:text-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
      >
        ← Back to Learn
      </button>
      <h1 className="text-xl font-semibold text-heading">{trackLabel(track.track)}</h1>
      <TrackTopics topics={track.topics} />
    </div>
  );
}

function LearnView() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const trackParam = searchParams.get("track");

  const [state, setState] = useState<ViewState>({ status: "loading" });
  const [focusTrack, setFocusTrackValue] = useState<string | null>(null);
  const [pendingTrack, setPendingTrack] = useState<string | null>(null);
  const [focusError, setFocusError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setState({ status: "loading" });
    try {
      const [curriculum, focus] = await Promise.all([getCurriculum(), getFocusTrack()]);
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
  }, [router]);

  useEffect(() => {
    load();
  }, [load]);

  const handleSetFocus = useCallback(
    async (track: string) => {
      const previous = focusTrack;
      setFocusTrackValue(track);
      setPendingTrack(track);
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
    return <CurriculumTreeSkeleton />;
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
    return <TrackDetail track={matched} onBack={() => router.push("/learn")} />;
  }

  const allStats = TRACK_ORDER.map((slug) => tracks.find((t) => t.track === slug))
    .filter((t): t is Track => t !== undefined)
    .map(computeStats);
  const availableStats = orderStats(
    allStats.filter((s) => s.lessonsReady > 0),
    focusTrack,
  );
  const comingSoon = allStats.filter((s) => s.lessonsReady === 0);

  return (
    <div className="space-y-8">
      {focusError && (
        <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
          <span>{focusError}</span>
          <button
            type="button"
            onClick={() => setFocusError(null)}
            className="text-xs font-medium underline underline-offset-2"
          >
            Dismiss
          </button>
        </div>
      )}

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {availableStats.map((stats) => (
          <TrackCard
            key={stats.track}
            stats={stats}
            isFocus={stats.track === focusTrack}
            onSetFocus={handleSetFocus}
            settingFocus={pendingTrack !== null}
          />
        ))}
      </div>

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
    <Suspense fallback={<CurriculumTreeSkeleton />}>
      <LearnView />
    </Suspense>
  );
}

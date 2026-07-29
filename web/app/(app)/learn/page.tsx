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

interface TrackCardEntry extends TrackStats {
  tree: Track;
}

// Carries the raw Track alongside its derived TrackStats so the render pass
// never needs to re-find a track by slug (a lookup that can only ever
// succeed here, since every entry was built from `tracks` in the first
// place — the dead "not found" branch that used to guard that lookup is
// gone because the case it guarded against can't happen).
function computeCardEntries(tracks: Track[]): TrackCardEntry[] {
  return tracks.map((tree) => {
    let totalConcepts = 0;
    let lessonsReady = 0;
    let chapterCount = 0;
    for (const topic of tree.topics) {
      chapterCount += topic.chapters.length;
      for (const chapter of topic.chapters) {
        const { available, total } = countAvailableLessons(chapter.concepts);
        totalConcepts += total;
        lessonsReady += available;
      }
    }
    return { track: tree.track, label: trackLabel(tree.track), totalConcepts, lessonsReady, chapterCount, tree };
  });
}

function ComingSoonRow({ label }: { label: string }) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-subtle/60 px-4 py-3 text-sm text-faint">
      <span>{label}</span>
      <span className="text-xs">Coming soon</span>
    </div>
  );
}

function TrackDetail({ track, progress }: { track: Track; progress: ProgressByKey | null }) {
  return (
    <div className="space-y-6">
      <LinkButton href="/learn" variant="ghost">
        ← Back to Learn
      </LinkButton>
      <h2 className="text-xl font-semibold text-heading">{trackLabel(track.track)}</h2>
      <TrackTopics topics={track.topics} progress={progress} />
    </div>
  );
}

function ProgressUnavailableBanner() {
  return (
    <div role="alert" className="rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger-strong">
      Could not load your reading progress — the read markers, chapter read counts, and track progress bars are
      hidden until this loads. Everything else on this page is unaffected.
    </div>
  );
}

function LearnView() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const trackParam = searchParams.get("track");

  const [state, setState] = useState<ViewState>({ status: "loading" });
  // null = "hasn't loaded" (or failed to load) — the one flag doubles as
  // both the availability signal and the data itself, so a track card or
  // chapter row can't end up with a real progress map but the "unavailable"
  // banner still showing, or vice versa (R1b: the two used to be separate
  // useState calls that had to be kept in sync by hand).
  const [progressIndex, setProgressIndex] = useState<ProgressByKey | null>(null);
  const { focusTrack, setFocusTrackValue, pendingTrack, error: focusError, setFocus: handleSetFocus, dismiss: dismissFocusError } =
    useFocusTrack();

  // getFocusTrack/getProgress already degrade non-auth failures internally
  // (see lib/api.ts) — only getCurriculum() failing can reject this
  // Promise.all, so a progress-fetch failure never takes the whole page down.
  const load = useCallback(async () => {
    setState({ status: "loading" });
    try {
      const [curriculum, progressResult, focus] = await Promise.all([getCurriculum(), getProgress(), getFocusTrack()]);
      setProgressIndex(progressResult.kind === "ok" ? indexProgress(progressResult.entries) : null);
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
        {progressIndex === null && <ProgressUnavailableBanner />}
        <TrackDetail track={matched} progress={progressIndex} />
      </div>
    );
  }

  const allEntries = computeCardEntries(sortTracksByDisplayOrder(tracks));
  const availableEntries = pinFocusFirst(
    allEntries.filter((e) => e.lessonsReady > 0),
    focusTrack,
  );
  const comingSoon = allEntries.filter((e) => e.lessonsReady === 0);

  return (
    <div className="space-y-8">
      {progressIndex === null && <ProgressUnavailableBanner />}

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

      {availableEntries.length === 0 ? (
        <EmptyState
          icon={<BookIcon />}
          message="No lessons are imported yet — check back once a batch is written."
        />
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {availableEntries.map((entry) => (
            <TrackCard
              key={entry.track}
              stats={entry}
              readStats={trackReadStats(entry.tree, progressIndex)}
              isFocus={entry.track === focusTrack}
              onSetFocus={handleSetFocus}
              saving={pendingTrack === entry.track}
              busy={pendingTrack !== null}
            />
          ))}
        </div>
      )}

      {comingSoon.length > 0 && (
        <section className="space-y-2">
          <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Coming soon</h2>
          <div className="space-y-2">
            {comingSoon.map((entry) => (
              <ComingSoonRow key={entry.track} label={entry.label} />
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

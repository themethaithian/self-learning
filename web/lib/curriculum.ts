import type { Concept, ProgressEntry, ProgressState, Track } from "./api";
import { sortTracksByDisplayOrder } from "./trackMeta";

export function countAvailableLessons(concepts: Concept[]): { available: number; total: number } {
  return {
    available: concepts.filter((concept) => concept.has_lesson).length,
    total: concepts.length,
  };
}

export function pinFocusFirst<T extends { track: string }>(items: T[], focusTrack: string | null): T[] {
  if (!focusTrack) return items;
  const focusIndex = items.findIndex((s) => s.track === focusTrack);
  if (focusIndex === -1) return items;
  const focus = items[focusIndex];
  return [focus, ...items.slice(0, focusIndex), ...items.slice(focusIndex + 1)];
}

export interface LessonBreadcrumb {
  track: string;
  chapterTitle: string;
  conceptTitle: string;
}

export type NextLessonResult =
  | { kind: "next"; track: string; topic: string; concept: string; title: string }
  | { kind: "end-of-track" }
  | { kind: "not-found" };

interface FlatConcept {
  track: string;
  topic: string;
  chapterTitle: string;
  concept: string;
  title: string;
  hasLesson: boolean;
  estMinutes: number | null;
}

// The API already returns topics/chapters/concepts sorted by position, so
// flattening preserves reading order without re-sorting here.
function flattenTrack(track: Track): FlatConcept[] {
  const flat: FlatConcept[] = [];
  for (const topic of track.topics) {
    for (const chapter of topic.chapters) {
      for (const concept of chapter.concepts) {
        flat.push({
          track: track.track,
          topic: topic.slug,
          chapterTitle: chapter.title,
          concept: concept.slug,
          title: concept.title,
          hasLesson: concept.has_lesson,
          estMinutes: concept.est_minutes,
        });
      }
    }
  }
  return flat;
}

export function locateLessonBreadcrumb(tracks: Track[], topicSlug: string, conceptSlug: string): LessonBreadcrumb | null {
  for (const track of tracks) {
    const found = flattenTrack(track).find((c) => c.topic === topicSlug && c.concept === conceptSlug);
    if (found) return { track: found.track, chapterTitle: found.chapterTitle, conceptTitle: found.title };
  }
  return null;
}

export function findNextLesson(tracks: Track[], topicSlug: string, conceptSlug: string): NextLessonResult {
  for (const track of tracks) {
    const flat = flattenTrack(track);
    const fromIndex = flat.findIndex((c) => c.topic === topicSlug && c.concept === conceptSlug);
    if (fromIndex === -1) continue;
    const next = flat.slice(fromIndex + 1).find((c) => c.hasLesson);
    return next
      ? { kind: "next", track: next.track, topic: next.topic, concept: next.concept, title: next.title }
      : { kind: "end-of-track" };
  }
  return { kind: "not-found" };
}

export type ProgressByKey = Record<string, ProgressState>;

export function progressKey(topicSlug: string, conceptSlug: string): string {
  return `${topicSlug}::${conceptSlug}`;
}

export function indexProgress(entries: ProgressEntry[]): ProgressByKey {
  const index: ProgressByKey = {};
  for (const entry of entries) index[progressKey(entry.topic, entry.concept)] = entry.state;
  return index;
}

export interface TrackProgressStats {
  track: string;
  done: number;
  total: number;
}

export function computeTrackProgress(track: Track, progressByKey: ProgressByKey): TrackProgressStats {
  let total = 0;
  let done = 0;
  for (const item of flattenTrack(track)) {
    if (!item.hasLesson) continue;
    total += 1;
    if ((progressByKey[progressKey(item.topic, item.concept)] ?? "not_started") === "passed") done += 1;
  }
  return { track: track.track, total, done };
}

export type NextUpResult =
  | {
      kind: "next";
      track: string;
      topic: string;
      concept: string;
      chapterTitle: string;
      title: string;
      estMinutes: number | null;
      state: "not_started" | "in_progress";
    }
  | { kind: "all-done" }
  | { kind: "no-lessons" };

// "all-done" (every has_lesson concept anywhere is passed) and "no-lessons"
// (no track has a lesson at all) are both empty results but mean opposite
// things to the user — collapsing them into one null/done state would show
// "nice work, you finished everything" when really nothing has been
// imported yet.
export function pickNextUp(tracks: Track[], progressByKey: ProgressByKey, focusTrack: string | null): NextUpResult {
  const ordered = pinFocusFirst(sortTracksByDisplayOrder(tracks), focusTrack);
  let anyLessonExists = false;

  for (const track of ordered) {
    for (const item of flattenTrack(track)) {
      if (!item.hasLesson) continue;
      anyLessonExists = true;
      const state = progressByKey[progressKey(item.topic, item.concept)] ?? "not_started";
      if (state === "passed") continue;
      return {
        kind: "next",
        track: item.track,
        topic: item.topic,
        concept: item.concept,
        chapterTitle: item.chapterTitle,
        title: item.title,
        estMinutes: item.estMinutes,
        state,
      };
    }
  }

  return anyLessonExists ? { kind: "all-done" } : { kind: "no-lessons" };
}

// Extracted out of the page component so the dedupe rule (don't list the
// track pickNextUp already surfaced as Next Up) is covered by a plain unit
// test instead of living untested inside JSX.
export function buildOtherTracks(
  tracks: Track[],
  progressByKey: ProgressByKey,
  focusTrack: string | null,
  nextUp: NextUpResult,
): TrackProgressStats[] {
  return pinFocusFirst(sortTracksByDisplayOrder(tracks), focusTrack)
    .map((track) => computeTrackProgress(track, progressByKey))
    .filter((stats) => stats.total > 0)
    .filter((stats) => !(nextUp.kind === "next" && stats.track === nextUp.track));
}

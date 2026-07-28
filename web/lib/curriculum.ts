import type { Concept, Track } from "./api";
import type { TrackStats } from "../components/TrackCard";

export function countAvailableLessons(concepts: Concept[]): { available: number; total: number } {
  return {
    available: concepts.filter((concept) => concept.has_lesson).length,
    total: concepts.length,
  };
}

export function pinFocusFirst(stats: TrackStats[], focusTrack: string | null): TrackStats[] {
  if (!focusTrack) return stats;
  const focusIndex = stats.findIndex((s) => s.track === focusTrack);
  if (focusIndex === -1) return stats;
  const focus = stats[focusIndex];
  return [focus, ...stats.slice(0, focusIndex), ...stats.slice(focusIndex + 1)];
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

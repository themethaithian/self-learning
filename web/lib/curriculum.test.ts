import { describe, expect, it } from "vitest";
import type { Chapter, Concept, Topic, Track } from "./api";
import type { TrackStats } from "../components/TrackCard";
import { countAvailableLessons, findNextLesson, locateLessonBreadcrumb, pinFocusFirst } from "./curriculum";

function concept(slug: string, hasLesson: boolean, position: number): Concept {
  return { slug, title: `Concept ${slug}`, position, has_lesson: hasLesson, est_minutes: hasLesson ? 5 : null };
}

function chapter(slug: string, position: number, concepts: Concept[]): Chapter {
  return { slug, title: `Chapter ${slug}`, position, concepts };
}

function topic(slug: string, position: number, chapters: Chapter[]): Topic {
  return { slug, title: `Topic ${slug}`, position, chapters };
}

function track(name: string, topics: Topic[]): Track {
  return { track: name, topics };
}

describe("countAvailableLessons", () => {
  it("counts concepts with a lesson out of the total", () => {
    const concepts = [concept("a", true, 1), concept("b", false, 2), concept("c", true, 3)];
    expect(countAvailableLessons(concepts)).toEqual({ available: 2, total: 3 });
  });

  it("handles an empty list", () => {
    expect(countAvailableLessons([])).toEqual({ available: 0, total: 0 });
  });
});

describe("pinFocusFirst", () => {
  const stats: TrackStats[] = [
    { track: "ddd", label: "DDD", totalConcepts: 5, lessonsReady: 5, chapterCount: 2 },
    { track: "aws", label: "AWS", totalConcepts: 5, lessonsReady: 5, chapterCount: 2 },
    { track: "go", label: "Go", totalConcepts: 5, lessonsReady: 5, chapterCount: 2 },
  ];

  it("returns the list unchanged when there is no focus track", () => {
    expect(pinFocusFirst(stats, null)).toEqual(stats);
  });

  it("moves the focus track to the front, keeping the rest in order", () => {
    expect(pinFocusFirst(stats, "go").map((s) => s.track)).toEqual(["go", "ddd", "aws"]);
  });

  it("returns the list unchanged when the focus track isn't present", () => {
    expect(pinFocusFirst(stats, "dsa")).toEqual(stats);
  });
});

describe("locateLessonBreadcrumb", () => {
  const tracks = [track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1)])])])];

  it("finds the track and chapter title for a known concept", () => {
    expect(locateLessonBreadcrumb(tracks, "t1", "k1")).toEqual({
      track: "ddd",
      chapterTitle: "Chapter c1",
      conceptTitle: "Concept k1",
    });
  });

  it("returns null for a concept slug that isn't in the tree", () => {
    expect(locateLessonBreadcrumb(tracks, "t1", "missing")).toBeNull();
  });

  it("resolves to the requesting topic's chapter when two topics share a concept slug", () => {
    // Both topics contain a concept slugged "shared" in different chapters.
    // Matching on concept slug alone (dropping the topic check) would return
    // t1's "Chapter A" for every lookup, since Array.find stops at the first
    // hit in flatten order — this fixture only passes if topic is honored too.
    const sharedSlugTracks = [
      track("ddd", [
        topic("t1", 1, [chapter("chapter-a", 1, [concept("shared", true, 1)])]),
        topic("t2", 2, [chapter("chapter-b", 1, [concept("shared", true, 1)])]),
      ]),
    ];
    expect(locateLessonBreadcrumb(sharedSlugTracks, "t2", "shared")).toEqual({
      track: "ddd",
      chapterTitle: "Chapter chapter-b",
      conceptTitle: "Concept shared",
    });
  });
});

describe("findNextLesson", () => {
  it("returns the next concept in the same chapter when it has a lesson", () => {
    const tracks = [
      track("ddd", [
        topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1), concept("k2", true, 2), concept("k3", true, 3)])]),
      ]),
    ];
    expect(findNextLesson(tracks, "t1", "k1")).toEqual({
      kind: "next",
      track: "ddd",
      topic: "t1",
      concept: "k2",
      title: "Concept k2",
    });
  });

  it("crosses a chapter boundary to the first available concept of the next chapter", () => {
    const tracks = [
      track("ddd", [
        topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1)]), chapter("c2", 2, [concept("k2", true, 1)])]),
      ]),
    ];
    expect(findNextLesson(tracks, "t1", "k1")).toEqual({
      kind: "next",
      track: "ddd",
      topic: "t1",
      concept: "k2",
      title: "Concept k2",
    });
  });

  it("skips a run of several lesson-less concepts across chapter and topic boundaries", () => {
    const tracks = [
      track("ddd", [
        topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1), concept("k2", false, 2), concept("k3", false, 3)])]),
        topic("t2", 2, [chapter("c2", 1, [concept("k4", false, 1), concept("k5", true, 2)])]),
      ]),
    ];
    expect(findNextLesson(tracks, "t1", "k1")).toEqual({
      kind: "next",
      track: "ddd",
      topic: "t2",
      concept: "k5",
      title: "Concept k5",
    });
  });

  it("reaches end-of-track at the last available concept in the whole track", () => {
    const tracks = [track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1), concept("k2", false, 2)])])])];
    expect(findNextLesson(tracks, "t1", "k1")).toEqual({ kind: "end-of-track" });
  });

  it("reaches end-of-track when no concept in the track has a lesson", () => {
    const tracks = [
      track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", false, 1), concept("k2", false, 2)])])]),
    ];
    expect(findNextLesson(tracks, "t1", "k1")).toEqual({ kind: "end-of-track" });
  });

  it("returns not-found for a concept slug that isn't in the tree at all", () => {
    const tracks = [track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1)])])])];
    expect(findNextLesson(tracks, "unknown-topic", "unknown-concept")).toEqual({ kind: "not-found" });
  });

  it("stops at end-of-track rather than spilling into the next track's first lesson", () => {
    // track1's only lesson is its last concept; track2 starts with a lesson.
    // Looping per track (correct) finds nothing left in track1 -> end-of-track.
    // Flattening all tracks together (a real mutation this repo saw pass 13/13
    // before this fixture existed) would instead walk straight into track2.
    const tracks = [
      track("track1", [
        topic("t1", 1, [chapter("c1", 1, [concept("last-in-track1", true, 1), concept("no-lesson", false, 2)])]),
      ]),
      track("track2", [topic("t2", 1, [chapter("c2", 1, [concept("first-in-track2", true, 1)])])]),
    ];
    expect(findNextLesson(tracks, "t1", "last-in-track1")).toEqual({ kind: "end-of-track" });
  });
});

import { describe, expect, it } from "vitest";
import type { Chapter, Concept, ProgressEntry, ProgressState, Topic, Track } from "./api";
import type { TrackStats } from "../components/TrackCard";
import {
  computeTrackProgress,
  countAvailableLessons,
  findNextLesson,
  indexProgress,
  locateLessonBreadcrumb,
  pickNextUp,
  pinFocusFirst,
  progressKey,
  type ProgressByKey,
} from "./curriculum";

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

function progressMap(entries: Array<[topic: string, concept: string, state: ProgressState]>): ProgressByKey {
  const map: ProgressByKey = {};
  for (const [topicSlug, conceptSlug, state] of entries) map[progressKey(topicSlug, conceptSlug)] = state;
  return map;
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

describe("pickNextUp", () => {
  it("returns the focus track's next concept ahead of an earlier-ranked track", () => {
    const tracks = [
      track("go", [topic("t-go", 1, [chapter("c-go", 1, [concept("g1", true, 1)])])]),
      track("ddd", [topic("t-ddd", 1, [chapter("c-ddd", 1, [concept("d1", true, 1)])])]),
    ];
    expect(pickNextUp(tracks, {}, "go")).toEqual({
      kind: "next",
      track: "go",
      topic: "t-go",
      concept: "g1",
      chapterTitle: "Chapter c-go",
      title: "Concept g1",
      estMinutes: 5,
      state: "not_started",
    });
  });

  it("falls through to the next track in display order when the focus track is fully passed", () => {
    // Tracks are handed in scrambled order on purpose — a fallback that just
    // walked the input array (rather than sorting by TRACK_DISPLAY_ORDER)
    // would land on "aws" here, since it comes right after "ddd" in this array.
    const tracks = [
      track("aws", [topic("t-aws", 1, [chapter("c-aws", 1, [concept("a1", true, 1)])])]),
      track("distsys", [topic("t-dist", 1, [chapter("c-dist", 1, [concept("s1", true, 1)])])]),
      track("ddd", [topic("t-ddd", 1, [chapter("c-ddd", 1, [concept("d1", true, 1)])])]),
    ];
    const progress = progressMap([["t-ddd", "d1", "passed"]]);
    expect(pickNextUp(tracks, progress, "ddd")).toEqual({
      kind: "next",
      track: "distsys",
      topic: "t-dist",
      concept: "s1",
      chapterTitle: "Chapter c-dist",
      title: "Concept s1",
      estMinutes: 5,
      state: "not_started",
    });
  });

  it("starts from the first track in display order when no focus track is set", () => {
    const tracks = [
      track("go", [topic("t-go", 1, [chapter("c-go", 1, [concept("g1", true, 1)])])]),
      track("ddd", [topic("t-ddd", 1, [chapter("c-ddd", 1, [concept("d1", true, 1)])])]),
    ];
    expect(pickNextUp(tracks, {}, null)).toMatchObject({ kind: "next", track: "ddd", concept: "d1" });
  });

  it("returns an in_progress concept ahead of a later not_started one in the same track", () => {
    const tracks = [
      track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1), concept("k2", true, 2)])])]),
    ];
    const progress = progressMap([["t1", "k1", "in_progress"]]);
    expect(pickNextUp(tracks, progress, null)).toEqual({
      kind: "next",
      track: "ddd",
      topic: "t1",
      concept: "k1",
      chapterTitle: "Chapter c1",
      title: "Concept k1",
      estMinutes: 5,
      state: "in_progress",
    });
  });

  it("skips concepts without a lesson", () => {
    const tracks = [
      track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", false, 1), concept("k2", true, 2)])])]),
    ];
    expect(pickNextUp(tracks, {}, null)).toMatchObject({ kind: "next", concept: "k2" });
  });

  it("skips a track that has zero lessons entirely", () => {
    const tracks = [
      track("aws", [topic("t-aws", 1, [chapter("c-aws", 1, [concept("a1", false, 1), concept("a2", false, 2)])])]),
      track("ddd", [topic("t-ddd", 1, [chapter("c-ddd", 1, [concept("d1", true, 1)])])]),
    ];
    expect(pickNextUp(tracks, {}, null)).toMatchObject({ kind: "next", track: "ddd", concept: "d1" });
  });

  it("returns all-done when every lesson in every track is passed", () => {
    const tracks = [
      track("ddd", [topic("t-ddd", 1, [chapter("c-ddd", 1, [concept("d1", true, 1)])])]),
      track("go", [topic("t-go", 1, [chapter("c-go", 1, [concept("g1", true, 1)])])]),
    ];
    const progress = progressMap([
      ["t-ddd", "d1", "passed"],
      ["t-go", "g1", "passed"],
    ]);
    expect(pickNextUp(tracks, progress, null)).toEqual({ kind: "all-done" });
  });

  it("returns no-lessons when no track has any lesson at all", () => {
    const tracks = [track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", false, 1)])])])];
    expect(pickNextUp(tracks, {}, null)).toEqual({ kind: "no-lessons" });
  });

  it("treats every concept as not_started when the progress map is empty", () => {
    const tracks = [track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1)])])])];
    expect(pickNextUp(tracks, {}, null)).toMatchObject({ state: "not_started" });
  });

  it("still includes a track absent from TRACK_DISPLAY_ORDER instead of dropping it", () => {
    const tracks = [
      track("ddd", [topic("t-ddd", 1, [chapter("c-ddd", 1, [concept("d1", true, 1)])])]),
      track("mystery-track", [topic("t-mystery", 1, [chapter("c-mystery", 1, [concept("m1", true, 1)])])]),
    ];
    const progress = progressMap([["t-ddd", "d1", "passed"]]);
    expect(pickNextUp(tracks, progress, null)).toMatchObject({ kind: "next", track: "mystery-track", concept: "m1" });
  });
});

describe("computeTrackProgress", () => {
  it("counts passed lessons out of the lessons available in the track", () => {
    const t = track("ddd", [
      topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1), concept("k2", true, 2), concept("k3", false, 3)])]),
    ]);
    const progress = progressMap([["t1", "k1", "passed"]]);
    expect(computeTrackProgress(t, progress)).toEqual({ track: "ddd", done: 1, total: 2 });
  });

  it("returns zero total for a track with no lessons", () => {
    const t = track("ddd", [topic("t1", 1, [chapter("c1", 1, [concept("k1", false, 1)])])]);
    expect(computeTrackProgress(t, {})).toEqual({ track: "ddd", done: 0, total: 0 });
  });

  it("does not count an in_progress lesson as done", () => {
    const t = track("ddd", [
      topic("t1", 1, [chapter("c1", 1, [concept("k1", true, 1), concept("k2", true, 2), concept("k3", true, 3)])]),
    ]);
    const progress = progressMap([
      ["t1", "k1", "passed"],
      ["t1", "k2", "in_progress"],
    ]);
    expect(computeTrackProgress(t, progress)).toEqual({ track: "ddd", done: 1, total: 3 });
  });
});

describe("indexProgress", () => {
  it("indexes entries by topic and concept, preserving state", () => {
    const entries: ProgressEntry[] = [
      { topic: "t1", concept: "k1", state: "passed", last_read_at: null, first_passed_at: null },
      { topic: "t2", concept: "k2", state: "in_progress", last_read_at: null, first_passed_at: null },
    ];
    expect(indexProgress(entries)).toEqual({
      [progressKey("t1", "k1")]: "passed",
      [progressKey("t2", "k2")]: "in_progress",
    });
  });

  it("keys by topic then concept, not the reverse", () => {
    const entries: ProgressEntry[] = [
      { topic: "alpha", concept: "beta", state: "passed", last_read_at: null, first_passed_at: null },
    ];
    const index = indexProgress(entries);
    expect(index[progressKey("alpha", "beta")]).toBe("passed");
    expect(index[progressKey("beta", "alpha")]).toBeUndefined();
  });

  it("returns an empty map for an empty entry list", () => {
    expect(indexProgress([])).toEqual({});
  });
});

describe("progressKey", () => {
  it("cannot collide when characters shift between topic and concept, unlike plain concatenation", () => {
    expect(progressKey("ab", "c")).not.toBe(progressKey("a", "bc"));
  });
});

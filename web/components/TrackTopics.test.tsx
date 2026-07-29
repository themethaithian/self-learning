// @vitest-environment jsdom
import { describe, expect, it } from "vitest";
import { fireEvent, render, screen, within } from "@testing-library/react";
import { TrackTopics } from "./TrackTopics";
import type { Chapter, Concept, ProgressState, Topic } from "@/lib/api";
import { progressKey, type ProgressByKey } from "@/lib/curriculum";

function concept(slug: string, hasLesson: boolean, position: number): Concept {
  return { slug, title: `Concept ${slug}`, position, has_lesson: hasLesson, est_minutes: hasLesson ? 5 : null };
}
function chapter(slug: string, title: string, position: number, concepts: Concept[]): Chapter {
  return { slug, title, position, concepts };
}

// 2 chapters, mixing all three states plus a lesson-less concept, matching
// the same shape as the curriculum.ts fixtures — a single-chapter fixture
// would let a chapter/track-level mix-up slip through unnoticed.
const topics: Topic[] = [
  {
    slug: "t1",
    title: "Topic",
    position: 1,
    chapters: [
      chapter("c1", "Model-Driven Foundations", 1, [
        concept("k1", true, 1),
        concept("k2", true, 2),
        concept("k3", true, 3),
      ]),
      chapter("c2", "Building Blocks", 2, [concept("k4", true, 1), concept("k5", false, 2)]),
      chapter("c3", "Aggregates", 3, [concept("k6", false, 1), concept("k7", false, 2)]),
    ],
  },
];

function progressMap(entries: Array<[string, string, ProgressState]>): ProgressByKey {
  const map: ProgressByKey = {};
  for (const [topic, conceptSlug, state] of entries) map[progressKey(topic, conceptSlug)] = state;
  return map;
}

function openChapter(name: RegExp) {
  const button = screen.getByRole("button", { name });
  fireEvent.click(button);
  return button;
}

describe("TrackTopics — accordion open/close", () => {
  it("hides concept rows (and their markers) from the accessibility tree until the chapter button is clicked, then reveals them", () => {
    // getByText matches raw text content regardless of visibility, so it
    // can't tell a real browser's collapsed accordion from a broken one —
    // getByRole is the query that actually respects the hidden attribute
    // (dom-accessibility-api excludes hidden subtrees from role queries).
    const progress = progressMap([["t1", "k1", "passed"]]);
    render(<TrackTopics topics={topics} progress={progress} />);

    expect(screen.queryAllByRole("img")).toHaveLength(0);

    openChapter(/Model-Driven Foundations/);

    expect(screen.getAllByRole("img")).toHaveLength(1);
  });

  it("flips aria-expanded from false to true when opened", () => {
    render(<TrackTopics topics={topics} progress={null} />);
    const button = screen.getByRole("button", { name: /Model-Driven Foundations/ });
    expect(button.getAttribute("aria-expanded")).toBe("false");
    fireEvent.click(button);
    expect(button.getAttribute("aria-expanded")).toBe("true");
  });

  it("keeps the concept list mounted (just hidden) while closed, so aria-controls always points at a real element", () => {
    render(<TrackTopics topics={topics} progress={null} />);
    const button = screen.getByRole("button", { name: /Model-Driven Foundations/ });
    const listId = button.getAttribute("aria-controls");
    const list = listId ? document.getElementById(listId) : null;
    expect(list).not.toBeNull();
    expect((list as HTMLElement).hidden).toBe(true);
    // The 3 concept <li>s stay in the DOM while closed too — visibility is
    // the hidden attribute's job, not conditional rendering of the rows.
    expect(list!.children).toHaveLength(3);
    fireEvent.click(button);
    expect((list as HTMLElement).hidden).toBe(false);
    expect(list!.children).toHaveLength(3);
  });

  it("still gives every concept row a marker slot even when its state is not_started (S7 spacer)", () => {
    // No progress entries at all — every concept in this chapter is
    // not_started, so none gets a visible marker, but each row still needs
    // its slot reserved or the titles go ragged. Scoped to this chapter's
    // own list — every other chapter on the page stays mounted-but-hidden
    // and would otherwise inflate the count.
    render(<TrackTopics topics={topics} progress={{}} />);
    const button = openChapter(/Model-Driven Foundations/);
    const list = document.getElementById(button.getAttribute("aria-controls")!)!;
    expect(within(list).getAllByTestId("concept-marker-slot")).toHaveLength(3);
  });
});

describe("TrackTopics", () => {
  it("marks a passed concept and an in_progress concept with distinct, named markers", () => {
    const progress = progressMap([
      ["t1", "k1", "passed"],
      ["t1", "k2", "in_progress"],
    ]);
    render(<TrackTopics topics={topics} progress={progress} />);
    openChapter(/Model-Driven Foundations/);

    expect(screen.getByRole("img", { name: "Read" })).toBeTruthy();
    expect(screen.getByRole("img", { name: "In progress" })).toBeTruthy();
    // Only k1 and k2 have a state worth marking; k3 (not_started) gets none.
    expect(screen.getAllByRole("img")).toHaveLength(2);
  });

  it("never swaps which marker means passed and which means in_progress", () => {
    const progress = progressMap([["t1", "k1", "in_progress"]]);
    render(<TrackTopics topics={topics} progress={progress} />);
    openChapter(/Model-Driven Foundations/);

    expect(screen.queryByRole("img", { name: "Read" })).toBeNull();
    expect(screen.getByRole("img", { name: "In progress" })).toBeTruthy();
  });

  it("renders no markers anywhere when progress hasn't loaded, even for a concept that was passed before", () => {
    render(<TrackTopics topics={topics} progress={null} />);
    openChapter(/Model-Driven Foundations/);
    expect(screen.queryAllByRole("img")).toHaveLength(0);
  });

  it("shows the chapter read count as 'read', never relabelled as 'ready'", () => {
    const progress = progressMap([["t1", "k1", "passed"]]);
    render(<TrackTopics topics={topics} progress={progress} />);
    const button = openChapter(/Model-Driven Foundations/);
    expect(within(button).getByText("1/3 read")).toBeTruthy();
    expect(within(button).queryByText(/ready/)).toBeNull();
  });

  it("falls back to the old n/N ready header when progress hasn't loaded", () => {
    render(<TrackTopics topics={topics} progress={null} />);
    const button = openChapter(/Model-Driven Foundations/);
    expect(within(button).getByText("3/3 ready")).toBeTruthy();
    expect(within(button).queryByText(/read$/)).toBeNull();
  });

  it("shows only the ready fallback, never a vacuous 0/0 read, for a chapter with zero lessons", () => {
    const progress = progressMap([["t1", "k1", "passed"]]);
    render(<TrackTopics topics={topics} progress={progress} />);
    const button = openChapter(/Aggregates/);
    expect(within(button).getByText("0/2 ready")).toBeTruthy();
    expect(within(button).queryByText(/\bread\b/)).toBeNull();
  });

  it("shows both the read line and the ready line when a chapter has lesson-less concepts and progress is loaded", () => {
    // c2 has withLesson=1, planned=2 (k5 has no lesson) — a real gap, so
    // both ratios should be visible with clearly different labels.
    const progress = progressMap([["t1", "k4", "passed"]]);
    render(<TrackTopics topics={topics} progress={progress} />);
    const button = openChapter(/Building Blocks/);
    expect(within(button).getByText("1/1 read")).toBeTruthy();
    expect(within(button).getByText("1/2 ready")).toBeTruthy();
  });
});

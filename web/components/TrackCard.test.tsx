// @vitest-environment jsdom
import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { TrackCard, type TrackStats } from "./TrackCard";
import type { TrackReadStats } from "@/lib/curriculum";

const baseStats: TrackStats = {
  track: "ddd",
  label: "Domain-Driven Design",
  totalConcepts: 33,
  lessonsReady: 4,
  chapterCount: 5,
};

function renderCard(readStats: TrackReadStats) {
  return render(
    <TrackCard stats={baseStats} readStats={readStats} isFocus={false} onSetFocus={vi.fn()} saving={false} busy={false} />,
  );
}

describe("TrackCard", () => {
  it("shows a progress bar and the read count when progress is available", () => {
    renderCard({ read: 1, available: 4 });
    expect(screen.getByRole("progressbar")).toBeTruthy();
    expect(screen.getByText(/^1\/4 lessons read/)).toBeTruthy();
  });

  it("never swaps read and available in the printed count", () => {
    renderCard({ read: 1, available: 4 });
    expect(screen.queryByText(/^4\/1 lessons read/)).toBeNull();
  });

  it("hides the bar and the read count entirely when progress hasn't loaded", () => {
    renderCard({ read: null, available: 4 });
    expect(screen.queryByRole("progressbar")).toBeNull();
    expect(screen.queryByText(/\d+\/\d+ lessons? read/)).toBeNull();
    expect(screen.getByText("4 lessons ready")).toBeTruthy();
  });

  it("surfaces how many more concepts are planned beyond what's readable today", () => {
    renderCard({ read: 4, available: 4 });
    expect(screen.getByText(/4\/4 lessons read/)).toBeTruthy();
    expect(screen.getByText(/29 more planned/)).toBeTruthy();
  });

  it("omits the 'more planned' clause once every planned concept has a lesson", () => {
    render(
      <TrackCard
        stats={{ ...baseStats, totalConcepts: 4 }}
        readStats={{ read: 4, available: 4 }}
        isFocus={false}
        onSetFocus={vi.fn()}
        saving={false}
        busy={false}
      />,
    );
    expect(screen.getByText("4/4 lessons read")).toBeTruthy();
    expect(screen.queryByText(/more planned/)).toBeNull();
  });

  it("a full bar still shows the planned-concepts line so it never reads as curriculum-finished", () => {
    renderCard({ read: 4, available: 4 });
    const bar = screen.getByRole("progressbar");
    expect(bar.getAttribute("aria-valuenow")).toBe("100");
    expect(screen.getByText(/33 concepts planned/)).toBeTruthy();
  });
});

import { describe, expect, it } from "vitest";
import { sortTracksByDisplayOrder } from "./trackMeta";

describe("sortTracksByDisplayOrder", () => {
  it("orders known tracks per TRACK_DISPLAY_ORDER regardless of input order", () => {
    const items = [{ track: "go" }, { track: "ddd" }, { track: "aws" }];
    expect(sortTracksByDisplayOrder(items).map((i) => i.track)).toEqual(["ddd", "aws", "go"]);
  });

  it("appends a track absent from TRACK_DISPLAY_ORDER to the end, not the start", () => {
    const items = [{ track: "mystery-track" }, { track: "go" }, { track: "ddd" }];
    expect(sortTracksByDisplayOrder(items).map((i) => i.track)).toEqual(["ddd", "go", "mystery-track"]);
  });

  it("does not mutate the input array", () => {
    const items = [{ track: "go" }, { track: "ddd" }];
    sortTracksByDisplayOrder(items);
    expect(items.map((i) => i.track)).toEqual(["go", "ddd"]);
  });
});

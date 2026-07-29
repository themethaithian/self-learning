// @vitest-environment jsdom
import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { ProgressBar } from "./ProgressBar";

describe("ProgressBar", () => {
  it("gets its accessible name from the required label prop", () => {
    render(<ProgressBar read={2} available={4} label="Domain-Driven Design reading progress" />);
    expect(screen.getByRole("progressbar", { name: "Domain-Driven Design reading progress" })).toBeTruthy();
  });

  it("computes the percentage from read/available itself", () => {
    render(<ProgressBar read={2} available={4} label="x" />);
    expect(screen.getByRole("progressbar").getAttribute("aria-valuenow")).toBe("50");
  });

  it("never divides by zero when available is 0", () => {
    render(<ProgressBar read={0} available={0} label="x" />);
    expect(screen.getByRole("progressbar").getAttribute("aria-valuenow")).toBe("0");
  });

  it("clamps instead of exceeding 100 when read is somehow greater than available", () => {
    render(<ProgressBar read={10} available={4} label="x" />);
    expect(screen.getByRole("progressbar").getAttribute("aria-valuenow")).toBe("100");
  });
});

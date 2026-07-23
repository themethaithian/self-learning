"use client";

import { useEffect, useRef, useState } from "react";

let mermaidInitialized = false;

export function Mermaid({ code }: { code: string }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function render() {
      try {
        const { default: mermaid } = await import("mermaid");
        if (!mermaidInitialized) {
          mermaid.initialize({
            startOnLoad: false,
            theme: "neutral",
            securityLevel: "strict",
            flowchart: { useMaxWidth: false },
            sequence: { useMaxWidth: false },
            state: { useMaxWidth: false },
            er: { useMaxWidth: false },
            class: { useMaxWidth: false },
            gantt: { useMaxWidth: false },
            pie: { useMaxWidth: false },
            journey: { useMaxWidth: false },
            gitGraph: { useMaxWidth: false },
            mindmap: { useMaxWidth: false },
            timeline: { useMaxWidth: false },
            c4: { useMaxWidth: false },
            sankey: { useMaxWidth: false },
            xyChart: { useMaxWidth: false },
            quadrantChart: { useMaxWidth: false },
            requirement: { useMaxWidth: false },
          });
          mermaidInitialized = true;
        }
        const id = `mermaid-${Math.random().toString(36).slice(2)}`;
        const { svg } = await mermaid.render(id, code);
        if (!cancelled && containerRef.current) {
          containerRef.current.innerHTML = svg;
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Failed to render diagram");
        }
      }
    }

    render();
    return () => {
      cancelled = true;
    };
  }, [code]);

  if (error) {
    return (
      <pre className="my-4 overflow-x-auto rounded-xl border border-subtle bg-page p-4 font-mono text-xs text-danger">
        Diagram render failed: {error}
      </pre>
    );
  }

  return (
    <div className="my-4 overflow-x-auto rounded-xl border border-subtle bg-page p-4">
      <div ref={containerRef} className="w-fit mx-auto" />
    </div>
  );
}

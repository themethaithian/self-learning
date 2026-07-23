"use client";

import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { Components } from "react-markdown";
import type { Element, Text } from "hast";
import { Mermaid } from "@/components/Mermaid";

function hastText(node: Element | undefined): string {
  if (!node) return "";
  return node.children
    .map((child) => (child.type === "text" ? (child as Text).value : hastText(child as Element)))
    .join("");
}

function codeLanguage(preNode: Element | undefined): string | undefined {
  const codeNode = preNode?.children.find((child): child is Element => child.type === "element" && child.tagName === "code");
  const classNames = codeNode?.properties?.className;
  const classes = Array.isArray(classNames) ? classNames.map(String) : [];
  return classes.find((c) => c.startsWith("language-"))?.replace("language-", "");
}

const components: Components = {
  h1: (props) => <h1 className="mt-8 text-2xl font-semibold text-heading first:mt-0" {...props} />,
  h2: (props) => <h2 className="mt-8 text-xl font-semibold text-heading" {...props} />,
  h3: (props) => <h3 className="mt-6 text-lg font-semibold text-heading" {...props} />,
  p: (props) => <p className="mt-4 leading-[1.8] text-body first:mt-0" {...props} />,
  ul: (props) => <ul className="mt-4 list-disc space-y-2 pl-6 text-body" {...props} />,
  ol: (props) => <ol className="mt-4 list-decimal space-y-2 pl-6 text-body" {...props} />,
  li: (props) => <li className="leading-[1.8]" {...props} />,
  blockquote: (props) => <blockquote className="mt-4 border-l-2 border-accent/40 pl-4 text-muted italic" {...props} />,
  a: (props) => <a className="text-accent underline underline-offset-2 hover:text-accent-hover" {...props} />,
  strong: (props) => <strong className="font-semibold text-heading" {...props} />,
  table: (props) => (
    <div className="mt-4 overflow-x-auto">
      <table className="w-full border-collapse text-sm" {...props} />
    </div>
  ),
  thead: (props) => (
    <thead className="border-b border-subtle text-left text-xs font-medium uppercase tracking-wide text-faint" {...props} />
  ),
  th: (props) => <th className="px-3 py-2" {...props} />,
  td: (props) => <td className="border-b border-subtle px-3 py-2 align-top" {...props} />,
  pre: ({ node, children, ...props }) => {
    const language = codeLanguage(node);
    if (language === "mermaid") {
      return <Mermaid code={hastText(node).trim()} />;
    }
    return (
      <pre className="mt-4 overflow-x-auto rounded-xl border border-subtle bg-page p-4 font-mono text-sm text-body" {...props}>
        {children}
      </pre>
    );
  },
  code: ({ className, children, ...props }) => (
    <code className={`rounded bg-page px-1.5 py-0.5 font-mono text-[0.9em] text-body ${className ?? ""}`} {...props}>
      {children}
    </code>
  ),
};

export function LessonBody({ markdown }: { markdown: string }) {
  return (
    <div className="max-w-[68ch] font-thai text-[17px] leading-[1.8] text-body">
      <Markdown remarkPlugins={[remarkGfm]} components={components}>
        {markdown}
      </Markdown>
    </div>
  );
}

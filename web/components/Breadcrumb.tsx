import Link from "next/link";

export interface Crumb {
  label: string;
  href?: string;
}

// Wraps onto multiple lines at narrow widths rather than truncating — a long
// chapter title loses meaning if cut mid-word, and the reader page has
// vertical room to spare above the lesson body.
export function Breadcrumb({ items }: { items: Crumb[] }) {
  return (
    <nav aria-label="Breadcrumb">
      <ol className="flex flex-wrap items-baseline gap-x-2 gap-y-1 text-sm">
        {items.map((item, index) => {
          const isLast = index === items.length - 1;
          return (
            <li key={index} className="flex min-w-0 items-baseline gap-x-2">
              {index > 0 && (
                <span aria-hidden className="text-faint">
                  ›
                </span>
              )}
              {item.href ? (
                // A link can never be the current page, so aria-current never
                // applies here — only isLast among the no-href crumbs below
                // decides that, which also keeps a single-crumb degraded
                // breadcrumb (always a link, e.g. "Learn") from being
                // mislabeled as the page the user is actually on.
                <Link
                  href={item.href}
                  className="rounded font-medium text-accent underline-offset-2 hover:text-accent-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
                >
                  {item.label}
                </Link>
              ) : (
                <span aria-current={isLast ? "page" : undefined} className={isLast ? "text-heading" : "text-muted"}>
                  {item.label}
                </span>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
}

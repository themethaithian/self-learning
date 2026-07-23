"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import type { ReactNode } from "react";
import { NAV_ITEMS } from "@/lib/nav";

function NavLink({ href, label, active }: { href: string; label: string; active: boolean }) {
  return (
    <Link
      href={href}
      className={`flex items-center rounded-full px-4 py-2 text-sm font-medium transition-colors duration-150 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2 ${
        active ? "bg-accent-tint text-accent-strong" : "text-muted hover:bg-page hover:text-body"
      }`}
    >
      {label}
    </Link>
  );
}

export function AppShell({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const activeItem = NAV_ITEMS.find((item) => pathname.startsWith(item.href));
  const title = activeItem?.label ?? "Self-Improve";

  return (
    <div className="min-h-screen bg-page md:pl-56">
      <aside className="fixed inset-y-0 left-0 hidden w-56 flex-col border-r border-subtle bg-surface p-4 md:flex">
        <div className="px-2 pb-6 text-sm font-semibold text-heading">Self-Improve</div>
        <nav className="flex flex-col gap-1">
          {NAV_ITEMS.map((item) => (
            <NavLink key={item.href} href={item.href} label={item.label} active={item === activeItem} />
          ))}
        </nav>
      </aside>

      <header className="sticky top-0 z-10 border-b border-subtle bg-surface/90 px-6 py-4 backdrop-blur">
        <h1 className="text-xl font-semibold text-heading">{title}</h1>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-8 pb-24 md:pb-8">{children}</main>

      <nav className="fixed inset-x-0 bottom-0 z-10 flex items-center justify-around border-t border-subtle bg-surface px-2 py-2 md:hidden">
        {NAV_ITEMS.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={`rounded-full px-3 py-1.5 text-xs font-medium transition-colors duration-150 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2 ${
              item === activeItem ? "bg-accent-tint text-accent-strong" : "text-muted"
            }`}
          >
            {item.label}
          </Link>
        ))}
      </nav>
    </div>
  );
}

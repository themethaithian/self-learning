export interface NavItem {
  href: string;
  label: string;
}

// Dashboard/Drill/DSA/Test/Tickets are placeholder routes (phase 4/5) — kept
// on disk but dropped from the nav so it never links to a page that's just
// "coming soon" (see docs/tickets/ux-today.md).
export const NAV_ITEMS: NavItem[] = [
  { href: "/today", label: "Today" },
  { href: "/learn", label: "Learn" },
];

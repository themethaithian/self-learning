export interface NavItem {
  href: string;
  label: string;
}

export const NAV_ITEMS: NavItem[] = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/learn", label: "Learn" },
  { href: "/drill", label: "Drill" },
  { href: "/dsa", label: "DSA" },
  { href: "/test", label: "Test" },
  { href: "/tickets", label: "Tickets" },
];

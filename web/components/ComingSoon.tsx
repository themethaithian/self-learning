import { EmptyState } from "@/components/EmptyState";
import { BookIcon } from "@/components/icons";

export function ComingSoon({ feature }: { feature: string }) {
  return <EmptyState icon={<BookIcon />} message={`${feature} is coming soon.`} />;
}

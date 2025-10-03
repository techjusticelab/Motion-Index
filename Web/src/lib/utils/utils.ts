import { format } from 'date-fns';

export function formatDate(dateString: string): string {
  if (!dateString) return 'N/A';
  try {
    return format(new Date(dateString), 'MMM d, yyyy');
  } catch (err) {
    return dateString;
  }
}

export function filterOptions(allOptions: string[] = [], searchInput: string): string[] {
  if (!searchInput) {
    return [...allOptions];
  }

  const searchLower = searchInput.toLowerCase();
  return allOptions.filter((item) => item.toLowerCase().includes(searchLower));
}
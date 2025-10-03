import { format } from 'date-fns';

export function formatDate(dateString: string): string {
  if (!dateString) return 'N/A';
  try {
    return format(new Date(dateString), 'MMM d, yyyy');
  } catch (err) {
    return dateString;
  }
}

export function formatDateTime(dateString: string): string {
  if (!dateString) return 'N/A';
  try {
    return format(new Date(dateString), 'MMM d, yyyy h:mm a');
  } catch (err) {
    return dateString;
  }
}

export function formatFileSize(bytes: number | string): string {
  if (!bytes || bytes === 0) return 'N/A';
  
  const size = typeof bytes === 'string' ? parseInt(bytes) : bytes;
  if (isNaN(size)) return 'N/A';
  
  const units = ['B', 'KB', 'MB', 'GB'];
  let unitIndex = 0;
  let fileSize = size;
  
  while (fileSize >= 1024 && unitIndex < units.length - 1) {
    fileSize /= 1024;
    unitIndex++;
  }
  
  return `${fileSize.toFixed(unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`;
}

export function formatPercentage(value: number | string): string {
  if (value === undefined || value === null) return 'N/A';
  
  const percent = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(percent)) return 'N/A';
  
  // If value is between 0 and 1, treat as decimal, otherwise as percentage
  const displayValue = percent <= 1 ? percent * 100 : percent;
  return `${displayValue.toFixed(1)}%`;
}

export function formatCount(count: number | string): string {
  if (!count || count === 0) return '0';
  
  const num = typeof count === 'string' ? parseInt(count) : count;
  if (isNaN(num)) return '0';
  
  return num.toLocaleString();
}

export function formatDocumentType(docType: string): string {
  if (!docType) return 'Unknown';
  
  // Convert snake_case to Title Case
  return docType
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

export function formatPersonName(person: any): string {
  if (typeof person === 'string') return person;
  if (typeof person === 'object' && person?.name) return person.name;
  return 'N/A';
}

export function formatCourtName(court: any): string {
  if (typeof court === 'string') return court;
  if (typeof court === 'object') {
    const parts = [];
    if (court.court_name) parts.push(court.court_name);
    if (court.jurisdiction) parts.push(`(${court.jurisdiction})`);
    return parts.join(' ') || 'N/A';
  }
  return 'N/A';
}

export function formatCaseInfo(caseInfo: any): { number?: string; name?: string; id?: string } {
  if (typeof caseInfo === 'object' && caseInfo) {
    return {
      number: caseInfo.case_number,
      name: caseInfo.case_name,
      id: caseInfo.case_id
    };
  }
  return {};
}

export function filterOptions(allOptions: string[] = [], searchInput: string): string[] {
  if (!searchInput) {
    return [...allOptions];
  }

  const searchLower = searchInput.toLowerCase();
  return allOptions.filter((item) => item.toLowerCase().includes(searchLower));
}

export function isEmptyValue(value: any): boolean {
  return value === undefined || 
         value === null || 
         value === '' || 
         (Array.isArray(value) && value.length === 0) ||
         (typeof value === 'object' && Object.keys(value).length === 0);
}
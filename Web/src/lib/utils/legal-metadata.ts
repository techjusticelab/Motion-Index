import type { Document, Judge, CourtInfo, Party, Attorney, Charge } from '$lib/types';

/**
 * Utility functions for extracting and formatting legal metadata
 * Specifically designed for public defender use cases
 */

/**
 * Extract judge name and title from Judge object or fallback to string
 */
export function extractJudgeName(judge: Judge | string | null | undefined): string {
	if (!judge) return '';
	
	if (typeof judge === 'string') {
		return judge;
	}
	
	if (typeof judge === 'object' && judge.name) {
		return judge.title ? `${judge.title} ${judge.name}` : judge.name;
	}
	
	return '';
}

/**
 * Extract court name and jurisdiction from CourtInfo object or fallback to string
 */
export function extractCourtInfo(court: CourtInfo | string | null | undefined): { name: string; jurisdiction: string } {
	if (!court) return { name: '', jurisdiction: '' };
	
	if (typeof court === 'string') {
		return { name: court, jurisdiction: '' };
	}
	
	if (typeof court === 'object') {
		const name = court.court_name || '';
		const jurisdiction = court.jurisdiction || court.district || court.county || '';
		return { name, jurisdiction };
	}
	
	return { name: '', jurisdiction: '' };
}

/**
 * Extract and format party information for display
 */
export function formatParties(parties: Party[] | null | undefined): string {
	if (!parties || !Array.isArray(parties) || parties.length === 0) {
		return '';
	}
	
	// Focus on defendant and prosecutor for public defenders
	const defendants = parties.filter(p => p.role?.toLowerCase().includes('defendant'));
	const prosecutors = parties.filter(p => p.role?.toLowerCase().includes('prosecutor') || p.role?.toLowerCase().includes('state'));
	
	const result = [];
	if (defendants.length > 0) {
		result.push(`Defendant: ${defendants[0].name}`);
	}
	if (prosecutors.length > 0) {
		result.push(`Prosecutor: ${prosecutors[0].name}`);
	}
	
	return result.join(' • ');
}

/**
 * Extract attorney information relevant to public defenders
 */
export function formatAttorneys(attorneys: Attorney[] | null | undefined): string {
	if (!attorneys || !Array.isArray(attorneys) || attorneys.length === 0) {
		return '';
	}
	
	// Focus on defense attorneys
	const defenseAttorneys = attorneys.filter(a => 
		a.role?.toLowerCase().includes('defense') || 
		a.role?.toLowerCase().includes('public defender')
	);
	
	if (defenseAttorneys.length > 0) {
		return `Defense: ${defenseAttorneys[0].name}`;
	}
	
	// If no defense attorney found, show first attorney
	return attorneys[0].name;
}

/**
 * Format charges for quick assessment
 */
export function formatCharges(charges: Charge[] | null | undefined): string[] {
	if (!charges || !Array.isArray(charges) || charges.length === 0) {
		return [];
	}
	
	return charges.slice(0, 3).map(charge => {
		const parts = [];
		if (charge.statute) parts.push(charge.statute);
		if (charge.grade) parts.push(charge.grade);
		if (charge.description) parts.push(charge.description.substring(0, 50) + (charge.description.length > 50 ? '...' : ''));
		return parts.join(' - ');
	});
}

/**
 * Determine document urgency/priority for public defenders
 */
export function getDocumentPriority(document: Document): 'high' | 'medium' | 'low' {
	const docType = document.doc_type?.toLowerCase() || '';
	const subject = document.metadata?.subject?.toLowerCase() || '';
	const status = document.metadata?.status?.toLowerCase() || '';
	
	// High priority: suppression motions, dismissals, favorable rulings
	if (
		docType.includes('motion_to_suppress') ||
		docType.includes('motion_to_dismiss') ||
		subject.includes('suppress') ||
		subject.includes('dismiss') ||
		status.includes('granted') ||
		status.includes('favorable')
	) {
		return 'high';
	}
	
	// Medium priority: other motions, orders
	if (
		docType.includes('motion') ||
		docType.includes('order') ||
		docType.includes('ruling')
	) {
		return 'medium';
	}
	
	return 'low';
}

/**
 * Calculate importance score for sorting documents (higher = more important)
 */
export function calculateDocumentImportanceScore(document: Document): number {
	let score = 0;
	
	const docType = document.doc_type?.toLowerCase() || '';
	const subject = document.metadata?.subject?.toLowerCase() || '';
	const status = document.metadata?.status?.toLowerCase() || '';
	const legalTags = document.metadata?.legal_tags?.map(tag => tag.toLowerCase()) || [];
	
	// Document type scoring (most important for public defenders)
	if (docType.includes('motion_to_suppress') || subject.includes('suppress')) {
		score += 1000; // Highest priority - suppression motions are crucial
	}
	if (docType.includes('motion_to_dismiss') || subject.includes('dismiss')) {
		score += 900; // Very high priority
	}
	if (docType.includes('motion_summary_judgment') || subject.includes('summary judgment')) {
		score += 800;
	}
	if (docType.includes('motion_in_limine') || subject.includes('in limine')) {
		score += 700;
	}
	if (docType.includes('motion_to_compel') || subject.includes('compel')) {
		score += 600;
	}
	if (docType.includes('order') || docType.includes('ruling')) {
		score += 500; // Orders and rulings are important for precedent
	}
	if (docType.includes('motion')) {
		score += 400; // Other motions
	}
	if (docType.includes('brief')) {
		score += 300;
	}
	if (docType.includes('transcript')) {
		score += 200;
	}
	
	// Status scoring (outcomes matter a lot)
	if (status.includes('granted') || status.includes('favorable')) {
		score += 500; // Favorable outcomes are very valuable
	}
	if (status.includes('denied') || status.includes('unfavorable')) {
		score += 300; // Still valuable to know what doesn't work
	}
	if (status.includes('pending')) {
		score += 100;
	}
	
	// Legal tags scoring
	const highValueTags = [
		'fourth amendment', 'suppression', 'search and seizure', 'miranda', 
		'due process', 'discovery', 'brady', 'prosecutorial misconduct',
		'ineffective assistance', 'sentencing', 'plea bargain'
	];
	
	legalTags.forEach(tag => {
		if (highValueTags.some(hvTag => tag.includes(hvTag))) {
			score += 200;
		}
	});
	
	// Charges scoring (serious charges get higher priority)
	const charges = document.metadata?.charges || [];
	charges.forEach(charge => {
		const description = charge.description?.toLowerCase() || '';
		const statute = charge.statute?.toLowerCase() || '';
		
		if (
			description.includes('felony') || 
			statute.includes('felony') ||
			charge.grade?.toLowerCase().includes('felony')
		) {
			score += 150;
		}
		if (
			description.includes('misdemeanor') || 
			statute.includes('misdemeanor') ||
			charge.grade?.toLowerCase().includes('misdemeanor')
		) {
			score += 75;
		}
	});
	
	// Court level scoring (higher courts = more precedential value)
	const courtInfo = extractCourtInfo(document.metadata?.court);
	if (courtInfo.name.toLowerCase().includes('supreme')) {
		score += 300;
	}
	if (courtInfo.name.toLowerCase().includes('appellate') || courtInfo.name.toLowerCase().includes('appeal')) {
		score += 200;
	}
	if (courtInfo.name.toLowerCase().includes('superior')) {
		score += 100;
	}
	
	// Recency bonus (more recent documents get slight boost)
	const relevantDate = getRelevantDate(document);
	const docDate = new Date(relevantDate.date);
	const now = new Date();
	const daysDiff = (now.getTime() - docDate.getTime()) / (1000 * 60 * 60 * 24);
	
	if (daysDiff < 30) {
		score += 50; // Very recent
	} else if (daysDiff < 90) {
		score += 30; // Recent
	} else if (daysDiff < 365) {
		score += 10; // Within a year
	}
	
	// Length penalty for very long documents (harder to review quickly)
	const pages = document.metadata?.pages || 0;
	if (pages > 100) {
		score -= 50;
	} else if (pages > 50) {
		score -= 20;
	}
	
	return score;
}

/**
 * Sort documents by importance for public defenders
 */
export function sortDocumentsByImportance(documents: Document[]): Document[] {
	return [...documents].sort((a, b) => {
		const scoreA = calculateDocumentImportanceScore(a);
		const scoreB = calculateDocumentImportanceScore(b);
		return scoreB - scoreA; // Sort descending (highest score first)
	});
}

/**
 * Get the most relevant date for display
 */
export function getRelevantDate(document: Document): { date: string; label: string } {
	const metadata = document.metadata;
	
	// Prioritize filing date or event date
	if (metadata?.filing_date) {
		return { date: metadata.filing_date, label: 'Filed' };
	}
	
	if (metadata?.event_date) {
		return { date: metadata.event_date, label: 'Event' };
	}
	
	if (metadata?.timestamp) {
		return { date: metadata.timestamp, label: 'Date' };
	}
	
	return { date: document.created_at, label: 'Created' };
}

/**
 * Get motion type/outcome status for visual indicators
 */
export function getMotionStatus(document: Document): { type: string; outcome?: string; color: string } {
	const docType = document.doc_type?.toLowerCase() || '';
	const status = document.metadata?.status?.toLowerCase() || '';
	const subject = document.metadata?.subject?.toLowerCase() || '';
	
	let type = 'Document';
	let outcome = undefined;
	let color = 'neutral';
	
	// Determine motion type
	if (docType.includes('motion_to_suppress') || subject.includes('suppress')) {
		type = 'Suppression Motion';
		color = 'blue';
	} else if (docType.includes('motion_to_dismiss') || subject.includes('dismiss')) {
		type = 'Motion to Dismiss';
		color = 'purple';
	} else if (docType.includes('motion')) {
		type = 'Motion';
		color = 'indigo';
	} else if (docType.includes('order')) {
		type = 'Order';
		color = 'green';
	} else if (docType.includes('ruling')) {
		type = 'Ruling';
		color = 'green';
	}
	
	// Determine outcome
	if (status.includes('granted') || status.includes('favorable')) {
		outcome = 'Granted';
		color = 'green';
	} else if (status.includes('denied') || status.includes('unfavorable')) {
		outcome = 'Denied';
		color = 'red';
	} else if (status.includes('pending')) {
		outcome = 'Pending';
		color = 'yellow';
	}
	
	return { type, outcome, color };
}

/**
 * Format case information for display
 */
export function formatCaseInfo(document: Document): { name: string; number: string } {
	const metadata = document.metadata;
	
	// Try structured case info first
	if (metadata?.case) {
		return {
			name: metadata.case.case_name || '',
			number: metadata.case.case_number || ''
		};
	}
	
	// Fall back to legacy fields
	return {
		name: metadata?.case_name || '',
		number: metadata?.case_number || ''
	};
}
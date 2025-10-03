/**
 * Network connectivity and Supabase service checker
 * Helps diagnose authentication connection issues
 */

export type NetworkCheckResult = {
    isOnline: boolean;
    supabaseReachable: boolean;
    possibleIssues: string[];
    recommendations: string[];
};

/**
 * Check basic network connectivity and Supabase accessibility
 */
export async function checkNetworkConnectivity(supabaseUrl: string): Promise<NetworkCheckResult> {
    const result: NetworkCheckResult = {
        isOnline: navigator.onLine,
        supabaseReachable: false,
        possibleIssues: [],
        recommendations: []
    };

    // Check if browser reports offline status
    if (!result.isOnline) {
        result.possibleIssues.push('Browser reports offline status');
        result.recommendations.push('Check your internet connection');
        return result;
    }

    try {
        // Try to reach Supabase health endpoint
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 5000);

        const response = await fetch(`${supabaseUrl}/rest/v1/`, {
            method: 'HEAD',
            signal: controller.signal,
            mode: 'cors'
        });

        clearTimeout(timeoutId);
        
        if (response.ok || response.status === 401 || response.status === 403) {
            // 401/403 means we reached Supabase but auth failed (expected)
            result.supabaseReachable = true;
        } else {
            result.possibleIssues.push(`Supabase returned status ${response.status}`);
        }
    } catch (error) {
        result.supabaseReachable = false;
        
        if (error instanceof Error) {
            const errorMessage = error.message.toLowerCase();
            
            if (errorMessage.includes('aborted') || errorMessage.includes('timeout')) {
                result.possibleIssues.push('Connection timeout to Supabase');
                result.recommendations.push('Check your internet speed and stability');
            } else if (errorMessage.includes('blocked') || errorMessage.includes('cors')) {
                result.possibleIssues.push('Request blocked (likely by ad blocker or browser extension)');
                result.recommendations.push('Disable ad blockers and privacy extensions');
                result.recommendations.push('Try incognito/private browsing mode');
            } else if (errorMessage.includes('network') || errorMessage.includes('fetch')) {
                result.possibleIssues.push('Network error reaching Supabase');
                result.recommendations.push('Check firewall and network restrictions');
                result.recommendations.push('Verify Supabase URLs are not blocked');
            } else {
                result.possibleIssues.push(`Unknown error: ${error.message}`);
            }
        }
    }

    // Add general recommendations if Supabase is not reachable
    if (!result.supabaseReachable) {
        result.recommendations.push('Add *.supabase.co to your allowlist');
        result.recommendations.push('Try a different browser or device');
        result.recommendations.push('Contact your network administrator if on corporate network');
    }

    return result;
}

/**
 * Get user-friendly error message based on network check results
 */
export function getNetworkErrorMessage(checkResult: NetworkCheckResult): string {
    if (!checkResult.isOnline) {
        return 'You appear to be offline. Please check your internet connection.';
    }

    if (!checkResult.supabaseReachable) {
        const issues = checkResult.possibleIssues.join(', ');
        const solutions = checkResult.recommendations
            .map(rec => `• ${rec}`)
            .join('\n');
        
        return `Cannot reach authentication service.\n\nPossible causes: ${issues}\n\nTry these solutions:\n${solutions}`;
    }

    return 'Network connectivity appears normal.';
}

/**
 * Browser-specific detection for common ad blockers
 */
export function detectAdBlockers(): string[] {
    const detectedBlockers: string[] = [];
    
    // Check for common ad blocker indicators
    if (typeof window !== 'undefined') {
        // uBlock Origin
        if (window.navigator.userAgent.includes('uBlock') || 
            document.querySelector('script[src*="ublock"]')) {
            detectedBlockers.push('uBlock Origin');
        }
        
        // AdBlock Plus
        if (window.navigator.userAgent.includes('Adblock') ||
            document.querySelector('link[href*="adblock"]')) {
            detectedBlockers.push('AdBlock Plus');
        }
        
        // Check for blocked resource indicators
        const blockedElements = document.querySelectorAll('[style*="display: none"]');
        if (blockedElements.length > 10) {
            detectedBlockers.push('Unknown ad blocker (detected blocked elements)');
        }
    }
    
    return detectedBlockers;
}
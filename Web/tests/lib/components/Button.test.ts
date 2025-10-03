import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import Button from '$lib/components/ui/Button.svelte';

describe('Button Component', () => {
	it('renders with default props', () => {
		const { getByRole } = render(Button, {
			props: { children: 'Click me' }
		});
		
		const button = getByRole('button');
		expect(button).toBeInTheDocument();
		expect(button).toHaveTextContent('Click me');
	});

	it('handles click events', async () => {
		const handleClick = vi.fn();
		const { getByRole } = render(Button, {
			props: { 
				children: 'Click me',
				onclick: handleClick
			}
		});
		
		const button = getByRole('button');
		await fireEvent.click(button);
		
		expect(handleClick).toHaveBeenCalledTimes(1);
	});

	it('can be disabled', () => {
		const { getByRole } = render(Button, {
			props: { 
				children: 'Click me',
				disabled: true
			}
		});
		
		const button = getByRole('button');
		expect(button).toBeDisabled();
	});

	it('shows loading state', () => {
		const { getByRole } = render(Button, {
			props: { 
				children: 'Click me',
				loading: true
			}
		});
		
		const button = getByRole('button');
		expect(button).toBeDisabled();
		// Could add more specific loading indicator checks here
	});
});
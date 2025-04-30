<script lang="ts">
	import { page } from '$app/stores';

	const email = $page.url.searchParams.get('email') || '';
</script>

<div class="flex min-h-screen flex-col justify-center bg-gray-50 py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<div class="bg-white px-4 py-8 text-center shadow sm:rounded-lg sm:px-10">
			<div
				class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-100"
			>
				<svg
					class="h-6 w-6 text-green-600"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M5 13l4 4L19 7"
					/>
				</svg>
			</div>
			<h2 class="mb-2 text-2xl font-bold text-gray-900">Check your email</h2>
			<p class="mb-6 text-gray-600">
				We've sent a confirmation link to <strong>{email}</strong>.<br />
				Please check your email and click the link to verify your account.
			</p>
			<div class="mt-4">
				<a href="/auth/login" class="font-medium text-indigo-600 hover:text-indigo-500">
					Return to sign in
				</a>
			</div>
		</div>
	</div>
</div>

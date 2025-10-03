<script lang="ts">
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';
</script>

<div
	class="min-h-screen bg-neutral-50 px-4 py-12 sm:px-6 lg:px-8"
	in:fade={{ duration: 600, easing: cubicOut }}
	out:fade
>
	<div class="mx-auto max-w-3xl" in:fly={{ y: 20, duration: 700, easing: cubicOut }}>
		<div class="text-center" in:fly={{ y: -10, duration: 700, delay: 200, easing: cubicOut }}>
			<h1
				class="text-3xl font-extrabold text-neutral-900 sm:text-4xl"
				in:slide={{ duration: 600, delay: 300 }}
			>
				Help & Support
			</h1>
			<p class="mt-3 text-lg text-neutral-500" in:fade={{ duration: 600, delay: 400 }}>
				Get assistance with Motion Index or learn how you can contribute to the project
			</p>
		</div>

		<div class="mt-12 grid gap-8 lg:grid-cols-2" in:fade={{ duration: 500, delay: 500 }}>
			<!-- Contact Information -->
			<div
				class="rounded-lg bg-white p-8 shadow"
				in:fly={{ x: -20, y: 10, duration: 700, delay: 600, easing: cubicOut }}
			>
				<div
					class="flex items-center justify-center"
					in:scale={{ start: 0.95, duration: 600, delay: 700, easing: elasticOut }}
				>
					<div class="rounded-full bg-primary-100 p-3">
						<svg
							class="h-6 w-6 text-primary-600"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
							/>
						</svg>
					</div>
				</div>
				<h2
					class="mt-4 text-center text-xl font-bold text-neutral-900"
					in:slide={{ duration: 500, delay: 800 }}
				>
					Contact Us
				</h2>
				<p class="mt-2 text-center text-neutral-500" in:slide={{ duration: 500, delay: 900 }}>
					Need help with Motion Index? Our team is here to assist you.
				</p>
				<div
					class="mt-6 space-y-4"
					in:fly={{ y: 10, duration: 600, delay: 1000, easing: cubicOut }}
				>
					<div class="flex items-center justify-center">
						<a
							href="mailto:info@techjusticelab.org"
							class="rounded-md bg-primary-50 px-4 py-2 text-primary-700 hover:bg-primary-100"
						>
							info@techjusticelab.org
						</a>
					</div>
					<div
						class="flex items-center justify-center"
						in:fly={{ y: 10, duration: 600, delay: 1100, easing: cubicOut }}
					>
						<a
							href="https://techjusticelab.org"
							target="_blank"
							rel="noopener noreferrer"
							class="rounded-md bg-primary-50 px-4 py-2 text-primary-700 hover:bg-primary-100"
						>
							techjusticelab.org
						</a>
					</div>
				</div>
			</div>

			<!-- GitHub Contribution -->
			<div
				class="rounded-lg bg-white p-8 shadow"
				in:fly={{ x: 20, y: 10, duration: 700, delay: 600, easing: cubicOut }}
			>
				<div
					class="flex items-center justify-center"
					in:scale={{ start: 0.95, duration: 600, delay: 700, easing: elasticOut }}
				>
					<div class="rounded-full bg-primary-100 p-3">
						<svg
							class="h-6 w-6 text-primary-600"
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="currentColor"
						>
							<path
								fill-rule="evenodd"
								d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
								clip-rule="evenodd"
							/>
						</svg>
					</div>
				</div>
				<h2
					class="mt-4 text-center text-xl font-bold text-neutral-900"
					in:slide={{ duration: 500, delay: 800 }}
				>
					Contribute to Motion Index
				</h2>
				<p class="mt-2 text-center text-neutral-500" in:slide={{ duration: 500, delay: 900 }}>
					Motion Index is an open-source project. We welcome contributions from developers of all
					skill levels.
				</p>
				<div
					class="mt-6 space-y-4"
					in:fly={{ y: 10, duration: 600, delay: 1000, easing: cubicOut }}
				>
					<a
						href="https://github.com/TheArctesian/Motion-Index"
						target="_blank"
						rel="noopener noreferrer"
						class="flex w-full items-center justify-center rounded-md border border-transparent bg-primary-600 px-6 py-3 text-base font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-900 focus:ring-offset-2"
						in:scale={{ start: 0.97, duration: 600, delay: 1100, easing: backOut }}
					>
						<svg
							class="mr-2 h-5 w-5"
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="currentColor"
						>
							<path
								fill-rule="evenodd"
								d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
								clip-rule="evenodd"
							/>
						</svg>
						View on GitHub
					</a>
				</div>
			</div>
		</div>

		<!-- FAQ Section -->
		<div class="mt-16" in:fly={{ y: 20, duration: 700, delay: 1200, easing: cubicOut }}>
			<h2
				class="text-center text-2xl font-bold text-neutral-900"
				in:slide={{ duration: 500, delay: 1300 }}
			>
				Frequently Asked Questions
			</h2>
			<div class="mt-8 space-y-6">
				{#each [{ question: 'What is Motion Index?', answer: 'Motion Index is a legal document repository and search tool designed to help legal professionals organize, categorize, and efficiently find case-related documents.' }, { question: 'How can I contribute to the project?', answer: 'You can contribute by forking the repository on GitHub, making improvements, and submitting pull requests. We welcome contributions to code, documentation, and testing.' }, { question: 'Is Motion Index open source?', answer: 'Yes, Motion Index is an open-source project developed by Tech Justice Lab. The code is available on GitHub and we encourage community contributions.' }, { question: 'How do I report bugs or request features?', answer: 'You can report bugs or request features by opening an issue on our GitHub repository or by contacting us directly at info@techjusticelab.org.' }] as faq, i}
					<div
						class="rounded-lg bg-white p-6 shadow"
						in:fly={{ y: 15, duration: 600, delay: 1400 + i * 100, easing: cubicOut }}
					>
						<h3
							class="text-lg font-medium text-neutral-900"
							in:slide={{ duration: 500, delay: 1450 + i * 100 }}
						>
							{faq.question}
						</h3>
						<p class="mt-2 text-neutral-600" in:fade={{ duration: 500, delay: 1500 + i * 100 }}>
							{faq.answer}
						</p>
					</div>
				{/each}
			</div>
		</div>

		<!-- Project Structure -->
		<div class="mt-16" in:fly={{ y: 20, duration: 700, delay: 1800, easing: cubicOut }}>
			<h2
				class="text-center text-2xl font-bold text-neutral-900"
				in:slide={{ duration: 500, delay: 1900 }}
			>
				Project Structure
			</h2>
			<div
				class="mt-8 overflow-hidden rounded-lg bg-white p-6 shadow"
				in:scale={{ start: 0.97, duration: 600, delay: 2000, easing: cubicOut }}
			>
				<p class="mb-4 text-neutral-600" in:fade={{ duration: 500, delay: 2100 }}>
					Motion Index is organized into the following main folders:
				</p>
				<ul class="space-y-4">
					{#each [{ name: 'API', description: 'Backend API services and endpoints' }, { name: 'Database', description: 'Files relating to interfacing with the databases' }, { name: 'Docs', description: 'Project documentation' }, { name: 'Documents', description: 'Local document storage for development' }, { name: 'Web', description: 'Frontend web application built with Svelte' }] as folder, i}
						<li
							class="flex items-start"
							in:fly={{ x: -10, duration: 500, delay: 2200 + i * 100, easing: cubicOut }}
						>
							<div class="mr-3 mt-1 flex-shrink-0">
								<svg
									class="h-5 w-5 text-primary-500"
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 20 20"
									fill="currentColor"
								>
									<path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
								</svg>
							</div>
							<div>
								<p class="font-medium text-neutral-900">{folder.name}</p>
								<p class="text-sm text-neutral-600">{folder.description}</p>
							</div>
						</li>
					{/each}
				</ul>
			</div>
		</div>

		<!-- Footer -->
		<div class="mt-16 text-center" in:fade={{ duration: 600, delay: 2700 }}>
			<p class="text-neutral-500">
				Motion Index is a project by <a
					href="https://techjusticelab.org"
					class="font-medium text-primary-600 hover:text-primary-500">Tech Justice Lab</a
				>
			</p>
			<p class="mt-1 text-sm text-neutral-400" in:fade={{ duration: 600, delay: 2800 }}>
				We're dedicated to improving access to justice through technology
			</p>
		</div>
	</div>
</div>

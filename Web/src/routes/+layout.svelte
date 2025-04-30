<script>
	let { children } = $props();
	import '../app.css';
</script>

<div class="min-h-screen bg-gray-50">
	<header class="bg-indigo-600 shadow-md">
		<div class="container mx-auto px-4 py-4">
			<div class="flex flex-col items-center justify-between sm:flex-row">
				<h1 class="text-2xl font-bold text-white">Motion Index</h1>
				<nav class="mt-3 sm:mt-0">
					<ul class="flex space-x-6 text-white">
						<li><a href="/" class="transition hover:text-indigo-200">Search</a></li>
						<li><a href="/upload" class="transition hover:text-indigo-200">Upload</a></li>

						<li>
							<a href="/account" class="transition hover:text-indigo-200"> Account</a>
						</li>
						<li>
							<a
								href="/help"
								class="rounded bg-white p-2 text-indigo-600 transition hover:bg-indigo-800 hover:text-white"
								>Help</a
							>
						</li>
					</ul>
				</nav>
			</div>
		</div>
	</header>

	<main>
		{@render children()}
	</main>

	<footer class="mt-12 border-t bg-gray-100">
		<div class="container mx-auto px-4 py-6 text-center text-sm text-gray-600">
			<p>Motion Index Search - Legal Document Repository</p>
		</div>
	</footer>
</div>

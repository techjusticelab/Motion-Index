<script lang="ts">
	import '../app.css';
	
	let { children } = $props();
</script>

<div class="min-h-screen bg-gray-50">
  <header class="bg-indigo-600 shadow-md">
    <div class="container mx-auto px-4 py-4">
      <h1 class="text-2xl font-bold text-white">Motion Index</h1>
    </div>
  </header>
  
  <main>
    {@render children()}
  </main>
  
  <footer class="bg-gray-100 border-t mt-12">
    <div class="container mx-auto px-4 py-6 text-center text-gray-600 text-sm">
      <p>Motion Index Search - Legal Document Repository</p>
    </div>
  </footer>
</div>

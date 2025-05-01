script>
	import { supabase } from '$lib/supabaseClient';

	let email = '';
	let password = '';

	async function handleLogin() {
		const { error } = await supabase.auth.signInWithPassword({
			email,
			password
		});

		if (error) alert(error.message);
	}

	async function handleSignUp() {
		const { error } = await supabase.auth.signUp({
			email,
			password
		});

		if (error) alert(error.message);
	}
</script>

<form on:submit|preventDefault>
	<input bind:value={email} type="email" placeholder="Email" />
	<input bind:value={password} type="password" placeholder="Password" />
	<button on:click={handleLogin}>Sign In</button>
	<button on:click={handleSignUp}>Sign Up</button>
</form>

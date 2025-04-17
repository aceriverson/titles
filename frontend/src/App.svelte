<script>
	import { onMount } from 'svelte';
	import Map from './Map.svelte';
	import SideBar from './SideBar.svelte';
	import 'remixicon/fonts/remixicon.css';
	import { user } from './stores.js';

	onMount(async () => {
		// Check for the token in the URL and save it to localStorage if it exists
		const urlParams = new URLSearchParams(window.location.search);
		const token = urlParams.get('token');
		if (token) {
			localStorage.setItem('token', token);
			history.replaceState(null, '', window.location.pathname);
		}

		const fetchUser = async () => {
			const token = localStorage.getItem('token');
			try {
				const response = await fetch('/api/user', {
					headers: {
						Authorization: `Bearer ${token}`
					}
				});
				const data = await response.json();
				return data;
			} catch {
				localStorage.setItem('token', null);
				return null;
			}
		};

		user.set(await fetchUser());
		console.log($user);
	});
</script>

<main>
	<SideBar />
	<Map />
</main>

<style>
	main {
		text-align: left;
		margin: 0 auto;
	}
</style>

<script>
	import { onMount } from 'svelte';
	import Map from './Map.svelte';
	import SideBar from './SideBar.svelte';
	import 'remixicon/fonts/remixicon.css';
	import { user } from './stores.js';

	onMount(async () => {
		const fetchUser = async () => {
			try {
				const response = await fetch('/api/polygons/user', {
					credentials: 'include'
				});
				const data = await response.json();
				return data;
			} catch {
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

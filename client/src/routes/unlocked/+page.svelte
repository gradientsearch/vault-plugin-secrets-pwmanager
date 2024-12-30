<script lang="ts">
	import Header from './components/header.svelte';
	import Password from './components/password.svelte';
	import Sidebar from './components/sidebar.svelte';
	import Vault from './components/vault.svelte';

	let bg = ['bg-sky-500', 'bg-slate-500', 'bg-pink-400', 'bg-green-400'];
</script>

<div class="flex h-full">
	<Sidebar></Sidebar>
	<div class="h-full w-full flex-col">
		<Header></Header>
		<div class="flex h-[calc(100vh-48px)] w-full">
			<Vault></Vault>
			<Password></Password>
		</div>
	</div>
</div>

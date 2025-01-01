<script lang="ts">
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { onMount } from 'svelte';
	import Header from './components/header.svelte';
	import Password from './components/password.svelte';
	import Sidebar from './components/sidebar.svelte';
	import Vault from './components/vault.svelte';
	let kp: KeyPair | undefined;

	onMount(async () => {
		kp = await storedKeyPair.get();
	});
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

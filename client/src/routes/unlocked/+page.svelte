<script lang="ts">
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { onMount } from 'svelte';
	import Header from './components/header/header.svelte';
	import Password from './components/password.svelte';
	import Sidebar from './components/sidebar.svelte';
	import Vault from './components/vault.svelte';
	import { getAPI, type Api } from '$lib/api';
	import type { Zarf } from './models/zarf';
	import type { PasswordListService } from './services/passwordList.service';
	import type { PasswordItem } from './models/input';

	let kp: KeyPair | undefined;
	let zarf: Zarf | null = $state(null);
	let api: Api | undefined;
	let selectedVault: PasswordList | null = $state(null);
	let passwordListService: PasswordListService | null = $state(null);

	let selectedPasswordItem: PasswordItem | undefined = $state();
	let passwordItems: PasswordItem[] = $state([]);

	function updateSelectedPasswordItem(passwordItem: PasswordItem) {
		selectedPasswordItem = passwordItem;
	}

	$effect(() => {
		selectedVault;
	});

	onMount(async () => {
		kp = await storedKeyPair.get();
		api = getAPI();
		zarf = {
			Keypair: kp,
			Api: api
		};
	});
</script>

<div class="flex h-full">
	{#if zarf !== null}
		<Sidebar bind:selectedVault></Sidebar>
		<div class="h-full w-full flex-col">
			<Header bind:passwordListService></Header>
			<div class="flex h-[calc(100vh-48px)] w-full">
				<Vault
					bind:zarf
					bind:passwordListService
					bind:selectedVault
					bind:selectedPasswordItem
					bind:passwordItems
				></Vault>
				<Password bind:selectedPasswordItem></Password>
			</div>
		</div>
	{/if}
</div>

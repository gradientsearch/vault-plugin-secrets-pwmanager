<script lang="ts">
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { onMount } from 'svelte';
	import Header from './components/header/header.svelte';
	import PasswordColumn from './layout/passwordColumn.svelte';
	import SidebarColumn from './layout/sidebarColumn.svelte';
	import PasswordListColumn from './layout/passwordListColumn.svelte';
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
		<SidebarColumn bind:selectedVault></SidebarColumn>
		<div class="h-full w-full flex-col">
			<Header bind:passwordListService></Header>
			<div class="flex h-[calc(100vh-48px)] w-full">
				<PasswordListColumn
					bind:zarf
					bind:passwordListService
					bind:selectedVault
					bind:selectedPasswordItem
					bind:passwordItems
				></PasswordListColumn>
				<PasswordColumn bind:selectedPasswordItem></PasswordColumn>
			</div>
		</div>
	{/if}
</div>

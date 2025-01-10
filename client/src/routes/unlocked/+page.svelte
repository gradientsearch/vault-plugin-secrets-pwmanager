<script lang="ts">
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { onMount } from 'svelte';
	import Header from './layout/header.svelte';
	import PasswordColumn from './layout/passwordColumn.svelte';
	import SidebarColumn from './layout/sidebarColumn.svelte';
	import PasswordItemsColumn from './layout/passwordItemsColumn.svelte';
	import { getAPI, type Api } from '$lib/api';
	import type { Zarf } from './models/zarf';
	import type { PasswordItemsService } from './services/passwordItems.service';
	import type { Entry } from './models/input';

	let kp: KeyPair | undefined;
	let zarf: Zarf | null = $state(null);
	let api: Api | undefined;
	let bundle: Bundle | null = $state(null);
	let passwordListService: PasswordItemsService | null = $state(null);

	let selectedPasswordItem: Entry | undefined = $state();
	let passwordItems: Entry[] = $state([]);

	$effect(() => {
		bundle;
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
		<SidebarColumn bind:bundle></SidebarColumn>
		<div class="h-full w-full flex-col">
			<Header bind:passwordListService></Header>
			<div class="flex h-[calc(100vh-48px)] w-full">
				<PasswordItemsColumn
					bind:zarf
					bind:passwordListService
					bind:bundle
					bind:selectedPasswordItem
					bind:passwordItems
				></PasswordItemsColumn>
				<PasswordColumn bind:selectedPasswordItem></PasswordColumn>
			</div>
		</div>
	{/if}
</div>

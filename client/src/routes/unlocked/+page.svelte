<script lang="ts">
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { onMount } from 'svelte';
	import Header from './layout/header.svelte';
	import PasswordColumn from './layout/passwordColumn.svelte';
	import SidebarColumn from './layout/sidebarColumn.svelte';
	import BundleColumn from './layout/bundleColumn.svelte'
	import { getAPI, type Api } from '$lib/api';
	import type { Zarf } from './models/zarf';
	import type { EntriesService } from './services/entry.service';
	import type { Entry } from './models/input';

	let kp: KeyPair | undefined;
	let zarf: Zarf | null = $state(null);
	let api: Api | undefined;
	let bundle: Bundle | null = $state(null);
	let passwordListService: EntriesService | null = $state(null);

	let selectedEntry: Entry | undefined = $state();
	let entries: Entry[] = $state([]);

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
				<BundleColumn
					bind:zarf
					bind:passwordListService
					bind:bundle
					bind:selectedEntry
					bind:entries
				></BundleColumn>
				<PasswordColumn bind:selectedEntry></PasswordColumn>
			</div>
		</div>
	{/if}
</div>

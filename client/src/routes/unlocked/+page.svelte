<script lang="ts">
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { onMount } from 'svelte';
	import HeaderView from './layout/headerView.svelte';
	import EntryView from './layout/entryView.svelte';
	import SidebarView from './layout/sidebarView.svelte';
	import BundleView from './layout/bundleView.svelte'
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
		<SidebarView bind:bundle></SidebarView>
		<div class="h-full w-full flex-col">
			<HeaderView bind:passwordListService></HeaderView>
			<div class="flex h-[calc(100vh-48px)] w-full">
				<BundleView
					bind:zarf
					bind:passwordListService
					bind:bundle
					bind:selectedEntry
					bind:entries
				></BundleView>
				<EntryView bind:selectedEntry></EntryView>
			</div>
		</div>
	{/if}
</div>

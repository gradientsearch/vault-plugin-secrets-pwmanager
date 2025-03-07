<script lang="ts">
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { onMount } from 'svelte';
	import HeaderView from './views/headerView.svelte';
	import EntryView from './views/entryView.svelte';
	import SidebarView from './views/sidebarView.svelte';
	import BundleView from './views/bundleView.svelte';
	import { getAPI, type Api } from '$lib/api';
	import type { Zarf } from './models/zarf';
	import type { BundleService } from './services/bundle.service';
	import type { Entry, Metadata } from './models/entry';
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	let kp: KeyPair | undefined;
	let zarf: Zarf | undefined = $state(undefined);
	let api: Api | undefined;
	let bundle: Bundle | null = $state(null);
	let bundleService: BundleService | undefined = $state(undefined);
	let selectedEntryMetadata: Metadata | undefined = $state();
	let entries: Entry[] = $state([]);

	onMount(async () => {
		kp = await storedKeyPair.get();

		api = getAPI();

		if (kp === undefined || api === undefined) {
			goto(`${base}/locked`);
			return;
		}

		zarf = {
			Keypair: kp,
			Api: api
		};
	});
	let clientHeight = $state(0)
</script>

<div bind:clientHeight class="flex h-full">
	{#if zarf !== undefined}
		<SidebarView bind:bundle bind:zarf {clientHeight}></SidebarView>
		<div class="h-full w-full flex-col">
			<HeaderView bind:bundleService></HeaderView>
			<div class="flex h-[calc(100vh-48px)] w-full">
				<BundleView
					bind:zarf
					bind:bundleService
					bind:bundle
					bind:selectedEntryMetadata
					bind:entries
				></BundleView>
				<EntryView bind:selectedEntryMetadata bind:bundleService></EntryView>
			</div>
		</div>
	{/if}
</div>

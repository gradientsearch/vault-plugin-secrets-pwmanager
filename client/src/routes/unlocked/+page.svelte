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
	import type { BundleMetadata } from './models/bundle/vault/metadata';

	let kp: KeyPair | undefined;
	let zarf: Zarf | undefined = $state(undefined);
	let api: Api | undefined;
	let bundle: Bundle | null = $state(null);
	let bundleService: BundleService | undefined = $state(undefined);
	let selectedEntryMetadata: Metadata | undefined = $state();
	let bundleMetadata: BundleMetadata = $state({
		entries: [],
		bundleName: '',
		version: 0
	});

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
</script>

<div class="flex h-[100vh]">
	{#if zarf !== undefined}
		<SidebarView bind:bundle bind:zarf></SidebarView>
		<div class="h-full w-full flex-col">
			<HeaderView bind:bundleMetadata bind:bundleService></HeaderView>
			<div class="flex h-[calc(100vh-48px)] w-full">
				<BundleView
					bind:zarf
					bind:bundleMetadata
					bind:bundleService
					bind:bundle
					bind:selectedEntryMetadata
				></BundleView>
				<EntryView bind:bundleMetadata bind:selectedEntryMetadata bind:bundleService></EntryView>
			</div>
		</div>
	{/if}
</div>

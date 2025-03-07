<script lang="ts">
	import { untrack } from 'svelte';
	import { newPasswordEntry as newPasswordEntry, type Entry, type Metadata } from '../models/entry';
	import { KVBundleService, type BundleService } from '../services/bundle.service';
	import { base } from '$app/paths';
	import type { BundleMetadata } from '../models/bundle/vault/metadata';

	let headerHeight = $state(0);
	let errorMessage: string | undefined = $state(undefined);

	let {
		bundle = $bindable(),
		zarf = $bindable(),
		bundleService = $bindable<BundleService>(),
		selectedEntryMetadata = $bindable<Entry>(),
		entries: entries = $bindable()
	} = $props();

	$effect(() => {
		bundle;

		untrack(() => {
			selectedEntryMetadata = undefined;
			setBundleService();
		});
	});

	async function onEntriesChanged() {
		let [metadata, err] = await bundleService.getMetadata();
		if (err !== undefined) {
			//TODO Toast error
			return;
		}
		entries = metadata.entries.reverse();

		let checkSelectedEntry = metadata.entries.filter((e: Metadata) => {
			return selectedEntryMetadata?.ID === e.ID;
		});

		if (checkSelectedEntry.length === 0 && selectedEntryMetadata !== undefined) {
			selectedEntryMetadata = undefined;
		} else if (entries.length > 0) {
			selectedEntryMetadata = entries[0];
		}
	}

	async function setBundleService() {
		if (bundle?.Type === 'bundle') {
			bundleService = new KVBundleService(zarf, bundle, onEntriesChanged);
			let err = await bundleService.init();
			if (err !== undefined) {
				errorMessage = err;
				alert(err);
			}
			let [vm, err2] = await bundleService.getMetadata();

			if (err2 !== undefined) {
				errorMessage = 'Error loading entries';
			}
			entries = (vm as BundleMetadata).entries.reverse();
		}
	}

	function onSelectedEntry(e: Entry) {
		selectedEntryMetadata = e;
	}
</script>

<div
	class="relative w-full max-w-96 overflow-x-hidden overflow-y-scroll border-2 border-r-8 border-border_primary bg-page_faint"
>
	<header
		bind:clientHeight={headerHeight}
		class="absolute top-0 flex w-full border-b-2 border-border_primary p-2"
	>
		<h1 class="text-base capitalize">{bundle?.Name} {bundle?.Type}</h1>
	</header>

	<span style="min-height: {headerHeight}px;" class="flex flex-1"></span>
	{#if entries}
		{#each entries as e}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				onclick={() => {
					onSelectedEntry(e);
				}}
				style="top: {headerHeight}px"
				class=" flex w-[100%] text-left text-base hover:scale-110 hover:bg-surface_interactive_hover hover:ps-3.5 hover:shadow-md {e.ID ===
				selectedEntryMetadata?.ID
					? 'bg-blue-100'
					: ''}"
			>
				<button class="flex w-full flex-row items-center p-4">
					<span class="pe-3 text-3xl">🔑</span>
					<div class="flex flex-col text-start">
						<span class="text-base font-bold text-foreground_strong"> {e.Name}</span>
						<span class="text-sm text-foreground_faint"> {e.Value}</span>
					</div>
				</button>
			</div>
		{/each}
	{/if}
</div>

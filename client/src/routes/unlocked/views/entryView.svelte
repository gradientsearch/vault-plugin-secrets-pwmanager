<script lang="ts">
	import { untrack } from 'svelte';
	import { type Metadata, type Entry, MODE } from '../models/entry';
	import { getPasswordComponent } from '../components/entries/components';
	import type { BundleService } from '../services/bundle.service';

	let {
		selectedEntryMetadata = $bindable<Metadata>(),
		bundleService = $bindable<BundleService>()
	} = $props();

	let copyOfSelectedEntry: Entry | undefined = $state();
	let errMessage: string | undefined = $state(undefined);
	let mode: MODE = $state(MODE.VIEW);

	$effect(() => {
		selectedEntryMetadata;
		console.log('selectedEntryMetadata Entryview: ', selectedEntryMetadata);
		if (selectedEntryMetadata) {
			untrack(() => {
				(async () => {
					console.log('untrack');
					console.log('async');
					let [e, err] = await bundleService.getEntry(selectedEntryMetadata);
					if (err !== undefined) {
						console.log(err);
						errMessage = `error getting entry: ${err}`;
					}

					copyOfSelectedEntry = e;
				})();
			});
		}
	});
</script>

<div class=" w-full border-t-2 border-border_primary">
	{#if copyOfSelectedEntry}
		{@const Component = getPasswordComponent(copyOfSelectedEntry?.Type)}
		<div class="h-[100%] w-full">
			<Component bind:entry={copyOfSelectedEntry} bind:bundleService bind:mode={mode}></Component>
		</div>
	{/if}
</div>

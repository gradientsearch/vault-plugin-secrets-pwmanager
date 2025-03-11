<script lang="ts">
	import { untrack } from 'svelte';
	import { type Metadata, type Entry, MODE } from '../models/entry';
	import { getPasswordComponent } from '../components/entries/components';
	import type { BundleService } from '../services/bundle.service';
	import type { BundleMetadata } from '../models/bundle/vault/metadata';

	let {
		selectedEntryMetadata = $bindable<Metadata>(),
		bundleService = $bindable<BundleService>(),
		bundleMetadata = $bindable<BundleMetadata>()
	} = $props();

	let originalEntry: Entry | undefined = $state();
	let copyOfSelectedEntry: Entry | undefined = $state();
	let errMessage: string | undefined = $state(undefined);
	let mode: MODE = $state(MODE.VIEW);

	$effect(() => {
		selectedEntryMetadata;
		if (selectedEntryMetadata) {
			untrack(() => {
				(async () => {
					let [e, err] = await bundleService.getEntry(selectedEntryMetadata);
					if (err !== undefined) {
						errMessage = `error getting entry: ${err}`;
					}
					originalEntry = e;
					copyOfSelectedEntry = e;
				})();
				mode = MODE.VIEW;
			});
		} else {
			copyOfSelectedEntry = undefined;
			mode = MODE.VIEW;
		}
	});

	function cancel() {
		copyOfSelectedEntry = JSON.parse(JSON.stringify(originalEntry))
	}
</script>

<div id="entry-view" class=" w-full overflow-y-scroll border-t-2 border-border_primary">
	{#if copyOfSelectedEntry}
		{@const Component = getPasswordComponent(copyOfSelectedEntry?.Type)}
		<div class="h-[100%] w-full">
			<Component bind:entry={copyOfSelectedEntry} bind:bundleMetadata bind:bundleService bind:mode cancel={cancel}
			></Component>
		</div>
	{/if}
</div>

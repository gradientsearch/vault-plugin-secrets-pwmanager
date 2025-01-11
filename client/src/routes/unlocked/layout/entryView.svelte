<script lang="ts">
	import { untrack } from 'svelte';
	import type { Entry } from '../models/entry';
	import { getPasswordComponent } from '../components/entries/components';

	let { selectedEntry = $bindable() } = $props();

	let copyOfSelectedEntry: Entry | undefined = $state();

	$effect(() => {
		selectedEntry;
		if (selectedEntry) {
            untrack(()=>{
                copyOfSelectedEntry = JSON.parse(JSON.stringify(selectedEntry));
            })

		}
	});
</script>

<div class=" w-full border-t-2 border-border_primary bg-page_faint">
	{#if copyOfSelectedEntry}
		{@const Component = getPasswordComponent(copyOfSelectedEntry?.Type)}
		<div class="w-full">
			<Component bind:entry={copyOfSelectedEntry}></Component>
		</div>
	{/if}
</div>

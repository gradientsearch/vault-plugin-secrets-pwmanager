<script lang="ts">
	import { untrack } from 'svelte';
	import type { Entry } from '../models/input';
	import { getPasswordComponent } from '../components/passwordItems/components';

	let { selectedPasswordItem = $bindable() } = $props();

	let copyOfSelectedPasswordItem: Entry | undefined = $state();

	$effect(() => {
		selectedPasswordItem;
		if (selectedPasswordItem) {
            untrack(()=>{
                copyOfSelectedPasswordItem = JSON.parse(JSON.stringify(selectedPasswordItem));
            })

		}
	});
</script>

<div class=" w-full border-t-2 border-border_primary bg-page_faint">
	{#if copyOfSelectedPasswordItem}
		{@const Component = getPasswordComponent(copyOfSelectedPasswordItem?.Type)}
		<div class="w-full">
			<Component bind:passwordItem={copyOfSelectedPasswordItem}></Component>
		</div>
	{/if}
</div>

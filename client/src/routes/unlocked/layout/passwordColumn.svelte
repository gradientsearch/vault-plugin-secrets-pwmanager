<script lang="ts">
	import { untrack } from 'svelte';
	import type { PasswordItem } from '../models/input';
	import { getPasswordComponent } from '../components/passwordItems/components';

	let { selectedPasswordItem = $bindable() } = $props();

	let copyOfSelectedPasswordItem: PasswordItem | undefined = $state();

	$effect(() => {
		selectedPasswordItem;
		if (selectedPasswordItem) {
            untrack(()=>{
                copyOfSelectedPasswordItem = JSON.parse(JSON.stringify(selectedPasswordItem));
                			
            console.log(copyOfSelectedPasswordItem)  
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

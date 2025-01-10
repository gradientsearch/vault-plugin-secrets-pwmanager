<script lang="ts">
	import { getPasswordComponent } from './passwordItems/components';

	let { selectedPasswordItem = $bindable() } = $props();

	$effect(() => {
		selectedPasswordItem;
        console.log(selectedPasswordItem)
	});
</script>

<div class=" w-full  border-t-2 border-border_primary bg-page_faint">
	{#if selectedPasswordItem}
		{@const Component = getPasswordComponent(selectedPasswordItem?.Type)}
		<div class="w-full">
			<Component bind:passwordItem={selectedPasswordItem}></Component>
		</div>
	{/if}
</div>

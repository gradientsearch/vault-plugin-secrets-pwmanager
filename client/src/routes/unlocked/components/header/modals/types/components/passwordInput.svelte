<script>
	import ItemInput from './itemInput.svelte';
	let { label, placeholder, value = $bindable(), reveal = $bindable() } = $props();
</script>

{#if !reveal}
	<ItemInput type={'password'} {label} {placeholder} bind:value></ItemInput>
{:else}
	<ItemInput type={'text'} {label} {placeholder} bind:value></ItemInput>
{/if}

<script>
	import ItemInput from './itemInput.svelte';
	let { label, placeholder, value = $bindable(), idx, last, reveal = $bindable() } = $props();
</script>

{#if !reveal}
	<ItemInput type={'password'} {label} {placeholder} bind:value {idx} {last}></ItemInput>
{:else}
	<ItemInput type={'text'} {label} {placeholder} bind:value {idx} {last}></ItemInput>
{/if}

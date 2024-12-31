<script lang="ts">
	import { base } from '$app/paths';

	var { href = '', fn = null, disabled = false, children } = $props();
</script>

{#if href.length > 0}
	<a
		class="mt-6 inline-block w-[76px] rounded border bg-blue-200 py-2 text-center text-sm font-bold text-foreground_high_contrast hover:bg-blue-400"
		{href}
	>
		{@render children?.()}
	</a>
{:else if fn != null}
	<button
		class="mt-6 inline-block w-[76px] rounded border bg-blue-200 py-2 text-center text-sm font-bold text-foreground_high_contrast hover:bg-blue-400"
		onclick={fn}
		disabled={disabled}
	>
		{@render children?.()}
	</button>
{/if}

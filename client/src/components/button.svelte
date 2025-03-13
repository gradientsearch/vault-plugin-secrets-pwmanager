<script lang="ts">
	import { base } from '$app/paths';

	var { href = '', fn = null, disabled = false, primary = true, children, margin = 'mt-6' } = $props();
</script>

{#if href.length > 0}
	<a
		class="{margin} inline-block w-[76px] rounded border {primary ? 'bg-blue-200 text-foreground_high_contrast hover:bg-blue-400': 'bg-surface_faint border-border_strong text-foreground_primary' } py-2 text-center text-sm font-bold "
		{href}
	>
		{@render children?.()}
	</a>
{:else if fn != null}
	<button
		class="{margin} inline-block w-[76px] rounded border {primary ? 'bg-blue-200 text-foreground_high_contrast hover:bg-blue-400': 'bg-surface_faint  border-border_strong text-foreground_primary hover:bg-foreground_high_contrast' } py-2 text-center text-sm font-bol"
		onclick={fn}
		disabled={disabled}
	>
		{@render children?.()}
	</button>
{/if}

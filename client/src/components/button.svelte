<script lang="ts">
	import { base } from '$app/paths';

	var { href="", fn = null, children } = $props();
</script>

{#if href.length > 0}
<a
	class="bg-blue-200 text-foreground_high_contrast hover:bg-blue-400 mt-6 inline-block w-[76px] rounded border py-2 text-center text-sm font-bold"
	href="{href}"
>
	{@render children?.()}
</a>
{:else if fn != null}
<button
	class="bg-blue-200 text-foreground_high_contrast hover:bg-blue-400 mt-6 inline-block w-[76px] rounded border py-2 text-center text-sm font-bold"
	onclick={fn}
>
	{@render children?.()}
</button>
{/if}
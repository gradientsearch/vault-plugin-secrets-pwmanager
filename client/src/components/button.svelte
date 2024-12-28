<script lang="ts">
	import { base } from '$app/paths';

	var { href, children } = $props();
</script>

<a
	class="bg-blue-200 text-foreground_high_contrast hover:bg-blue-400 mt-6 inline-block w-[76px] rounded border py-2 text-center text-sm font-bold"
	href="{href}"
>
	{@render children?.()}
</a>

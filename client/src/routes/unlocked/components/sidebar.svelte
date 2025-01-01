<script lang="ts">
	let {children}= $props()
</script>

<div
	class="bg-sky-400 hidden w-full content-center text-center md:max-w-64 md:bg-green-500 lg:block"
>
	{@render children?.()}
</div>

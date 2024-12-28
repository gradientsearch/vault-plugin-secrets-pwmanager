<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<div class="h-screen w-full text-3xl">
	<div class="flex h-full flex-1 flex-col">
		<div class="flex-1">
			{@render children()}
		</div>
	</div>
</div>

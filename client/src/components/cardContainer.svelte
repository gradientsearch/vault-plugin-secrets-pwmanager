<script lang="ts">
	var { className="", children } = $props();
</script>

<div class="border-border_faint relative block overflow-hidden border p-3 shadow-lg {className}">
	{@render children?.()}
</div>

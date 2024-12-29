<script lang="ts">
	var { className = '', children } = $props();
</script>

<div class="relative block overflow-hidden border border-border_faint p-3 shadow-lg {className}">
	{@render children?.()}
</div>

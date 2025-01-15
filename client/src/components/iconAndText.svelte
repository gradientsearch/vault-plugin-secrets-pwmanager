<script>
	import Icon from './icon.svelte';

	var { className } = $props();
</script>

<div class="flex-row justify-center {className}">
	<div class="flex justify-center">
		<Icon className="brand-icon-large"></Icon>
	</div>
	<p class="text-center text-xs">
		Vault Password Manager <a
			class="text-blue-200 underline"
			href="https://github.com/gradientsearch/vault-plugin-secrets-pwmanager">plugin</a
		>
	</p>
</div>

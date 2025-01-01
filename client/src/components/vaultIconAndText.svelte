<script>
	import VaultIcon from './vaultIcon.svelte';

	var { className } = $props();
</script>

<div class="flex-row justify-center {className}">
	<div class="flex justify-center">
		<VaultIcon className="brand-icon-large"></VaultIcon>
	</div>
	<p class="text-center text-xs">
		Vault Password Manager <a
			class="text-blue-200 underline"
			href="https://github.com/gradientsearch/vault-plugin-secrets-pwmanager">plugin</a
		>
	</p>
</div>

<script lang="ts">
	import { onMount } from 'svelte';
	import type { PasswordItem } from '../models/input';
	import type { PasswordListService } from '../models/passwordList.service';

	let passwordItems: PasswordItem[] = $state([]);
	let { selectedVault = $bindable(), zarf = $bindable() } = $props();
    let passwordListService: PasswordListService
	onMount(() => {});

    function setPasswordListService(){
        if (selectedVault?.Type === 'vault') {
            
        }
    }

	function getListOfPasswords() {
		// TODO define the type to return
		return zarf?.Api?.getPasswordListMetadata(selectedVault);
	}

	$effect(() => {
		selectedVault;

		getListOfPasswords();
	});
</script>

<div class="w-full max-w-96 overflow-y-scroll border-2 border-border_primary bg-page_faint">
	<header class="sticky top-0 border-b-2 border-border_primary p-2">
		<h1 class="text-base">{selectedVault?.Name} Vault</h1>
	</header>
</div>

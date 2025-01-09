<script lang="ts">
	import { Api, getAPI } from '$lib/api';
	import { onMount } from 'svelte';

	let { selectedVault = $bindable() } = $props();
	let api: Api | undefined;

	onMount(() => {
		api = getAPI();
	});


    function getListOfPasswords(){
        if (api === undefined){
            return
        }

        api.getPasswordListMetadata(selectedVault)
    }

	$effect(() => {
        selectedVault
        getListOfPasswords()

    });
</script>

<div class="w-full max-w-96 overflow-y-scroll border-2 border-border_primary bg-page_faint">
	<header class="sticky top-0 border-b-2 border-border_primary p-2">
		<h1 class="text-base">{selectedVault?.Name} Vault</h1>
	</header>
</div>

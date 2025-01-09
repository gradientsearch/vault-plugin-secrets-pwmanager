<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import type { PasswordItem } from '../models/input';
	import {
		VaultPasswordListService,
		type PasswordListService
	} from '../models/passwordList.service';

	let passwordItems: PasswordItem[] = $state([]);
	let {
		selectedVault = $bindable(),
		zarf = $bindable(),
		passwordListService = $bindable()
	} = $props();
	onMount(() => {});

	let onVaultAddFn = (pi: PasswordItem[]) => {
        passwordItems = pi
    };

	function setPasswordListService() {
		if (selectedVault?.Type === 'vault') {
			passwordListService = new VaultPasswordListService(zarf, selectedVault, onVaultAddFn);
			passwordItems = passwordListService.get();
		}
	}

	$effect(() => {
		selectedVault;
		untrack(() => {
			setPasswordListService();
		});
	});
</script>

<div class="w-full max-w-96 overflow-y-scroll border-2 border-border_primary bg-page_faint">
	<header class="sticky top-0 border-b-2 border-border_primary p-2">
		<h1 class="text-base">{selectedVault?.Name} {selectedVault?.Type}</h1>
	</header>

    {#each passwordItems as i}
    <div>
        <button>
            {i.Metadata.Name}
        </button>
    </div>
    {/each}
</div>

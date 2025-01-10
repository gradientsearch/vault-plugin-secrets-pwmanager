<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import { newPasswordItem, type PasswordItem } from '../models/input';
	import {
		VaultPasswordListService,
		type PasswordListService
	} from '../models/passwordList.service';
	import { base } from '$app/paths';

	let passwordItems: PasswordItem[] = $state([]);
	let {
		selectedVault = $bindable(),
		zarf = $bindable(),
		passwordListService = $bindable()
	} = $props();
	onMount(() => {});

	let onVaultAddFn = (pi: PasswordItem[]) => {
		passwordItems = pi;
	};

	function setPasswordListService() {
		if (selectedVault?.Type === 'vault') {
			passwordListService = new VaultPasswordListService(zarf, selectedVault, onVaultAddFn);
			passwordItems = passwordListService.get();
			let pi = newPasswordItem();
			pi.Metadata.Name = 'My Secret Password';
			pi.Metadata.Value = 'stephen';
			passwordItems.push(pi);
			passwordItems = passwordItems;
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
		<div class="flex w-full text-base p-4 hover:bg-surface_interactive_hover    ">
            <button class="flex flex-row">
            <img class="h-8 p2" src="{base}/icons/key.svg" alt="key icon">
			<div class="flex flex-col text-start">
                <span class='text-base font-bold text-foreground_strong'> {i.Metadata.Name}</span>            
                <span class='text-sm text-foreground_faint'> {i.Metadata.Value}</span>            
                
            </div>
            
			</button>
		</div>
	{/each}
</div>

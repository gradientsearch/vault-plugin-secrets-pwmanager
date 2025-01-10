<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import { newPasswordItem, type PasswordItem } from '../models/input';
	import { VaultPasswordListService } from '../services/passwordList.service';
	import { base } from '$app/paths';
	import Password from './password.svelte';

let headerHeight = $state(0)

	let {
		selectedVault = $bindable(),
		zarf = $bindable(),
		passwordListService = $bindable(),
		selectedPasswordItem = $bindable(),
        passwordItems = $bindable(),
	} = $props();

	$effect(() => {
		selectedVault;
		untrack(() => {
			setPasswordListService();
		});
	});

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
            pi.Core.Items[0].Value = 'stephen'
            pi.Core.Items[1].Value = 'super-secret-password'
            pi.Type = 'password'
            pi.Name = 'My Secret Password';
            
			passwordItems.push(pi);
			passwordItems = passwordItems;
		}
	}

	function onSelectedPasswordItem(passwordItem: PasswordItem) {
		selectedPasswordItem = passwordItem;
	}
</script>

<div
	class="relative  w-full max-w-96 overflow-y-scroll border-2 border-r-8 border-border_primary bg-page_faint"
>
	<header bind:clientHeight={headerHeight}  class="absolute top-0 flex w-full border-b-2 border-border_primary p-2">
		<h1 class="text-base">{selectedVault?.Name} {selectedVault?.Type}</h1>
	</header>

    <span style="min-height: {headerHeight}px;" class="flex flex-1"></span>
	{#each passwordItems as i}
    
		<div style="top: {headerHeight}px" class="flex w-full p-4 text-base hover:bg-surface_interactive_hover">
			<button
				onclick={() => {
					onSelectedPasswordItem(i);
				}}
				class="flex w-full flex-row"
			>
				<img class="p2 h-8" src="{base}/icons/key.svg" alt="key icon" />
				<div class="flex flex-col text-start">
					<span class="text-base font-bold text-foreground_strong"> {i.Metadata.Name}</span>
					<span class="text-sm text-foreground_faint"> {i.Metadata.Value}</span>
				</div>
			</button>
		</div>
	{/each}
</div>


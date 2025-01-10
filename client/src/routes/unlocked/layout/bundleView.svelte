<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import { newEntry as newPasswordEntry, type Entry } from '../models/input';
	import { VaultEntriesService } from '../services/entry.service';
	import { base } from '$app/paths';

	let headerHeight = $state(0);

	let {
		bundle = $bindable(),
		zarf = $bindable(),
		passwordListService = $bindable(),
		selectedEntry: selectedEntry = $bindable(),
		entries: entries = $bindable()
	} = $props();

	$effect(() => {
		bundle;
		untrack(() => {
			setEntrysService();
		});
	});

	onMount(() => {});

	let onVaultAddFn = (es: Entry[]) => {
		entries = es;
	};

	async function setEntrysService() {
		if (bundle?.Type === 'vault') {
			passwordListService = new VaultEntriesService(zarf, bundle, onVaultAddFn);
			// TODO Grab PasswordList Decryption key, create it if it does not exist (policy will only allow the owner of the vault to do this). /keys/{{ identity.entity.id }}
			entries = await passwordListService.get();
			let e = newPasswordEntry();
			e.Metadata.Name = 'My Secret Password';
			e.Metadata.Value = 'stephen';
			e.Core.Items[0].Value = 'stephen';
			e.Core.Items[1].Value = 'super-secret-password';
			e.Type = 'password';
			e.Name = 'My Secret Password';

			entries.push(e);
			entries = entries;
		}
	}

	function onSelectedEntry(e: Entry) {
		selectedEntry = e;
	}
</script>

<div
	class="relative w-full max-w-96 overflow-y-scroll border-2 border-r-8 border-border_primary bg-page_faint"
>
	<header
		bind:clientHeight={headerHeight}
		class="absolute top-0 flex w-full border-b-2 border-border_primary p-2"
	>
		<h1 class="text-base">{bundle?.Name} {bundle?.Type}</h1>
	</header>

	<span style="min-height: {headerHeight}px;" class="flex flex-1"></span>
	{#each entries as e}
		<div
			style="top: {headerHeight}px"
			class="flex w-full p-4 text-base hover:bg-surface_interactive_hover"
		>
			<button
				onclick={() => {
					onSelectedEntry(e);
				}}
				class="flex w-full flex-row"
			>
				<img class="p2 h-8" src="{base}/icons/key.svg" alt="key icon" />
				<div class="flex flex-col text-start">
					<span class="text-base font-bold text-foreground_strong"> {e.Metadata.Name}</span>
					<span class="text-sm text-foreground_faint"> {e.Metadata.Value}</span>
				</div>
			</button>
		</div>
	{/each}
</div>

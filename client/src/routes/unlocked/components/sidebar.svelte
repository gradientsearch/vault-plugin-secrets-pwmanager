<script lang="ts">
	import { onMount } from 'svelte';
	import VaultIcon from '../../../components/vaultIcon.svelte';

	let {selectedVault = $bindable()} = $props()
	let privateVault: PasswordList | null = $state(null)
	onMount(() => {
		let info = localStorage.getItem('loginInfo')
		if (info !== null) {
			let infoObj =  JSON.parse(info)
			let p: PasswordList = {
				Type: 'vault',
				Path: `vaults/${infoObj['entityID']}/private`,
				Name: 'private'
			}
			privateVault = p
			selectedVault = p
		}
	})

	
</script>

<div class="bg-token_side_nav_color_surface_primary hidden w-full md:max-w-64 lg:block">
	<div class="flex-row px-3 py-4">
		<div class="flex content-start">
			<VaultIcon className="nav-header-icon"></VaultIcon>
		</div>
		<ul class="pt-3">
			<li
				class="text-token_side_nav_color_foreground_faint height-[36px] px-[8px] py-[9px] text-sm font-bold"
			>
				<div>Vaults</div>
			</li>

			{#if privateVault}
			<li
				class="text-token_side_nav_color_foreground_strong bg-token_side_nav_color_surface_interactive_active hover:bg-token_side_nav_color_surface_interactive_hover height-[36px] rounded-lg px-[8px] py-[9px] text-sm"
			>
				<div>{privateVault.Name}</div>
			</li>
			{/if}
		</ul>
	</div>
</div>

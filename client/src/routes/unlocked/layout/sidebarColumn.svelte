<script lang="ts">
	import { onMount } from 'svelte';
	import VaultIcon from '../../../components/vaultIcon.svelte';

	let { passwordBundle = $bindable() } = $props();
	let privateVault: PasswordBundle | null = $state(null);
	onMount(() => {
		let info = localStorage.getItem('loginInfo');
		if (info !== null) {
			let infoObj = JSON.parse(info);
			let p: PasswordBundle = {
				Type: 'vault',
				Path: `vaults/${infoObj['entityID']}/private`,
				Name: 'private',
				Owner: infoObj['entityID']
			};
			privateVault = p;
			passwordBundle = p;
		}
	});
</script>

<div class="hidden w-full bg-token_side_nav_color_surface_primary md:max-w-64 lg:block">
	<div class="flex-row px-3 py-4">
		<div class="flex content-start">
			<VaultIcon className="nav-header-icon"></VaultIcon>
		</div>
		<ul class="pt-3">
			<li
				class="height-[36px] px-[8px] py-[9px] text-sm font-bold text-token_side_nav_color_foreground_faint"
			>
				<div>Vaults</div>
			</li>

			{#if privateVault}
				<li
					class="height-[36px] rounded-lg bg-token_side_nav_color_surface_interactive_active px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_strong hover:bg-token_side_nav_color_surface_interactive_hover"
				>
					<div>{privateVault.Name}</div>
				</li>
			{/if}
		</ul>
	</div>
</div>

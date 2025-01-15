<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '../../../components/icon.svelte';

	let { bundle = $bindable() } = $props();
	let privateBundle: Bundle | null = $state(null);
	onMount(() => {
		let info = localStorage.getItem('loginInfo');
		if (info !== null) {
			let infoObj = JSON.parse(info);
			let p: Bundle = {
				Type: 'bundle',
				Path: `bundle/${infoObj['entityID']}/private`,
				Name: 'private',
				Owner: infoObj['entityID']
			};
			privateBundle = p;
			bundle = p;
		}
	});
</script>

<div class="hidden w-full bg-token_side_nav_color_surface_primary md:max-w-64 lg:block">
	<div class="flex-row px-3 py-4">
		<div class="flex content-start">
			<Icon className="nav-header-icon"></Icon>
		</div>
		<ul class="pt-3">
			<li
				class="height-[36px] px-[8px] py-[9px] text-sm font-bold text-token_side_nav_color_foreground_faint"
			>
				<div>Bundles</div>
			</li>

			{#if privateBundle}
				<li
					class="height-[36px] rounded-lg bg-token_side_nav_color_surface_interactive_active px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_strong hover:bg-token_side_nav_color_surface_interactive_hover"
				>
					<div>{privateBundle.Name}</div>
				</li>
			{/if}
		</ul>
	</div>
</div>

<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '../../../components/icon.svelte';

	let { bundle = $bindable(), zarf = $bindable() } = $props();
	let defaultBundle: Bundle | undefined = $state(undefined);
	onMount(() => {
		let info = localStorage.getItem('loginInfo');
		if (info !== null) {
			let infoObj = JSON.parse(info);
			let b: Bundle = {
				Type: 'bundle',
				Path: `bundles/data/${infoObj['entityID']}/${infoObj['entityID']}`,
				Name: 'personal',
				Owner: infoObj['entityID']
			};
			defaultBundle = b;
			bundle = b;
		}

		(async () => {
			// retrieve list of bundles
			// retrieve the name of the bundles and cache that in local storage
			// list the bundles/ add on click event
		})();
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

			{#if defaultBundle}
				<li
					class="height-[36px] rounded-lg bg-token_side_nav_color_surface_interactive_active px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_strong hover:bg-token_side_nav_color_surface_interactive_hover"
				>
					<div class="capitalize">{defaultBundle.Name}</div>
				</li>
			{/if}
		</ul>
	</div>
</div>

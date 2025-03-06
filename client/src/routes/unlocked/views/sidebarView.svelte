<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '../../../components/icon.svelte';
	import { KVBundleService } from '../services/bundle.service';

	let { bundle = $bindable(), zarf = $bindable() } = $props();
	let bundles: Bundle[]  = $state([]);

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
			bundles?.push(b);
			bundle = b;
		}

		(async () => {
			// retrieve list of bundles
			let [bs, err] = await KVBundleService.getBundles(zarf);
			if (err !== undefined) {
				console.log(`error listing bundles ${err}`);
			}

			if (bs !== undefined) {
				bundles = [...bundles, ...bs];
			}

			console.log(bundles);

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

			{#if bundles}
				{#each bundles as b}
					<li
						class="{b.Path === bundle.Path
							? 'bg-token_side_nav_color_surface_interactive_active'
							: ''} height-[36px] my-2 min-h-[36px] rounded-lg px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_strong hover:bg-token_side_nav_color_surface_interactive_hover"
					>
						{#if b.Name.length === 0}
							<div class="capitalize">bundle</div>
						{:else}
							<div class="capitalize">{b.Name}</div>
						{/if}
					</li>
				{/each}
			{/if}
		</ul>
	</div>
</div>

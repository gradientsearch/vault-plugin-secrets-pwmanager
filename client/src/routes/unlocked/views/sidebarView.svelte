<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import Icon from '../../../components/icon.svelte';
	import { KVBundleService } from '../services/bundle.service';
	import Modal from '../../../components/modal.svelte';
	import BundleModal from '../modals/bundle.svelte';

	let { bundle = $bindable(), zarf = $bindable(), clientHeight } = $props();
	let bundles: Bundle[] = $state([]);
	let showModal = $state(false);
	let newBundle: Bundle | undefined = $state();

	$effect(() => {
		newBundle;
		untrack(() => {
			console.log('new bundle', newBundle?.Name);
			if (newBundle !== undefined) {
				let copyBundle = Object.assign({}, newBundle);
				bundles.push(copyBundle);
				newBundle = undefined;
			}
		});
	});

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

	let headerHeight = $state(0);
</script>

<div
	style="max-height: {clientHeight}px"
	class="w-full bg-token_side_nav_color_surface_primary md:max-w-64 lg:block"
>
	<div class="flex-row">
		<div bind:clientHeight={headerHeight} class="flex-column flex content-start p-2">
			<Icon className="nav-header-icon"></Icon>
			<div
				class="height-[36px] px-[8px] py-[9px] text-sm font-bold text-token_side_nav_color_foreground_faint"
			>
				<div class="flex flex-row items-center justify-between">
					<span>Bundles</span>

					<button
						onclick={() => {
							showModal = true;
							console.log('clicked', showModal);
						}}
					>
						<span class="p-2 text-sm hover:bg-surface_interactive_hover hover:text-blue-300">
							Add Bundle +
						</span>
					</button>
				</div>
			</div>
		</div>
		<ul style="max-height: {clientHeight - headerHeight}px; height: {clientHeight - headerHeight}px" class="overflow-hidden overflow-y-scroll">
			{#if bundles}
				{#each bundles as b}
					<!-- svelte-ignore a11y_click_events_have_key_events -->
					<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
					<li
						onclick={() => {
							bundle = b;
						}}
						class="{b.Path === bundle.Path
							? 'bg-token_side_nav_color_surface_interactive_active'
							: ''} height-[36px] my-2 min-h-[36px] rounded-lg px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_strong hover:cursor-pointer hover:bg-token_side_nav_color_surface_interactive_hover"
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
<Modal bind:showModal>
	<BundleModal edit={true} bind:showModal bind:newBundle bind:zarf></BundleModal>
</Modal>

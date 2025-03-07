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
			if (newBundle !== undefined) {
				let copyBundle = Object.assign({}, newBundle);
				bundles.push(copyBundle);
				bundle = copyBundle;
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
			let [bs, err] = await KVBundleService.getBundles(zarf);
			if (err !== undefined) {
				console.log(`error listing bundles ${err}`);
			}

			if (bs !== undefined) {
				bundles = [...bundles, ...bs];
			}
		})();
	});
</script>

<div class="h-[100vh] w-full bg-token_side_nav_color_surface_primary md:max-w-64 lg:block">
	<div class="flex-row">
		<div class=" flex h-[64px] flex-row content-start px-1 py-2 bg-token_side_nav_color_surface_primary">
			<Icon className="nav-header-icon max-w-16"></Icon>
			<span class="flex flex-1"></span>
			<div
				class="height-[36px] px-[8px] py-[9px] text-sm font-bold text-token_side_nav_color_foreground_faint"
			></div>
		</div>
		<ul class="h-[calc(100vh-64px)] overflow-hidden overflow-y-scroll">
			<li
				class="height-[36px] my-2 flex min-h-[36px] flex-row px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_faint"
			>
				<span class="flex items-end">Bundles</span>
				<span class="flex flex-grow"></span>
				<button
					onclick={() => {
						showModal = true;
					}}
				>
					<span
						class="p-2 text-3xl hover:bg-surface_interactive_hover hover:text-blue-300"
					>
						+
					</span>
				</button>
			</li>
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
							: ''} height-[36px] my-2 min-h-[36px] px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_strong hover:cursor-pointer hover:bg-token_side_nav_color_surface_interactive_hover"
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

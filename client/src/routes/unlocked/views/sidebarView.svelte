<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import Icon from '../../../components/icon.svelte';
	import { KVBundleService } from '../services/bundle.service';
	import Modal from '../../../components/modal.svelte';
	import CreateModal from '../modals/createBundle.svelte';
	import { userService } from '../services/user.service';

	let { bundle = $bindable(), zarf = $bindable() } = $props();

	let bundles: Bundle[] = $state([]);
	let sharedBundles: Bundle[] = $state([]);
	let showModal = $state(false);
	let newBundle: Bundle | undefined = $state();
	let username: string = $state('');

	$effect(() => {
		newBundle;
		untrack(() => {
			addNewBundle();
		});
	});

	onMount(() => {
		username = userService.getUsername();
		addPersonalBundle();
		getBundles();
	});

	/**
	 * addNewBundle is called when a user successfully
	 * creates a new bundle.
	 */
	function addNewBundle() {
		if (newBundle !== undefined) {
			let copyBundle = Object.assign({}, newBundle);
			bundles.push(copyBundle);
			bundle = copyBundle;
			newBundle = undefined;
		}
	}

	/**
	 * addPersonalBundle creates a users personal bundle
	 */
	function addPersonalBundle() {
		let entityID = userService.getEntityID();
		if (entityID !== undefined) {
			let b: Bundle = {
				Type: 'bundle',
				Path: `bundles/data/${entityID}/${entityID}`,
				Name: 'personal',
				Owner: entityID,
				IsAdmin: true,
				ID: entityID,
				Users: []
			};
			bundles?.push(b);
			bundle = b;
		}
	}

	/**
	 * getBundles retrieves the users bundles from the pwmanager
	 */
	async function getBundles() {
		let [bs, err] = await KVBundleService.getBundles(zarf);
		if (err !== undefined) {
			console.log(`error listing bundles ${err}`);
		}

		if (bs !== undefined && bs.bundles) {
			bundles = [...bundles, ...bs.bundles];
		}

		if (bs !== undefined && bs.sharedBundles) {
			sharedBundles = [...bs.sharedBundles];
		}
	}

</script>

<div class="h-[100vh] w-full bg-token_side_nav_color_surface_primary md:max-w-64 lg:block">
	<div class="flex-row">
		<div
			class=" flex h-[64px] flex-row content-start bg-token_side_nav_color_surface_primary px-1 py-2"
		>
			<Icon className="nav-header-icon max-w-16"></Icon>
			<span class="flex flex-1"></span>
			<div
				class="height-[36px] px-[8px] py-[9px] text-sm font-bold text-token_side_nav_color_foreground_faint"
			>
				{username}
			</div>
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
					<span class="p-2 text-3xl hover:bg-surface_interactive_hover hover:text-blue-300">
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
						{#if b?.Name?.length === 0}
							<div class="capitalize">bundle</div>
						{:else}
							<div class="">{b?.Name}</div>
						{/if}
					</li>
				{/each}
			{/if}

			<li
				class="height-[36px] my-2 flex min-h-[36px] flex-row px-[8px] py-[9px] text-sm text-token_side_nav_color_foreground_faint"
			>
				<span class="flex items-end">Shared</span>
			</li>
			{#if sharedBundles}
				{#each sharedBundles as b}
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
						{#if b?.Name?.length === 0}
							<div class="capitalize">bundle</div>
						{:else}
							<div class="">{b?.Name}</div>
						{/if}
					</li>
				{/each}
			{/if}
		</ul>
	</div>
</div>
<Modal bind:showModal>
	<CreateModal edit={true} bind:showModal bind:newBundle bind:zarf></CreateModal>
</Modal>

<script lang="ts">
	import { base } from '$app/paths';
	import Password from './password.svelte';

	let { bundleService = $bindable(), showModal = $bindable() } = $props();

	let selectedType: any = $state(undefined);
	let containerHeight: number | undefined = $state();
	let headerHeight: number | undefined = $state();
	let types = [
		{
			name: 'password',
			src: '🔑',
			component: Password
		}
		// {
		// 	name: 'note',
		// 	src: 'note.png'
		// },
		// {
		// 	name: 'ssh',
		// 	src: 'ssh.png'
		// }
	];

	function onSelectedType(t: any) {
		selectedType = t;
	}
</script>

{#if selectedType === undefined}
	<header class="border-b-2 border-b-border_primary bg-page_primary p-4">
		<h1 class="text-start font-bold">Create Password</h1>
	</header>

	<div class="4xl:grid-cols-4 grid grid-cols-1 gap-4 p-4 sm:grid-cols-2 2xl:grid-cols-3">
		{#each types as t}
			<button
				onclick={() => onSelectedType(t)}
				class="relative flex min-h-32 flex-row items-center justify-items-start bg-page_faint shadow-lg hover:bg-surface_interactive_hover hover:text-foreground_action_hover"
			>
				<span class="p-3 text-3xl">{t.src}</span>
				<span class="left-0 right-0 flex-1 text-center md:absolute">{t.name}</span>
			</button>
		{/each}
	</div>
{:else}
	{@const Component = selectedType.component}
	<div bind:clientHeight={containerHeight} class="h-[100%]">
		<header bind:clientHeight={headerHeight}>
			<div class="flex flex-row justify-start border-b-2 border-border_primary p-4 font-bold">
				<div>New Item</div>
			</div>
		</header>
		<Component
			bind:bundleService
			bind:showModal
			clientHeight={containerHeight - headerHeight}
			cancel={() => {
				selectedType = undefined;
			}}
		/>
	</div>
{/if}

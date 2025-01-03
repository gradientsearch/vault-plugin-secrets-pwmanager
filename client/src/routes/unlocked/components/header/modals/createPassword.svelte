<script lang="ts">
	import { base } from '$app/paths';
	import Password from './types/password.svelte';

	let selectedType: any = $state(undefined);

	let types = [
		{
			name: 'password',
			src: 'key.svg',
			component: Password
		},
		{
			name: 'note',
			src: 'note.png'
		},
		{
			name: 'ssh',
			src: 'ssh.png'
		}
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
				<img class="h-16" src="{base}/icons/{t.src}" alt="{t.name} icon" />
				<span class="left-0 right-0 flex-1 text-center md:absolute">{t.name}</span>
			</button>
		{/each}
	</div>
{:else}
	{@const Component = selectedType.component}
	<div class="overflow-y-scroll">
		<Component />
	</div>
{/if}

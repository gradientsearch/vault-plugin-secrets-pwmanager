<script lang="ts">
	import { base } from '$app/paths';

	let types = [
		{
			name: 'password',
			src: 'key.svg'
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
</script>

<header class="border-b-2 border-b-border_primary bg-page_primary p-4">
	<h1 class="text-start font-bold">Create Password</h1>
</header>

<div class="grid 4xl:grid-cols-4 2xl:grid-cols-3 sm:grid-cols-2 grid-cols-1 gap-4 p-4">
    {#each types as t}
	<button
		class="relative flex min-h-32 flex-row items-center justify-items-start bg-page_faint shadow-lg hover:bg-surface_interactive_hover hover:text-foreground_action_hover"
	>
		<img class="h-16" src="{base}/icons/{t.src}" alt="{t.src} icon" />
		<span class="left-0 right-0 flex-1 text-center md:absolute">{t.name}</span>
	</button>
    {/each}
</div>

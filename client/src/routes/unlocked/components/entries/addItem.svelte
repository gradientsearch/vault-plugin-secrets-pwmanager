<script lang="ts">
	import { onMount } from 'svelte';

	let { fn = $bindable() } = $props();

	let showMenu = $state(false);
	let clientHeight = $state(0);
	let bottom = $state('');

	onMount(() => {});

	let inputTypes = [
		'text',
		'password',
		'date',
		'text',
		'password',
		'date',
		'text',
		'password',
		'date',
		'text',
		'password',
		'date'
	];
</script>

<div bind:clientHeight class="absolute right-0">
	<div class="z-10 flex h-full w-10 items-center justify-center hover:cursor-pointer">
		<button
			class="flex w-[100%] items-center justify-center text-base"
			onclick={(event) => {
				if (Math.abs(event.clientY - window.innerHeight) < 150) {
					bottom = 'bottom-0';
				} else {
					bottom = 'top-5';
				}
				showMenu = !showMenu;
			}}
			aria-label="menu"
		>
			➕
		</button>
	</div>
	{#if showMenu}
		<div
			style="transform: translate3d(0px, 0px, 0px);"
			class="border-gray-100 bg-white absolute end-0 z-50 w-32 rounded-md border bg-page_faint shadow-lg {bottom} right-8 h-32 overflow-y-scroll"
			role="menu"
		>
			{#each inputTypes as i}
				<div class="">
					<button
						onclick={() => {
							fn(i);
							showMenu = false;

							setTimeout(() => {
								let el = document.getElementById('entry-view');
								el?.scroll({ top: el.scrollHeight, behavior: 'smooth' });
							}, 150);
						}}
						class="text-gray-500 hover:text-gray-700 block w-full rounded-md rounded-b-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
						role="menuitem"
					>
						{i}
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>

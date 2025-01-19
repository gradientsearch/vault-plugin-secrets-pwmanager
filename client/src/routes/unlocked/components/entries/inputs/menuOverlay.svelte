<script lang="ts">
	import { onMount } from "svelte";

	let { inputHeight =$bindable(), type = $bindable<string>(), value} = $props();

	let showMenu = $state(false);

	let inputType: string | undefined = $state();
    let reveal = $state(false)

	onMount(() => {
		inputType = type;
	});

	function copyToClipBoard(text: string) {
		navigator.clipboard.writeText(text).then(
			function () {
				console.log('Async: Copying to clipboard was successful!');
			},
			function (err) {
				console.error('Async: Could not copy text: ', err);
			}
		);
	}
</script>

<div class="absolute right-0">
	<div
		style="min-height: {inputHeight}px"
		class="flex h-full w-10 items-center justify-center hover:cursor-pointer z-10"
	>
		<button
			style="min-height: {inputHeight}px;"
			class="flex w-[100%] items-center justify-center"
			onclick={() => {
				showMenu = !showMenu;
			}}
			aria-label="menu"
		>
			<svg
				version="1.1"
				id="svg2"
				width="3.3915606"
				height="14.342"
				viewBox="0 0 3.3915606 14.342"
				xmlns="http://www.w3.org/2000/svg"
			>
				<defs id="defs6" />
				<g id="g8" transform="translate(-6.6792989,-1.8920578)">
					<ellipse
						style="fill:#000000"
						id="path224"
						cx="8.3750792"
						cy="3.5880578"
						rx="1.6957803"
						ry="1.696"
					/>
					<ellipse
						style="fill:#000000"
						id="ellipse226"
						cx="8.3750792"
						cy="9.0630579"
						rx="1.6957803"
						ry="1.696"
					/>
					<ellipse
						style="fill:#000000"
						id="ellipse228"
						cx="8.3750792"
						cy="14.538058"
						rx="1.6957803"
						ry="1.696"
					/>
				</g>
			</svg>
		</button>
	</div>
	{#if showMenu}
		<div
			style="transform: translate3d(0px, -{inputHeight / 4}px, 0px);"
			class="border-gray-100 bg-white absolute end-0 w-32 rounded-md border bg-page_faint shadow-lg z-50"
			role="menu"
		>
			<div class="">
				<button
					onclick={() => {
						copyToClipBoard(value);
						showMenu = false;
					}}
					class="text-gray-500 hover:text-gray-700 block w-full rounded-md rounded-b-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
					role="menuitem"
				>
					Copy
				</button>

				{#if inputType === 'password'}
					{#if !reveal}
						<button
							onclick={() => {
								type = 'text';
								showMenu = false;
                                reveal = true
							}}
							class="text-gray-500 hover:bg-gray-50 hover:text-gray-700 block w-full rounded-md rounded-t-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
							role="menuitem"
						>
							Reveal
						</button>
					{:else}
						<button
							onclick={() => {
								type = 'password';
								showMenu = false;
                                reveal = false
							}}
							class="text-gray-500 hover:bg-gray-50 hover:text-gray-700 block w-full rounded-md rounded-t-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
							role="menuitem"
						>
							Conceal
						</button>
					{/if}
				{/if}
			</div>
		</div>
	{/if}
</div>

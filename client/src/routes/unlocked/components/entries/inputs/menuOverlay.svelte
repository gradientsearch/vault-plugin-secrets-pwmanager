<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import { MODE } from '../../../models/entry';
	import type { Input } from '../../../models/input';
	import { getSecureRandomString } from '$lib/helper';
	import Button from '../../../../../components/button.svelte';

	let {
		input = $bindable<Input>(),
		inputHeight = $bindable(),
		inputType = $bindable(),
		mode,
		isCore,
		idx,
		onDeleteItem,
		id
	} = $props();

	let showMenu = $state(false);
	let showGeneratePasswordMenu = $state(false);
	let reveal = $state(false);

	let generatedPassword = $state('');
	let generatedPasswordLength = $state(20);
	$effect(() => {
		id;
		untrack(() => {
			reveal = false;
		});
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

	function onGeneratePassword() {
		showMenu = false;
		showGeneratePasswordMenu = true;

		generatedPassword = getSecureRandomString(generatedPasswordLength);
	}

	function onUseGeneratedPassword() {
		input.Value = generatedPassword;
		showGeneratePasswordMenu = false;
	}
</script>

<div class="absolute right-0">
	<div
		style="min-height: {inputHeight}px"
		class="z-10 flex h-full w-10 items-center justify-center hover:cursor-pointer"
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
			style="transform: translate3d(-10px, -{inputHeight / 4}px, 0px);"
			class="border-gray-100 bg-white absolute end-0 z-50 w-64 rounded-md border bg-page_faint shadow-lg"
			role="menu"
		>
			<div class="">
				<button
					onclick={() => {
						copyToClipBoard(input.Value);
						showMenu = false;
					}}
					class="text-gray-500 hover:text-gray-700 block w-full rounded-md rounded-b-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
					role="menuitem"
				>
					Copy
				</button>

				{#if input.Type === 'password'}
					{#if !reveal}
						<button
							onclick={() => {
								inputType = 'text';
								showMenu = false;
								reveal = true;
							}}
							class="text-gray-500 hover:bg-gray-50 hover:text-gray-700 block w-full rounded-md rounded-t-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
							role="menuitem"
						>
							Reveal
						</button>
					{:else}
						<button
							onclick={() => {
								inputType = 'password';
								showMenu = false;
								reveal = false;
							}}
							class="text-gray-500 hover:bg-gray-50 hover:text-gray-700 block w-full rounded-md rounded-t-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
							role="menuitem"
						>
							Conceal
						</button>
					{/if}
				{/if}
				{#if mode === MODE.EDIT && input.Type === 'password'}
					<button
						onclick={() => {
							onGeneratePassword();
						}}
						class="text-gray-500 hover:text-gray-700 block w-full rounded-md rounded-b-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
						role="menuitem"
					>
						Generate Password
					</button>
				{/if}
				{#if !isCore && mode === MODE.EDIT}
					<button
						onclick={() => {
							onDeleteItem(idx);
							showMenu = false;
						}}
						class="text-gray-500 hover:text-gray-700 block w-full rounded-md rounded-b-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
						role="menuitem"
					>
						Delete
					</button>
				{/if}
			</div>
		</div>
	{/if}

	{#if showGeneratePasswordMenu}
		<div
			style="min-height: {inputHeight}px"
			class="flex-column z-10 flex h-full w-10 items-center justify-center hover:cursor-pointer"
		>
			<div
				style="transform: translate3d(-10px, -{inputHeight / 4}px, 0px);"
				class="border-gray-100 bg-white absolute end-0 z-50 w-64 rounded-md border bg-page_faint p-2 shadow-lg"
				role="menu"
			>
				<label for="generated">
					<span class="text-gray-700 text-sm font-medium"> Generated </span>

					<div class="flex flex-row">
						<input
							type="text"
							id="generated"
							bind:value={generatedPassword}
							class="border-gray-300 rounded shadow-sm sm:text-sm"
						/>
						<button onclick={onGeneratePassword} class="text-lg p-2">🔄</button>
					</div>
				</label>

				<div>
					<label for="length">
						<span class="text-gray-700 text-sm font-medium"> Length </span>

						<input
							type="number"
							id="length"
							bind:value={generatedPasswordLength}
							class="border-gray-300 rounded shadow-sm sm:text-sm"
						/>
					</label>
				</div>

				<div class="pt-4">
					<button
						class=" font-bol z-50
			 inline-block w-[76px] rounded
				border bg-blue-200 py-2 text-center text-sm hover:bg-blue-400"
						onclick={onUseGeneratedPassword}
					>
						Use
					</button>
					<button
						class="rounded z-50 w-[76px]
			 border border-border_strong bg-surface_faint py-2 text-center text-sm font-bold text-foreground_primary hover:bg-foreground_high_contrast"
						onclick={() => (showGeneratePasswordMenu = false)}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>

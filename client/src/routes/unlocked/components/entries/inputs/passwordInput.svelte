<!-- 
@component

## PasswordInput
This is the component used to show the password input element

### Props
- `label`: the label shown above the input text field i.e username, password
- `placeholder`: what the user sees before entering in a value to the input field
- `value`: bindable value of the user provided text
- `idx`: the section index of the input field 0 index and last index have rounded corners
- `last`: the section index of the 

### Example
<Password bind:entry></Password>

-->


<script lang="ts">
	import ItemInput from './itemInput.svelte';
	let reveal = $state(false);
	let inputHeight = $state(0);
	let { label, placeholder, value = $bindable(), idx, last, id } = $props();


	/** effect needed to conceal passwords when an entry changes
	 * this happens if the same input type is displayed in the same
	 * position the element is not recreated. This sets default inputType
	 * back to password when an entry id is updated.
	*/
	$effect(() => {
		id;

		inputType = 'password'
	})

	let viewState = $state('edit');
	let showMenu = $state(false);
	let inputType = $state('password');

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

<div class="relative flex flex-row">
	{#if viewState === 'edit'}
		<div bind:clientHeight={inputHeight} class=" w-full flex-1">
			<ItemInput type={inputType} {label} {placeholder} bind:value {idx} {last} id={id}></ItemInput>
		</div>

		<div class="absolute right-0"></div>
		<div class="absolute right-0 z-50">
			<div
				style="min-height: {inputHeight}px"
				class="flex h-full w-10 items-center justify-center hover:cursor-pointer"
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
					class="border-gray-100 bg-white absolute end-0 z-50 w-32 rounded-md border bg-page_faint shadow-lg"
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
							<button
								onclick={() => {
									inputType = 'text';
									showMenu = false;
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
								}}
								class="text-gray-500 hover:bg-gray-50 hover:text-gray-700 block w-full rounded-md rounded-t-none px-4 py-2 text-start text-sm hover:bg-surface_interactive_hover"
								role="menuitem"
							>
								Conceal
							</button>
						{/if}
					</div>
				</div>
			{/if}
		</div>
	{:else if viewState === 'new'}
		<div bind:clientHeight={inputHeight} class=" w-full flex-1">
			<ItemInput type={'password'} {label} {placeholder} bind:value {idx} {last} {id}></ItemInput>
		</div>
	{/if}
</div>

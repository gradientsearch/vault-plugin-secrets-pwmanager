<script lang="ts">
	import { onMount } from 'svelte';
	import { MODE } from '../../../models/entry';
	import MenuOverlay from './menuOverlay.svelte';
	import type { Input } from '../../../models/input';

	let {
		input = $bindable<Input>(),
		idx,
		last,
		id,
		mode = $bindable<MODE>(),
		isCore,
		onDeleteItem
	} = $props();

	let inputHeight = $state(0);
	let inputType: string | undefined = $state();


	$effect(() => {
		id;
		inputType = input.Type
	})
	onMount(() => {
		inputType = input.Type;
	});
</script>

<div class="relative flex flex-row">
	<div bind:clientHeight={inputHeight} class=" w-full flex-1">
		<label
			for="username"
			class="border-gray-200 focus-within:border-blue-600 focus-within:ring-blue-600 block overflow-hidden border-x border-b {idx ===
			0
				? 'rounded-t-md border-t'
				: ''} {last ? 'rounded-b-md' : ''} px-3 py-2 shadow-sm focus-within:ring-1"
		>
			<div class="text-gray-700 text-xs font-medium">
				<span class="w-100 flex">
					{input.Label}
				</span>

				<input
					id={input.Label + '-' + idx}
					style={input.Type === 'date' ? 'width: auto;' : ''}
					class="focus:border-transparent mt-1 w-full border-none p-0 focus:outline-none focus:ring-0 sm:text-sm"
					disabled={mode === MODE.VIEW}
					autocomplete="off"
					type={inputType}
					placeholder={input.Placeholder}
					bind:value={input.Value}
				/>
			</div>
		</label>
	</div>
	<MenuOverlay bind:input bind:inputHeight bind:inputType {mode} {isCore} {idx} {onDeleteItem} {id}></MenuOverlay>
</div>

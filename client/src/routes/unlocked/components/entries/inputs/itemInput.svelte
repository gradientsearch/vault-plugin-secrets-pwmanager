<script lang="ts">
	import { onMount } from 'svelte';
	import { MODE } from '../../../models/entry';
	import MenuOverlay from './menuOverlay.svelte';

	let {
		label,
		type,
		placeholder,
		value = $bindable<string>(),
		idx,
		last,
		id,
		mode = $bindable<MODE>(),
		isCore,
		onDeleteItem
	} = $props();

	let inputHeight = $state(0);
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
					{label}
				</span>
				<input
					style={type === 'date' ? 'width: auto;' : ''}
					autocomplete="off"
					{type}
					id={label + '-' + idx}
					{placeholder}
					class="focus:border-transparent mt-1 w-full border-none p-0 focus:outline-none focus:ring-0 sm:text-sm"
					bind:value
					disabled={mode === MODE.VIEW}
				/>
			</div>
		</label>
	</div>
	<MenuOverlay bind:inputHeight bind:type {value} {mode} {isCore} {idx} {onDeleteItem} {id}></MenuOverlay>
</div>

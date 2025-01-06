<script lang="ts">
	import { base } from '$app/paths';
	import Button from '../../../../../../components/button.svelte';
	import { newPasswordItem, type PasswordItem } from '../../../../models/input';
	import ItemInput from './components/itemInput.svelte';

	let vm: PasswordItem = $state(newPasswordItem());

	function onSave() {}

	function onBack() {}
</script>

<div class="flex h-[100%] flex-col">
	<header class="p-4">
		<div class="mt-3 flex flex-row">
			<img class="max-h-12" src="{base}/icons/key.svg" alt="key icon" />
			<input
				type="text"
				multiple
				class="form-input mt-1 block w-full"
				placeholder="Password"
				bind:value={vm.Name}
			/>
		</div>
	</header>

	<div class="block p-4">
		<div class="text-md grid min-w-96 grid-cols-1">
			{#each vm.Core.Items as v, idx}
				<ItemInput type={v.Type} label={v.Label} placeholder={v.Placeholder} bind:value={v.Value} {idx} last={vm.Core.Items.length - 1 === idx}
				></ItemInput>
			{/each}
		</div>
	</div>
	<span class="flex flex-1"></span>
	<div>
		<Button fn={onSave}>Save</Button>
		<Button fn={onBack}>Back</Button>
	</div>
</div>

<style>
	label {
		/* color: #1f2124;
		font-size: 14px;
		font-weight: 700;
		align-items: center;
		gap: 4px;
		width: min-content;
		min-width: 100%; */
	}
</style>

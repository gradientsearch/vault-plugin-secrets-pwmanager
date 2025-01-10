<script lang="ts">
	import { base } from '$app/paths';
	import { Api } from '$lib/api';
	import Button from '../../../../../../components/button.svelte';
	import { newPasswordItem, type Metadata, type PasswordItem } from '../../../../models/input';
	import { getComponent } from './types';

	let {
		passwordListService = $bindable(),
		clientHeight = $bindable(),
		showModal = $bindable(),
		cancel
	} = $props();
	let pi: PasswordItem = $state(newPasswordItem());

	async function onSave() {
		let meta: Metadata = {
			Name: pi.Name,
			Type: 'password',
			// Important to note that the 0 indexed value for Password item is username
			// if this were to to the 1 index that would be the password!
			Value: pi.Core.Items[0].Value
		};
		pi.Metadata = meta;
		let err = await passwordListService.add(pi);
		if (err !== undefined) {
			// alert user there was a problem saving the password
		} else {
			showModal = false
		}
	}

	function onCancel() {
		cancel();
	}
</script>

<div class="flex flex-col" style="height: {clientHeight}px;">
	<header class="p-4">
		<div class="mt-3 flex flex-row">
			<img class="max-h-12" src="{base}/icons/key.svg" alt="key icon" />
			<input
				type="text"
				multiple
				class="form-input mt-1 block w-full"
				placeholder="Password"
				bind:value={pi.Name}
			/>
		</div>
	</header>

	<div class="block p-4">
		<div class="text-md grid min-w-96 grid-cols-1">
			{#each pi.Core.Items as v, idx}
				{@const Component = getComponent(v.Type)}
				<Component
					type={v.Type}
					label={v.Label}
					placeholder={v.Placeholder}
					bind:value={v.Value}
					{idx}
					last={pi.Core.Items.length - 1 === idx}
				/>
				<!-- <ItemInput type={v.Type} label={v.Label} placeholder={v.Placeholder} bind:value={v.Value} {idx} last={vm.Core.Items.length - 1 === idx}
				></ItemInput> -->
			{/each}
		</div>
	</div>
	<span class="flex flex-1"></span>
	<div class="p-4">
		<Button fn={onSave}>Save</Button>
		<Button fn={onCancel} primary={false}>Cancel</Button>
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

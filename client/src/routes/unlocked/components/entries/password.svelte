<!-- 
@component
## Password
This is the component used to show the password Entry.

### Props
- `entry`: The generalized data structure pwmanager uses to define passwords i.e entry could be a password, login, secure note.

### Example
<Password bind:entry></Password>

-->

<script lang="ts">
	import { base } from '$app/paths';
	import Button from '../../../../components/button.svelte';
	import { MODE, type Metadata } from '../../models/entry';
	import type { BundleService } from '../../services/bundle.service';
	import { getInputComponent } from '../entries/components';

	let {
		entry = $bindable(),
		bundleService = $bindable<BundleService>(),
		mode = $bindable<MODE>(),
		cancel = () => {}
	} = $props();


	async function onSave() {
		let meta: Metadata = {
			Name: entry.Name,
			Type: 'password',
			// Important to note that the 0 indexed value for Password item is username
			// if this were to to the 1 index that would be the password!
			Value: entry.Core.Items[0].Value,
			ID: ''
		};
		entry.Metadata = meta;
		console.log(bundleService);
		let err = await bundleService.addEntry(entry);
		if (err !== undefined) {
			console.log('err: ', err);
			// alert user there was a problem saving the password
		} else {
		}
	}

	function onCancel() {
		cancel();
	}
</script>

<form class="flex h-[100%] flex-col">
	<header class="p-4">
		<div class="mt-3 flex flex-row">
			<img class="max-h-12" src="{base}/icons/key.svg" alt="key icon" />
			<input
				type="text"
				multiple
				class="form-input mt-1 block w-full"
				placeholder="Password"
				bind:value={entry.Name}
			/>
		</div>
	</header>

	<div class="block p-4">
		<div class="text-md grid min-w-96 grid-cols-1">
			{#each entry.Core.Items as v, idx}
				{@const Component = getInputComponent(v.Type)}
				<Component
					type={v.Type}
					label={v.Label}
					placeholder={v.Placeholder}
					bind:value={v.Value}
					{idx}
					last={entry.Core.Items.length - 1 === idx}
					id={entry.Metadata.ID}
					mode={mode}
				/>
			{/each}
		</div>
	</div>

	<span class="flex flex-1"></span>
	{#if mode === MODE.EDIT}
		<div class="p-4">
			<Button fn={onSave}>Save</Button>
			<Button fn={onCancel} primary={false}>Cancel</Button>
		</div>
	{/if}
</form>

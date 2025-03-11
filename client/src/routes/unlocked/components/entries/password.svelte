<!-- 
@component
## Password
This is the component used to show the password Entry.

### Props
- `entry`: The generalized data structure pwmanager uses to define passwords i.e entry could be a password, login, secure note.
- `bundleService`: the service responsible for making requests to the server.
- `mode`: sets the current mode i.e EDIT or VIEW.
- `cancel`: a function called if a user is editing an entry and then cancels. Used to reset entry.
- `save`: a function called on save entry.
### Example
<Password bind:entry bind:bundleService></Password>

-->

<script lang="ts">
	import Button from '../../../../components/button.svelte';
	import type { BundleMetadata } from '../../models/bundle/vault/metadata';
	import { MODE, type Entry, type Metadata } from '../../models/entry';
	import { DateInput, PasswordInput, TextInput, type Input } from '../../models/input';
	import type { BundleService } from '../../services/bundle.service';
	import { getInputComponent } from '../entries/components';
	import AddItem from './addItem.svelte';

	let {
		entry = $bindable(),
		bundleService = $bindable<BundleService>(),
		bundleMetadata = $bindable<BundleMetadata>(),
		mode = $bindable<MODE>(),
		cancel = $bindable(),
		save = () => {}
	} = $props();

	async function onSave() {
		let meta: Metadata = {
			Name: entry.Name,
			Type: 'password',
			// Important to note that the 0 indexed value for Password item is username
			// if this were to be the 1 index that would be the password!
			Value: entry.Core.Items[0].Value,
			ID: entry.Metadata.ID ? entry.Metadata.ID : ''
		};
		entry.Metadata = meta;

		let err = await bundleService.putEntry(entry, bundleMetadata);
		if (err !== undefined) {
			console.log('err: ', err);

			// alert user there was a problem saving the password
		} else {
			save();
			mode = MODE.VIEW;
		}
	}

	function onCancel() {
		cancel();
		mode = MODE.VIEW;
	}

	function addItem(itemType: string) {
		let input: Input;
		switch (itemType) {
			case 'password':
				input = new PasswordInput();
				break;
			case 'text':
				input = new TextInput();
				break;
			case 'date':
				input = new DateInput();
				break;
			default:
				return;
		}

		entry.More.Items.push(input);
	}

	function onDeleteItem(idx: number) {
		if (idx >= 0 && idx < entry.More.Items.length) {
			entry.More.Items.splice(idx, 1);
		}
	}
</script>

<form class="flex h-[100%] flex-col">
	<header class="p-4">
		<div class="mt-3 flex flex-row">
			<span class="pe-3 text-3xl">🔑</span>
			<input
				type="text"
				multiple
				class="form-input mt-1 block w-full"
				placeholder="Password"
				bind:value={entry.Name}
				disabled={mode === MODE.VIEW}
			/>
		</div>
	</header>

	<div class="block p-4">
		<div class="text-md grid min-w-96 grid-cols-1">
			{#each entry.Core.Items as v, idx}
				{@const Component = getInputComponent(v.Type)}
				<Component
					bind:input={entry.Core.Items[idx]}
					{idx}
					last={entry.Core.Items.length - 1 === idx}
					id={entry.Metadata.ID}
					{mode}
					isCore={true}
					onDeleteItem
				/>
			{/each}

			{#if entry.More.Items.length > 0}
				<span class="h-20"></span>
			{/if}

			{#each entry.More.Items as v, idx}
				{@const Component = getInputComponent(v.Type)}
				<Component
					bind:input={entry.More.Items[idx]}
					{idx}
					last={entry.More.Items.length - 1 === idx}
					id={entry.Metadata.ID}
					{mode}
					isCore={false}
					{onDeleteItem}
				/>
			{/each}
		</div>

		<div
			class="relative flex flex-row p-2"
			style="visibility: {mode === MODE.EDIT ? 'visible' : 'hidden'};"
		>
			<span class="flex-1"></span>
			<AddItem fn={addItem} {mode}></AddItem>
		</div>
	</div>

	<span class="flex flex-1"></span>
	{#if mode === MODE.EDIT}
		<div class="p-4">
			<Button fn={onSave}>Save</Button>
			<Button fn={onCancel} primary={false}>Cancel</Button>
		</div>
	{:else}
		<div class="p-4">
			<Button
				fn={() => {
					mode = MODE.EDIT;
				}}>Edit</Button
			>
			<Button
				primary={false}
				fn={() => {
					bundleService.deleteEntry(entry.Metadata.ID);
				}}>Delete</Button
			>
		</div>
	{/if}
</form>

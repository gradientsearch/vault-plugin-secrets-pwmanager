<script lang="ts">
	import Button from '../../../components/button.svelte';
	import { newEntry, type Metadata, type Entry } from '../models/input';
	import Password from '../components/entries/password.svelte';

	let {
		passwordListService = $bindable(),
		clientHeight = $bindable(),
		showModal = $bindable(),
		cancel
	} = $props();
	let entry: Entry = $state(newEntry());

	async function onSave() {
		let meta: Metadata = {
			Name: entry.Name,
			Type: 'password',
			// Important to note that the 0 indexed value for Password item is username
			// if this were to to the 1 index that would be the password!
			Value: entry.Core.Items[0].Value,
			Path: ''
		};
		entry.Metadata = meta;
		let err = await passwordListService.add(entry);
		if (err !== undefined) {
			// alert user there was a problem saving the password
		} else {
			showModal = false;
		}
	}

	function onCancel() {
		cancel();
	}
</script>

<div class="flex flex-col" style="height: {clientHeight}px;">
	<Password bind:entry></Password>
	<span class="flex flex-1"></span>
	<div class="p-4">
		<Button fn={onSave}>Save</Button>
		<Button fn={onCancel} primary={false}>Cancel</Button>
	</div>
</div>

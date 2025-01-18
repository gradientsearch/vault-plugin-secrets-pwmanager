<script lang="ts">
	import Button from '../../../components/button.svelte';
	import { newPasswordEntry, type Metadata, type Entry } from '../models/entry';
	import Password from '../components/entries/password.svelte';
	import type { BundleService } from '../services/bundle.service';

	let {
		bundleService = $bindable<BundleService>(),
		clientHeight = $bindable<Number>(),
		showModal = $bindable<Boolean>(),
		cancel
	} = $props();
	let entry: Entry = $state(newPasswordEntry());

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
		console.log(bundleService)
		let err = await bundleService.addEntry(entry);
		if (err !== undefined) {
			console.log('err: ', err)
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
	<Password bind:entry bind:bundleService cancel={cancel} state={'create'}></Password>
</div>

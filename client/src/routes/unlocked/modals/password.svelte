<script lang="ts">
	import Button from '../../../components/button.svelte';
	import { newPasswordEntry, type Metadata, type Entry, MODE } from '../models/entry';
	import Password from '../components/entries/password.svelte';
	import type { BundleService } from '../services/bundle.service';
	import type { BundleMetadata } from '../models/bundle/vault/metadata';

	let {
		bundleService = $bindable<BundleService>(),
		bundleMetadata = $bindable<BundleMetadata>(),
		clientHeight = $bindable<Number>(),
		showModal = $bindable<Boolean>(),
		cancel
	} = $props();
	let entry: Entry = $state(newPasswordEntry());

	async function onSave() {
		showModal = false;
	}
	function onCancel() {
		cancel();
	}
</script>

<div class="flex flex-col" style="height: {clientHeight}px;">
	<Password bind:entry bind:bundleMetadata bind:bundleService save={onSave} cancel={onCancel} mode={MODE.EDIT}
	></Password>
</div>

<script lang="ts">
	import Button from '../../../components/button.svelte';
	import { newPasswordEntry, type Metadata, type Entry, MODE } from '../models/entry';
	import Password from '../components/entries/password.svelte';
	import type { BundleService } from '../services/bundle.service';
	import Bundle from '../components/bundles/bundle.svelte';
	import { onMount } from 'svelte';

	let {
		showModal = $bindable<Boolean>(),
        zarf = $bindable(),
        newBundle = $bindable(),
		edit
	} = $props();

	async function onSave() {
		showModal = false;
	}
	function onCancel() {
		showModal = false;
	}
	let containerHeight: number | undefined = $state();
	let headerHeight: number | undefined = $state();

</script>

<div bind:clientHeight={containerHeight} class="h-[100%]">
	<header bind:clientHeight={headerHeight}>
		<div class="flex flex-row justify-start border-b-2 border-border_primary p-4 font-bold">
			<div>New Bundle</div>
		</div>
	</header>

	<div class="flex flex-col" style="height: {containerHeight - headerHeight}px;">
		<Bundle bind:zarf bind:newBundle save={onSave} cancel={onCancel} {edit} ></Bundle>
	</div>
</div>

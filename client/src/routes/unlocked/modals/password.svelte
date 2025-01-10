<script lang="ts">
	import Button from '../../../components/button.svelte';
	import { newPasswordItem, type Metadata, type PasswordItem } from '../models/input';
	import Password from '../components/passwordItems/password.svelte';

	let {
		passwordListService = $bindable(),
		clientHeight = $bindable(),
		showModal = $bindable(),
		cancel
	} = $props();
	let passwordItem: PasswordItem = $state(newPasswordItem());

	async function onSave() {
		let meta: Metadata = {
			Name: passwordItem.Name,
			Type: 'password',
			// Important to note that the 0 indexed value for Password item is username
			// if this were to to the 1 index that would be the password!
			Value: passwordItem.Core.Items[0].Value,
			Path: ''
		};
		passwordItem.Metadata = meta;
		let err = await passwordListService.add(passwordItem);
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
	<Password bind:passwordItem></Password>
	<span class="flex flex-1"></span>
	<div class="p-4">
		<Button fn={onSave}>Save</Button>
		<Button fn={onCancel} primary={false}>Cancel</Button>
	</div>
</div>

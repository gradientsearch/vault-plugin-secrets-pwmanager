<script lang="ts">
	import Button from '../../../../components/button.svelte';
	import { KVBundleService, type BundleService } from '../../services/bundle.service';

	let {
		newBundle = $bindable(),
		zarf = $bindable(),
		edit = false,
		cancel = $bindable(),
		save = () => {}
	} = $props();

	let isEdit = $state(edit);
	async function onSave() {
		let [bundle, err] = await KVBundleService.createBundle(zarf, name);
		if (err !== undefined) {
			console.log(`error creating bundle: ${err}`);
			return;
		}

		newBundle = bundle;
		save();
	}

	function onCancel() {
		cancel();
	}

	let name = $state('');
</script>

<div class="flex flex-row p-4">
	<div class="flex flex-row">
		<div class="text-md grid min-w-96 grid-cols-1 gap-6">
			<label class="block">
				<span class="text-gray-700 text-base font-medium">Name</span>
				<input
					type="text"
					multiple
					class="form-input mt-1 block w-full"
					placeholder="pwmanager"
					bind:value={name}
				/>
			</label>
		</div>
	</div>
</div>

<span class="flex flex-1"></span>
{#if isEdit}
	<div class="p-4">
		<Button fn={onSave}>Save</Button>
		<Button fn={onCancel} primary={false}>Cancel</Button>
	</div>
{:else}
	<div class="p-4">
		<Button
			fn={() => {
				isEdit = !isEdit;
			}}>Edit</Button
		>
	</div>
{/if}

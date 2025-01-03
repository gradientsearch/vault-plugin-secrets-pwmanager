<script lang="ts">
	import { base } from '$app/paths';
	import Button from '../../../../../../components/button.svelte';

	interface Metadata {
		uuid: string;
		name: string;
		type: string;
	}
	interface Password {
		metadata: Metadata;
	}

	let vm: Password = $state({
		metadata: {
			uuid: '',
			name: 'Password',
			type: 'password'
		}
	});

	let userName: string;

    function onSave(){

    }

    function onBack(){

    }
</script>


<header class="p-4">
	<div class="flex flex-row mt-3">
		<img class="max-h-12" src="{base}/icons/key.svg" alt="key icon" />
		<input
			type="text"
			multiple
			class="form-input mt-1 block w-full"
			placeholder="Password"
			bind:value={vm.metadata.name}
		/>
	</div>
</header>

<div class="p-4">
	<div class="text-md grid min-w-96 grid-cols-1 gap-6">
		<label class="block">
			<span class="text-gray-700">username</span>
			<input
				type="text"
				multiple
				class="form-input mt-1 block w-full"
				placeholder=""
				bind:value={userName}
			/>
		</label>
	</div>
</div>

<div>
    <Button fn={onSave}>Save</Button>
    <Button fn={onBack}>Back</Button>
</div>

<style>
	label {
		color: #1f2124;
		font-size: 14px;
		font-weight: 700;
		align-items: center;
		gap: 4px;
		width: min-content;
		min-width: 100%;
	}
</style>

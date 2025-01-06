<script lang="ts">
	import { base } from '$app/paths';
	import Button from '../../../../../../components/button.svelte';
	import type { Input, PasswordItem } from '../../../../models/input';

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

    let inputValues = {

    }

	function onSave() {}

	function onBack() {}

    
</script>

<div class="flex h-[100%] flex-col">
	<header class="p-4">
		<div class="mt-3 flex flex-row">
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

	<div class="block p-4">
		<div class="text-md grid min-w-96 grid-cols-1">
			<!--
  Heads up! 👋

  Plugins:
    - @tailwindcss/forms
-->


			<label
				for="username"
				class="border-gray-200 focus-within:border-blue-600 focus-within:ring-blue-600 block overflow-hidden rounded-md border px-3 py-2 shadow-sm focus-within:ring-1"
			>
				<div class="text-gray-700 text-xs font-medium">
					username

					<input
						type="email"
						id="UserEmail"
						placeholder="anthony@rhcp.com"
						class="focus:border-transparent mt-1 w-full border-none p-0 focus:outline-none focus:ring-0 sm:text-sm"
					/>
				</div>
			</label>
		</div>
	</div>
	<span class="flex flex-1"></span>
	<div>
		<Button fn={onSave}>Save</Button>
		<Button fn={onBack}>Back</Button>
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

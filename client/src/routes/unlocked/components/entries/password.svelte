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
	import { getInputComponent } from '../entries/components';

	let { entry = $bindable(), state = 'new' } = $props();
</script>

<form>
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
				/>
			{/each}
		</div>
	</div>
</form>

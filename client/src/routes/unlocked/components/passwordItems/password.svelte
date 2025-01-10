<script lang="ts">
	import { base } from '$app/paths';
	import { getInputComponent } from '../passwordItems/components';

	let { passwordItem = $bindable() } = $props();

</script>

<header class="p-4">
	<div class="mt-3 flex flex-row">
		<img class="max-h-12" src="{base}/icons/key.svg" alt="key icon" />
		<input
			type="text"
			multiple
			class="form-input mt-1 block w-full"
			placeholder="Password"
			bind:value={passwordItem.Name}
		/>
	</div>
</header>

<div class="block p-4">
	<div class="text-md grid min-w-96 grid-cols-1">
		{#each passwordItem.Core.Items as v, idx}
			{@const Component = getInputComponent(v.Type)}
			<Component
				type={v.Type}
				label={v.Label}
				placeholder={v.Placeholder}
				bind:value={v.Value}
				{idx}
				last={passwordItem.Core.Items.length - 1 === idx}
			/>
		{/each}
	</div>
</div>

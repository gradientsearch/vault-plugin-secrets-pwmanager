<script lang="ts">
	import { on } from 'svelte/events';
	import { getPasswordComponent } from './passwordItems/types';

	let { passwordItem = $bindable() } = $props();

	$effect(() => {
		passwordItem;
	});
</script>

<div class=" w-full content-center border-t-2 border-border_primary bg-page_faint">
	<!-- {@const Component = getPasswordComponent(passwordItem.Type)}
	<Component
    ></Component> -->
</div>

<script>
	import { base } from '$app/paths';

	import Button from './button.svelte';
	import CardContainer from './cardContainer.svelte';

	var { title, body, button, className = '' } = $props();
</script>

<CardContainer {className}>
	<div class="sm:flex sm:justify-between sm:gap-4">
		<div>
			<h3 class="text-gray-900 text-lg font-bold sm:text-xl">{title}</h3>
		</div>
	</div>

	<div class="mt-4">
		<p class="text-gray-500 text-pretty text-sm">
			{body}
		</p>
	</div>

	<!-- Base -->

	{@render button()}
</CardContainer>

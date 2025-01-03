<script lang="ts">
	import Title from '../../../components/title.svelte';

	let showModal: boolean = $state(false);
	function createPassword() {
		showModal = true;
	}
	function closeModal() {
		showModal = false;
	}
	function modalClick(event: Event) {
		event.stopPropagation();
	}
</script>

<header class="sticky top-0 h-[48px] bg-surface_faint p-1">
	<div class="flex justify-end">
		<button onclick={createPassword}>
			<span class="m-4 p-2 text-sm hover:bg-surface_interactive_hover hover:text-blue-300">
				Create password +
			</span>
		</button>
	</div>
</header>

{#if showModal}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<button
		aria-label="modal-overlay"
		onclick={closeModal}
		class="absolute bottom-0 left-0 right-0 top-0 flex flex-1 bg-foreground_strong/30 hover:cursor-default"
	>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div
			onclick={modalClick}
			class="absolute bottom-0 left-0 right-0 top-0 mx-auto my-auto h-[100vh] w-[100vw] bg-foreground_high_contrast shadow-lg lg:h-[80vh] xl:w-[50vw]"
		>
			<h1 class="text-start font-bold">Create Password</h1>
		</div>
	</button>
{/if}

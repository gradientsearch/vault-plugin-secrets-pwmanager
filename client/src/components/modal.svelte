<script lang="ts">
	let { children,  showModal =$bindable() } = $props();

    function closeModal() {
		showModal = false;
	}

    function modalClick(event: Event) {
		event.stopPropagation();
	}

</script>

{#if showModal}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<button
		aria-label="modal-overlay"
		onclick={closeModal}
		class="absolute bottom-0 left-0 right-0 top-0 flex flex-1 bg-foreground_strong/30 hover:cursor-default"
	>
		<div
			onclick={modalClick}
			class="absolute bottom-0 left-0 right-0 top-0 mx-auto my-auto h-[100vh] w-[100vw] bg-foreground_high_contrast shadow-lg lg:h-[80vh] xl:w-[50vw]"
		>
			{@render children?.()}
		</div>
	</button>
{/if}

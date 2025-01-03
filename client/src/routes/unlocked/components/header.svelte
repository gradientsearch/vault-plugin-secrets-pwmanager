<script lang="ts">
	import Modal from '../../../components/modal.svelte';


	let showModal: boolean = $state(false);
	function createPassword() {
		showModal = true;
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

<Modal bind:fnOnClick={modalClick} bind:showModal={showModal}>
	<h1 class="text-start font-bold">Create Password</h1>
</Modal>

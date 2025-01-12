<script lang="ts">
	import Modal from '../../../components/modal.svelte';
	import CreatePassword from '../modals/createPassword.svelte';

	let { bundleService = $bindable() } = $props();
	let showModal: boolean = $state(false);
	function createPassword() {
		showModal = true;
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

<Modal bind:showModal>
	<CreatePassword bind:bundleService bind:showModal></CreatePassword>
</Modal>

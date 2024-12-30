<script lang="ts">
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import Button from '../components/button.svelte';
	import Card from '../components/card.svelte';
	import VaultIcon from '../components/vaultIcon.svelte';
	import VaultIconAndText from '../components/vaultIconAndText.svelte';

	//goto(`${base}/locked`)
</script>

<div class="flex h-full w-full justify-center">
	<div class="flex-row">
		<VaultIconAndText className="mt-14"></VaultIconAndText>
		<Card
			className="mt-6"
			title="Enroll New Device"
			body="Choose this option if you wish to enroll a new device with this password manager."
		>
			{#snippet button()}
				<Button href="{base}/enroll">Enroll</Button>
			{/snippet}
		</Card>

		<Card
			className="mt-10"
			title="Register as New User"
			body="Choose this option to register as a first-time user with this password manager."
		>
			{#snippet button()}
				<Button href="{base}/register">Register</Button>
			{/snippet}
		</Card>
	</div>
</div>

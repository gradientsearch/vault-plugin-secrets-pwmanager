<script>
	import Title from '../../components/title.svelte';
	import VaultIconAndText from '../../components/vaultIconAndText.svelte';
	import Form from './components/form.svelte';
</script>

<div class="flex h-full w-full justify-center">
	<div class="flex-row">
        
		<VaultIconAndText className="mt-14"></VaultIconAndText>
		<Title className="mb-2 mt-8">Register</Title>
        <Form></Form>
	</div>
</div>

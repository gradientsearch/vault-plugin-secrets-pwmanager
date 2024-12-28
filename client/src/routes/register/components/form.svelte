<script lang="ts">
	import Button from '../../../components/button.svelte';
	import CardContainer from '../../../components/cardContainer.svelte';

	class Register {
		mount: string = ""
		url: string = ""
		token: string = ""
		password: string = ""
		retypedPassword: string = ""
	}

    let register = new Register()

	function onRegister() {}
</script>

<div class="flex h-full w-full justify-center">
	<div class="flex-row">
		<CardContainer className="overflow-">
			<div class="text-md grid min-w-96 grid-cols-1 gap-6">
				<label class="block">
					<span class="text-gray-700">Mount</span>
					<input
						type="text"
						multiple
						class="form-input mt-1 block w-full"
						placeholder="pwmanager"
						bind:value={register.mount}
					/>
				</label>
				<label class="block">
					<span class="text-gray-700">Vault Address</span>
					<input
						type="url"
						multiple
						class="form-input mt-1 block w-full"
						placeholder="https://vault.example.com:8200"
						bind:value={register.url}
					/>
				</label>
				<label class="block">
					<span class="text-gray-700">Vault token</span>
					<input type="password" multiple class="form-input mt-1 block w-full" bind:value={register.token} />
				</label>
				<label class="block">
					<span class="text-gray-700">Password Manager Password</span>
					<input
						type="password"
						multiple
						class="form-input mt-1 block w-full"
						bind:value={register.password}
					/>
				</label>

				<label class="block">
					<span class="text-gray-700">Retype Password Manager Password</span>
					<input
						type="password"
						multiple
						class="form-input mt-1 block w-full"
						bind:value={register.retypedPassword}
					/>
				</label>
				<Button fn={onRegister}>Register</Button>
			</div>
		</CardContainer>
        <article class="max-w-sm text-wrap text-bre">
        <p class="break-all">{JSON.stringify(register)} </p>
    </article>
	</div>
    
</div>


<style>
	label {
		color: #1f2124;
		font-size: 14px;
		font-weight: 700;
		align-items: center;
		gap: 4px;
		width: min-content;
		min-width: 100%;
	}
</style>

<script lang="ts">
	import { buildUUK, bytesToHex } from '$lib/uuk';
	import Button from '../../../components/button.svelte';
	import CardContainer from '../../../components/cardContainer.svelte';
	import { onMount } from 'svelte';
	import Password from '../../unlocked/components/password.svelte';
	import { Api } from '$lib/api';

	class Register {
		mount: string = 'pwmanager';
		url: string = 'http://localhost:8200';
		token: string = '';
		password: string = '';
		retypedPassword: string = '';
		secretKey: Uint8Array = new Uint8Array();
	}

	let errorText: string | undefined = undefined;

	let register = new Register();
	let encoder = new TextEncoder();
	let isRegistering = false;
	async function onRegister() {
		isRegistering = true;
		let randomSeq = bytesToHex(crypto.getRandomValues(new Uint8Array(26)));
		let secretKey = `H1-${register.mount}-${randomSeq}`; // combination Secret ID - secret

		let api = new Api(register.token, register.url, register.mount);
		let tokenInfo = await api.tokenLookup();
		let entityID = tokenInfo['data']['entity_id'];
		let uuk = await buildUUK(
			encoder.encode(register.password),
			encoder.encode(register.mount),
			encoder.encode(secretKey),
			encoder.encode(entityID)
		);

		let err = await api.register(uuk);
		if (err != undefined) {
			errorText = err.message
		}

		isRegistering = false;
	}
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
					<input
						type="password"
						multiple
						class="form-input mt-1 block w-full"
						bind:value={register.token}
					/>
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

				{#if errorText}
					<p class="text-sm text-foreground_critical">{errorText}</p>
				{/if}

				{#if !isRegistering}
					<Button disabled={isRegistering} fn={onRegister}>Register</Button>
				{:else}
					<div class="flex items-center justify-center">
						<svg
							class="text-white h-3 w-3 animate-spin"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
						>
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
							></circle>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							></path>
						</svg>
						<span class="ml-1 text-sm">Processing</span>
					</div>
				{/if}
			</div>
		</CardContainer>
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

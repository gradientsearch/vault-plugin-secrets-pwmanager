<script lang="ts">
	import { buildUUK, bytesToHex } from '$lib/uuk';
	import Button from '../../../components/button.svelte';
	import CardContainer from '../../../components/cardContainer.svelte';
	import { onMount } from 'svelte';
	import Password from '../../unlocked/components/password.svelte';

	class Register {
		mount: string = 'pwmanager';
		url: string = 'http://localhost:8200';
		token: string = '';
		password: string = '';
		retypedPassword: string = '';
		secretKey: Uint8Array = new Uint8Array();
	}

	let register = new Register();
	let encoder = new TextEncoder();
	let isRegistering = false;
	async function onRegister() {
		isRegistering = true;
		let randomSeq = bytesToHex(crypto.getRandomValues(new Uint8Array(26)));
		let secretKey = `H1-${register.mount}-${randomSeq}`; // combination Secret ID - secret
		// TODO grab entity-id from token
		let uuk = await buildUUK(
			encoder.encode(register.password),
			encoder.encode(register.mount),
			encoder.encode(secretKey),
			new Uint8Array([1, 2, 3])
		);
		console.log(JSON.stringify(uuk));
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
				<Button fn={onRegister}>Register</Button>
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
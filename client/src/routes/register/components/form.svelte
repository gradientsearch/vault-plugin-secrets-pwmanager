<script lang="ts">
	// TODO refactor this to follow a mvvm pattern
	import { buildUUK, bytesToHex } from '$lib/uuk';
	import Button from '../../../components/button.svelte';
	import CardContainer from '../../../components/cardContainer.svelte';

	import { Api } from '$lib/api';
	import Title from '../../../components/title.svelte';
	import VaultIconAndText from '../../../components/vaultIconAndText.svelte';

	class Register {
		mount: string = 'pwmanager';
		url: string = 'http://localhost:8200';
		token: string = '';
		password: string = '';
		retypedPassword: string = '';
		secretKey: string = '';
		secretKeyDisplay: string = '';
	}

	let errorText: string | undefined = undefined;

	let register = new Register();
	let encoder = new TextEncoder();
	let isRegistering = false;
	let registered = false;
	async function onRegister() {
		if (register.password != register.retypedPassword) {
			errorText = 'passwords do not match';
			return;
		}
		isRegistering = true;
		let randomSeq = bytesToHex(crypto.getRandomValues(new Uint8Array(15)));
		let secretKeyDisplay = randomSeq.replace(/(.{6})/g, '$1-').slice(0, -1);
		register.secretKey = `H1-${register.mount}-${randomSeq}`; // combination Secret ID - secret
		register.secretKeyDisplay = `H1-${register.mount}-${secretKeyDisplay}`;

		let api = new Api(register.token, register.url, register.mount);
		let tokenInfo = await api.tokenLookup();
		let entityID = tokenInfo['data']['entity_id'];
		let uuk = await buildUUK(
			encoder.encode(register.password),
			encoder.encode(register.mount),
			encoder.encode(register.secretKey),
			encoder.encode(entityID)
		);

		let err = await api.register(uuk);
		if (err != undefined) {
			errorText = err.message;
		} else {
			registered = true;
		}

		isRegistering = false;
		let encSecretKey = bytesToHex(new TextEncoder().encode(register.secretKey));
		localStorage.setItem('secretkey', encSecretKey);
	}
</script>

{#if !registered}
	<div class="flex h-full w-full justify-center">
		<div class="flex-row">
			<VaultIconAndText className="mt-14"></VaultIconAndText>
			<Title className="mb-2 mt-8">Register</Title>
			<CardContainer className="overflow-rounded-3xl">
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
{:else}
	<div class="flex h-full w-full justify-center">
		<div class="max-w-2xl flex-row gap-3">
			<VaultIconAndText className="mt-14"></VaultIconAndText>

			<div
				class="mt-6 rounded-3xl bg-red-200 p-3 text-center text-base font-bold text-foreground_high_contrast"
			>
				Important Notice: Print and Keep This Page Safe
			</div>

			<CardContainer className="bg-red-100 mt-10">
				<div class="text-md grid min-w-96 grid-cols-1 gap-6">
					<label class="block">
						<span class="text-gray-700">Mount</span>
						<input
							type="text"
							multiple
							disabled
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
							disabled
							class="form-input mt-1 block w-full"
							placeholder="https://vault.example.com:8200"
							bind:value={register.url}
						/>
					</label>
					<label class="block">
						<span class="text-gray-700">Secret Key</span>
						<input
							type="text"
							multiple
							disabled
							class="form-input mt-1 block w-full"
							bind:value={register.secretKeyDisplay}
						/>
					</label>

					<label class="block">
						<span class="text-gray-700">Password Manager Password</span>
						<input type="text" multiple class="form-input mt-1 block w-full" />
					</label>
				</div></CardContainer
			>

			<p class="mt-3 text-sm">
				This page contains your Secret Key and a textbox to fill in your Password Manager Secret,
				both of which are required to unlock your passwords and secrets. Without these, there is no
				way to access your vault.
			</p>

			<p class="text-sm">
				To ensure you don't lose access, we highly recommend that you print this page and store it
				in a secure place. Do not share your Secret Key, Password Manager Secret, or password with
				anyone you do not trust.
			</p>
			<p class="text-sm">Store it safely! You are the only one who can unlock your vault.</p>
		</div>
	</div>
{/if}

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

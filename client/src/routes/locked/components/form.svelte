<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { Api } from '$lib/api';
	import {  decryptEncPriKey } from '$lib/uuk';
	import Button from '../../../components/button.svelte';
	import CardContainer from '../../../components/cardContainer.svelte';
	import Title from '../../../components/title.svelte';

	import IconAndText from '../../../components/iconAndText.svelte';
	import { storedKeyPair, type KeyPair } from '$lib/asym_key_store';
	import { hexToBytes } from '$lib/helper';

    
	class SignIn {
		mount: string = 'pwmanager';
		url: string = 'http://localhost:8200';
		token: string = '';
		password: string = '';
	}

	let signIn = new SignIn();
	let errorText: string | undefined;

	let isSigningIn = false;
    
	async function onSignIn() {
        
		isSigningIn = true;
		let secretKeyHex = localStorage.getItem('secretkey');

		if (secretKeyHex == null) {
			goto(`${base}/`);
			return;
		}

		let secretKey = new TextDecoder().decode(hexToBytes(secretKeyHex));
		
		let api = new Api(signIn.token, signIn.url, signIn.mount);
		let tokenInfo = await api.tokenLookup();
		let entityID = tokenInfo['data']['entity_id'];

		// TODO update when using different auth methods
		let username = tokenInfo['data']['meta']['username']

		let [uuk,err] = await api.uuk(entityID);
		if (err != undefined || uuk == undefined) {
			errorText = 'error retrieving UUK';
			isSigningIn = false;
			return;
		}

        let encoder = new TextEncoder();
		encoder.encode(signIn.password);
		encoder.encode(signIn.mount);
		encoder.encode(secretKey);
		encoder.encode(entityID);

		let [prikey, pubkey] = await decryptEncPriKey(
			uuk,
			encoder.encode(signIn.password),
			encoder.encode(signIn.mount),
			encoder.encode(secretKey),
			encoder.encode(entityID)
		);

		let keypair: KeyPair = {
			PriKey: prikey,
			PubKey: pubkey
		};

		storedKeyPair.set(keypair);
		localStorage.setItem('loginInfo', JSON.stringify({token: signIn.token, url: signIn.url, mount: signIn.mount, entityID: entityID, username: username}))
		
		goto(`${base}/unlocked`);
		isSigningIn = false;
	}
</script>

<div class="flex h-full w-full justify-center">
	<div class="flex-row">
		<IconAndText className="mt-14"></IconAndText>
		<Title className="mb-2 mt-8">Sign in to PwManger</Title>
		<CardContainer className="overflow-rounded-3xl">
			<div class="text-md grid min-w-96 grid-cols-1 gap-6">
				<label class="block">
					<span class="text-gray-700">Mount</span>
					<input
						type="text"
						multiple
						class="form-input mt-1 block w-full"
						placeholder="pwmanager"
						bind:value={signIn.mount}
					/>
				</label>
				<label class="block">
					<span class="text-gray-700">Vault Address</span>
					<input
						type="url"
						multiple
						class="form-input mt-1 block w-full"
						placeholder="https://vault.example.com:8200"
						bind:value={signIn.url}
					/>
				</label>
				<label class="block">
					<span class="text-gray-700">Vault token</span>
					<input
						type="password"
						multiple
						class="form-input mt-1 block w-full"
						bind:value={signIn.token}
					/>
				</label>
				<label class="block">
					<span class="text-gray-700">Password Manager Password</span>
					<input
						type="password"
						multiple
						class="form-input mt-1 block w-full"
						bind:value={signIn.password}
					/>
				</label>

				{#if errorText}
					<p class="text-sm text-foreground_critical">{errorText}</p>
				{/if}

				{#if !isSigningIn}
					<Button disabled={isSigningIn} fn={onSignIn}>Sign in</Button>
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

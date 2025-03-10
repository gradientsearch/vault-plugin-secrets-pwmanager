<script lang="ts">
	import { onMount } from 'svelte';
	import Button from '../../../../components/button.svelte';
	import { KVBundleService, type BundleService } from '../../services/bundle.service';

	let {
		bundle = $bindable(),
		zarf = $bindable(),
		cancel = $bindable(),
		save = () => {}
	} = $props();

	interface user {
		name: string;
		isAdmin: boolean;
		caps: {};
	}

	let caps = {
		create: true,
		read: true,
		update: true,
		patch: true,
		delete: true,
		list: true
	};

	async function onSave() {
		let bundleUsers: BundleUser[] = [];
		users.forEach((u) => {
			let ucap:any = [];
			Object.keys(u.Capabilities).forEach((c) => {
				if (u.Capabilities[c]) {
					ucap.push(c);
				}
			});
			u.Capabilities = ucap.join(',');
			bundleUsers.push(u);
		});

		let [pubkeys, err] = await KVBundleService.updateSharedBundleUsers(
			zarf,
			bundle.Owner,
			bundle.ID,
			bundleUsers
		);
		if (err !== undefined) {
			console.log(`error updating bundle users: ${err}`);
			return;
		}

		console.log(pubkeys);

		//update bundle metadata
		save();
	}

	function onCancel() {
		cancel();
	}

	let bundleUser = {
		EntityName: '',
		EntityId: '',
		Capabilities: {},
		IsAdmin: false
	};

	let bundleName = $state('');
	let bu: BundleUser = $state(Object.assign({}, bundleUser));

	let users: BundleUser[] = $state([]);
	onMount(() => {
		bundleName = bundle?.Name;
		console.log(bundle);

		users = [...bundle.Users];
	});
</script>

<div class="flex flex-row p-4">
	<div class="flex flex-row">
		<div class="flex flex-col">
			<div class="text-md grid min-w-96 grid-cols-1 gap-6">
				<label class="block">
					<span class="text-gray-700 text-base font-medium">Name</span>
					<input
						type="text"
						multiple
						class="form-input mt-1 block w-full"
						placeholder="pwmanager"
						bind:value={bundleName}
					/>
				</label>
			</div>

			<div class="flex flex-row">
				<div class="text-md grid min-w-96 grid-cols-1 gap-6">
					<label class="block">
						<span class="text-gray-700 text-base font-medium">Users</span>

						<div class="flex flex-row">
							<input
								type="text"
								multiple
								class="form-input mt-1 block w-full"
								placeholder="bob"
								bind:value={bu.EntityName}
							/>
							<div class="flex flex-row">
								{#each Object.keys(caps) as c}
									<label class="flex cursor-pointer items-start ps-4">
										<div class="flex items-center pe-1">
											&#8203;&#8203;
											<input
												onclick={() => {
													bu.Capabilities[c] = !bu.Capabilities[c];
													if (bu.Capabilities[c] === false) {
														bu.IsAdmin = false;
													}
												}}
												bind:checked={bu.Capabilities[c]}
												type="checkbox"
												class="border-gray-300 rounded-sm"
											/>
										</div>

										<div>
											<strong class="text-gray-900 pe-2 text-center text-lg font-medium">
												{c}</strong
											>
										</div>
									</label>
								{/each}
								<label class="flex cursor-pointer items-start ps-4">
									<div class="flex items-center pe-1">
										&#8203;&#8203;
										<input
											onclick={() => {
												if (!bu.IsAdmin) {
													Object.keys(caps).forEach((c) => {
														bu.Capabilities[c] = true;
													});
												}
												bu.IsAdmin = !bu.IsAdmin;
											}}
											bind:checked={bu.IsAdmin}
											type="checkbox"
											class="border-gray-300 rounded-sm"
										/>
									</div>

									<div>
										<strong class="text-gray-900 pe-2 text-lg font-medium">admin</strong>
									</div>
								</label>
							</div>

							<div class="px-2">
								<Button
									margin="m-0"
									fn={() => {
										users.push(bu);
										bu = Object.assign({}, bundleUser);
									}}>Add</Button
								>
							</div>
						</div>
					</label>

					{#each users as u, idx}
						<div class="flex flex-row">
							<input
								type="text"
								multiple
								disabled
								class="form-input mt-1 block w-full border-none"
								placeholder="bob"
								bind:value={u.EntityName}
							/>
							<div class="flex flex-row">
								{#each Object.keys(caps) as c}
									<label class="flex cursor-pointer items-start ps-4">
										<div class="flex items-center pe-1">
											&#8203;&#8203;
											<input
												onclick={() => {
													u.Capabilities[c] = !u.Capabilities[c];
													if (u.Capabilities[c] === false) {
														u.IsAdmin = false;
													}
												}}
												bind:checked={u.Capabilities[c]}
												type="checkbox"
												class="border-gray-300 rounded-sm"
											/>
										</div>

										<div>
											<strong class="text-gray-900 pe-2 text-center text-lg font-medium">
												{c}</strong
											>
										</div>
									</label>
								{/each}
								<label class="flex cursor-pointer items-start ps-4">
									<div class="flex items-center pe-1">
										&#8203;&#8203;
										<input
											onclick={() => {
												if (!u.IsAdmin) {
													Object.keys(caps).forEach((c) => {
														u.Capabilities[c] = true;
													});
												}
												u.IsAdmin = !u.IsAdmin;
											}}
											bind:checked={u.IsAdmin}
											type="checkbox"
											class="border-gray-300 rounded-sm"
										/>
									</div>

									<div>
										<strong class="text-gray-900 pe-2 text-lg font-medium">admin</strong>
									</div>
								</label>
							</div>

							<div class="px-2">
								<Button
									margin="m-0"
									primary={false}
									fn={() => {
										users.splice(idx, 1);
										users = users;
									}}>delete {idx}</Button
								>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>

<span class="flex flex-1"></span>

<div class="p-4">
	<Button fn={onSave}>Save</Button>
	<Button fn={onCancel} primary={false}>Cancel</Button>
</div>

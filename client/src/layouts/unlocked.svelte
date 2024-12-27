<script lang="ts">
	import Header from "../components/unlocked/header.svelte";
	import Password from "../components/unlocked/password.svelte";
	import Vault from "../components/unlocked/vault.svelte";

	let bg = ['bg-sky-500', 'bg-slate-500', 'bg-pink-400', 'bg-green-400'];
</script>

<div class="flex h-full">
	<div class="{bg[0]} w-full md:bg-green-500 md:max-w-64 lg:block hidden  content-center text-center">(sidebar)</div>
	<div class="h-full w-full flex-col">
		<Header></Header>
		<div class="flex h-[calc(100vh-48px)] w-full">
			<Vault></Vault>
			<Password></Password>
		</div>
	</div>
</div>
